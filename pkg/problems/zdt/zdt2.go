package zdt

import (
	"math"
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// ZDT2Solution implements the ZDT2 problem.
// f1(x) = x₀
// g(x) = 1 + (9/(n-1)) * ∑₍ᵢ₌₁₎ⁿ₋₁ xᵢ
// f2(x) = g * (1 - (x₀/g)²)
type ZDT2Solution struct {
	Dimensions       int       `json:"dimensions"`
	X                []float64 `json:"x"`
	CachedObjectives []float64 `json:"objectives"`
}

// RandomZDT2Solution creates a random solution with all decision variables in [0, 1].
func RandomZDT2Solution(dimensions int) problems.Solution {
	x := make([]float64, dimensions)
	for i := range x {
		x[i] = rand.Float64()
	}
	return &ZDT2Solution{Dimensions: dimensions, X: x}
}

func (s *ZDT2Solution) Objectives() []float64 {
	if len(s.CachedObjectives) > 0 {
		return s.CachedObjectives
	}
	n := float64(s.Dimensions)
	f1 := s.X[0]
	sum := 0.0
	for i := 1; i < s.Dimensions; i++ {
		sum += s.X[i]
	}
	g := 1.0 + (9.0/(n-1))*sum
	f2 := g * (1 - math.Pow(f1/g, 2))
	s.CachedObjectives = []float64{f1, f2}
	return s.CachedObjectives
}

func (s *ZDT2Solution) Fitness() float64 {
	obj := s.Objectives()
	return obj[0] + obj[1]
}

func (s *ZDT2Solution) Crossover(other problems.Solution) []problems.Solution {
	otherSol, ok := other.(*ZDT2Solution)
	if !ok || s.Dimensions != otherSol.Dimensions {
		return []problems.Solution{s}
	}
	child1 := make([]float64, s.Dimensions)
	child2 := make([]float64, s.Dimensions)
	for i := range s.Dimensions {
		if rand.Float64() < 0.5 {
			child1[i] = s.X[i]
			child2[i] = otherSol.X[i]
		} else {
			child1[i] = otherSol.X[i]
			child2[i] = s.X[i]
		}
	}
	return []problems.Solution{
		&ZDT2Solution{Dimensions: s.Dimensions, X: child1},
		&ZDT2Solution{Dimensions: s.Dimensions, X: child2},
	}
}

func (s *ZDT2Solution) Mutate(rate float64) problems.Solution {
	mutant := make([]float64, s.Dimensions)
	copy(mutant, s.X)
	for i := range s.Dimensions {
		if rand.Float64() < rate {
			delta := rand.NormFloat64() * 0.1
			mutant[i] += delta
			if mutant[i] < 0 {
				mutant[i] = 0
			}
			if mutant[i] > 1 {
				mutant[i] = 1
			}
		}
	}
	return &ZDT2Solution{Dimensions: s.Dimensions, X: mutant}
}
