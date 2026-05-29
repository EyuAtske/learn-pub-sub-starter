package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	gamelogic.PrintServerHelp()
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
	serverLoop:
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("sending a pause message")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		case "resume":
			fmt.Println("sending a resume message")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		case "quit":
			fmt.Println("exiting the server")
			break serverLoop
		default:
			fmt.Println("Incorrect command")
		}
	}
	fmt.Println("Connection was successfull")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("program is shutting down")
	con.Close()
}
