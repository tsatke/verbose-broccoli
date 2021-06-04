package app

import "net/http"

func (suite *AppSuite) TestHealthcheck() {
	suite.
		Get("/healthcheck").
		ExpectJSON(http.StatusOK, M{
			"success": true,
		})
}
