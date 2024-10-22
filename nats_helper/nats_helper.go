package nats_helper

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

type NatsMessage struct {
	UserId int64  `json:"userId"`
	Text   string `json:"natsMessage"`
}

type Command struct {
	UserId    int64    `json:"userId"`
	Arguments []string `json:"arguments"`
}

type NatsListenerHandlerFunction func(*nats.Msg)

type NatsListenerHandler struct {
	Function NatsListenerHandlerFunction
}

//goland:noinspection GoUnusedExportedFunction
func PublishToNATS(queue string, message []byte) error {
	nc, err := connect()
	if err != nil {
		log.Printf("[PublishToNATS] Помилка під'єднання до NATS (черга \"%s\"): %v", err, queue)
		return err
	}
	defer nc.Close()

	if err = nc.Publish(queue, message); err != nil {
		log.Printf("[PublishToNATS] Помилка публікації в чергу \"%s\" NATS: %v", queue, err)
		return err
	}

	log.Printf("[PublishToNATS] Повідомлення надіслано в чергу \"%s\" NATS", queue)
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func StartNatsListener(queue string, handler *NatsListenerHandler) error {
	nc, err := connect()
	if err != nil {
		log.Printf("[StartNatsListener] Помилка під'єднання до NATS (черга \"%s\"): %v", err, queue)
		return nil
	}
	if _, err = nc.Subscribe(queue, func(msg *nats.Msg) { handler.Function(msg) }); err != nil {
		log.Printf("[StartNatsListener] Помилка підписки до черги \"%s\" в NATS: %v", queue, err)
		return err
	}

	if err = nc.Flush(); err != nil {
		log.Printf("[StartNatsListener] Помилка flash в черзі \"%s\" NATS : %v", queue, err)
		return err
	}

	if err = nc.LastError(); err != nil {
		log.Printf("[StartNatsListener] Помилка для черги \"%s\" в NATS: %v", queue, err)
		return err
	}

	log.Printf("[StartNatsListener] Підписка до черги \"%s\" вдала", queue)
	return nil
}

func connect() (*nats.Conn, error) {
	ip := os.Getenv("BROKER_IP")
	port := os.Getenv("BROKER_PORT")
	natsUrl := fmt.Sprintf("nats://%s:%s", ip, port)

	log.Printf("[connect] Під'єднуюся до NATS: %s", natsUrl)

	for i := 0; i < 5; i++ {
		nc, err := nats.Connect(natsUrl)
		if err != nil {
			log.Printf("[NatsListener] Error connecting to NATS (%d try): %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}

		log.Printf("[connect] Під'єднався до NATS: %s", natsUrl)
		return nc, nil
	}
	return nil, fmt.Errorf("[connect] Неможливо під'єднатися до NATS: %s", natsUrl)
}
