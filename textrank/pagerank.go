package textrank

import (
	"fmt"
	"math/cmplx"

	"gonum.org/v1/gonum/mat"
)

const defaultDamping = 0.85

func buildAdjacencyMatrix(numNodes int, edges []edge) *mat.Dense {
	m := mat.NewDense(numNodes, numNodes, nil)
	for _, edge := range edges {
		m.Set(edge.start, edge.end, edge.score)
		m.Set(edge.end, edge.start, edge.score)
	}

	// Calculate weighted adjacency
	output := mat.NewDense(numNodes, numNodes, nil)
	for i := 0; i < numNodes; i++ {
		sum := 0.0
		for k := 0; k < numNodes; k++ {
			sum += m.At(i, k)
		}
		if sum != 0.0 {
			for j := 0; j < numNodes; j++ {
				weighted := m.At(i, j) / sum
				output.Set(i, j, weighted)
			}
		}
	}

	// fmt.Println(output)
	return output
}

func buildProbMatrix(numNodes int) *mat.Dense {
	inv := float64(1 / float64(numNodes))
	m := mat.NewDense(numNodes, numNodes, nil)
	for i := 0; i < numNodes; i++ {
		for j := 0; j < numNodes; j++ {
			m.Set(i, j, inv)
		}
	}
	return m
}

func buildPageRankMatrix(numNodes int, edges []edge, damping float64) *mat.Dense {
	adj := buildAdjacencyMatrix(numNodes, edges)
	prob := buildProbMatrix(numNodes)
	var sadj mat.Dense
	sadj.Scale(damping, adj)
	var sprob mat.Dense
	sprob.Scale(1-damping, prob)
	var ret mat.Dense
	ret.Add(&sadj, &sprob)
	return &ret
}

func pageRank(numNodes int, edges []edge) ([]float64, error) {
	m := buildPageRankMatrix(numNodes, edges, defaultDamping)
	var eig mat.Eigen
	ok := eig.Factorize(m, mat.EigenLeft)
	if !ok {
		return nil, fmt.Errorf("unable to compute eigenvectors")
	}
	var ev mat.CDense
	eig.LeftVectorsTo(&ev)

	output := make([]float64, numNodes)
	for i := 0; i < numNodes; i++ {
		output[i] = cmplx.Abs(ev.At(i, 0))
	}

	return output, nil
}
