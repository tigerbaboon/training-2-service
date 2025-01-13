package user

import (
	userdto "app/app/modules/user/dto"
	"app/app/modules/user/entity"
	hashpassword "app/helper/hashPassword"
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *bun.DB
}

func newService(db *bun.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (service *UserService) User() *userdto.UserResponse {
	return &userdto.UserResponse{}
}

func (service *UserService) CreateUser(ctx context.Context, req *userdto.UserDTORequest) error {
	hashPassword, err := hashpassword.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entity.UserEntity{
		Username:  req.Username,
		Password:  hashPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	_, err = service.db.NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) Login(ctx context.Context, reqlogin userdto.UserLoginRequest) (*userdto.UserResponse, error) {
	var user entity.UserEntity
	err := service.db.NewSelect().
		Model(&user).
		Where("username = ?", reqlogin.Username).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqlogin.Password))
	if err != nil {
		return nil, err
	}

	sampleSecret := []byte(viper.GetString("MY_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID,
		"userType": "user",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(sampleSecret)
	if err != nil {
		return nil, err
	}

	userResponse := &userdto.UserResponse{}
	copier.Copy(userResponse, &user)

	userResponse.Token = tokenString

	return userResponse, nil
}

func (service *UserService) UpdateUser(ctx context.Context, username string, req *userdto.UserUpdateRequest) error {
	var user entity.UserEntity
	err := service.db.NewSelect().
		Model(&user).
		Where("username = ?", username).
		Scan(ctx)
	if err != nil {
		return err
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	user.UpdatedAt = time.Now().Unix()

	_, err = service.db.NewUpdate().
		Model(&user).
		Where("username = ?", username).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) GetUser(ctx context.Context, userx string) (*userdto.UserResponseAll, error) {
	var user entity.UserEntity

	err := service.db.NewSelect().
		Model(&user).
		Where("id = ?", userx).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	userResponse := &userdto.UserResponseAll{}
	copier.Copy(userResponse, &user)

	return userResponse, nil

}

func (service *UserService) GetUserByID(ctx context.Context, userID string) (*userdto.UserResponseAll, error) {
	var user userdto.UserResponseAll

	err := service.db.NewSelect().
		TableExpr("userx as u").
		Column("u.id", "u.username", "u.email", "u.first_name", "u.last_name").
		Where("u.id = ?", userID).
		Scan(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
