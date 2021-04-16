package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sync"
	"time"
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

type ConcurrencyData struct {
	Position int
	Distance float64
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

	start := time.Now()
	Linear(nodePool, c)
	log.Printf("Linear Ellapsed %v\n", time.Since(start))
	PrintCycle(c)

	nodePool2 := AssembleNodePool(&inputFIle)
	c2 := &Cycle{
		NodesIncluded: 0,
		Start:         nil,
	}
	time.Sleep(5 * time.Second)

	start = time.Now()
	Concurrent(nodePool2, c2)
	log.Printf("Ellapsed %v\n", time.Since(start))

	PrintCycle(c2)

	return
}

func PrintCycle(c *Cycle) {
	n := c.Start
	for i := 0; i < c.NodesIncluded; i++ {
		fmt.Print(n.NodeID, "->")
		n = n.Next
	}
	fmt.Print(n.NodeID, "\n")

	n = c.Start
	for i := 0; i < c.NodesIncluded; i++ {
		fmt.Println("Node:", n.NodeID, "Distance to next: ", n.DistanceToOtherNodes[n.Next.NodeID])
		n = n.Next
	}

}

func Linear(nodePool []*Node, c *Cycle) {
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

}

func Concurrent(nodePool []*Node, c *Cycle) {
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

			var wg sync.WaitGroup
			j := 0
			dChannel := make(chan ConcurrencyData)

			for currentNode != c.Start {
				wg.Add(1)
				//j = 0
				go func(pos int, n, cn *Node, group *sync.WaitGroup) {
					defer group.Done()
					dik := n.DistanceToOtherNodes[cn.NodeID]
					dkj := n.DistanceToOtherNodes[cn.Next.NodeID]
					dij := cn.DistanceToOtherNodes[cn.Next.NodeID]

					d := dik + dkj - dij
					dChannel <- ConcurrencyData{
						Position: pos,
						Distance: d,
					}

				}(j, nodePool[i], currentNode, &wg)
				j++

				currentNode = currentNode.Next
			}

			go func(g *sync.WaitGroup, channel chan ConcurrencyData) {
				g.Wait()
				close(dChannel)
			}(&wg, dChannel)

			for d := range dChannel {
				if d.Distance < distance {
					distance = d.Distance
					position = d.Position
				}
			}

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
