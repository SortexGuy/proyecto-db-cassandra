package recommendations

import (
	"container/heap"
	"fmt"

	"github.com/SortexGuy/proyecto-db-cassandra/src/movies"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

//"github.com/SortexGuy/proyecto-db-cassandra/src/movies"

type Recommendation struct {
	UserID int64   `json:"user_id"`
	Movies []int64 `json:"recommendations"`
}

// Nodo personalizado que implementa la interfaz graph.Node
type Node struct {
	id int64
}

// Método ID para cumplir con graph.Node
func (n Node) ID() int64 {
	return n.id
}

func CreateGraph(users []int64, movies []int64, relations []movies.MovieByUser) *simple.UndirectedGraph {
	g := simple.NewUndirectedGraph()

	// Mapas para buscar nodos por ID
	//nodes := make(map[int64]Node)
	nodes := make(map[int64]graph.Node)
	for _, user_id := range users {
		user_id = -1 * user_id
		n := Node{id: user_id}
		nodes[user_id] = n
		g.AddNode(n)
	}

	for _, movie_id := range movies {
		n := Node{id: movie_id}
		nodes[movie_id] = n
		g.AddNode(n)
	}

	// Agrega aristas (relaciones)

	for _, relation := range relations {
		userID := -1 * relation.UserID
		movieID := relation.MovieID
		if userNode, ok := nodes[userID]; ok {
			if movieNode, ok := nodes[movieID]; ok {

				g.SetEdge(simple.Edge{F: userNode, T: movieNode})

			} else {
				fmt.Printf("Advertencia: Película con ID %d no encontrada.\n", movieID)
			}
		} else {
			fmt.Printf("Advertencia: Usuario con ID %d no encontrado.\n", userID)
		}
	}

	return g
}

type Item struct {
	ID    int64
	Score float64
}

func HybridRecommendation(g *simple.UndirectedGraph, initialUserID int64, steps int, lambda float64, num_recommendations int) Recommendation {
	resources, maxResources, heat, maxHeat := propagate(g, initialUserID, steps)
	// Normalizar ambos resultados y combinar
	hybridScores := make(map[int64]float64)
	nodes := g.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		if node.ID() >= 0 { // Solo considerar objetos (IDs positivos)
			heatValue, heatExists := heat[node.ID()]
			if heatExists && !heatValue.previousInteraction && heatValue.value != 0 {
				hybridScores[node.ID()] = lambda * (heatValue.value / maxHeat)
			}
			resourceValue, resourceExists := resources[node.ID()]
			if resourceExists && !resourceValue.previousInteraction && resourceValue.value != 0 {
				hybridScores[node.ID()] += (1 - lambda) * (resourceValue.value / maxResources)
			}
		}

	}
	var recommendations Recommendation

	// Asigna valores a los campos correctamente
	recommendations.UserID = initialUserID
	recommendations.Movies = sortMapDescending(hybridScores, num_recommendations)

	return recommendations
}

type ResourceInfo struct {
	value               float64
	previousInteraction bool
}

