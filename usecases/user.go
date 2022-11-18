package usecases

import (
	"eam-auth-go-ms/models"
	"eam-auth-go-ms/iteractor"
)

type IUserUseCase interface {
	Exist(id uint) bool
	Get(id uint) (*models.User, error)
	GetByEmail(email string, org string) (*models.User, error)
	GetByAlias(alias string, org string) (*models.User, error)
	Update(user *models.User) (string, error)
	UpdateByAdmin(user *models.User) error
	//Register(handler string, data *pb.RegAndLoginReq) (*models.User, error)
	VerifyEmail(id uint) error
	ChangeAlias(id uint, alias string, org string) error
	ChangeEmail(id uint, email string, org string) error
	GetByExtraKeyValue(organizationId uint, key, value string) (*models.User, error)
	//List(opt shared.ListOptions, filter *pb.UserFilterModel) (*[]models.User, error)
	//Count(filter *pb.UserFilterModel) int64
	VerifyByToken(token string) error
	ChangePassword(userId uint64, pass string, oldPass string) error
	Create(user *models.User) error
}

type UserUseCase struct {
	JwtIteractor       iteractor.JwtIteractor
	UserIteractor      iteractor.UserIteractor
}

func newUserUseCase(JwtIteractor iteractor.JwtIteractor, UserIteractor iteractor.UserIteractor) IUserUseCase {
	return &UserUseCase{
		JwtIteractor: JwtIteractor,
		UserIteractor: UserIteractor,
	}
}