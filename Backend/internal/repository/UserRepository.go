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

	//more function for seller
	CreateBankAccount(bankAccount domain.BankAccount) (domain.BankAccount, error)

	//cart
	FindCartItems(uId uint) ([]domain.Cart, error)
	FindCartItem(uId uint, pId uint) (domain.Cart, error)
	CreateCart(cart domain.Cart) error
	UpdateCart(cart domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(uId uint) error

	//order
	CreateOrder(o domain.Order) error
	FindOrders(uId uint) ([]domain.Order, error)
	FindOrderByid(id uint, uId uint) (domain.Order, error)

	//profile
	CreateProfile(e domain.Address) error
	UpdateProfile(e domain.Address) error
}

// Private struct
type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) CreateOrder(o domain.Order) error {
	err := r.db.Create(&o).Error
	return err
}

func (r *userRepository) FindOrders(uId uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Where("user_id = ?", uId).Find(&orders).Error
	return orders, err
}

func (r *userRepository) FindOrderByid(id uint, uId uint) (domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("items").Where("id=? AND user_id = ?", id, uId).Find(&order).Error
	return order, err
}

func (r *userRepository) CreateProfile(e domain.Address) error {
	err := r.db.Create(&e).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateProfile(e domain.Address) error {
	err := r.db.Where("userId = ?", e.UserId).Updates(&e).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindCartItems(uId uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Where("user_id = ?", uId).Find(&carts).Error
	return carts, err
}

func (r *userRepository) FindCartItem(uId uint, pId uint) (domain.Cart, error) {
	var carts domain.Cart
	err := r.db.Where("user_id = ? AND id = ?", uId, pId).Find(&carts).Error
	return carts, err
}

func (r *userRepository) CreateCart(cart domain.Cart) error {
	err := r.db.Create(&cart).Error
	if err != nil {
		return err
	}
	return err
}

func (r *userRepository) UpdateCart(c domain.Cart) error {
	var cart domain.Cart
	err := r.db.Model(&cart).Where("id = ?", cart.ID).Updates(&c).Error
	if err != nil {
		return err
	}
	return err
}

func (r *userRepository) DeleteCartById(id uint) error {
	err := r.db.Where("id = ?", id).Delete(&domain.Cart{}).Error
	if err != nil {
		return err
	}
	return err
}

func (r *userRepository) DeleteCartItems(uId uint) error {
	err := r.db.Where("user_id = ?", uId).Delete(&domain.Cart{}).Error
	return err
}

func (r *userRepository) CreateBankAccount(e domain.BankAccount) (domain.BankAccount, error) {
	return e, r.db.Create(&e).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	fmt.Println("User successfully saved:", user)
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
	if err := r.db.Where("email = ?", email).Preload("Address").First(&user).Error; err != nil {
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
