package request

type CreateBalanceRequest struct {
	OwnerID            string `json:"ownerID" binding:"required"`
	OrganizationRandId string `json:"organizationRandId" binding:"required"`
	Currency           string `json:"currency" binding:"required"`
}
