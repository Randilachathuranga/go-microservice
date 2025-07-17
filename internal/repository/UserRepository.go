package repository

import (
	"fmt"
	"go-ecommerce-app/internal/domain"
	"gorm.io/gorm"
)

// Public interface
type UserRepository interface {
	CreateUser(user domain.User) (domain.User, error)
	FindUserById(userId uint) (domain.User, error)
	FindUserByEmail(email string) (domain.User, error)
	UpdateUser(id uint, user domain.User) (domain.User, error)
	DeleteUser(user domain.User) error
	ViewUser(user domain.User) (domain.User, error)
	ViewAllUsers() ([]domain.User, error)
}

// Private struct
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return domain.User{}, err
		fmt.Println("database connection established", user)
	}
	return user, nil
}

func (r *userRepository) FindUserById(userId uint) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, userId).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) FindUserByEmail(email string) (domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(id uint, user domain.User) (domain.User, error) {
	if err := r.db.Model(&domain.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return domain.User{}, err
	}
	return r.FindUserById(id)
}

func (r *userRepository) DeleteUser(user domain.User) error {
	return r.db.Delete(&user).Error
}

func (r *userRepository) ViewUser(user domain.User) (domain.User, error) {
	var result domain.User
	if err := r.db.First(&result, user.ID).Error; err != nil {
		return domain.User{}, err
	}
	return result, nil
}

func (r *userRepository) ViewAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
