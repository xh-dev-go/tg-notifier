package init_core

import (
	"errors"
	"flag"
	"fmt"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagBool"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagString"
	"os"
	"strings"
)

var listOfEnvs = []string{
	"X_TIMEZONE",
	strings.ToUpper("X_" + tgToken),
	strings.ToUpper("X_" + kafkaHost),
	strings.ToUpper("X_" + kafkaWithCredential),
	strings.ToUpper("X_" + kafkaUsername),
	strings.ToUpper("X_" + kafkaPassword),
	strings.ToUpper("X_" + kafkaConsumerGroup),
	strings.ToUpper("X_" + kafkaSourceTopic),
	strings.ToUpper("X_" + kafkaResultTopic),
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
	configWithEnv       = "config-with-env"
)

var listOfArgs = []string{
	configWithEnv,
	tgToken,
	kafkaHost,
	kafkaWithCredential,
	kafkaUsername,
	kafkaPassword,
}

var withEnv *flagBool.BoolParam
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

func envConfig() {

	var boolValue bool

	valueT := GetOsEnv(strings.ToUpper("X_"+tgToken), "")
	token.SetValue(&valueT)

	valueH := GetOsEnv(strings.ToUpper("X_"+kafkaHost), "")
	host.SetValue(&valueH)

	valueWC := GetOsEnv(strings.ToUpper("X_"+kafkaWithCredential), "")
	if valueWC == "true" {
		boolValue = true
	} else {
		boolValue = false
	}
	withCredential.SetValue(&boolValue)

	valueUN := GetOsEnv(strings.ToUpper("X_"+kafkaUsername), "")
	un.SetValue(&valueUN)

	valuePWD := GetOsEnv(strings.ToUpper("X_"+kafkaPassword), "")
	password.SetValue(&valuePWD)

	valueCG := GetOsEnv(strings.ToUpper("X_"+kafkaConsumerGroup), DefaultConsumerGroup)
	consumerGroup.SetValue(&valueCG)

	valueST := GetOsEnv(strings.ToUpper("X_"+kafkaSourceTopic), "")
	sourceTopic.SetValue(&valueST)

	valueRT := GetOsEnv(strings.ToUpper("X_"+kafkaResultTopic), "")
	resultTopic.SetValue(&valueRT)
}

const DefaultConsumerGroup = "go-telegram-bot-sender"

func CmdPreparation() {
	token = flagString.New(tgToken, "the token used for sending notification").BindCmd()
	host = flagString.New(kafkaHost, "the host to kafka").BindCmd()
	withCredential = flagBool.NewDefault(kafkaWithCredential, "does the kafka connection with a credential", false).BindCmd()
	un = flagString.New(kafkaUsername, "does the kafka connection username").BindCmd()
	password = flagString.New(kafkaPassword, "does the kafka connection password").BindCmd()
	consumerGroup = flagString.NewDefault(kafkaConsumerGroup, DefaultConsumerGroup, "consumer group id").BindCmd()
	sourceTopic = flagString.New(kafkaSourceTopic, "source topic which provide the source of notification").BindCmd()
	resultTopic = flagString.New(kafkaResultTopic, "result topic which result event submitting to").BindCmd()
	withEnv = flagBool.NewDefault(configWithEnv, "configuration with environment variables", false).BindCmd()
	flag.Parse()

	println("-----------------")
	for _, v := range os.Args {
		println(v)
	}
	println("-----------------")
	if withEnv.Value() {
		envConfig()
	}

	verification()

	fmt.Println("=======================")
	fmt.Printf("Config with env: %t\n", withEnv.Value())
	fmt.Printf("Host: %s\n", host.Value())
	fmt.Printf("WithCredential: %t\n", withCredential.Value())
	if withCredential.Value() {
		fmt.Printf("Username: %s\n", un.Value())
	}
	fmt.Printf("Consumer group: %s\n", consumerGroup.Value())
	fmt.Printf("Source Topic: %s\n", sourceTopic.Value())
	fmt.Printf("Result Topic: %s\n", resultTopic.Value())
}
