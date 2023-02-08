package init_core

import (
	"errors"
	"flag"
	"fmt"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagBool"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagString"
	"os"
)

var listOfEnvs = []string{
	"X_TIMEZONE",
}

const (
	tgToken             = "tg-token"
	kafkaHost           = "kafka-host"
	kafkaWithCredential = "kafka-with-credential"
	kafkaUsername       = "kafka-username"
	kafkaPassword       = "kafka-password"
	kafkaConsumerGroup  = "kafka-consumer-group"
	kafkaSourceTopic    = "kafka-source-topic"
	kafkaResultTopic    = "kafka-result-topic"
)

var listOfArgs = []string{
	tgToken,
	kafkaHost,
	kafkaWithCredential,
	kafkaUsername,
	kafkaPassword,
}

var token *flagString.StringParam
var host *flagString.StringParam
var withCredential *flagBool.BoolParam
var un *flagString.StringParam
var password *flagString.StringParam
var consumerGroup *flagString.StringParam
var sourceTopic *flagString.StringParam
var resultTopic *flagString.StringParam

func verification() {
	if token.Value() == "" {
		panic(errors.New("telegram bot token not exists"))
	}
	if host.Value() == "" {
		panic(errors.New("kafka host not exists"))
	}
	if consumerGroup.Value() == "" {
		panic(errors.New("consumer group not exists"))
	}
	if sourceTopic.Value() == "" {
		panic(errors.New("source topic not exists"))
	}
	if resultTopic.Value() == "" {
		panic(errors.New("result topic not exists"))
	}
	if withCredential.Value() {
		if un.Value() == "" {
			panic(errors.New("kafka username is empty"))
		}
		if password.Value() == "" {
			panic(errors.New("kafka password is empty"))
		}
	}
}

func GetTgToken() string {
	return token.Value()
}
func GetKafkaHost() string {
	return host.Value()
}
func IsWithCredential() bool {
	return withCredential.Value()
}
func GetUsername() string {
	if IsWithCredential() {
		return un.Value()
	} else {
		panic(errors.New("no username should be required"))
	}
}
func GetPassword() string {
	if IsWithCredential() {
		return password.Value()
	} else {
		panic(errors.New("no password should be required"))
	}
}
func GetConsumerGroup() string {
	return consumerGroup.Value()
}
func GetSourceTopic() string {
	return sourceTopic.Value()
}
func GetResultTopic() string {
	return resultTopic.Value()
}

func GetOsEnv(key, defaultValue string) string {
	for _, v := range listOfEnvs {
		if v == key {
			value := os.Getenv(key)
			if value == "" {
				return defaultValue
			} else {
				return value
			}
		}
	}

	panic(errors.New(fmt.Sprintf("key[%s] is not allowed env variable", key)))
}

func CmdPreparation() {
	token = flagString.New(tgToken, "the token used for sending notification").BindCmd()
	host = flagString.New(kafkaHost, "the host to kafka").BindCmd()
	withCredential = flagBool.NewDefault(kafkaWithCredential, "does the kafka connection with a credential", false).BindCmd()
	un = flagString.New(kafkaUsername, "does the kafka connection username").BindCmd()
	password = flagString.New(kafkaPassword, "does the kafka connection password").BindCmd()
	consumerGroup = flagString.NewDefault(kafkaConsumerGroup, "go-telegram-bot-sender", "consumer group id").BindCmd()
	sourceTopic = flagString.New(kafkaSourceTopic, "source topic which provide the source of notification").BindCmd()
	resultTopic = flagString.New(kafkaResultTopic, "result topic which result event submitting to").BindCmd()
	flag.Parse()

	verification()
}
