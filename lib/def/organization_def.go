package def

import "errors"

var OrganizationMismatch = errors.New("Organization mismatch")
var OrganizationNotFound = errors.New("Organization not found")
var DuplicateSlug = errors.New("Duplicate slug")
var DuplicateName = errors.New("Duplicate name")

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}
