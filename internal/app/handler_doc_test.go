package app

import (
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (suite *AppSuite) TestPostDocumentNoLogin() {
	suite.
		Request("POST", "/doc").
		Body(M{
			"filename": "myfile",
			"size":     1234,
		}).
		Expect(http.StatusUnauthorized, M{
			"message": "not logged in",
			"success": false,
		})
}

func (suite *AppSuite) TestPostDocument() {
	user := suite.login()

	testUUID := uuid.New()
	suite.app.genUUID = func() uuid.UUID {
		return testUUID
	}
	clock := SingleTimestampClock{time.Now()}
	suite.app.clock = clock

	suite.
		Request("POST", "/doc").
		Body(M{
			"filename": "myfile",
		}).
		Expect(http.StatusOK, M{
			"success": true,
			"id":      testUUID.String(),
		})

	// check that document is created correctly
	doc, err := suite.app.documents.Get(DocID(testUUID.String()))
	suite.NoError(err)
	suite.Equal(DocumentHeader{
		ID:      DocID(testUUID.String()),
		Name:    "myfile",
		Owner:   user,
		Created: clock.Timestamp,
		Updated: time.Time{},
	}, doc)

	// check that permissions are set up correctly
	acl, err := suite.app.documents.ACL(DocID(testUUID.String()))
	suite.NoError(err)
	suite.Equal(ACL{
		Permissions: map[string]Permission{
			user: {
				Username: user,
				Read:     true,
				Write:    true,
				Delete:   true,
				Share:    true,
			},
		},
	}, acl)
}

func (suite *AppSuite) TestPostContent() {
	user := suite.login()
	_ = user

	data := []byte("hello")
	testUUID := uuid.New()
	suite.app.genUUID = func() uuid.UUID {
		return testUUID
	}
	clock := SingleTimestampClock{time.Now()}
	suite.app.clock = clock

	// create required document header
	suite.
		Request("POST", "/doc").
		Body(M{
			"filename": "myfile",
		}).
		Expect(http.StatusOK, M{
			"success": true,
			"id":      testUUID.String(),
		})

	// post content
	suite.
		Request("POST", "/doc/"+testUUID.String()+"/content").
		File("file", "ignored", data).
		Expect(http.StatusOK, M{
			"success": true,
		})

	// check that content is stored correctly
	rdc, err := suite.app.objects.Read(DocID(testUUID.String()))
	suite.NoError(err)

	content, err := io.ReadAll(rdc)
	suite.NoError(err)
	suite.NoError(rdc.Close())

	suite.Equal(data, content)
}
