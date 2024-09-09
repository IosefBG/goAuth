package services

import (
	"backendGoAuth/internal/entities"
	"backendGoAuth/internal/repositories"
	"time"
)

type AdminService interface {
	GetAllUsers() ([]entities.User, error)
	EditUser(user entities.User) error
	DeleteUser(userID int) error
}

type adminService struct {
	repo repositories.UserRepository
}

func NewAdminService(repo repositories.UserRepository) AdminService {
	return &adminService{repo}
}

func (s *adminService) GetAllUsers() ([]entities.User, error) {
	return s.repo.GetAllUsers()
}

func (s *adminService) EditUser(user entities.User) error {
	user.UpdatedAt = time.Now()
	return s.repo.EditUser(user)
}

func (s *adminService) DeleteUser(userID int) error {
	return s.repo.DeleteUser(userID)
}
