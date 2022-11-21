package repositories

import (
	"errors"

	"github.com/shkiperko0/auth-go-ms/models"

	// "fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "gorm.io/datatypes"
	"gorm.io/gorm"
)

//type ListOptions struct{}
//type MUserFilter struct{}

type IUserRepository interface {
	// GetStructureCount(id uint, currentDepth *int, maxDepth *int) int64
	// FirstLineCount(id *uint) int64
	// ListIds() (*[]uint, error)
	// Create(user *models.User) error
	// Update(user *models.User) error
	GetById(id uint) (*models.User, error)
	// GetByAlias(alias string, org uint) (*models.User, error)
	// GetByEmail(email string, organizationId uint) (*models.User, error)
	// GetByExtraKeyValue(organizationId uint, key, value string) (*models.User, error)
	// List(opt ListOptions, filter *MUserFilter) (*[]models.User, error)
	// Count(filter *MUserFilter) int64
	// GetUserByField(organizationId uint, field string, value string) (*models.User, error)
	// Delete(id uint64) error
}

func (r *UserRepository) GetById(id uint) (*models.User, error) {
	var item *models.User
	res := r.DB.First(&item, id)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "user.userNotFound")
	}

	return item, res.Error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) IUserRepository {
	return &UserRepository{
		DB: DB,
	}
}

// func (r *UserRepository) Delete(id uint64) error {
// 	result := r.DB.Delete(&models.User{}, id)
// 	return result.Error
// }

// func (r *UserRepository) GetByAlias(alias string, org uint) (*models.User, error) {
// 	var item *models.User
// 	res := r.DB.Where("alias = ? AND organization_id = ?", alias, org).First(&item)
// 	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 		return nil, status.Errorf(codes.NotFound, "user.userNotFound")
// 	}

// 	return item, res.Error
// }

// func (r *UserRepository) GetByExtraKeyValue(organizationId uint, key, value string) (*models.User, error) {
// 	var item models.User
// 	res := r.DB.Where("organization_id = ?", organizationId).First(&item, datatypes.JSONQuery("extra").Equals(value, key))
// 	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 		return nil, status.Errorf(codes.NotFound, "user.userNotFound")
// 	}

// 	return &item, res.Error
// }

// func (r *UserRepository) ListIds() (*[]uint, error) {
// 	var items []uint
// 	res := r.DB.Table("users").Select("id").Find(&items)
// 	return &items, res.Error
// }

// func (r *UserRepository) Create(user *models.User) error {
// 	result := r.DB.Create(&user)
// 	return result.Error
// }

// func (r *UserRepository) Update(user *models.User) error {
// 	result := r.DB.Save(&user)
// 	return result.Error
// }

// func (r *UserRepository) GetByEmail(email string, organizationId uint) (*models.User, error) {
// 	var item *models.User
// 	res := r.DB.Where("email = ? AND organization_id = ?", email, organizationId).First(&item)
// 	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 		return nil, status.Errorf(codes.NotFound, "user.userNotFound")
// 	}

// 	return item, res.Error
// }

// func (r *UserRepository) GetUserByField(organizationId uint, field string, value string) (*models.User, error) {
// 	var user models.User

// 	res := r.DB.Where(fmt.Sprintf("organization_id = ? AND %s = ?", field), organizationId, value).First(&user)

// 	if res.Error != nil {
// 		return nil, res.Error
// 	}

// 	return &user, res.Error
// }
