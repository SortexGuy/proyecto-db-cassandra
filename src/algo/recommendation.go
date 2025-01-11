package algo

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/graph/simple"
)

// Nodo personalizado que implementa la interfaz graph.Node
type Node struct {
	id int64
}

// Método ID para cumplir con graph.Node
func (n Node) ID() int64 {
	return n.id
}

type Item struct {
	ID    int64
	Score float64
}

func HybridRecommendation(g *simple.UndirectedGraph, initialUserID int64, steps int, lambda float64) []Item {
	// Obtener los resultados de las dos estrategias
	heat, max_heat := heatPropagation(g, initialUserID, steps)
	resources, max_resources := resourceAllocation(g, initialUserID, steps)

	// Normalizar ambos resultados y combinar
	hybridScores := make(map[int64]float64)
	nodes := g.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		if node.ID() >= 0 { // Solo considerar objetos (IDs positivos)
			heatValue, heatExists := heat[node.ID()]
			if heatExists && !heatValue.previousInteraction && heatValue.value != 0 {
				hybridScores[node.ID()] = lambda * (heatValue.value / max_heat)
			}
			resourceValue, resourceExists := resources[node.ID()]
			if resourceExists && !resourceValue.previousInteraction && resourceValue.value != 0 {
				hybridScores[node.ID()] += (1 - lambda) * (resourceValue.value / max_resources)
			}
		}

	}
	return sortMapDescending(hybridScores)
}

type ResourceInfo struct {
	value               float64
	previousInteraction bool
}

func heatPropagation(g *simple.UndirectedGraph, initialUserID int64, steps int) (map[int64]*ResourceInfo, float64) {
	// Inicializar recursos en los objetos conectados al usuario inicial
	heat := initializeMap(g, initialUserID)
	degrees := calculateDegrees(g)
	max := 0.0
	// Propagación del calor por los pasos indicados
	for step := 0; step < steps; step++ {
		// Paso 1: Pelicula → Usuario
		userHeat := make(map[int64]*ResourceInfo)
		nodes := g.Nodes()
		for nodes.Next() {
			node := nodes.Node()
			if node.ID() < 0 { // Usuarios (IDs negativos)
				sum := 0.0
				movies := g.From(node.ID())
				for movies.Next() {
					movie := movies.Node()
					sum += heat[movie.ID()].value / float64(degrees[node.ID()])
				}
				userHeat[node.ID()] = &ResourceInfo{
					value:               sum,
					previousInteraction: false,
				}
			}
		}

		// Paso 2: Usuario → Pelicula
		newHeat := make(map[int64]*ResourceInfo)
		nodes = g.Nodes()
		for nodes.Next() {
			node := nodes.Node()
			if node.ID() >= 0 { // Peliculas (IDs positivos)
				sum := 0.0
				users := g.From(node.ID())
				for users.Next() {
					user := users.Node()
					sum += userHeat[user.ID()].value / float64(degrees[node.ID()])
				}
				newHeat[node.ID()] = &ResourceInfo{
					value:               sum,
					previousInteraction: heat[node.ID()].previousInteraction,
				}
				if step == steps-1 && sum > max {
					max = sum
				}
			} else {
				newHeat[node.ID()] = &ResourceInfo{
					value:               0.0,
					previousInteraction: false,
				}
			}
		}

		heat = newHeat
	}

	return heat, max
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

func resourceAllocation(g *simple.UndirectedGraph, initialUserID int64, steps int) (map[int64]*ResourceInfo, float64) {
	// Inicializar recursos en los objetos conectados al usuario inicial
	resources := initializeMap(g, initialUserID)
	degrees := calculateDegrees(g)
	max := 0.0
	// Propagación de recursos por los pasos indicados
	for step := 0; step < steps; step++ {
		// Paso 1: Distribuir recursos de peliculas a usuarios
		userResources := make(map[int64]*ResourceInfo)
		nodes := g.Nodes()
		for nodes.Next() {
			node := nodes.Node()
			if node.ID() < 0 { // ID negativo para usuarios
				sum := 0.0
				movies := g.From(node.ID())
				for movies.Next() {
					movie := movies.Node()
					sum += resources[movie.ID()].value / float64(degrees[movie.ID()])
				}
				userResources[node.ID()] = &ResourceInfo{
					value:               sum,
					previousInteraction: false,
				}
			}
		}

		// Paso 2: Distribuir recursos de usuarios a objetos
		newResources := make(map[int64]*ResourceInfo)
		nodes = g.Nodes()
		for nodes.Next() {
			node := nodes.Node()
			if node.ID() >= 0 { // ID positivo para objetos
				sum := 0.0
				users := g.From(node.ID())
				for users.Next() {
					user := users.Node()
					sum += userResources[user.ID()].value / float64(degrees[user.ID()])
				}
				newResources[node.ID()] = &ResourceInfo{
					value:               sum,
					previousInteraction: resources[node.ID()].previousInteraction,
				}
				if step == steps-1 && sum > max {
					max = sum
				}
			} else {
				newResources[node.ID()] = &ResourceInfo{
					value:               0.0,
					previousInteraction: false,
				}
			}
		}

		resources = newResources
	}

	return resources, max
}

func sortMapDescending(hybridScores map[int64]float64) []Item {
	items := make([]Item, 0, len(hybridScores))
	for id, score := range hybridScores {
		items = append(items, Item{ID: id, Score: score})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Score > items[j].Score // Orden descendente
	})

	return items
}

type UserMovie struct {
	UserID  int64
	MovieID int64
}

func createGraph(users []int64, movies []int64, relations []UserMovie) *simple.UndirectedGraph {
	g := simple.NewUndirectedGraph()

	// Mapas para buscar nodos por ID
	nodes := make(map[int64]Node)

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

/*func test() {

	g := createGraph(users, movies, relations)

	resources := HybridRecommendation(g, -1, 1, 0.5)

	// Imprimir resultados
	fmt.Println("Distribución de recursos:")

	for _, item := range resources {
		fmt.Printf("ID: %d, Score: %.7f\n", item.ID, item.Score)

	}

}*/
