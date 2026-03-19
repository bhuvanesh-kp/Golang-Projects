package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarber  int
	BarberDoneChan  chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string){
	shop.NumberOfBarber++

	go func ()  {
		isSleeping := false
		msg := fmt.Sprintf("%s goes to the waiting room to check for clients.", barber)
		color.Green(msg)

		for {
			// if there is no client , the barber goes to sleep
			if len(shop.ClientsChan) == 0 {
				msg = fmt.Sprintf("There is nothing to do, so %s takes a nap.", barber)
				color.Yellow(msg)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientsChan

			if shopOpen{
				if isSleeping{
					msg = fmt.Sprintf("%s wake %s up.", client, barber)
					color.Blue(msg)
					isSleeping =  false
				}
				shop.cutHair(barber, client)
			}else{
				// shop is closed , so send the barber home and close this go routine
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string){
	msg := fmt.Sprintf("%s is cutting %s's hair", barber, client)
	color.Green(msg)

	time.Sleep(shop.HairCutDuration)

	msg = fmt.Sprintf("%s is finished cutting %s's hair", barber, client)
	color.Green(msg)
}

func (shop *BarberShop) sendBarberHome(barber string){
	msg := fmt.Sprintf("%s is going home.\n ", barber)
	color.White(msg)
	shop.BarberDoneChan <- true
}

func (shop *BarberShop) closeShopForDay(){
	msg := "Closing shop for the day"
	color.Red(msg)

	close(shop.ClientsChan)
	shop.Open = false

	for a := 1;a <= shop.NumberOfBarber;a++{
		<-shop.BarberDoneChan
	}

	close(shop.BarberDoneChan)

	color.Cyan("--------------------------")
	msg = "The barber shop is closed for the day , everyone has gone home"
	color.White(msg)
}

func (shop *BarberShop) addClient(client string){
	msg := fmt.Sprintf("%s arrives", client)
	color.Cyan(msg)

	if shop.Open{
		select{
		case shop.ClientsChan <-client:
			msg = fmt.Sprintf("%s takes a seat in the waiting room.", client)
			color.Cyan(msg)
		default:
			msg = fmt.Sprintf("Waiting room is full , so %s leaves!", client)
			color.Red(msg)
		}
	}else{
		msg = fmt.Sprintf("The shop is already closed , so %s leaves!", client)
		color.Red(msg)
	}
}