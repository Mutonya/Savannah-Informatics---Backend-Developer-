package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/domain/models"
	"github.com/Mutonya/Savanah/internal/domain/repositories"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	mock     sqlmock.Sqlmock
	repo     repositories.CustomerRepository
	testCust *models.Customer
}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	var (
		db  *gorm.DB
		err error
	)

	// Create mock database connection
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.mock = mock
	suite.repo = repositories.NewCustomerRepository(db)

	suite.testCust = &models.Customer{
		Model:     gorm.Model{ID: 1},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+1234567890",
		OAuthID:   "oauth123",
	}
}

func (suite *CustomerRepositoryTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestCustomerRepositorySuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}

func (suite *CustomerRepositoryTestSuite) TestGetByOAuthID_Success() {
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "phone", "oauth_id"}).
		AddRow(suite.testCust.ID, suite.testCust.FirstName, suite.testCust.LastName,
			suite.testCust.Email, suite.testCust.Phone, suite.testCust.OAuthID)

	suite.mock.ExpectQuery(`SELECT \* FROM "customers" WHERE oauth_id = \$1`).
		WithArgs(suite.testCust.OAuthID).
		WillReturnRows(rows)

	customer, err := suite.repo.GetByOAuthID(suite.testCust.OAuthID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testCust.ID, customer.ID)
	assert.Equal(suite.T(), suite.testCust.Email, customer.Email)
}

func (suite *CustomerRepositoryTestSuite) TestGetByOAuthID_NotFound() {
	suite.mock.ExpectQuery(`SELECT \* FROM "customers" WHERE oauth_id = \$1`).
		WithArgs("nonexistent").
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := suite.repo.GetByOAuthID("nonexistent")

	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *CustomerRepositoryTestSuite) TestGetByID_Success() {
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email"}).
		AddRow(suite.testCust.ID, suite.testCust.FirstName, suite.testCust.LastName, suite.testCust.Email)

	suite.mock.ExpectQuery(`SELECT \* FROM "customers" WHERE "customers"."id" = \$1`).
		WithArgs(suite.testCust.ID).
		WillReturnRows(rows)

	customer, err := suite.repo.GetByID(suite.testCust.ID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.testCust.ID, customer.ID)
}

func (suite *CustomerRepositoryTestSuite) TestCreate_Success() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "customers"`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			suite.testCust.FirstName,
			suite.testCust.LastName,
			suite.testCust.Email,
			suite.testCust.Phone,
			"", // address
			suite.testCust.OAuthID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	err := suite.repo.Create(suite.testCust)

	assert.NoError(suite.T(), err)
}

func (suite *CustomerRepositoryTestSuite) TestUpdate_Success() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`UPDATE "customers"`).
		WithArgs(
			sqlmock.AnyArg(), // updated_at
			suite.testCust.FirstName,
			suite.testCust.LastName,
			suite.testCust.Email,
			suite.testCust.Phone,
			"", // address
			suite.testCust.OAuthID,
			suite.testCust.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Update(suite.testCust)

	assert.NoError(suite.T(), err)
}

func (suite *CustomerRepositoryTestSuite) TestDelete_Success() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`DELETE FROM "customers"`).
		WithArgs(suite.testCust.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Delete(suite.testCust.ID)

	assert.NoError(suite.T(), err)
}
