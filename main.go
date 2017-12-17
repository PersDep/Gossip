package main

import (
	"./graph"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func print(graph graph.Graph) {
	for node, nodes := range graph {
		fmt.Println(node)
		fmt.Println(nodes)
	}
}

func main() {
	var nodesAmount, port, minDegree, maxDegree int
	var errors [4]error

	nodesAmount, errors[0] = strconv.Atoi("50")
	port, errors[1] = strconv.Atoi("4400")
	minDegree, errors[2] = strconv.Atoi("5")
	maxDegree, errors[3] = strconv.Atoi("7")

	for _, error := range errors {
		if error != nil {
			panic(error)
		}
	}

	rand.Seed(int64(time.Now().Nanosecond()))

	graph := graph.Generate(nodesAmount, minDegree, maxDegree, port)
	print(graph)
}
