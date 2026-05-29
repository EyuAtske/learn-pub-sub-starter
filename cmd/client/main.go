package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"
	con, err := amqp.Dial(connStr)
	if err != nil{
		fmt.Println("Connection faild: %w", err)
	}
	defer con.Close()
	username, err:= gamelogic.ClientWelcome()
	queueName := fmt.Sprintf(routing.PauseKey + "." + username)
	queuetype := pubsub.Transient
	_, _, err = pubsub.DeclareAndBind(con, routing.ExchangePerilDirect, queueName, routing.PauseKey, queuetype)
	if err != nil{
		fmt.Print(err.Error())
		return
	}
	fmt.Println("Connection was successfull")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("program is shutting down")
	con.Close()
}
