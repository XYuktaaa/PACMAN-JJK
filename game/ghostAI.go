// ghostAI.go
package main

import (
	"container/heap"
	"math"
)

type Node struct {
	X, Y     int
	G, H, F  float64
	Parent   *Node
	Index    int // for priority queue
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].F < pq[j].F }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := x.(*Node)
	n.Index = len(*pq)
	*pq = append(*pq, n)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := old[len(old)-1]
	*pq = old[:len(old)-1]
	return n
}

func heuristic(a, b Node) float64 {
	dx := math.Abs(float64(a.X - b.X))
	dy := math.Abs(float64(a.Y - b.Y))
	return dx + dy // Manhattan distance
}

func findPath(level [][]int, startX, startY, endX, endY int) []Node {
	start := &Node{X: startX, Y: startY}
	goal := Node{X: endX, Y: endY}

	openSet := make(PriorityQueue, 0)
	heap.Init(&openSet)
	heap.Push(&openSet, start)

	gScore := make(map[[2]int]float64)
	gScore[[2]int{start.X, start.Y}] = 0

	visited := make(map[[2]int]bool)

	for openSet.Len() > 0 {
		current := heap.Pop(&openSet).(*Node)
		if current.X == goal.X && current.Y == goal.Y {
			return reconstructPath(current)
		}

		visited[[2]int{current.X, current.Y}] = true

		for _, dir := range [][2]int{{0,1},{1,0},{0,-1},{-1,0}} {
			nx, ny := current.X+dir[0], current.Y+dir[1]
			if ny < 0 || ny >= len(level) || nx < 0 || nx >= len(level[0]) || level[ny][nx] == 1 {
				continue
			}
			if visited[[2]int{nx, ny}] {
				continue
			}

			tempG := gScore[[2]int{current.X, current.Y}] + 1
			neighbor := &Node{X: nx, Y: ny, Parent: current}
			if oldG, ok := gScore[[2]int{nx, ny}]; !ok || tempG < oldG {
				gScore[[2]int{nx, ny}] = tempG
				neighbor.G = tempG
				neighbor.H = heuristic(*neighbor, goal)
				neighbor.F = neighbor.G + neighbor.H
				heap.Push(&openSet, neighbor)
			}
		}
	}
	return nil // no path found
}

func reconstructPath(n *Node) []Node {
	var path []Node
	for n != nil {
		path = append([]Node{*n}, path...)
		n = n.Parent
	}
	return path
}

