package models

type Application struct {
	ID             uint32
	DisplayName    *string       `json:"displayName" gorm:"not null,size:128"`
	Description    string        `json:"description" gorm:"size:128"`
	OrganizationId *uint         `json:"organizationId"`
	Organization   *Organization `json:"organization"`
	ClientId       string        `json:"clientId" gorm:"size:128"`
}
