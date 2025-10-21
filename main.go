package paystore

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"paystore/config"
	"paystore/lib/def"
	"paystore/lib/helper"
	"paystore/lib/request"
	vendorModel "paystore/user/vendors/model"
	vendorRequest "paystore/user/vendors/request"
)

func main() {
	writeDB := helper.CreatePostgresConnection(
		os.Getenv("DB_WRITE_HOST"), os.Getenv("DB_WRITE_PORT"), os.Getenv("DB_WRITE_USER"),
		os.Getenv("DB_WRITE_PASSWORD"), os.Getenv("DB_WRITE_NAME"), os.Getenv("DB_WRITE_SSLMODE"))
	defer writeDB.Close()
	readDB := helper.CreatePostgresConnection(
		os.Getenv("DB_READ_HOST"), os.Getenv("DB_READ_PORT"), os.Getenv("DB_READ_USER"),
		os.Getenv("DB_READ_PASSWORD"), os.Getenv("DB_READ_NAME"), os.Getenv("DB_READ_SSLMODE"))
	defer readDB.Close()
	redis := helper.ConnectRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_USER"),
		os.Getenv("REDIS_PASS"), false)

	config := config.DefaultConfig(os.Getenv("VENDOR_TABLE_NAME"))

	client := New(writeDB, readDB, redis, config)
	httpClient := NewHTTPHandler(client)

	app := fiber.New()
	app.Post("/webhook/xendit", httpClient.ReceivePayment)

	app.Listen(":" + os.Getenv("PORT"))
}

type HTTPCommandHandler struct {
	paystoreClient *PaystoreClient
}

func (h *HTTPCommandHandler) CreateBalance(c *fiber.Ctx) error {
	var requestBody request.CreateBalanceRequest

	return nil
}

func (h *HTTPCommandHandler) CreateOrganization(c *fiber.Ctx) error {
	var requestBody request.CreateOrganizationRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return helper.ReturnErrorResponse(c, fiber.StatusBadRequest, err, "E4001")
	}

	return nil
}

func (h *HTTPCommandHandler) CreatePin(c *fiber.Ctx) error {

	return nil
}

func (h *HTTPCommandHandler) RequestWithdraw(c *fiber.Ctx) error {

	return nil
}

func (h *HTTPCommandHandler) ReceivePayment(c *fiber.Ctx) error {
	// TODO: Implement your custom payment processing logic here
	// This function serves as a template for handling payment webhooks

	var requestBody vendorRequest.XenditReceivePayment
	if err := c.BodyParser(&requestBody); err != nil {
		return err
	}

	vendorItem := vendorModel.NewVendor()
	vendorItem.ID = requestBody.ID
	vendorItem.ExternalID = requestBody.ExternalID
	vendorItem.UserID = requestBody.UserID
	vendorItem.Status = requestBody.Status
	vendorItem.Amount = requestBody.Amount
	vendorItem.PaidAmount = requestBody.PaidAmount
	vendorItem.AdjustedReceivedAmount = requestBody.AdjustedReceivedAmount
	vendorItem.FeesPaidAmount = requestBody.FeesPaidAmount
	vendorItem.PaidAt = requestBody.PaidAt
	vendorItem.Created = requestBody.Created
	vendorItem.Updated = requestBody.Updated
	vendorItem.Currency = requestBody.Currency
	vendorItem.BankCode = requestBody.BankCode
	vendorItem.PaymentMethod = requestBody.PaymentMethod
	vendorItem.PaymentChannel = requestBody.PaymentChannel
	vendorItem.PaymentDestination = requestBody.PaymentDestination

	receivePaymentRequest := request.ReceivePaymentRequest{
		AccountUUID:    vendorItem.ExternalID,
		Amount:         vendorItem.Amount,
		VendorRecordID: vendorItem.ID,
	}

	newPayment, errCreatePayment := h.paystoreClient.ReceivePayment(receivePaymentRequest, vendorItem)
	if errCreatePayment != nil {
		statusResponse := fiber.StatusInternalServerError
		appCode := "E5001"
		if errCreatePayment == def.BalanceNotFound {
			statusResponse = fiber.StatusNotFound
			appCode = "E4001"
		}
		if errCreatePayment == def.OrganizationNotFound {
			statusResponse = fiber.StatusNotFound
			appCode = "E4002"
		}

		return helper.ReturnErrorResponse(c, statusResponse, errCreatePayment, appCode)
	}

	return c.JSON(map[string]string{"uuid": newPayment.GetUUID()})
}

func NewHTTPHandler(paystoreClient *PaystoreClient) *HTTPCommandHandler {
	return &HTTPCommandHandler{
		paystoreClient: paystoreClient,
	}
}

type HTTPFetcherHandler struct {
	paystoreFetcher *PaystoreFetcher
}
