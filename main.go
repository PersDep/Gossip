package main

import (
	"./graph"
	"fmt"
	"math/rand"
	"os"
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
	var nodesAmount, port, minDegree, maxDegree, ttl int
	var errors [5]error

	nodesAmount, errors[0] = strconv.Atoi(os.Args[1])
	port, errors[1] = strconv.Atoi(os.Args[2])
	minDegree, errors[2] = strconv.Atoi(os.Args[3])
	maxDegree, errors[3] = strconv.Atoi(os.Args[4])
	ttl, errors[4] = strconv.Atoi(os.Args[5])

	for _, error := range errors {
		if error != nil {
			panic(error)
		}
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	graph := graph.Generate(nodesAmount, minDegree, maxDegree, port)

	quitChan := make(chan struct{})
	killChan := make(chan struct{}, nodesAmount)

	for i := 0; i < nodesAmount; i++ {
		go runGossipNerwork(i, graph, quitChan, killChan, ttl, i == 0)
	}

	<-quitChan
	for i := 0; i < nodesAmount; i++ {
		killChan <- struct{}{}
	}

	time.Sleep(time.Millisecond * 100)
}
