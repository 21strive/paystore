package def

import "errors"

var OrganizationMismatch = errors.New("Organization mismatch")
var OrganizationNotFound = errors.New("Organization not found")

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}
