package operation

import (
	"context"
	"paystore/lib/def"
	pb "paystore/protos"
)

// GRPCHandler
/*
	- CreateBalance
	- CreatePayment
	- FinalizedPayment
	- CreateWithdraw
	- FinalizedWithdraw
*/
func pbToGoPaymentStatus(pbStatus pb.PaymentStatus) def.PaymentStatus {
	switch pbStatus {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		return def.PaymentStatusPending
	case pb.PaymentStatus_PAYMENT_STATUS_PAID:
		return def.PaymentStatusPaid
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return def.PaymentStatusFailed
	default:
		return def.PaymentStatusPending // or handle error
	}
}

func pbToGoWithdrawStatus(pbStatus pb.PaymentStatus) def.WithdrawStatus {
	switch pbStatus {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		return def.StatusPending
	case pb.PaymentStatus_PAYMENT_STATUS_PAID:
		return def.StatusSuccess
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return def.StatusFailed
	default:
		return def.StatusPending // or handle error
	}
}

type GRPCHandler struct {
	pb.UnimplementedPaystoreServer
	paystoreClient *PaystoreClient
}

func (grpc *GRPCHandler) CreateBalance(ctx context.Context, in *pb.CreateBalanceRequest) (*pb.CreatedResponse, error) {
	balance, errCreate := grpc.paystoreClient.CreateBalance(in.ExternalID, in.Currency, in.OrganizationSlug)
	if errCreate != nil {
		if errCreate == def.OrganizationNotFound {
			return nil, def.OrganizationNotFound
		}
		return nil, errCreate
	}

	return &pb.CreatedResponse{ID: balance.GetUUID()}, nil
}

func (grpc *GRPCHandler) CreatePayment(ctx context.Context, in *pb.CreatePaymentRequest) (*pb.CreatedResponse, error) {
	payment, errCreate := grpc.paystoreClient.CreatePayment(in.AccountUUID, in.Amount, in.VendorRecordId)
	if errCreate != nil {
		return nil, errCreate
	}

	return &pb.CreatedResponse{ID: payment.GetUUID()}, nil
}

func (grpc *GRPCHandler) FinalizedPayment(ctx context.Context, in *pb.FinalizedPaymentRequest) (*pb.EmptyResponse, error) {
	errFinalized := grpc.paystoreClient.FinalizedPayment(in.AccountUUID, in.PaymentUUID, pbToGoPaymentStatus(in.PaymentStatus))
	if errFinalized != nil {
		return nil, errFinalized
	}

	return &pb.EmptyResponse{}, nil
}

func (grpc *GRPCHandler) CreateWithdraw(ctx context.Context, in *pb.CreateWithdrawRequest) (*pb.CreatedResponse, error) {
	withdraw, errCreate := grpc.paystoreClient.CreateWithdraw(in.AccountUUID, in.Amount, in.VendorRecordId)
	if errCreate != nil {
		return nil, errCreate
	}

	return &pb.CreatedResponse{ID: withdraw.GetUUID()}, nil
}

func (grpc *GRPCHandler) FinalizedWithdraw(ctx context.Context, in *pb.FinalizedWithdrawRequest) (*pb.EmptyResponse, error) {
	errFinalized := grpc.paystoreClient.FinalizedWithdraw(in.WithdrawUUID, in.WithdrawUUID, pbToGoWithdrawStatus(in.WithdrawStatus))
	if errFinalized != nil {
		return nil, errFinalized
	}

	return &pb.EmptyResponse{}, nil
}

func NewGRPCHandler(paystoreClient *PaystoreClient) *GRPCHandler {
	return &GRPCHandler{
		paystoreClient: paystoreClient,
	}
}
