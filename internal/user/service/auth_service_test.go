package service

import (
	"context"
	"testing"

	"github.com/PauloGuillen/gostosobookings/internal/user/repository"
)

func TestAuthService_Login(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				repo: tt.fields.repo,
			}
			got, err := s.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