func propagate(g *simple.UndirectedGraph, initialUserID int64, steps int) (map[int64]*ResourceInfo, float64, map[int64]*ResourceInfo, float64) {
	// Inicializar recursos en los objetos conectados al usuario inicial
	resources := initializeMap(g, initialUserID)
	heat := initializeMap(g, initialUserID)
	degrees := calculateDegrees(g)
	inverseDegrees := calculateInverseDegree(degrees)
	maxResources, maxHeat := 0.0, 0.0

	userResources, userHeat := make(map[int64]*ResourceInfo), make(map[int64]*ResourceInfo)
	newResources, newHeat := make(map[int64]*ResourceInfo), make(map[int64]*ResourceInfo)

	for step := 0; step < steps; step++ {
		// Paso 1: Distribuir recursos de peliculas a usuarios
		nodes := g.Nodes()
		for nodes.Next() {
			nodeID := nodes.Node().ID()
			if nodeID < 0 { // ID negativo para usuarios
				sumResources, sumHeat := 0.0, 0.0
				movies := g.From(nodeID)
				for movies.Next() {
					movieID := movies.Node().ID()
					sumResources += resources[movieID].value * inverseDegrees[movieID]
					sumHeat += heat[movieID].value * inverseDegrees[nodeID]
				}

				userResources[nodeID] = &ResourceInfo{
					value:               sumResources,
					previousInteraction: false,
				}
				userHeat[nodeID] = &ResourceInfo{
					value:               sumHeat,
					previousInteraction: false,
				}

			}
		}

		// Paso 2: Distribuir recursos de usuarios a objetos
		nodes = g.Nodes()
		for nodes.Next() {
			nodeID := nodes.Node().ID()
			if nodeID >= 0 { // ID positivo para objetos
				sumResources, sumHeat := 0.0, 0.0
				users := g.From(nodeID)
				for users.Next() {
					userID := users.Node().ID()
					sumResources += userResources[userID].value * inverseDegrees[userID]
					sumHeat += userHeat[userID].value * inverseDegrees[nodeID]
				}
				newResources[nodeID] = &ResourceInfo{
					value:               sumResources,
					previousInteraction: resources[nodeID].previousInteraction,
				}
				newHeat[nodeID] = &ResourceInfo{
					value:               sumHeat,
					previousInteraction: heat[nodeID].previousInteraction,
				}

				if step == steps-1 {
					if sumResources > maxResources {
						maxResources = sumResources
					}
					if sumHeat > maxHeat {
						maxHeat = sumHeat
					}
				}

			} else {
				newResources[nodeID] = &ResourceInfo{
					value:               0.0,
					previousInteraction: false,
				}
				newHeat[nodeID] = &ResourceInfo{
					value:               0.0,
					previousInteraction: false,
				}
			}
		}

		resources = newResources
		heat = newHeat

		for k := range userResources {
			delete(userResources, k)
			delete(userHeat, k)
			delete(newResources, k)
			delete(newHeat, k)
		}
	}

	return resources, maxResources, heat, maxHeat
}
func initializeMap(g *simple.UndirectedGraph, initialUserID int64) map[int64]*ResourceInfo {
	temp := make(map[int64]*ResourceInfo)

	nodes := g.Nodes()
	for nodes.Next() {
		nodeID := nodes.Node().ID()
		temp[nodeID] = &ResourceInfo{value: 0, previousInteraction: false}

		// Verificar la vecindad dentro del mismo bucle
		if g.HasEdgeBetween(nodeID, initialUserID) {
			temp[nodeID].value = 1
			temp[nodeID].previousInteraction = true
		}
	}

	return temp
}

// Función para calcular el grado de cada nodo en el grafo
func calculateDegrees(g *simple.UndirectedGraph) map[int64]int {
	degrees := make(map[int64]int)
	nodes := g.Nodes()
	for nodes.Next() {
		node := nodes.Node()

		// Contar el número de nodos conectados
		count := 0
		connections := g.From(node.ID())
		for connections.Next() {
			count++
		}

		degrees[node.ID()] = count
	}
	return degrees
}

func calculateInverseDegree(degrees map[int64]int) map[int64]float64 {
	inverseDegrees := make(map[int64]float64)
	for nodeID, degree := range degrees {
		if degree > 0 { // Evitar división por cero
			inverseDegrees[nodeID] = 1.0 / float64(degree)
		} else {
			inverseDegrees[nodeID] = 0.0
		}
	}
	return inverseDegrees
}

type ItemHeap []Item

func (h ItemHeap) Len() int           { return len(h) }
func (h ItemHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h ItemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ItemHeap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}

func (h *ItemHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func sortMapDescending(hybridScores map[int64]float64, numRecommendations int) []int64 {
	// Usar un min-heap para mantener solo los numRecommendations mayores
	h := &ItemHeap{}
	heap.Init(h)

	for id, score := range hybridScores {
		if h.Len() < numRecommendations {
			heap.Push(h, Item{ID: id, Score: score})
		} else if score > (*h)[0].Score {
			heap.Pop(h)
			heap.Push(h, Item{ID: id, Score: score})
		}
	}

	result := make([]int64, numRecommendations)
	for i := numRecommendations - 1; i >= 0; i-- {
		result[i] = (heap.Pop(h).(Item)).ID
	}

	return result
}

/*
func test() {

	g := createGraph(users, movies, relations)
	resources := HybridRecommendation(g, -1, 1, 0.5, 10)

	fmt.Printf("Imprimir resultados\n")

	for _, item := range resources {
		fmt.Printf("ID: %d\n", item)
	}
}*/
