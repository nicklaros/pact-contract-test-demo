package test

import (
	"fmt"
	"log"
	"os"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

func getPactBrokerUsername() string {
	url := os.Getenv("PACT_BROKER_USERNAME")
	if url == "" {
		url = "pact"
	}

	return url
}

func getPactBrokerPassword() string {
	url := os.Getenv("PACT_BROKER_PASSWORD")
	if url == "" {
		url = "mysecret"
	}

	return url
}

func publishPact() error {
	var dir, _ = os.Getwd()
	var pactDir = fmt.Sprintf("%s/pacts", dir)

	p := dsl.Publisher{}

	commitHash, err := getCommitHash()
	if err != nil {
		return err
	}

	err = p.Publish(types.PublishRequest{
		PactURLs:        []string{pactDir},
		PactBroker:      getPactBrokerHost(),
		BrokerUsername:  getPactBrokerUsername(),
		BrokerPassword:  getPactBrokerPassword(),
		ConsumerVersion: commitHash,
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

	return err
}

func main() {
	publishPact()
}
