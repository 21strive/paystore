package request

type ReceivePaymentRequest struct {
	AccountUUID    string `json:"account_uuid" binding:"required"`
	Amount         int64  `json:"amount" binding:"required"`
	VendorRecordID string `json:"vendor_record_id" binding:"required"`
}
