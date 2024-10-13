package nats_helper

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

type NatsMessageHandler func(*nats.Msg)

type NatsListener struct {
	Handler NatsMessageHandler
}

//goland:noinspection GoUnusedExportedFunction
func PublishToNATS(natsQueue string, message []byte) error {
	nc, err := connect()
	if err != nil {
		return err
	}
	defer nc.Close()

	err = nc.Publish(natsQueue, message)
	if err != nil {
		return err
	}

	log.Printf("[PublishToNATS] Повідомлення надіслано в NATS")
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func StartNatsListener(queueName string, function *NatsListener) *nats.Conn {
	nc, err := connect()
	if err != nil {
		log.Printf("[NatsListener] Error connecting to NATS: %v", err)
		return nil
	}
	if _, err = nc.Subscribe(queueName, func(msg *nats.Msg) { function.Handler(msg) }); err != nil {
		log.Printf("[NatsListener] Error subscribing to NATS queue: %v", err)
		return nil
	}
	err = nc.Flush()
	if err != nil {
		log.Printf("[NatsListener] Error during NATS flush: %v", err)
		return nil
	}
	if err = nc.LastError(); err != nil {
		log.Printf("[NatsListener] NATS subscription error: %v", err)
	}

	log.Println("[NatsListener] NATS subscription successful")
	return nc
}

func connect() (*nats.Conn, error) {
	ip := os.Getenv("BROKER_IP")
	port := os.Getenv("BROKER_PORT")
	natsUrl := fmt.Sprintf("nats://%s:%s", ip, port)

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
