package main

import (
	"container/heap"
	"math"
	"fmt"
)

type Node struct {
	X, Y     int
	G, H, F  float64
	Parent   *Node
	Index    int // for priority queue
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { 
	// Primary sort by F score, secondary by H score (tie-breaker)
	if pq[i].F == pq[j].F {
		return pq[i].H < pq[j].H
	}
	return pq[i].F < pq[j].F 
}

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

// Manhattan distance heuristic (consistent and admissible)
func heuristic(a, b Node) float64 {
	dx := math.Abs(float64(a.X - b.X))
	dy := math.Abs(float64(a.Y - b.Y))
	return dx + dy
}

// Improved A* pathfinding with better performance and error handling
func findPath(level [][]int, startX, startY, endX, endY int) []Node {
    fmt.Printf("findPath called: start(%d,%d) -> end(%d,%d)\n", startX, startY, endX, endY)
    
    // Validate input bounds
    if startY < 0 || startY >= len(level) || startX < 0 || startX >= len(level[0]) ||
       endY < 0 || endY >= len(level) || endX < 0 || endX >= len(level[0]) {
        fmt.Printf("Out of bounds: start(%d,%d) end(%d,%d), level size %dx%d\n", 
                   startX, startY, endX, endY, len(level[0]), len(level))
        return nil
    }

    fmt.Printf("Start tile value: %d, End tile value: %d\n", level[startY][startX], level[endY][endX])

    // Check if start or end is a wall
    if level[startY][startX] == 1 {
        fmt.Printf("Start position (%d,%d) is a wall!\n", startX, startY)
        return nil
    }
    if level[endY][endX] == 1 {
        fmt.Printf("End position (%d,%d) is a wall!\n", endX, endY)
        return nil
    }

    // If already at target, return single-node path
    if startX == endX && startY == endY {
        return []Node{{X: startX, Y: startY}}
    }

    start := &Node{X: startX, Y: startY, G: 0}
    goal := Node{X: endX, Y: endY}
    start.H = heuristic(*start, goal)
    start.F = start.G + start.H

    openSet := make(PriorityQueue, 0)
    heap.Init(&openSet)
    heap.Push(&openSet, start)

    gScore := make(map[[2]int]float64)
    gScore[[2]int{start.X, start.Y}] = 0

    closedSet := make(map[[2]int]bool)
    inOpenSet := make(map[[2]int]*Node)
    inOpenSet[[2]int{start.X, start.Y}] = start

    directions := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
    maxIterations := len(level) * len(level[0]) * 2
    iterations := 0

    for openSet.Len() > 0 && iterations < maxIterations {
        iterations++
        
        current := heap.Pop(&openSet).(*Node)
        currentKey := [2]int{current.X, current.Y}
        
        delete(inOpenSet, currentKey)
        
        if current.X == goal.X && current.Y == goal.Y {
            fmt.Printf("Path found after %d iterations!\n", iterations)
            return reconstructPath(current)
        }

        closedSet[currentKey] = true

        for _, dir := range directions {
            nx, ny := current.X + dir[0], current.Y + dir[1]
            neighborKey := [2]int{nx, ny}

            if ny < 0 || ny >= len(level) || nx < 0 || nx >= len(level[0]) {
                continue
            }
            
            // IMPORTANT: Allow movement through empty spaces (0), pellets (2), and power pellets (4)
            if level[ny][nx] != 0 && level[ny][nx] != 2 && level[ny][nx] != 4 {
                continue
            }

            if closedSet[neighborKey] {
                continue
            }

            tentativeG := current.G + 1.0

            if existingNode, exists := inOpenSet[neighborKey]; exists {
                if tentativeG < existingNode.G {
                    existingNode.G = tentativeG
                    existingNode.F = existingNode.G + existingNode.H
                    existingNode.Parent = current
                    heap.Fix(&openSet, existingNode.Index)
                }
            } else {
                neighbor := &Node{
                    X: nx, 
                    Y: ny, 
                    G: tentativeG,
                    Parent: current,
                }
                neighbor.H = heuristic(*neighbor, goal)
                neighbor.F = neighbor.G + neighbor.H
                
                heap.Push(&openSet, neighbor)
                inOpenSet[neighborKey] = neighbor
                gScore[neighborKey] = tentativeG
            }
        }
    }

    fmt.Printf("No path found after %d iterations\n", iterations)
    return nil
}
// Reconstruct path from goal to start
func reconstructPath(node *Node) []Node {
	path := make([]Node, 0)
	current := node
	
	for current != nil {
		path = append([]Node{{X: current.X, Y: current.Y}}, path...)
		current = current.Parent
	}
	
	return path
}

// Smooth path by removing unnecessary waypoints (optional optimization)
func smoothPath(level [][]int, path []Node) []Node {
	if len(path) <= 2 {
		return path
	}

	smoothed := make([]Node, 0, len(path))
	smoothed = append(smoothed, path[0])

	i := 0
	for i < len(path)-1 {
		j := len(path) - 1
		
		// Find the furthest node we can reach directly
		for j > i+1 {
			if hasLineOfSight(level, path[i], path[j]) {
				break
			}
			j--
		}
		
		smoothed = append(smoothed, path[j])
		i = j
	}

	return smoothed
}

// Check if there's a clear line of sight between two nodes
func hasLineOfSight(level [][]int, start, end Node) bool {
	dx := end.X - start.X
	dy := end.Y - start.Y
	
	steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))
	if steps == 0 {
		return true
	}

	stepX := float64(dx) / float64(steps)
	stepY := float64(dy) / float64(steps)

	for i := 0; i <= steps; i++ {
		x := start.X + int(math.Round(float64(i)*stepX))
		y := start.Y + int(math.Round(float64(i)*stepY))
		
		if y < 0 || y >= len(level) || x < 0 || x >= len(level[0]) || level[y][x] == 1 {
			return false
		}
	}
	
	return true
}

// Alternative simple pathfinding for when A* might be overkill
func findSimplePath(level [][]int, startX, startY, endX, endY int) []Node {
	// Simple greedy approach - move towards target when possible
	path := make([]Node, 0)
	x, y := startX, startY
	path = append(path, Node{X: x, Y: y})

	maxSteps := 100 // Prevent infinite loops
	steps := 0

	for (x != endX || y != endY) && steps < maxSteps {
		steps++
		moved := false

		// Prefer horizontal movement first
		if x != endX {
			newX := x
			if endX > x {
				newX = x + 1
			} else {
				newX = x - 1
			}
			
			if newX >= 0 && newX < len(level[0]) && level[y][newX] == 0 {
				x = newX
				path = append(path, Node{X: x, Y: y})
				moved = true
			}
		}

		// Then try vertical movement
		if !moved && y != endY {
			newY := y
			if endY > y {
				newY = y + 1
			} else {
				newY = y - 1
			}
			
			if newY >= 0 && newY < len(level) && level[newY][x] == 0 {
				y = newY
				path = append(path, Node{X: x, Y: y})
				moved = true
			}
		}

		// If we can't move towards target, we're stuck
		if !moved {
			break
		}
	}

	if x == endX && y == endY {
		return path
	}
	
	// Fallback to A* if simple path fails
	return findPath(level, startX, startY, endX, endY)
}
