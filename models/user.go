package models

type User struct {
	ID    uint
	Alias *string `json:"alias,omitempty" gorm:"uniqueIndex, size:64"`
	//ParentId       *uint   `json:"parentId" gorm:"index"`
	Email     string  `json:"email" gorm:"size:256"`
	Verified  bool    `json:"verified"`
	Password  string  `json:"-,omitempty"`
	FirstName *string `json:"firstName,omitempty" gorm:"size:128"`
	//MiddleName     *string `json:"middleName,omitempty" gorm:"size:128"`
	Lastname *string `json:"lastame,omitempty" gorm:"size:128"`
	//Phone          *string `json:"phone,omitempty" gorm:"index,size:32"`
	//Lang           string  `json:"lang" gorm:"size:8"`
	Role    string `json:"role" gorm:"size:64"`
	Blocked bool   `json:"blocked"`
	//OrganizationId uint    `json:"organizationId"`
}
