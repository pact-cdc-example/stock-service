package product_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gofiber/fiber/v2"
	"github.com/pact-cdc-example/stock-service/app/product"
	"github.com/pact-cdc-example/stock-service/pkg/cerr"
	"github.com/pact-cdc-example/stock-service/pkg/httpclient"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

const (
	getProductByIDPath     = "/api/v1/products/%s"
	consumerStockService   = "StockService"
	providerProductService = "ProductService"
)

type ProductConsumerTestSuite struct {
	suite.Suite
	pact          *dsl.Pact
	client        product.Client
	pactServerURL string
}

func TestProductConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductConsumerTestSuite))
}

func (s *ProductConsumerTestSuite) SetupSuite() {
	s.initPact()

	s.client = product.NewClient(&product.NewClientOpts{
		HTTPClient: httpclient.New(),
		BaseURL:    s.pactServerURL,
	})
}

func (s *ProductConsumerTestSuite) TearDownSuite() {
	defer s.pact.Teardown()
}

func (s *ProductConsumerTestSuite) TestGivenGetProductByIDReqThenItShouldReturnProductNotFoundErrWhenProductWithGivenIDNotExists() {
	givenProductID := gofakeit.UUID()

	s.pact.AddInteraction().
		Given("i get product not found error when the product with given id does not exists").
		UponReceiving("A request for product with a non exist product id").
		WithRequest(dsl.Request{
			Method:  http.MethodGet,
			Path:    dsl.String(fmt.Sprintf(getProductByIDPath, givenProductID)),
			Headers: map[string]dsl.Matcher{},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusBadRequest,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"code":    20001,
				"message": "Product not found.",
			},
		})

	var test = func() error {
		_, err := s.client.GetProductByID(context.Background(), givenProductID)
		return err
	}

	err := s.pact.Verify(test)

	s.Equal(err, cerr.Bag{Code: 20001, Message: "Product not found."})
}

func (s *ProductConsumerTestSuite) TestGivenGetProductByIDReqThenItShouldReturnProductWhenProductWithGivenIDExists() {
	givenProductID := gofakeit.UUID()

	givenProduct := product.Product{
		ID:        givenProductID,
		Name:      gofakeit.Name(),
		Code:      gofakeit.Word(),
		Color:     gofakeit.Color(),
		CreatedAt: gofakeit.Date(),
		UpdatedAt: gofakeit.Date(),
		Price:     gofakeit.Price(10, 100),
		ImageURL:  gofakeit.ImageURL(200, 100),
		Type:      gofakeit.Word(),
	}

	s.pact.AddInteraction().
		Given("i get product with given id").
		UponReceiving("A request for product with a exist product id").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String(fmt.Sprintf(getProductByIDPath, givenProductID)),
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"id":         dsl.Like(givenProduct.ID),
				"name":       dsl.Like(givenProduct.Name),
				"code":       dsl.Like(givenProduct.Code),
				"color":      dsl.Like(givenProduct.Color),
				"created_at": dsl.Like(givenProduct.CreatedAt),
				"updated_at": dsl.Like(givenProduct.UpdatedAt),
				"price":      dsl.Like(givenProduct.Price),
				"image_url":  dsl.Like(givenProduct.ImageURL),
				"type":       dsl.Like(givenProduct.Type),
			},
		})

	var test = func() error {
		_, err := s.client.GetProductByID(context.Background(), givenProductID)
		return err
	}

	err := s.pact.Verify(test)

	s.Nil(err)
}

func (s *ProductConsumerTestSuite) initPact() {
	s.pact = &dsl.Pact{
		Host:     "localhost",
		Consumer: consumerStockService,
		Provider: providerProductService,

		DisableToolValidityCheck: true,
		PactFileWriteMode:        "overwrite",
		LogDir:                   "./pacts/logs",
	}
	//it must be used otherwise it could not create pact file
	s.pact.Setup(true)

}
