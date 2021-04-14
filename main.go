package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

type InputFile struct {
	Matrix [][]struct {
		//DistanceInMeters    int     `json:"distance_in_meters"`
		TravelTimeInMinutes float64 `json:"travel_time_in_minutes"`
	} `json:"matrix"`
}

type Node struct {
	Next                 *Node
	Before               *Node
	DistanceToOtherNodes []float64
	NodeID               int
}

type Cycle struct {
	NodesIncluded int
	Start         *Node
}

func main() {

	file, err := ioutil.ReadFile("1000_cm.json")
	if err != nil {
		panic(err)
	}

	var inputFIle InputFile
	err = json.Unmarshal(file, &inputFIle)
	if err != nil {
		panic(err)
	}

	nodePool := AssembleNodePool(&inputFIle)
	c := &Cycle{
		NodesIncluded: 0,
		Start:         nil,
	}

	for i := 0; i < len(nodePool); i++ {
		if nodePool[i] == nil {
			continue
		}
		// I can assume when the cycle is empty, the pool is full, so I can get the first element
		if c.NodesIncluded == 0 {
			c.Start = nodePool[i]
			nodePool[i] = nil

			c.Start.Next = c.Start
			c.Start.Before = c.Start

			c.NodesIncluded++

			continue
		} else if c.NodesIncluded == 1 {
			/*n := c.Start
			shortestNodeDistance := math.MaxFloat64

			selected := 0
			for j := range n.DistanceToOtherNodes {
				if n.DistanceToOtherNodes[j] == 0 {
					continue
				}

				if n.DistanceToOtherNodes[j] < shortestNodeDistance {
					selected = j
					shortestNodeDistance = n.DistanceToOtherNodes[j]
					c.Start.Next = nodePool[j]
					c.Start.Before = nodePool[j]

					nodePool[j].Next = c.Start
					nodePool[j].Before = c.Start
				}
			}

			nodePool[selected] = nil
			c.NodesIncluded++*/
			c.Start.Next = nodePool[i]
			c.Start.Before = nodePool[i]

			nodePool[i].Next = c.Start
			nodePool[i].Before = c.Start

			nodePool[i] = nil
			c.NodesIncluded++
		} else {
			//dik+dkj-dij
			currentNode := c.Start.Next
			position := 0
			distance := math.MaxFloat64

			j := 0
			for currentNode != c.Start {
				//j = 0
				dik := nodePool[i].DistanceToOtherNodes[currentNode.NodeID]
				dkj := nodePool[i].DistanceToOtherNodes[currentNode.Next.NodeID]
				dij := currentNode.DistanceToOtherNodes[currentNode.Next.NodeID]

				d := dik + dkj - dij

				if d < distance {
					distance = d
					position = j
				}

				j++

				currentNode = currentNode.Next
			}

			fmt.Println("Should be inserted in:", position)

			currentNode = c.Start
			for j = 0; j < position; j++ {
				currentNode = currentNode.Next
			}

			tmpN := currentNode.Next
			tmpN.Before = nodePool[i]
			currentNode.Next = nodePool[i]
			nodePool[i].Next = tmpN
			nodePool[i].Before = currentNode

			nodePool[i] = nil
			c.NodesIncluded++
		}
	}

	n := c.Start
	for i := 0; i < c.NodesIncluded; i++ {
		fmt.Print(n.NodeID, "->")
		n = n.Next
	}

	n = c.Start
	for i := 0; i < c.NodesIncluded; i++ {
		fmt.Println("Node:", n.NodeID, "Distance to next: ", n.DistanceToOtherNodes[n.Next.NodeID])
		n = n.Next
	}

	return
}

func AssembleNodePool(inputFIle *InputFile) []*Node {
	nodePool := make([]*Node, len(inputFIle.Matrix))
	// Create node pool
	for i := range inputFIle.Matrix {
		distances := make([]float64, len(inputFIle.Matrix[i]))
		for j := range inputFIle.Matrix[i] {
			distances[j] = inputFIle.Matrix[i][j].TravelTimeInMinutes
		}

		n := &Node{
			Next:                 nil,
			Before:               nil,
			DistanceToOtherNodes: distances,
			NodeID:               i,
		}

		nodePool[i] = n
	}
	return nodePool
}
