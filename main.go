package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	a := &Node{Name: "a", X: 1, Y: 1, SensorTypes: []string{"temp", "hum"}, Receiver: make(chan Data, 1)}
	b := &Node{Name: "b", X: 1, Y: 2, SensorTypes: []string{"temp", "hum"}, Receiver: make(chan Data, 1)}
	c := &Node{Name: "c", X: 2, Y: 1, SensorTypes: []string{"temp", "hum"}, Receiver: make(chan Data, 1)}
	a.Neibours = []*Node{b, c}
	b.Neibours = []*Node{a, c}
	c.Neibours = []*Node{a, b}
	go a.Agent()
	go b.Agent()
	go c.Agent()
	http.HandleFunc("/", greet)
	log.Printf("%v\n", http.ListenAndServe(":8081", nil))
}

type Node struct {
	Name        string
	X, Y        int
	SensorTypes []string
	Neibours    []*Node
	Receiver    chan Data
}

type Data struct {
	Interest  string
	DataPoint int
}

func (node *Node) Agent() {
	fmt.Printf("%+v\n", *node)
	i := 0
	for {
		select {
		case data := <-node.Receiver:
			fmt.Printf("node: %s data: %v\n", node.Name, data)
		case <-time.After(1 * time.Millisecond):
			node.Neibours[i%len(node.Neibours)].Receiver <- Data{"xy", i}
			i++
		}
	}
}
