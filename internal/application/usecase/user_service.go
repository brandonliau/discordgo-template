package usecase

import (
	"discordgo-template/internal/domain/user"

	"github.com/google/uuid"
)

type UserService struct {
	userRepository user.UserRepository
}

func NewUserService(userRepository user.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// --- Get User ---
// question: [UserID string] or [ID string]
type GetUserRequest struct {
	ID uuid.UUID
}

type GetUserResult struct {
	User *user.User
}

func (s *UserService) Get(req GetUserRequest) (*GetUserResult, error) {
	usr, err := s.userRepository.Get(req.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserResult{User: usr}, nil
}

// --- Get All Users ---
type GetAllUsersRequest struct{}

type GetAllUsersResult struct {
	Users []*user.User
}

func (s *UserService) GetAll(req GetAllUsersRequest) (*GetAllUsersResult, error) {
	users, err := s.userRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return &GetAllUsersResult{Users: users}, nil
}

// --- User Join ---
// question: [UserID string] or [ID string]
type UserJoinRequest struct{}

type UserJoinResult struct {
	User *user.User
}

func (s *UserService) Join(req UserJoinRequest) (*UserJoinResult, error) {
	usr := user.NewUser()

	if err := s.userRepository.Create(usr); err != nil {
		return nil, err
	}

	return &UserJoinResult{User: usr}, nil
}

// --- User Leave ---
// question: [UserID string] or [ID string]
type UserLeaveRequest struct {
	ID uuid.UUID
}

type UserLeaveResult struct {
	User *user.User
}

func (s *UserService) Leave(req UserLeaveRequest) (*UserLeaveResult, error) {
	usr, err := s.userRepository.Get(req.ID)
	if err != nil {
		return nil, err
	}

	if err := s.userRepository.Delete(usr); err != nil {
		return nil, err
	}

	return &UserLeaveResult{User: usr}, nil
}
