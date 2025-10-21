package request

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}
