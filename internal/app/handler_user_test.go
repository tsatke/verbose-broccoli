package app

import "net/http"

func (suite *AppSuite) TestLoginSuccess() {
	suite.NoError(suite.app.users.CreateUser("foo", "bar"))

	suite.
		Request("POST", "/user/login").
		Body(M{
			"username": "foo",
			"password": "bar",
		}).
		Expect(http.StatusOK, M{
			"success": true,
		})
}

func (suite *AppSuite) TestLoginFailed() {
	suite.
		Request("POST", "/user/login").
		Body(M{
			"username": "foo",
			"password": "bar",
		}).
		Expect(http.StatusUnauthorized, M{
			"success": false,
			"message": "invalid credentials",
		})
}

func (suite *AppSuite) TestLoginTwice() {
	suite.NoError(suite.app.users.CreateUser("foo", "bar"))

	suite.
		Request("POST", "/user/login").
		Body(M{
			"username": "foo",
			"password": "bar",
		}).
		Expect(http.StatusOK, M{
			"success": true,
		})

	suite.
		Request("POST", "/user/login").
		Body(M{
			"username": "foo",
			"password": "bar",
		}).
		Expect(http.StatusOK, M{
			"success": true,
			"message": "already logged in",
		})
}

func (suite *AppSuite) TestLoginLogoutLogin() {
	suite.NoError(suite.app.users.CreateUser("foo", "bar"))

	suite.
		Request("POST", "/user/login").
		Body(M{
			"username": "foo",
			"password": "bar",
		}).
		Expect(http.StatusOK, M{
			"success": true,
		})

	suite.
		Request("GET", "/user/logout").
		Expect(http.StatusOK, M{
			"success": true,
		})

	suite.
		Request("POST", "/user/login").
		Body(M{
			"username": "foo",
			"password": "bar",
		}).
		Expect(http.StatusOK, M{
			"success": true,
		})
}

func (suite *AppSuite) TestLogoutNotLoggedIn() {
	suite.
		Request("GET", "/user/logout").
		Expect(http.StatusUnauthorized, M{
			"success": false,
			"message": "not logged in",
		})
}
