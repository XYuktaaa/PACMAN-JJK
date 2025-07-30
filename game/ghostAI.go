package main

import(
    "container/heap"
    "math"
    
)
type Node struct{
    X, Y   int
    G, H   int
    Parent *Node
    Index  int
}

func (n *Node) F() float{
    return n.G + n.H
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
    return len(pq)
}
func (pq PriorityQueue) Less(i, j int)bool {
    return pq[i].F()< pq[j].F()
}

