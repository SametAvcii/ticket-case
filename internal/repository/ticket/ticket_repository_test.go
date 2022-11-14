package repository

import (
	"context"
	"errors"
	"log"
	"regexp"
	models "samet-avci/gowit/models/ticket"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository ITicketRepository
}

func (s *Suite) SetupSuite() {
	var err error

	sqlDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	if err != nil {
		log.Fatalf("[sqlmock new] %s", err)
	}

	// create dialector
	dialector := mysql.New(mysql.Config{
		Conn:       sqlDB,
		DriverName: "mysql",
	})

	// a SELECT VERSION() query will be run when gorm opens the database
	// so we need to expect that here
	columns := []string{"version"}
	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows(columns).FromCSVString("1"),
	)

	// open the database
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("[gorm open] %s", err)
	}

	s.mock = mock
	s.repository = NewTicketRepository(db)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_Repository_IsDuplicate() {
	var (
		ID          uint = 1
		name             = "deneme"
		description      = "deneme"
		allocation  uint = 10
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `tickets` WHERE name = ? ORDER BY `tickets`.`id` LIMIT 1")).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "desc", "allocation"}).
			AddRow(ID, name, description, allocation))

	ctx := context.Background()
	IsDuplicate := s.repository.IsDuplicate(ctx, name)

	require.Equal(s.T(), IsDuplicate, true)
}

func (s *Suite) Test_Repository_IsDuplicateError() {
	var (
		name = "deneme"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `tickets` WHERE name = ? ORDER BY `tickets`.`id` LIMIT 1")).
		WithArgs(name).
		WillReturnError(errors.New("error test"))
	ctx := context.Background()
	IsDuplicate := s.repository.IsDuplicate(ctx, name)

	require.Equal(s.T(), IsDuplicate, false)
}
func (s *Suite) Test_Repository_CreateTicket() {

	req := models.Ticket{
		Name:       "deneme",
		Desc:       "deneme",
		Allocation: 10,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tickets` (`name`,`desc`,`allocation`) VALUES (?,?,?)")).
		WithArgs(req.Name, req.Desc, req.Allocation).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	ctx := context.Background()
	err := s.repository.CreateTicket(ctx, &req)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_Repository_CreateTicketError() {
	var (
		name        = "deneme"
		description = "deneme"
		allocation  = 10
	)

	req := &models.Ticket{
		Name:       name,
		Desc:       description,
		Allocation: uint(allocation),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tickets` (`name`,`desc`,`allocation`) VALUES (?,?,?)")).
		WithArgs(name, description, allocation).
		WillReturnError(errors.New("error test"))
	s.mock.ExpectRollback()

	ctx := context.Background()
	err := s.repository.CreateTicket(ctx, req)

	require.Error(s.T(), err)
}

func (s *Suite) Test_Repository_GetTicketByID() {
	var (
		ID          uint = 1
		name             = "deneme"
		description      = "deneme"
		allocation  uint = 10
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `tickets` WHERE id = ? ORDER BY `tickets`.`id` LIMIT 1")).
		WithArgs(ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "desc", "allocation"}).
			AddRow(ID, name, description, allocation))

	ctx := context.Background()
	res, err := s.repository.GetTicketByID(ctx, int(ID))

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(models.Ticket{ID: ID, Name: name, Desc: description, Allocation: allocation}, res))
}

func (s *Suite) Test_Repository_GetTicketByIDError() {
	var ID uint = 1

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `tickets` WHERE id = ? ORDER BY `tickets`.`id` LIMIT 1")).
		WithArgs(ID).
		WillReturnError(errors.New("error test"))

	ctx := context.Background()
	_, err := s.repository.GetTicketByID(ctx, int(ID))

	require.Error(s.T(), err)
}

func (s *Suite) Test_Repository_UpdateAllocation() {
	var (
		ID         uint = 1
		allocation uint = 10
	)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `tickets` SET `allocation`=? WHERE id = ?")).
		WithArgs(allocation, ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	ctx := context.Background()
	err := s.repository.UpdateAllocation(ctx, allocation, ID)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_Repository_UpdateAllocationError() {
	var (
		ID         uint = 1
		allocation uint = 10
	)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `tickets` SET `allocation`=? WHERE id = ?")).
		WithArgs(allocation, ID).
		WillReturnError(errors.New("test error"))
	s.mock.ExpectRollback()

	ctx := context.Background()
	err := s.repository.UpdateAllocation(ctx, allocation, ID)

	require.Error(s.T(), err)
}

func (s *Suite) Test_Repository_SaveSoldTicket() {

	req := models.SoldTicket{
		UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		Quantity: 10,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sold_tickets` (`user_id`,`quantity`) VALUES (?,?)")).
		WithArgs(req.UserID, req.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	ctx := context.Background()
	err := s.repository.SaveSoldTicket(ctx, req)

	require.NoError(s.T(), err)
}
func (s *Suite) Test_Repository_SaveSoldTicketError() {

	req := models.SoldTicket{
		UserID:   uuid.MustParse("406c1d05-bbb2-4e94-b183-7d208c2692e1"),
		Quantity: 10,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sold_tickets` (`user_id`,`quantity`) VALUES (?,?)")).
		WithArgs(req.UserID, req.Quantity).
		WillReturnError(errors.New("error test"))
	s.mock.ExpectRollback()

	ctx := context.Background()
	err := s.repository.SaveSoldTicket(ctx, req)

	require.Error(s.T(), err)
}
