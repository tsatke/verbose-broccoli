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

func (suite *AppSuite) TestLoginWithChallenge() {
	suite.createUser("testuser", "testpass")
	mock := new(MockAuthService)
	suite.app.auth = mock

	mock.
		On("Login", "testuser", "testpass").
		Return(LoginResult{
			Success:   true,
			Challenge: "testchallenge",
		}, nil)
	mock.
		On("AnswerChallenge", "testuser", "testchallenge", "testresponse").
		Return(LoginResult{
			Success: true,
		}, nil)

	suite.
		Request("POST", "/auth/login").
		Body(M{
			"username": "testuser",
			"password": "testpass",
		}).
		Expect(http.StatusOK, M{
			"success":   true,
			"challenge": "testchallenge",
		})
	suite.
		Request("POST", "/auth/challenge").
		Body(M{
			"username":        "testuser",
			"challenge":       "testchallenge",
			"client_response": "testresponse",
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

func (suite *AppSuite) TestLogout() {
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

	suite.
		Request("GET", "/auth/logout").
		Expect(http.StatusOK, M{
			"success": true,
		})

	suite.
		Request("GET", "/auth/logout").
		Expect(http.StatusUnauthorized, M{
			"success": false,
			"message": "not logged in",
		})
}
