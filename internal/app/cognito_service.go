package app

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/tsatke/verbose-broccoli/internal/app/config"
)

//go:generate mockery --inpackage --testonly --case snake --name cognitoIdentityProviderAPI --filename cognito_service_mock_test.go

type cognitoIdentityProviderAPI interface {
	InitiateAuth(context.Context, *cip.InitiateAuthInput, ...func(*cip.Options)) (*cip.InitiateAuthOutput, error)
	RespondToAuthChallenge(context.Context, *cip.RespondToAuthChallengeInput, ...func(*cip.Options)) (*cip.RespondToAuthChallengeOutput, error)
}

type CognitoService struct {
	poolID   string
	clientID string

	idProvider cognitoIdentityProviderAPI
}

func NewCognitoService(cfg config.Config) *CognitoService {
	return &CognitoService{
		poolID:   cfg.GetString(config.AWSCognitoPoolID),
		clientID: cfg.GetString(config.AWSCognitoClientID),

		idProvider: cip.NewFromConfig(cfg.AWS),
	}
}

func (s *CognitoService) Login(user, pass string) (LoginResult, error) {
	resp, err := s.idProvider.InitiateAuth(context.Background(), &cip.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": user,
			"PASSWORD": pass,
		},
		ClientId: aws.String(s.clientID),
	})
	if err != nil {
		return LoginResult{}, fmt.Errorf("initiate auth: %w", err)
	}

	if resp.AuthenticationResult != nil {
		return LoginResult{
			Done:    true,
			Success: true,
			Token:   *resp.AuthenticationResult.IdToken,
		}, nil
	}

	return LoginResult{
		Done:      false,
		Success:   true,
		Challenge: string(resp.ChallengeName),
	}, nil
}

func (s *CognitoService) AnswerChallenge(user, challenge, payload string) (LoginResult, error) {
	resp, err := s.idProvider.RespondToAuthChallenge(context.Background(), &cip.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameType(challenge),
		ClientId:      aws.String(s.clientID),
		ChallengeResponses: map[string]string{
			challenge: payload,
		},
	})

	if err != nil {
		return LoginResult{}, fmt.Errorf("respond to auth: %w", err)
	}

	return LoginResult{
		Done:    true,
		Success: true,
		Token:   *resp.AuthenticationResult.IdToken,
	}, nil
}

func (s *CognitoService) TokenValid(s2 string) bool {
	panic("implement me")
}
