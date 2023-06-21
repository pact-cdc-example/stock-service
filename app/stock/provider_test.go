package stock_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pact-cdc-example/stock-service/app/stock"
	"github.com/pact-cdc-example/stock-service/pkg/server"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
	"time"
)

const (
	pactBrokerLocalURL = "http://localhost"
)

type PactSettings struct {
	Host            string
	ProviderName    string
	BrokerBaseURL   string
	BrokerUsername  string // Basic authentication
	BrokerPassword  string // Basic authentication
	ConsumerName    string
	ConsumerVersion string // a git sha, semantic version number
	ConsumerTag     string // dev, staging, prod
	ProviderVersion string
}

func (s *PactSettings) getPactURL(useLocal bool) string {
	var pactURL string

	if s.ConsumerVersion == "" {
		pactURL = fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/latest/master.json", s.BrokerBaseURL, s.ProviderName, s.ConsumerName)
	} else {
		pactURL = fmt.Sprintf("%s/pacts/provider/%s/consumer/%s/version/%s.json", s.BrokerBaseURL, s.ProviderName, s.ConsumerName, s.ConsumerVersion)
	}

	return pactURL
}

type ProviderTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	pactSettings *PactSettings
	ctx          context.Context
	l            *logrus.Logger
	app          server.Server
	mockRepo     *stock.MockRepository
	serverPort   string
}

func TestProvider(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}

func (s *ProviderTestSuite) SetupSuite() {
	s.l, _ = test.NewNullLogger()
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = stock.NewMockRepository(s.ctrl)

	stockService := stock.NewService(&stock.NewServiceOpts{
		R: s.mockRepo,
		L: s.l,
	})

	stockHandler := stock.NewHandler(&stock.NewHandlerOpts{
		S: stockService,
		L: s.l,
	})

	sp, err := utils.GetFreePort()
	s.Nil(err)

	s.serverPort = fmt.Sprintf("%d", sp)

	s.app = server.New(&server.NewServerOpts{
		Port: s.serverPort,
	}, []server.RouteHandler{
		stockHandler,
	})

	//err = createProductTableOnDB(postgreDB)
	s.Nil(err)

	go func() {
		if serverErr := s.app.Run(); serverErr != nil {
			fmt.Println("serverErr", serverErr)
		}
	}()

	_ = os.Setenv("CONSUMER_NAME", "BasketService")
	_ = os.Setenv("CONSUMER_TAG", "dev")
	_ = os.Setenv("GIT_SHORT_HASH", "4.0.2")
	_ = os.Setenv("CONSUMER_VERSION", "4.0.2")
	s.pactSettings = &PactSettings{
		Host:            "localhost",
		ProviderName:    "StockService",
		ConsumerName:    os.Getenv("CONSUMER_NAME"),
		ConsumerVersion: os.Getenv("CONSUMER_VERSION"),
		BrokerBaseURL:   pactBrokerLocalURL,
		ConsumerTag:     os.Getenv("CONSUMER_TAG"),
		ProviderVersion: os.Getenv("GIT_SHORT_HASH"),
	}
	time.Sleep(3 * time.Second)
}

func (s *ProviderTestSuite) TestProvider() {
	pact := &dsl.Pact{
		Host:                     s.pactSettings.Host,
		Provider:                 s.pactSettings.ProviderName,
		Consumer:                 s.pactSettings.ConsumerName,
		DisableToolValidityCheck: true,
	}

	providerBaseURL := fmt.Sprintf("http://%s:%s", s.pactSettings.Host, s.serverPort)

	verifyRequest := types.VerifyRequest{
		ProviderBaseURL:            providerBaseURL,
		PactURLs:                   []string{s.pactSettings.getPactURL(true)},
		BrokerURL:                  s.pactSettings.BrokerBaseURL,
		Tags:                       []string{s.pactSettings.ConsumerTag},
		BrokerUsername:             s.pactSettings.BrokerUsername,
		BrokerPassword:             s.pactSettings.BrokerPassword,
		FailIfNoPactsFound:         true,
		PublishVerificationResults: true,
		ProviderVersion:            s.pactSettings.ProviderVersion,
		StateHandlers:              map[string]types.StateHandler{},
		BeforeEach:                 nil,
		AfterEach:                  nil,
		TagWithGitBranch:           false,
	}
	defer pact.Teardown()
	verifyResponses, err := pact.VerifyProvider(s.T(), verifyRequest)
	s.Nil(err)

	if err != nil {
		log.Println(err)
	}

	log.Printf("%d pact tests run", len(verifyResponses))
}
