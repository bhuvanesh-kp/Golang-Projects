package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10		
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {
	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// print starting message
	color.Cyan("sleeping barber problem")
	color.Cyan("--------------------------")

	// create channels if we need any
	clientChannel := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarber:  0,
		BarberDoneChan:  doneChan,
		ClientsChan:     clientChannel,
		Open:            true,
	}

	color.Cyan("The shop is open ryt now")

	// add barbers
	shop.addBarber("james bond")

	// adding more barbers to ensure clients are satisfied with service
	shop.addBarber("kratos")
	shop.addBarber("tony stark")
	shop.addBarber("john wick")
	shop.addBarber("max verstappen")

	// start the barbershop as a go routine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func ()  {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients
	i := 1

	go func ()  {
		for {
			// get a random number with average arrival rate
			randomMillsecond := rand.Int() % (2 * arrivalRate)
			select{
			case <- shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMillsecond)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}	
	}()

	// block untill the barbershop is closed
	time.Sleep(7 * time.Second)
}
