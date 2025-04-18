package sga

import (
	"math/rand/v2"
	"sort"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type Algorithm struct {
	algos.GeneticAlgorithm
	params         Params
	generations    int
	population     []problems.Solution
	eliteSize      int
	matingPoolSize int
}

type Step struct {
	algos.GeneticAlgorithmStep
	Generation int `json:"generation"`
}

func NewAlgorithm(problem problems.Problem, timeout time.Duration, params Params, logger algos.ProgressLoggerProvider) *Algorithm {
	return &Algorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, timeout, logger),
		params:           params,
		generations:      0,
		eliteSize:        int(float64(params.PopulationSize) * params.ElitePercentile),
		matingPoolSize:   int(float64(params.PopulationSize) * params.MatingPoolPercentile),
	}
}

func (alg *Algorithm) Run() {
	alg.InitPopulation()
	fitness := 0.0
	for time.Since(alg.StartTimestamp) < alg.Timeout {
		alg.Evolve()
		alg.generations++
		bestFitness := alg.Solution.Fitness()
		if fitness != bestFitness {
			fitness = bestFitness
			var pareto [][]float64
			for _, solution := range alg.population {
				pareto = append(pareto, solution.Objectives())
			}
			alg.LogStep(Step{
				GeneticAlgorithmStep: algos.GeneticAlgorithmStep{Elapsed: time.Since(alg.StartTimestamp), Solution: alg.Solution, ParetoFront: pareto},
				Generation:           alg.generations,
			})
		}
	}
}

func (alg *Algorithm) InitPopulation() {
	pop := make([]problems.Solution, alg.params.PopulationSize)
	for i := range pop {
		pop[i] = alg.Problem.RandomSolution()
	}
	alg.population = pop
}

func (alg *Algorithm) Evolve() {
	alg.evaluateGeneration()

	newPopulation := make([]problems.Solution, 0, alg.params.PopulationSize)

	// perform elitism
	newPopulation = append(newPopulation, alg.population[:alg.eliteSize]...)

	// generate rest of the population from mating pool
	for len(newPopulation) < alg.params.PopulationSize {
		p1Ind := rand.IntN(alg.matingPoolSize)
		p2Ind := rand.IntN(alg.matingPoolSize)
		if p1Ind == p2Ind {
			continue
		}
		parent1 := alg.population[p1Ind]
		parent2 := alg.population[p2Ind]
		children := parent1.Crossover(parent2)
		for i := range children {
			children[i] = children[i].Mutate(alg.params.MutationRate)
		}
		newPopulation = append(newPopulation, children...)
	}

	alg.population = newPopulation
}

func (alg *Algorithm) evaluateGeneration() {
	sort.Slice(alg.population, func(i, j int) bool {
		return alg.population[i].Fitness() < alg.population[j].Fitness()
	})
	alg.Solution = alg.population[0]
}
