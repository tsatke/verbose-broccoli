package app

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestCognitoServiceSuite(t *testing.T) {
	suite.Run(t, new(CognitoServiceTestSuite))
}

type CognitoServiceTestSuite struct {
	suite.Suite

	service  *CognitoService
	clientID string
	poolID   string

	client *mockCognitoIdentityProviderAPI
}

func (suite *CognitoServiceTestSuite) SetupTest() {
	suite.clientID = "client-id"
	suite.poolID = "pool-id"
	suite.client = &mockCognitoIdentityProviderAPI{}
	suite.service = &CognitoService{
		poolID:     suite.poolID,
		clientID:   suite.clientID,
		idProvider: suite.client,
	}
}

func (suite *CognitoServiceTestSuite) TestLogin() {
	suite.client.
		On("InitiateAuth",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *cognitoidentityprovider.InitiateAuthInput) bool {
				return suite.EqualValues("USER_PASSWORD_AUTH", i.AuthFlow) &&
					suite.Equal(map[string]string{
						"USERNAME": "testuser",
						"PASSWORD": "testpass",
					}, i.AuthParameters) &&
					suite.Equal(suite.clientID, *i.ClientId)
			}),
		).
		Return(&cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &types.AuthenticationResultType{
				IdToken: aws.String("testtoken"),
			},
		}, nil).
		Once()

	res, err := suite.service.Login("testuser", "testpass")
	suite.NoError(err)
	suite.Equal(LoginResult{
		Done:      true,
		Success:   true,
		Challenge: "",
		Token:     "testtoken",
	}, res)
}

func (suite *CognitoServiceTestSuite) TestLoginWithChallenge() {
	suite.client.
		On("InitiateAuth",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *cognitoidentityprovider.InitiateAuthInput) bool {
				return suite.EqualValues("USER_PASSWORD_AUTH", i.AuthFlow) &&
					suite.Equal(map[string]string{
						"USERNAME": "testuser",
						"PASSWORD": "testpass",
					}, i.AuthParameters) &&
					suite.Equal(suite.clientID, *i.ClientId)
			}),
		).
		Return(&cognitoidentityprovider.InitiateAuthOutput{
			ChallengeName: "testchallenge",
		}, nil).
		Once()

	res, err := suite.service.Login("testuser", "testpass")
	suite.NoError(err)
	suite.Equal(LoginResult{
		Done:      false,
		Success:   true,
		Challenge: "testchallenge",
		Token:     "",
	}, res)
}
