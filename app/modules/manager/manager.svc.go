package manager

import (
	managerdto "app/app/modules/manager/dto"
	managerent "app/app/modules/manager/ent"
	hashpassword "app/helper/hashPassword"
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type ManagerService struct {
	db *bun.DB
}

func newService(db *bun.DB) *ManagerService {
	return &ManagerService{
		db: db,
	}
}

func (svc *ManagerService) CreateManager(ctx context.Context, req *managerdto.ManagerDTORequest) (*managerdto.ManagerDTOResponse, error) {
	hashpassword, err := hashpassword.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	manager := &managerent.Managers{
		Username:    req.Username,
		Password:    hashpassword,
		ManagerName: req.ManagerName,
	}

	_, err = svc.db.
		NewInsert().
		Model(manager).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	managerResponse := &managerdto.ManagerDTOResponse{}
	copier.Copy(managerResponse, manager)

	return managerResponse, nil
}

func (svc *ManagerService) UpdateManager(ctx context.Context, username string, req *managerdto.ManagerUpdateRequest) (*managerdto.ManagerDTOResponse, error) {
	var manager managerent.Managers
	err := svc.db.NewSelect().
		Model(&manager).
		Where("username = ?", username).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	manager.ManagerName = req.ManagerName

	_, err = svc.db.NewUpdate().
		Model(&manager).
		Where("username = ?", username).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	managerResponse := &managerdto.ManagerDTOResponse{}
	copier.Copy(managerResponse, &manager)

	return managerResponse, nil
}

func (svc *ManagerService) LoginManager(ctx context.Context, req *managerdto.ManagerLoginRequest) (*managerdto.ManagerLoginResponse, error) {
	var manager managerent.Managers
	err := svc.db.NewSelect().
		Model(&manager).
		Where("username = ?", req.Username).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(manager.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	sampleSecret := []byte(viper.GetString("MY_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   manager.ID,
		"userType": "manager",
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(sampleSecret)
	if err != nil {
		return nil, err
	}

	managerResponse := &managerdto.ManagerLoginResponse{}
	copier.Copy(managerResponse, &manager)

	managerResponse.Token = tokenString

	return managerResponse, nil
}

func (svc *ManagerService) GetManagerByID(ctx context.Context, mngID string) (*managerdto.ManagerDTOResponse, error) {
	var manager managerent.Managers
	err := svc.db.NewSelect().
		Model(&manager).
		Where("id = ?", mngID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	managerResponse := &managerdto.ManagerDTOResponse{}
	copier.Copy(managerResponse, &manager)

	return managerResponse, nil
}
