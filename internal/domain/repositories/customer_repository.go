package repositories

import (
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/domain/models"
)

type CustomerRepositorys interface {
	GetByOAuthID(oauthID string) (*models.Customer, error)
	GetByID(id uint) (*models.Customer, error)
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
}

type customerRepositorys struct {
	db *gorm.DB
}

//func NewCustomerRepository(db *gorm.DB) CustomerRepository {
//	return &customerRepository{db: db}
//}

func (r *customerRepository) GetByOAuthIDs(oauthID string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.Where("oauth_id = ?", oauthID).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByIDs(id uint) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Creates(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) Updates(customer *models.Customer) error {
	return r.db.Save(customer).Error
}
