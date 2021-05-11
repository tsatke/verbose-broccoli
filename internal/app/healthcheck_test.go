package app

import "net/http"

func (suite *AppSuite) TestHealthcheck() {
	suite.
		Request("GET", "/healthcheck").
		Expect(http.StatusOK, M{
			"success": true,
		})
}
