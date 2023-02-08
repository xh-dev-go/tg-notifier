package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"os"
	InitCore "tg-notification/init-core"
	TgSend "tg-notification/tg-send"
	"time"
)

type NotificationEventDto struct {
	Topic            string `json:"topic"`
	Partition        int    `json:"partition"`
	Offset           int64  `json:"offset"`
	CompleteTime     string `json:"completeTime"`
	ProcessingPeriod int64  `json:"processingPeriod"`
	Result           bool   `json:"result"`
	Err              string `json:"err"`
}

type MessageModel struct {
	Receiver int64  `json:"receiver"`
	Encoding string `json:"encoding"`
	Message  string `json:"message"`
}

func (model MessageModel) receiverStr() string {
	return fmt.Sprintf("%d", model.Receiver)
}

func (model MessageModel) decode() (MessageModel, error) {
	if model.Encoding == "none" {
		return model, nil
	} else if model.Encoding == "base64" {
		bytes, err := base64.StdEncoding.DecodeString(model.Message)
		if err != nil {
			return model, err
		}

		return MessageModel{
			Message:  string(bytes),
			Receiver: model.Receiver,
			Encoding: "none",
		}, nil
	} else {
		return model, errors.New(fmt.Sprintf("Encoding[%s] not support", model.Encoding))
	}

}

func Kafka() {
	var w *kafka.Writer
	var r *kafka.Reader
	if InitCore.IsWithCredential() {
		mechanism, err := scram.Mechanism(scram.SHA256, InitCore.GetUsername(), InitCore.GetPassword())
		if err != nil {
			log.Fatalln(err)
		}

		dialer := &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		}

		w = &kafka.Writer{
			Addr:     kafka.TCP(InitCore.GetKafkaHost()),
			Balancer: &kafka.LeastBytes{},
			Transport: &kafka.Transport{
				SASL: mechanism,
				TLS:  &tls.Config{},
			},
		}

		r = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{InitCore.GetKafkaHost()},
			GroupID: InitCore.GetConsumerGroup(),
			Topic:   InitCore.GetSourceTopic(),
			Dialer:  dialer,
		})

	}

	for {
		var result = false
		var errMsg = ""
		var messageModel MessageModel
		var startTime time.Time
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			errMsg = err.Error()
			goto sendEvent
		}
		startTime = time.Now()

		if err = json.Unmarshal(m.Value, &messageModel); err != nil {
			errMsg = err.Error()
			goto sendEvent
		}
		messageModel, err = messageModel.decode()
		if err != nil {
			errMsg = err.Error()
			goto sendEvent
		}

		err = TgSend.Send(messageModel.receiverStr(), messageModel.Message, InitCore.GetTgToken())
		if err != nil {
			errMsg = err.Error()
			goto sendEvent
		}

		result = true

	sendEvent:

		now := time.Now()
		elapse := now.Sub(startTime)

		dto := NotificationEventDto{
			Topic:            m.Topic,
			Partition:        m.Partition,
			Offset:           m.Offset,
			CompleteTime:     now.Format(time.RFC3339),
			ProcessingPeriod: elapse.Milliseconds(),
			Result:           result,
			Err:              errMsg,
		}
		source := fmt.Sprintf("%s/%d/%d", m.Topic, m.Partition, m.Offset)

		dtoByte, err := json.Marshal(dto)

		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Topic: InitCore.GetResultTopic(),
				Value: dtoByte,
				Headers: []kafka.Header{
					{
						Key:   "version",
						Value: []byte("1.0.0"),
					},
					{
						Key:   "source",
						Value: []byte(source),
					},
				},
			},
		)
		if err != nil {
			var std = log.New(os.Stderr, "", log.LstdFlags)
			std.Output(2, fmt.Sprintln(err))
		}
	}
}

func main() {
	InitCore.CmdPreparation()
	Kafka()
}

func init() {
	tz := InitCore.GetOsEnv("X_TIMEZONE", "Asia/Hong_Kong")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
	time.Local = loc
}
