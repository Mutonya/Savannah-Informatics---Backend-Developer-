package repositories

import (
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/domain/models"
)

type CustomerRepository interface {
	GetByOAuthID(oauthID string) (*models.Customer, error) // get customer by OAuthId
	GetByID(id uint) (*models.Customer, error)             // get user by ID
	Create(customer *models.Customer) error                // create customer
	Update(customer *models.Customer) error                // update the customer
	Delete(id uint) error                                  // delete user with ID
}

type customerRepository struct {
	db *gorm.DB // db connection
}

//	initialize repo
//
// constructor Injection for DI
// hold the db privately {Encapsulation}
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// Primary lookup during authentication
func (r *customerRepository) GetByOAuthID(oauthID string) (*models.Customer, error) {
	var customer models.Customer // Initialize empty customer struct

	// Query: SELECT * FROM customers WHERE oauth_id = ? LIMIT 1

	if err := r.db.Where("oauth_id = ?", oauthID).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByID(id uint) (*models.Customer, error) {
	var customer models.Customer // Initialize empty customer struct

	// Query: SELECT * FROM customers WHERE id = ? LIMIT 1 {First automatically add LIMIT 1
	// First finds the first record ordered by primary key, matching given conditions conds
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// Create a new customer {Persist new customer to database}
// Create inserts value, returning the inserted data's primary key in value's id
// Auto-populates ID, CreatedAt, UpdatedAt
// Validates struct tags (e.g., gorm:"not null")
func (r *customerRepository) Create(customer *models.Customer) error {
	// Query: INSERT INTO customers (...) VALUES (...)
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(customer *models.Customer) error {
	// Save updates value in database. If value doesn't contain a matching primary key, value is inserted.
	// Query: UPDATE customers SET ... WHERE id = ?
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uint) error {
	// Delete deletes value matching given conditions.
	//If value contains primary key it is included in the conditions. If
	// value includes a deleted_at field, then
	//Delete performs a soft delete instead by setting deleted_at with the current
	// time if null.
	// Query: DELETE FROM customers WHERE id = ?
	return r.db.Delete(&models.Customer{}, id).Error
}
