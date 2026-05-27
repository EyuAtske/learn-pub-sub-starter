package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"
	con, err := amqp.Dial(connStr)
	if err != nil{
		fmt.Println("Connection faild: %w", err)
	}
	defer con.Close()
	ch, err := con.Channel()
	if err != nil{
		fmt.Print("coundn't create channel")
		return
	}
	pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	fmt.Println("Connection was successfull")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("program is shutting down")
	con.Close()
}
