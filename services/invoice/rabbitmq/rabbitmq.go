package transports

import (
	"os"

	"gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"github.com/ARmo-BigBang/kit/log"

	"github.com/streadway/amqp"
)

type agent struct {
	Queue   amqp.Queue
	Channel *amqp.Channel
	Logger  log.Logger
	Invoice service.Invoice
}

func MakeAgent(client *amqp.Connection, qeueuName string, logger log.Logger, invoice service.Invoice) agent {
	// Connect to message broker
	ch, err := client.Channel()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// Insure queue is created or create it
	q, err := ch.QueueDeclare(
		qeueuName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	return agent{
		Queue:   q,
		Channel: ch,
		Logger:  logger,
		Invoice: invoice,
	}
}

// Register Consumers Here
func (ag *agent) RegisterConsumers() {
	// Register consumer for queue
	msgs, err := ag.Channel.Consume(
		ag.Queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)

	if err != nil {
		ag.Logger.Error(err)
		os.Exit(1)
	}

	go func() {
		for d := range msgs {
			ag.Logger.Info("Received a message: %s", d.Body)
			// Call service method to process message
			// this will internally call service.DoSomething()
			// a.doSomething(d.Body)
		}
	}()
}
