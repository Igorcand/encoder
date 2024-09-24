package queue

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

type RabbitMQ struct {
	User 				string
	Password 			string
	Host 				string
	Port 				string
	VHost 				string
	ConsumerQueueName 	string
	ConsumerName 		string
	AutoAck				bool
	Args				amqp.Table
	Channel				*amqp.Channel
}

func NewRabbitMQ() *RabbitMQ{
	rabbitMQArgs := amqp.Table{}
	rabbitMQArgs["x-dead-letter-exchange"] = os.Getenv("RABBITMQ_DLX")

	rabbitMQ := RabbitMQ{
		User: 				os.Getenv("RABBITMQ_DEFAULT_USER"),
		Password: 			os.Getenv("RABBITMQ_DEFAULT_PASSWORD"),
		Host: 				os.Getenv("RABBITMQ_DEFAULT_HOST"),
		Port: 				os.Getenv("RABBITMQ_DEFAULT_PORT"),
		VHost: 				os.Getenv("RABBITMQ_DEFAULT_VHOST"),
		ConsumerQueueName: 	os.Getenv("RABBITMQ_CONSUMER_QUEUE_NAME"),
		ConsumerName: 		os.Getenv("RABBITMQ_CONSUMER_NAME"),
		AutoAck: 			false,
		Args: 				rabbitMQArgs,
	}
	return &rabbitMQ
}

func (r *RabbitMQ) Connect() *amqp.Channel{
	dsn := "amqp://" + r.User + ":" + r.Password + "@" + r.Host + ":" + r.Port + r.VHost
	conn, err := amqp.Dial(dsn)
	failOnError(err, "failed to connect to RabbitMQ")

	r.Channel, err = conn.Channel()
	failOnError(err, "failed to open a RabbitMQ")
	return r.Channel
}

func (r *RabbitMQ) Consume(messageChannel chan amqp.Delivery){
	q, err := r.Channel.QueueDeclare(
		r.ConsumerQueueName,	// name
		true,					// durable
		false,					// delete when unused
		false,					// exclusive
		false,					//  no-wait
		r.Args,					// arguments
	)

	failOnError(err, "failed to declare queue")

	incomingMessage, err := r.Channel.Consume(
		q.Name,				// queue
		r.ConsumerName,		// consumer
		r.AutoAck, 			// auto-ack
		false,				// exclusive
		false,				// no-local
		false, 				// no-wait
		nil,				// args
	)
	failOnError(err, "failed to register consumer")

	go func(){
		for message := range incomingMessage{
			log.Println("Incoming new message")
			messageChannel <- message
		} 
		log.Println("RabbitMQ channel closed")
		close(messageChannel)

	}()


}

func (r *RabbitMQ) Notify(message string, contentType string, exchange string, routingKey string) error{
	err := r.Channel.Publish(
		exchange,		// exchange
		routingKey, 	// routing key
		false, 			//mandatory
		false,			// immediate
		amqp.Publishing{
			ContentType: 	contentType,
			Body:			[]byte(message),
		})
	
	if err != nil{
		return err
	}

	return nil
}


func failOnError(err error, msg string){
	if err != nil{
		log.Fatalf("%s: %s", msg, err)
	}
}

