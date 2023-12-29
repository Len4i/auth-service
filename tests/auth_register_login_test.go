package tests

import (
	"testing"
	"time"

	aaav1 "github.com/Len4i/aaa/gen/go/aaa"
	"github.com/Len4i/auth-service/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID          = 999
	emtyAppID      = 0
	notExistAppID  = 998
	notExistUserID = 99999
	appSecret      = "test-secret"
	adminUserID    = 999
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePass(passDefaultLen)

	respReg, err := s.AuthClient.Register(ctx, &aaav1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := s.AuthClient.Login(ctx, &aaav1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["user_id"].(float64)))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.Equal(t, email, claims["email"].(string))

	const deltaSec = 1

	assert.InDelta(t, loginTime.Add(s.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSec)
}

func TestRegister_DoubleRegistration(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePass(passDefaultLen)

	respReg, err := s.AuthClient.Register(ctx, &aaav1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	_, err = s.AuthClient.Login(ctx, &aaav1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	_, err = s.AuthClient.Register(ctx, &aaav1.RegisterRequest{
		Email:    email,
		Password: "other-password",
	})
	require.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Internal desc = user already exists")

	// Try login again with original password
	_, err = s.AuthClient.Login(ctx, &aaav1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)
}

func TestRegister_IncorrectInput(t *testing.T) {
	ctx, s := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "empty email",
			email:       "",
			password:    randomFakePass(passDefaultLen),
			expectedErr: "rpc error: code = InvalidArgument desc = email is required",
		},
		{
			name:        "invalid email",
			email:       "this is not an email",
			password:    randomFakePass(passDefaultLen),
			expectedErr: "rpc error: code = InvalidArgument desc = email is not valid",
		},
		{
			name:        "existing user, incorrect password",
			email:       "admin-user@localhost.com",
			password:    randomFakePass(passDefaultLen),
			expectedErr: "rpc error: code = Internal desc = user already exists",
		},
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := s.AuthClient.Register(ctx, &aaav1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}

func TestLogin_IncorrectInput(t *testing.T) {
	ctx, s := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int
		expectedErr string
	}{
		{
			name:        "not existing email",
			email:       "notexistingemail@incognito.com",
			password:    randomFakePass(passDefaultLen),
			appID:       appID,
			expectedErr: "rpc error: code = Internal desc = internal error",
		},
		{
			name:        "empty email",
			email:       "",
			password:    randomFakePass(passDefaultLen),
			appID:       appID,
			expectedErr: "rpc error: code = InvalidArgument desc = email is required",
		},
		{
			name:        "invalid email",
			email:       "this is not an email",
			password:    randomFakePass(passDefaultLen),
			appID:       appID,
			expectedErr: "rpc error: code = InvalidArgument desc = email is not valid",
		},
		{
			name:        "existing user, incorrect password",
			email:       "admin-user@localhost.com",
			password:    randomFakePass(passDefaultLen),
			appID:       appID,
			expectedErr: "rpc error: code = Internal desc = internal error",
		},
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "rpc error: code = InvalidArgument desc = password is required",
		},
		{
			name:        "empty appID",
			email:       gofakeit.Email(),
			password:    randomFakePass(passDefaultLen),
			appID:       emtyAppID,
			expectedErr: "rpc error: code = Internal desc = internal error",
		},
		{
			name:        "not existing AppID",
			email:       gofakeit.Email(),
			password:    randomFakePass(passDefaultLen),
			appID:       notExistAppID,
			expectedErr: "rpc error: code = Internal desc = internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := s.AuthClient.Login(ctx, &aaav1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    int32(tt.appID),
			})
			require.Error(t, err)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}

func TestIsAdmin_NotAdmin(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePass(passDefaultLen)

	respReg, err := s.AuthClient.Register(ctx, &aaav1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respNotAdmin, err := s.AuthClient.IsAdmin(ctx, &aaav1.IsAdminRequest{
		UserId: respReg.GetUserId(),
	})
	require.NoError(t, err)
	assert.False(t, respNotAdmin.GetIsAdmin())

}

func TestIsAdmin_IsAdmin(t *testing.T) {
	ctx, s := suite.New(t)

	respAdmin, err := s.AuthClient.IsAdmin(ctx, &aaav1.IsAdminRequest{
		UserId: adminUserID,
	})
	require.NoError(t, err)
	assert.True(t, respAdmin.GetIsAdmin())
}

func TestIsAdmin_InvalidUserID(t *testing.T) {
	ctx, s := suite.New(t)

	_, err := s.AuthClient.IsAdmin(ctx, &aaav1.IsAdminRequest{
		UserId: 0,
	})
	require.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = userID is required")
}
func TestIsAdmin_EmptyRequest(t *testing.T) {
	ctx, s := suite.New(t)

	_, err := s.AuthClient.IsAdmin(ctx, &aaav1.IsAdminRequest{})
	require.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = userID is required")

}

func TestIsAdmin_NotExistUserID(t *testing.T) {
	ctx, s := suite.New(t)

	_, err := s.AuthClient.IsAdmin(ctx, &aaav1.IsAdminRequest{})
	require.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = userID is required")

	_, err = s.AuthClient.IsAdmin(ctx, &aaav1.IsAdminRequest{
		UserId: notExistUserID,
	})
	require.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Internal desc = internal error")
}

func randomFakePass(len int) string {
	return gofakeit.Password(true, true, true, true, true, len)
}
