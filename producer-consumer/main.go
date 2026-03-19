package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NumberOfPizzas int = 10

var pizzaMade, pizzaCancelled, total int

type Producer struct {
	data chan PizzaOrdered
	quit chan chan error
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

type PizzaOrdered struct {
	pizzaNumber int
	message     string
	success     bool
}


func Pizzaria(pizzaMaker *Producer) {
	// running pizza job and keep track of pizza made
	// run until a quit notification is received 
	// try to make pizzas

	for ;;{
		//currentPizza := makePizza()
	}
}

func main() {
	fmt.Println("Creating a producer consumer problem simulation using pizza shop simulation")

	//seed the random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("pizza shop is open for bussiness")

	pizzaJob := &Producer{
		data: make(chan PizzaOrdered),
		quit: make(chan chan error),
	}

	// running producer job in background
	go Pizzaria(pizzaJob)
}
