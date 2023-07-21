package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"pact-contract-test-demo/common/pointer"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

func NewPactTestSuite(provider string) (*PactTestSuite, error) {
	commitHash, err := getCommitHash()
	if err != nil {
		return nil, err
	}

	branch, err := getBranch()
	if err != nil {
		return nil, err
	}

	ci := strings.ToLower(os.Getenv("CI")) == "true"

	return &PactTestSuite{
		CommitHash:         commitHash,
		Branch:             branch,
		ServiceHostBaseURL: getServiceHostBaseURL(),
		BrokerURL:          getPactBrokerHost(),
		Provider:           provider,
		PactURLs:           getPactURLs(),
		CI:                 ci,
	}, nil
}

func getServiceHostBaseURL() string {
	port := os.Getenv("HTTP_SERVER_PORT")
	if port == "" {
		port = "8087"
	}

	return fmt.Sprintf("http://localhost:%s", port)
}

func getServiceHost() string {
	return fmt.Sprintf("%s/rpc", getServiceHostBaseURL())
}

func getPactBrokerHost() string {
	url := os.Getenv("PACT_BROKER_URL")
	if url == "" {
		// TODO: change to http://localhost:9292 for default value when devops done adding PACT_BROKER_URL env variable in wakanda pipeline
		url = "https://pact.srcli.xyz.dmmy.me"
	}

	return url
}

func getCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.CombinedOutput()

	return strings.TrimSpace(string(output)), err
}

func getBranch() (string, error) {
	branch := os.Getenv("BRANCH_NAME")
	if branch != "" {
		return branch, nil
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()

	return strings.TrimSpace(string(output)), err
}

func getPactURLs() []string {
	pactURL := os.Getenv("PACT_URL")

	if pactURL == "" {
		return nil
	}

	return []string{
		pactURL,
	}
}

type PactTestSuite struct {
	CommitHash         string
	Branch             string
	ServiceHostBaseURL string
	BrokerURL          string
	Provider           string
	PactURLs           []string
	CI                 bool
}

func (suite *PactTestSuite) VerifyContract(t *testing.T, stateHandlers types.StateHandlers) error {
	pact := &dsl.Pact{
		Provider: suite.Provider,
	}

	verifyRequest := types.VerifyRequest{
		ProviderBaseURL:            suite.ServiceHostBaseURL,
		ProviderVersion:            suite.CommitHash,
		ProviderBranch:             suite.Branch,
		PublishVerificationResults: suite.CI,
		StateHandlers:              stateHandlers,
	}

	if suite.PactURLs != nil {
		verifyRequest.PactURLs = suite.PactURLs
	} else {
		verifyRequest.BrokerURL = suite.BrokerURL
		verifyRequest.EnablePending = true

		// To enable the WIP pacts feature, you need to set the includeWipPactsSince field to the date from which you want to start using this feature. eg "2020-10-31".
		// The date is required so that you don't suddenly start verifying 100 past feature pacts in your build all of a sudden.
		//
		// ref: https://docs.pact.io/pact_broker/advanced_topics/wip_pacts
		verifyRequest.IncludeWIPPactsSince = pointer.Pointer(time.Date(2023, 07, 05, 00, 00, 00, 0, time.UTC))

		verifyRequest.ConsumerVersionSelectors = []types.ConsumerVersionSelector{
			{
				MainBranch: true,
			},
			{
				DeployedOrReleased: true,
			},
		}
	}

	fmt.Println(fmt.Printf("VERIFY CONTRACT: (provider_version=%s, provider_branch=%s, pact_urls=%s)", suite.CommitHash, suite.Branch, strings.Join(suite.PactURLs, ",")))

	// Verify the Provider using published contracts on Pact Broker.
	_, err := pact.VerifyProvider(t, verifyRequest)

	return err
}
