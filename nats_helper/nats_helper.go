package nats_helper

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

type natsBotText struct {
	UserId int64  `json:"user_id"`
	Text   string `json:"text"`
}

type natsBotCommand struct {
	UserId    int64    `json:"user_id"`
	Arguments []string `json:"arguments"`
}

type natsBotFile struct {
	UserId   int64  `json:"user_id"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

type NatsListenerHandlerFunction func()

type NatsListenerHandler struct {
	Function NatsListenerHandlerFunction
}

//goland:noinspection GoUnusedExportedFunction
func PublishTextMessage(queue string, userId int64, text string) error {
	msg := natsBotText{
		UserId: userId,
		Text:   text,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[PublishTextMessage] Помилка: %s", err)
		return err
	}

	if err = publishMessageToNats(queue, jsonData); err == nil {
		log.Printf("[PublishToNATS] Повідомлення надіслано в чергу \"%s\" NATS", queue)
		return nil
	} else {
		log.Printf("[PublishToNATS] Помилка публікації в чергу \"%s\" NATS: %v", queue, err)
		return err
	}
}

//goland:noinspection GoUnusedExportedFunction
func PublishCommandMessage(queue string, userId int64, message []string) error {
	msg := natsBotCommand{
		UserId:    userId,
		Arguments: message,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[PublishCommandMessage] Помилка: %s", err)
		return err
	}

	if err = publishMessageToNats(queue, jsonData); err == nil {
		log.Printf("[PublishCommandMessage] Повідомлення надіслано в чергу \"%s\" NATS", queue)
		return nil
	} else {
		log.Printf("[PublishCommandMessage] Помилка публікації в чергу \"%s\" NATS: %v", queue, err)
		return err
	}
}

//goland:noinspection GoUnusedExportedFunction
func PublishFileInfoMessage(queue string, userId int64, fileId string, fileName string, fileSize int64, mimeType string) error {
	msg := natsBotFile{
		UserId:   userId,
		FileId:   fileId,
		FileName: fileName,
		Size:     fileSize,
		MimeType: mimeType,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[PublishFileInfoMessage] Помилка: %s", err)
		return err
	}

	if err = publishMessageToNats(queue, jsonData); err == nil {
		log.Printf("[PublishFileInfoMessage] Повідомлення надіслано в чергу \"%s\" NATS", queue)
		return nil
	} else {
		log.Printf("[PublishFileInfoMessage] Помилка публікації в чергу \"%s\" NATS: %v", queue, err)
		return err
	}
}

//goland:noinspection GoUnusedExportedFunction

//goland:noinspection GoUnusedExportedFunction
func ParseNatsBotText(data []byte) (int64, string, error) {
	var msg natsBotText
	if err := json.Unmarshal(data, &msg); err == nil {
		userId := msg.UserId
		text := msg.Text
		log.Printf("[ParseNatsBotText] Отримано текст \"%s\" для користувача %d", text, userId)
		return userId, text, nil
	} else {
		log.Printf("[ParseNatsBotText] Помилка при розборі повідомлення з NATS: %v", err)
		return 0, "", err
	}
}

//goland:noinspection GoUnusedExportedFunction
func ParseNatsBotCommand(data []byte) (int64, []string, error) {
	var msg natsBotCommand
	if err := json.Unmarshal(data, &msg); err == nil {
		userId := msg.UserId
		arguments := msg.Arguments
		log.Printf("[ParseNatsBotText] Отримано аргументи команди \"%v\" для користувача %d", arguments, userId)
		return userId, arguments, nil
	} else {
		log.Printf("[ParseNatsBotText] Помилка при розборі повідомлення з NATS: %v", err)
		return 0, nil, err
	}
}

//goland:noinspection GoUnusedExportedFunction
func ParseNatsBotFile(data []byte) (int64, string, string, int64, string, error) {
	var msg natsBotFile
	if err := json.Unmarshal(data, &msg); err == nil {
		userId := msg.UserId
		fileId := msg.FileId
		fileName := msg.FileName
		size := msg.Size
		mimeType := msg.MimeType
		log.Printf("[ParseNatsBotFile] Отримано файл \"%s\" для користувача %d", fileName, userId)
		return userId, fileId, fileName, size, mimeType, nil
	} else {
		log.Printf("[ParseNatsBotFile] Помилка при розборі повідомлення з NATS: %v", err)
		return 0, "", "", 0, "", err
	}
}

//goland:noinspection GoUnusedExportedFunction
func StartNatsListener(queue string, handler *NatsListenerHandler) error {
	nc, err := connect()
	if err != nil {
		log.Printf("[StartNatsListener] Помилка під'єднання до NATS (черга \"%s\"): %v", err, queue)
		return nil
	}
	if _, err = nc.Subscribe(queue, func(msg *nats.Msg) { handler.Function() }); err != nil {
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

func publishMessageToNats(queue string, message []byte) error {
	if nc, err := connect(); err == nil {
		defer nc.Close()
		return nc.Publish(queue, message)
	} else {
		log.Printf("[publishMessageToNats] Помилка під'єднання до NATS (черга \"%s\"): %v", err, queue)
		return err
	}
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
