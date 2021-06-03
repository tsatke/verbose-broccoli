package app

import "net/http"

func (suite *AppSuite) TestGetContentNoLogin() {
	suite.
		Request("GET", "/doc/abc/content").
		Expect(http.StatusUnauthorized, M{
			"message": "not logged in",
			"success": false,
		})
}

func (suite *AppSuite) TestGetContentNoDocument() {
	suite.login()

	suite.
		Request("GET", "/doc/abc/content").
		Expect(http.StatusBadRequest, M{
			"message": "no content for id",
			"success": false,
		})
}

func (suite *AppSuite) TestGetContent() {
	suite.login()
	suite.createContent("abc", []byte("hello"))

	suite.
		Request("GET", "/doc/abc/content").
		ExpectRaw(http.StatusOK, []byte("hello"))
}

func (suite *AppSuite) TestGetContentMultipleTimes() {
	suite.login()
	suite.createContent("abc", []byte("hello"))

	suite.
		Request("GET", "/doc/abc/content").
		ExpectRaw(http.StatusOK, []byte("hello"))
	suite.
		Request("GET", "/doc/abc/content").
		ExpectRaw(http.StatusOK, []byte("hello"))
	suite.
		Request("GET", "/doc/abc/content").
		ExpectRaw(http.StatusOK, []byte("hello"))
}
