package app

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

func TestAuroraIndexTestSuite(t *testing.T) {
	suite.Run(t, new(AuroraIndexTestSuite))
}

type AuroraIndexTestSuite struct {
	suite.Suite

	index *AuroraIndex
	mock  sqlmock.Sqlmock
	db    *sql.DB
}

func (suite *AuroraIndexTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	suite.NoError(err)

	suite.mock = mock
	suite.db = db
	suite.index = &AuroraIndex{suite.db}
}

func (suite *AuroraIndexTestSuite) TearDownTest() {
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *AuroraIndexTestSuite) TestCreate() {
	suite.mock.
		ExpectBegin()
	prepHeader := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size) VALUES ($1, $2, $3)`).
		WillBeClosed()
	prepHeader.
		ExpectExec().
		WithArgs("docID", "docName", 1234).
		WillReturnResult(sqlmock.NewResult(0, 1))
	prepACL := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_acls (doc_id, username, read, write, delete, share) VALUES ($1, $2, $3, $4, $5, $6)`).
		WillBeClosed()
	prepACL.
		ExpectExec().
		WithArgs("docID", "username", true, true, true, true).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.
		ExpectCommit()

	suite.NoError(suite.index.Create(DocumentHeader{
		ID:   "docID",
		Name: "docName",
		Size: 1234,
	}, ACL{
		Permissions: map[string]Permission{
			"username": {
				Username: "username",
				Read:     true,
				Write:    true,
				Delete:   true,
				Share:    true,
			},
		},
	}))
}

func (suite *AuroraIndexTestSuite) TestCreateFailHeader() {
	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	prepHeader := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size) VALUES ($1, $2, $3)`).
		WillBeClosed()
	prepHeader.
		ExpectExec().
		WithArgs("docID", "docName", 1234).
		WillReturnError(testErr)
	suite.mock.
		ExpectRollback()

	suite.ErrorIs(suite.index.Create(DocumentHeader{
		ID:   "docID",
		Name: "docName",
		Size: 1234,
	}, ACL{
		Permissions: map[string]Permission{
			"username": {
				Username: "username",
				Read:     true,
				Write:    true,
				Delete:   true,
				Share:    true,
			},
		},
	}), testErr)
}

func (suite *AuroraIndexTestSuite) TestCreateFailACL() {
	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	prepHeader := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size) VALUES ($1, $2, $3)`).
		WillBeClosed()
	prepHeader.
		ExpectExec().
		WithArgs("docID", "docName", 1234).
		WillReturnResult(sqlmock.NewResult(0, 1))
	prepACL := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_acls (doc_id, username, read, write, delete, share) VALUES ($1, $2, $3, $4, $5, $6)`).
		WillBeClosed()
	prepACL.
		ExpectExec().
		WithArgs("docID", "username", true, true, true, true).
		WillReturnError(testErr)
	suite.mock.
		ExpectRollback()

	suite.ErrorIs(suite.index.Create(DocumentHeader{
		ID:   "docID",
		Name: "docName",
		Size: 1234,
	}, ACL{
		Permissions: map[string]Permission{
			"username": {
				Username: "username",
				Read:     true,
				Write:    true,
				Delete:   true,
				Share:    true,
			},
		},
	}), testErr)
}

func (suite *AuroraIndexTestSuite) TestTx() {
	suite.mock.
		ExpectBegin()
	suite.mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.
		ExpectCommit()

	suite.NoError(suite.index.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec("SELECT")
		return err
	}))
}

func (suite *AuroraIndexTestSuite) TestFailBegin() {
	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin().
		WillReturnError(testErr)

	suite.ErrorIs(suite.index.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec("SELECT")
		return err
	}), testErr)
}

func (suite *AuroraIndexTestSuite) TestFailCommit() {
	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	suite.mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.
		ExpectCommit().
		WillReturnError(testErr)

	suite.ErrorIs(suite.index.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec("SELECT")
		return err
	}), testErr)
}
