package app

import "net/http"

func (suite *AppSuite) TestLogin() {
	suite.createUser("testuser", "testpass")

	suite.
		Post("/auth/login").
		BodyJSON(M{
			"username": "testuser",
			"password": "testpass",
		}).
		ExpectJSON(http.StatusOK, M{
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
		Post("/auth/login").
		BodyJSON(M{
			"username": "testuser",
			"password": "testpass",
		}).
		ExpectJSON(http.StatusOK, M{
			"success":   true,
			"challenge": "testchallenge",
		})
	suite.
		Post("/auth/challenge").
		BodyJSON(M{
			"username":        "testuser",
			"challenge":       "testchallenge",
			"client_response": "testresponse",
		}).
		ExpectJSON(http.StatusOK, M{
			"success": true,
		})
}

func (suite *AppSuite) TestLoginNoUser() {
	suite.
		Post("/auth/login").
		BodyJSON(M{
			"username": "testuser",
			"password": "testpass",
		}).
		ExpectJSON(http.StatusUnauthorized, M{
			"success": false,
			"message": "invalid credentials",
		})
}

func (suite *AppSuite) TestLoginInvalid() {
	suite.createUser("someuser", "somepass")

	suite.
		Post("/auth/login").
		BodyJSON(M{
			"username": "testuser",
			"password": "testpass",
		}).
		ExpectJSON(http.StatusUnauthorized, M{
			"success": false,
			"message": "invalid credentials",
		})
}

func (suite *AppSuite) TestLogout() {
	suite.createUser("testuser", "testpass")

	suite.
		Post("/auth/login").
		BodyJSON(M{
			"username": "testuser",
			"password": "testpass",
		}).
		ExpectJSON(http.StatusOK, M{
			"success": true,
		})

	suite.
		Get("/auth/logout").
		ExpectJSON(http.StatusOK, M{
			"success": true,
		})

	suite.
		Get("/auth/logout").
		ExpectJSON(http.StatusUnauthorized, M{
			"success": false,
			"message": "not logged in",
		})
}
