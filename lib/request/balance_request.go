package request

type CreateBalanceRequest struct {
	OwnerID          string `json:"ownerID" binding:"required"`
	OrganizationSlug string `json:"organizationSlug" binding:"required"`
	Currency         string `json:"currency" binding:"required"`
}
