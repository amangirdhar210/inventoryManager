package service

import (
	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/amangirdhar210/inventory-manager/internal/core/ports"
)

type authService struct {
	repo           ports.ManagerRepository
	tokenGenerator ports.TokenGenerator
}

func NewAuthService(repo ports.ManagerRepository, tokenGenerator ports.TokenGenerator) AuthService {
	return &authService{
		repo:           repo,
		tokenGenerator: tokenGenerator,
	}
}

func (s *authService) Login(email, password string) (string, error) {
	manager, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := manager.CheckPassword(password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.tokenGenerator.GenerateToken(manager)
	if err != nil {
		return "", domain.ErrTokenGeneration
	}

	return token, nil
}
