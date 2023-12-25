package jwt

import (
	"testing"
	"time"

	"github.com/Len4i/auth-service/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func TestNewToken(t *testing.T) {
	type args struct {
		user     models.User
		app      models.App
		duration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				user: models.User{
					ID:    1,
					Email: "mail1@buba.com",
				},
				app: models.App{
					ID:     12,
					Secret: "secret",
				},
				duration: 5 * time.Minute,
			},
			wantErr: false,
			want:    "mail1@buba.com",
		},
		{
			name: "still happy path",
			args: args{
				user: models.User{
					ID:    0,
					Email: "",
				},
				app: models.App{
					ID:     0,
					Secret: "",
				},
				duration: 5 * time.Minute,
			},
			wantErr: false,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewToken(tt.args.user, tt.args.app, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			token, err := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.args.app.Secret), nil
			})
			if err != nil {
				t.Error(err)
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				t.Error("token claims are not of type jwt.MapClaims")
			}
			if claims["email"] != tt.want {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
