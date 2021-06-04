package app

import (
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (suite *AppSuite) TestPostDocumentNoLogin() {
	suite.
		Post("/doc").
		BodyJSON(M{
			"filename": "myfile",
			"size":     1234,
		}).
		ExpectJSON(http.StatusUnauthorized, M{
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
		Post("/doc").
		BodyJSON(M{
			"filename": "myfile",
		}).
		ExpectJSON(http.StatusOK, M{
			"success": true,
			"id":      testUUID.String(),
		})

	// check that document is created correctly
	doc, err := suite.app.documents.Get(DocID(testUUID.String()))
	suite.NoError(err)

	suite.Equal(DocID(testUUID.String()), doc.ID)
	suite.Equal("myfile", doc.Name)
	suite.Equal(user, doc.Owner)
	suite.Truef(clock.Timestamp.Equal(doc.Created), "expected %v, but got %v", clock.Timestamp, doc.Created)
	suite.Truef(time.Time{}.Equal(doc.Updated), "expected %v, but got %v", time.Time{}, doc.Updated)

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
		Post("/doc").
		BodyJSON(M{
			"filename": "myfile",
		}).
		ExpectJSON(http.StatusOK, M{
			"success": true,
			"id":      testUUID.String(),
		})

	// post content
	suite.
		Post("/doc/"+testUUID.String()+"/content").
		File("file", "ignored", data).
		ExpectJSON(http.StatusOK, M{
			"success": true,
		})

	// check that content is stored correctly
	rdc, err := suite.app.objects.Read(DocID(testUUID.String()))
	suite.NoError(err)

	content, err := io.ReadAll(rdc)
	suite.NoError(err)
	suite.NoError(rdc.Close())

	suite.Equal(data, content)

	// check that the header information updates the updated field
	doc, err := suite.app.documents.Get(DocID(testUUID.String()))
	suite.NoError(err)

	suite.Equal(DocID(testUUID.String()), doc.ID)
	suite.Equal("myfile", doc.Name)
	suite.Equal(user, doc.Owner)
	suite.Truef(clock.Timestamp.Equal(doc.Created), "expected %v, but got %v", clock.Timestamp, doc.Created)
	suite.Truef(clock.Timestamp.Equal(doc.Updated), "expected %v, but got %v", clock.Timestamp, doc.Updated)
}
