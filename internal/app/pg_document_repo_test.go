package app

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

func TestPostgresDocumentRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresDocumentRepoTestSuite))
}

type PostgresDocumentRepoTestSuite struct {
	suite.Suite

	index *PostgresDocumentRepo
	mock  sqlmock.Sqlmock
	db    *sql.DB
}

func (suite *PostgresDocumentRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	suite.NoError(err)

	suite.mock = mock
	suite.db = db
	suite.index = &PostgresDocumentRepo{suite.db}
}

func (suite *PostgresDocumentRepoTestSuite) TearDownTest() {
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *PostgresDocumentRepoTestSuite) TestCreate() {
	docCreateTime := time.Now()

	suite.mock.
		ExpectBegin()
	prepHeader := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size, owner, created) VALUES ($1, $2, $3, $4, $5)`).
		WillBeClosed()
	prepHeader.
		ExpectExec().
		WithArgs("docID", "docName", 1234, "username", docCreateTime).
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
		ID:      "docID",
		Name:    "docName",
		Size:    1234,
		Owner:   "username",
		Created: docCreateTime,
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

func (suite *PostgresDocumentRepoTestSuite) TestCreateFailHeaderPrepare() {
	docCreateTime := time.Now()

	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size, owner, created) VALUES ($1, $2, $3, $4, $5)`).
		WillReturnError(testErr)
	suite.mock.
		ExpectRollback()

	suite.ErrorIs(suite.index.Create(DocumentHeader{
		ID:      "docID",
		Name:    "docName",
		Size:    1234,
		Owner:   "username",
		Created: docCreateTime,
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

func (suite *PostgresDocumentRepoTestSuite) TestCreateFailHeader() {
	docCreateTime := time.Now()

	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	prepHeader := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size, owner, created) VALUES ($1, $2, $3, $4, $5)`).
		WillBeClosed()
	prepHeader.
		ExpectExec().
		WithArgs("docID", "docName", 1234, "username", docCreateTime).
		WillReturnError(testErr)
	suite.mock.
		ExpectRollback()

	suite.ErrorIs(suite.index.Create(DocumentHeader{
		ID:      "docID",
		Name:    "docName",
		Size:    1234,
		Owner:   "username",
		Created: docCreateTime,
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

func (suite *PostgresDocumentRepoTestSuite) TestCreateFailACLPrepare() {
	docCreateTime := time.Now()

	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	prepHeader := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size, owner, created) VALUES ($1, $2, $3, $4, $5)`).
		WillBeClosed()
	prepHeader.
		ExpectExec().
		WithArgs("docID", "docName", 1234, "username", docCreateTime).
		WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.
		ExpectPrepare(`INSERT INTO au_document_acls (doc_id, username, read, write, delete, share) VALUES ($1, $2, $3, $4, $5, $6)`).
		WillReturnError(testErr)
	suite.mock.
		ExpectRollback()

	suite.ErrorIs(suite.index.Create(DocumentHeader{
		ID:      "docID",
		Name:    "docName",
		Size:    1234,
		Owner:   "username",
		Created: docCreateTime,
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

func (suite *PostgresDocumentRepoTestSuite) TestCreateFailACL() {
	docCreateTime := time.Now()

	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	prepHeader := suite.mock.
		ExpectPrepare(`INSERT INTO au_document_headers (doc_id, name, size, owner, created) VALUES ($1, $2, $3, $4, $5)`).
		WillBeClosed()
	prepHeader.
		ExpectExec().
		WithArgs("docID", "docName", 1234, "username", docCreateTime).
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
		ID:      "docID",
		Name:    "docName",
		Size:    1234,
		Owner:   "username",
		Created: docCreateTime,
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

func (suite *PostgresDocumentRepoTestSuite) TestTx() {
	suite.mock.
		ExpectBegin()
	suite.mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.
		ExpectCommit()

	suite.NoError(tx(suite.index.db, func(tx *sql.Tx) error {
		_, err := tx.Exec("SELECT")
		return err
	}))
}

func (suite *PostgresDocumentRepoTestSuite) TestFailBegin() {
	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin().
		WillReturnError(testErr)

	suite.ErrorIs(tx(suite.index.db, func(tx *sql.Tx) error {
		_, err := tx.Exec("SELECT")
		return err
	}), testErr)
}

func (suite *PostgresDocumentRepoTestSuite) TestFailCommit() {
	testErr := errors.New("test error")
	suite.mock.
		ExpectBegin()
	suite.mock.
		ExpectExec("SELECT").
		WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.
		ExpectCommit().
		WillReturnError(testErr)

	suite.ErrorIs(tx(suite.index.db, func(tx *sql.Tx) error {
		_, err := tx.Exec("SELECT")
		return err
	}), testErr)
}
