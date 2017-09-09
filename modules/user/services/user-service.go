package services

import (
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/wen-bing/go-enterprise-web-toolkit/core/db"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	s := UserService{}
	s.db = db
	return &s
}

func (s *UserService) CreateUser(user models.User) models.UserBasicObject {
	user.Id = db.GenerateModelId()
	bytePassword := []byte(user.Password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(passwordHash)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	//save to
	s.db.Create(user)
	obj := models.UserModel2BasicObject(user)
	return obj
}

func (s *UserService) GetUserById(id string) (models.UserBasicObject, error) {
	var user models.User
	var obj models.UserBasicObject
	e := s.db.First(&user, "id = ?", id).Error
	obj = models.UserModel2BasicObject(user)
	return obj, e
}

func (s *UserService) Login(userId, password string) (string, error) {
	user := models.User{}
	err := checkmail.ValidateFormat(userId)
	if err == nil {
		e := s.db.Where(&models.User{Email: userId}).First(&user).Error
		if e != nil {
			return "", e
		}
	} else {
		e := s.db.Where(&models.User{Phone: userId}).First(&user).Error
		if e != nil {
			return "", e
		}
	}
	err = verifyPassword(password, user)
	return user.Id, err
}

func verifyPassword(password string, user models.User) error {
	bytePassword := []byte(password)
	hashedPassword := []byte(user.Password)
	return bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
}
