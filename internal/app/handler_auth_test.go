package app

import "net/http"

func (suite *AppSuite) TestLogin() {
	suite.createUser("testuser", "testpass")

	suite.
		Request("POST", "/auth/login").
		Body(M{
			"username": "testuser",
			"password": "testpass",
		}).
		Expect(http.StatusOK, M{
			"success": true,
		})
}

func (suite *AppSuite) TestLoginNoUser() {
	suite.
		Request("POST", "/auth/login").
		Body(M{
			"username": "testuser",
			"password": "testpass",
		}).
		Expect(http.StatusUnauthorized, M{
			"success": false,
			"message": "invalid credentials",
		})
}

func (suite *AppSuite) TestLoginInvalid() {
	suite.createUser("someuser", "somepass")

	suite.
		Request("POST", "/auth/login").
		Body(M{
			"username": "testuser",
			"password": "testpass",
		}).
		Expect(http.StatusUnauthorized, M{
			"success": false,
			"message": "invalid credentials",
		})
}
