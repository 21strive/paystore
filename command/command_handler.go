package command

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
*/
type GRPCHandler struct {
	pb.UnimplementedPaystoreServer
	paystoreClient *PaystoreClient
}

func (grpc *GRPCHandler) CreateBalance(ctx context.Context, in *pb.CreateBalanceRequest) (*pb.RequestResponse, error) {

	balance, errCreate := grpc.paystoreClient.CreateBalance(in.ExternalID, in.Currency, in.OrganizationSlug)
	if errCreate != nil {
		if errCreate == def.OrganizationNotFound {
			return nil, def.OrganizationNotFound
		}
		return nil, errCreate
	}

	return &pb.RequestResponse{
		ID:     balance.GetUUID(),
		Status: pb.Status_SUCCES,
	}, nil
}
