package textrank

type edge struct {
	start int
	end   int
	score float64
}

type graph struct {
	edges     []edge
	nodeCount int
}

func NewGraphForSentences(sentences []sentence) graph {
	var g graph
	g.nodeCount = len(sentences)
	g.edges = g.computeEdges(sentences)
	return g
}

func (g *graph) PruneGraph() (graph, map[int]int) {
	neighbourCount := g.computeNeighborhood()
	mapOldtoNew := computeNodeMapping(neighbourCount)
	return graph{
		nodeCount: len(mapOldtoNew),
		edges:     mapEdges(mapOldtoNew, g.edges),
	}, mapOldtoNew
}

func (g *graph) computeEdges(sentences []sentence) []edge {
	output := make([]edge, 0, g.nodeCount*g.nodeCount/2)
	for i := 0; i < g.nodeCount; i++ {
		for j := i + 1; j < g.nodeCount; j++ {
			similarity := getSimilarity(sentences[i], sentences[j])
			if similarity != 0 {
				output = append(output, edge{i, j, similarity})
			}
		}
	}
	return output
}

func (g *graph) computeNeighborhood() map[int]int {
	output := make(map[int]int, g.nodeCount)
	for _, edge := range g.edges {
		output[edge.start]++
		output[edge.end]++
	}
	return output
}

func computeNodeMapping(neighborhood map[int]int) map[int]int {
	mapOldtoNew := make(map[int]int)
	newIndex := 0
	for i, neighborCount := range neighborhood {
		if neighborCount > 0 {
			mapOldtoNew[i] = newIndex
			newIndex++
		}
	}
	return mapOldtoNew
}

func mapEdges(mapOldtoNew map[int]int, edges []edge) []edge {
	output := make([]edge, len(edges))
	for i, e := range edges {
		output[i] = edge{mapOldtoNew[e.start], mapOldtoNew[e.end], e.score}
	}
	return output
}
