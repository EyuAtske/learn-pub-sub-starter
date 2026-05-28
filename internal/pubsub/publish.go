package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int 
const(
	Durable SimpleQueueType = 0
	Transient SimpleQueueType = 1
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error{
	data, err := json.Marshal(val)
	if err != nil{
		return err
	}
	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error){
	ch, _ := conn.Channel()
	queueT := true
	switch queueType{
	case Transient:
		queueT = false
	}
	queue, err := ch.QueueDeclare(queueName, queueT, !queueT, !queueT, false, nil)
	if err != nil{
		fmt.Errorf(err.Error())
	}
	ch.QueueBind(queueName, key, exchange, false, nil)
	return ch, queue, nil
	}