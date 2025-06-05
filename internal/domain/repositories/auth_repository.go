package repositories

import (
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/domain/models"
)

type CustomerRepository interface {
	GetByOAuthID(oauthID string) (*models.Customer, error)
	GetByID(id uint) (*models.Customer, error)
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
	Delete(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) GetByOAuthID(oauthID string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.Where("oauth_id = ?", oauthID).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Customer{}, id).Error
}
