package main

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/graph/simple"
)

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

func initialize_map(g *simple.UndirectedGraph, initialUserID int64) map[int64]*ResourceInfo {
	temp := make(map[int64]*ResourceInfo)

	// Inicializar todos los nodos en el mapa
	nodes := g.Nodes()
	for nodes.Next() {
		node := nodes.Node()
		temp[node.ID()] = &ResourceInfo{
			value:               0,
			previousInteraction: false,
		}
	}

	// Asignar valores iniciales a los vecinos del nodo inicial
	neighbors := g.From(initialUserID)
	//Hay ue agregar un comprobante para que no recorra varias veces el mismo nodo.
	for neighbors.Next() {
		neighbor := neighbors.Node()
		if temp[neighbor.ID()] != nil { // Verificar que el nodo exista en el mapa
			temp[neighbor.ID()].value = 1
			temp[neighbor.ID()].previousInteraction = true
		}
	}

	return temp
}

func resourceAllocation(g *simple.UndirectedGraph, initialUserID int64, steps int) (map[int64]*ResourceInfo, float64) {
	// Inicializar recursos en los objetos conectados al usuario inicial
	resources := initialize_map(g, initialUserID)
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

func heatPropagation(g *simple.UndirectedGraph, initialUserID int64, steps int) (map[int64]*ResourceInfo, float64) {
	// Inicializar recursos en los objetos conectados al usuario inicial
	heat := initialize_map(g, initialUserID)
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

func sortMapDescending(m map[int64]float64) map[int64]float64 {
	// Extraer las claves del mapa en un slice
	keys := make([]int64, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	// Ordenar las claves de mayor a menor según los valores
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]] // Comparación de mayor a menor
	})

	// Crear un nuevo mapa con las claves ordenadas
	sortedMap := make(map[int64]float64)
	for _, key := range keys {
		sortedMap[key] = m[key]
	}

	return sortedMap
}

func HybridRecommendation(g *simple.UndirectedGraph, initialUserID int64, steps int, lambda float64) map[int64]float64 {
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
			if heatExists && !heatValue.previousInteraction {
				hybridScores[node.ID()] = lambda * (heatValue.value / max_heat)
			}
			resourceValue, resourceExists := resources[node.ID()]
			if resourceExists && !resourceValue.previousInteraction {
				hybridScores[node.ID()] += (1 - lambda) * (resourceValue.value / max_resources)
			}
		}

	}
	return sortMapDescending(hybridScores)
}

// Nodo personalizado que implementa la interfaz graph.Node
type Node struct {
	id                   int64
	previous_interaction bool
	name                 []string
}

// Método ID para cumplir con graph.Node
func (n Node) ID() int64 {
	return n.id
}

type ResourceInfo struct {
	value               float64
	previousInteraction bool
}

func test() {
	// Crear un grafo simple no dirigido
	graph := simple.NewUndirectedGraph()

	// Crear nodos personalizados con IDs únicos
	u1 := Node{id: -1}
	u2 := Node{id: -2}
	u3 := Node{id: -3}
	u4 := Node{id: -4}
	o1 := Node{id: 1}
	o2 := Node{id: 2}
	o3 := Node{id: 3}
	o4 := Node{id: 4}
	o5 := Node{id: 5}

	// Agregar nodos al grafo
	graph.AddNode(u1)
	graph.AddNode(u2)
	graph.AddNode(o1)
	graph.AddNode(o2)
	graph.AddNode(o3)

	// Agregar aristas entre usuarios y objetos
	graph.SetEdge(graph.NewEdge(u1, o1))
	graph.SetEdge(graph.NewEdge(u1, o4))
	graph.SetEdge(graph.NewEdge(u2, o1))
	graph.SetEdge(graph.NewEdge(u2, o2))
	graph.SetEdge(graph.NewEdge(u2, o3))
	graph.SetEdge(graph.NewEdge(u2, o4))
	graph.SetEdge(graph.NewEdge(u3, o1))
	graph.SetEdge(graph.NewEdge(u3, o3))
	graph.SetEdge(graph.NewEdge(u4, o3))
	graph.SetEdge(graph.NewEdge(u4, o5))

	resources := HybridRecommendation(graph, u1.ID(), 1, 1)

	// Imprimir resultados
	fmt.Println("Distribución de recursos:")

	for id, value := range resources {
		fmt.Printf("Nodo %d: %.4f\n", id, value)
	}
}
