package main

import (
	"container/heap"
	"math"
	// "fmt"
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

// A* pathfinding implementation
func findPath(level [][]int, startX, startY, endX, endY int) []Node {
    // Validate bounds
    if startY < 0 || startY >= len(level) || startX < 0 || startX >= len(level[0]) ||
       endY < 0 || endY >= len(level) || endX < 0 || endX >= len(level[0]) {
        return nil
    }

    // Check if start or end is a wall
    if level[startY][startX] == TileWall || level[endY][endX] == TileWall {
        return nil
    }

    // If already at target
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

    closedSet := make(map[[2]int]bool)
    inOpenSet := make(map[[2]int]*Node)
    inOpenSet[[2]int{start.X, start.Y}] = start

    directions := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
    maxIterations := len(level) * len(level[0])
    iterations := 0

    for openSet.Len() > 0 && iterations < maxIterations {
        iterations++
        
        current := heap.Pop(&openSet).(*Node)
        currentKey := [2]int{current.X, current.Y}
        
        delete(inOpenSet, currentKey)
        
        if current.X == goal.X && current.Y == goal.Y {
            return reconstructPath(current)
        }

        closedSet[currentKey] = true

        for _, dir := range directions {
            nx, ny := current.X + dir[0], current.Y + dir[1]
            neighborKey := [2]int{nx, ny}

            if ny < 0 || ny >= len(level) || nx < 0 || nx >= len(level[0]) {
                continue
            }
            
            // Allow movement through walkable tiles
            if !isWalkableTile(level[ny][nx]) {
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
            }
        }
    }

    return nil
}

// Helper function to check if a tile is walkable
func isWalkableTile(tile int) bool {
    return tile == TileEmpty || tile == TilePellet || tile == TilePowerPellet
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
func findSimplePath(level [][]int, startX, startY, endX, endY, tileSize int) []Node {
    startTileX := startX / tileSize
    startTileY := startY / tileSize
    endTileX := endX / tileSize
    endTileY := endY / tileSize
    
    path := make([]Node, 0)
    x, y := startTileX, startTileY
    path = append(path, Node{X: x, Y: y})

    maxSteps := 50
    steps := 0

    for (x != endTileX || y != endTileY) && steps < maxSteps {
        steps++
        moved := false

        // Try horizontal movement first
        if x != endTileX {
            newX := x
            if endTileX > x {
                newX = x + 1
            } else {
                newX = x - 1
            }
            
            if newX >= 0 && newX < len(level[0]) && 
               (level[y][newX] == TileEmpty || level[y][newX] == TilePellet || level[y][newX] == TilePowerPellet) {
                x = newX
                path = append(path, Node{X: x, Y: y})
                moved = true
            }
        }

        // Try vertical movement
        if !moved && y != endTileY {
            newY := y
            if endTileY > y {
                newY = y + 1
            } else {
                newY = y - 1
            }
            
            if newY >= 0 && newY < len(level) && 
               (level[newY][x] == TileEmpty || level[newY][x] == TilePellet || level[newY][x] == TilePowerPellet) {
                y = newY
                path = append(path, Node{X: x, Y: y})
                moved = true
            }
        }

        if !moved {
            break
        }
    }

    return path
}
