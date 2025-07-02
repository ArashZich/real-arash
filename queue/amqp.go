package queue

import (
	"log"

	"gitag.ir/armogroup/armo/services/reality/config"
	"github.com/streadway/amqp"
)

func Connect() *amqp.Connection {
	var (
		rabbitMQUserName = config.AppConfig.RabbitMQUserName
		rabbitMQPassword = config.AppConfig.RabbitMQPassword
		rabbitMQHost     = config.AppConfig.RabbitMQHost
		rabbitMQPort     = config.AppConfig.RabbitMQPort
	)

	connStr := "amqp://" + rabbitMQUserName + ":" + rabbitMQPassword + "@" + rabbitMQHost + ":" + rabbitMQPort + "/"

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	return conn
}
