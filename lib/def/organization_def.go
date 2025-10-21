package def

import "errors"

var OrganizationMismatch = errors.New("Organization mismatch")

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}
