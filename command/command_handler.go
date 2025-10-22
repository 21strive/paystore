package command

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"paystore/lib/def"
	"paystore/lib/helper"
	"paystore/lib/request"
	pb "paystore/protos"
	vendorModel "paystore/user/vendors/model"
	vendorRequest "paystore/user/vendors/request"
)

// HTTPCommandHandler
/*
	- UpdatePIN
	- CreateWithdraw
*/
type HTTPCommandHandler struct {
	paystoreClient *PaystoreClient
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

// GRPCHandler
/*
	- CreateBalance
	- CreatePayment
	- FinalizedPayment
*/
type GRPCHandler struct {
	pb.UnimplementedPaystoreServer
	paystoreClient *PaystoreClient
}

func (grpc *GRPCHandler) CreateBalance(ctx context.Context, in *pb.CreateBalanceRequest) (*pb.ReqeustResponse, error) {
	createBalanceRequest := request.CreateBalanceRequest{
		OwnerID:          in.OwnerID,
		Currency:         in.Currency,
		OrganizationSlug: in.OrganizationSlug,
	}

	balance, errCreate := grpc.paystoreClient.CreateBalance(createBalanceRequest)
	if errCreate != nil {
		if errCreate == def.OrganizationNotFound {
			return nil, def.OrganizationNotFound
		}
		return nil, errCreate
	}

	return &pb.ReqeustResponse{
		ID:     balance.GetUUID(),
		Status: pb.Status_SUCCES,
	}, nil
}
