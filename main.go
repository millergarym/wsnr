package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	a := &Node{X: 1, Y: 1, SensorTypes: []string{"temp", "hum"}, Receiver: make(chan Data)}
	b := &Node{X: 1, Y: 2, SensorTypes: []string{"temp", "hum"}, Receiver: make(chan Data)}
	c := &Node{X: 2, Y: 1, SensorTypes: []string{"temp", "hum"}, Receiver: make(chan Data)}
	a.Neibours = []*Node{b, c}
	b.Neibours = []*Node{a, c}
	c.Neibours = []*Node{a, b}
	fmt.Printf("a: %+v\n", a)
	fmt.Printf("b: %+v\n", b)
	fmt.Printf("c: %+v\n", c)
	go a.Agent()
	go b.Agent()
	go c.Agent()
	http.HandleFunc("/", greet)
	log.Printf("%v\n", http.ListenAndServe(":8081", nil))
}

type Node struct {
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
	i := 0
	for {
		select {
		case data := <-node.Receiver:
			fmt.Printf("xy %d:%d data: %v\n", node.X, node.Y, data)
		case <-time.After(300 * time.Millisecond):
			node.Neibours[i%len(node.Neibours)].Receiver <- Data{"xy", i}
			i++
		}
	}
}
