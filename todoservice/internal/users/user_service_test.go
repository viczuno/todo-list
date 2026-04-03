package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateUser(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testID := "test-uuid-123"

	tests := []struct {
		name           string
		input          models.User
		setupMocks     func(repo *automock.UserRepository, uuid *automock.UUIDService, timeService *automock.TimeService)
		expectedID     string
		expectedErrMsg string
	}{
		{
			name: "success - creates user with all fields",
			input: models.User{
				Email:    "test@example.com",
				GithubID: "github-123",
				Role:     "admin",
			},
			setupMocks: func(repo *automock.UserRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.MatchedBy(func(user models.User) bool {
					return user.ID == testID &&
						user.Email == "test@example.com" &&
						user.GithubID == "github-123" &&
						user.Role == "admin" &&
						user.CreatedAt == mockTime &&
						user.UpdatedAt == mockTime
				})).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "success - creates user with minimal fields",
			input: models.User{
				Email: "minimal@example.com",
				Role:  "reader",
			},
			setupMocks: func(repo *automock.UserRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "success - creates user with reader role",
			input: models.User{
				Email:    "reader@example.com",
				GithubID: "github-456",
				Role:     "reader",
			},
			setupMocks: func(repo *automock.UserRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "success - creates user with writer role",
			input: models.User{
				Email:    "writer@example.com",
				GithubID: "github-789",
				Role:     "writer",
			},
			setupMocks: func(repo *automock.UserRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "error - repository create fails",
			input: models.User{
				Email:    "test@example.com",
				GithubID: "github-123",
				Role:     "admin",
			},
			setupMocks: func(repo *automock.UserRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return("", errors.New("database connection failed")).Once()
			},
			expectedID:     "",
			expectedErrMsg: "database connection failed",
		},
		{
			name: "error - duplicate email constraint",
			input: models.User{
				Email: "duplicate@example.com",
				Role:  "reader",
			},
			setupMocks: func(repo *automock.UserRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return("", errors.New("unique constraint violation")).Once()
			},
			expectedID:     "",
			expectedErrMsg: "unique constraint violation",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewUserRepository(t)
			uuidService := automock.NewUUIDService(t)
			timeService := automock.NewTimeService(t)

			tc.setupMocks(repo, uuidService, timeService)

			svc := users.NewService(repo, uuidService, timeService)

			// Act
			id, err := svc.CreateUser(ctx, tc.input)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
				assert.Empty(t, id)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}

func TestService_GetUser(t *testing.T) {
	ctx := context.Background()
	testID := "test-uuid-123"
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		userID         string
		setupMocks     func(repo *automock.UserRepository)
		expectedUser   models.User
		expectedErrMsg string
	}{
		{
			name:   "success - returns user by id",
			userID: testID,
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.User{
					ID:        testID,
					Email:     "test@example.com",
					GithubID:  "github-123",
					Role:      "admin",
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
				}, nil).Once()
			},
			expectedUser: models.User{
				ID:        testID,
				Email:     "test@example.com",
				GithubID:  "github-123",
				Role:      "admin",
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:   "success - returns user with reader role",
			userID: testID,
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.User{
					ID:    testID,
					Email: "reader@example.com",
					Role:  "reader",
				}, nil).Once()
			},
			expectedUser: models.User{
				ID:    testID,
				Email: "reader@example.com",
				Role:  "reader",
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - user not found",
			userID: "non-existent-id",
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, "non-existent-id").Return(models.User{}, errors.New("user not found")).Once()
			},
			expectedUser:   models.User{},
			expectedErrMsg: "user not found",
		},
		{
			name:   "error - database connection error",
			userID: testID,
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.User{}, errors.New("database connection refused")).Once()
			},
			expectedUser:   models.User{},
			expectedErrMsg: "database connection refused",
		},
		{
			name:   "error - empty id",
			userID: "",
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, "").Return(models.User{}, errors.New("invalid id")).Once()
			},
			expectedUser:   models.User{},
			expectedErrMsg: "invalid id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewUserRepository(t)
			tc.setupMocks(repo)

			svc := users.NewService(repo, nil, nil)

			// Act
			user, err := svc.GetUser(ctx, tc.userID)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	testID := "test-uuid-123"
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          models.User
		setupMocks     func(repo *automock.UserRepository)
		expectedErrMsg string
	}{
		{
			name: "success - updates all fields",
			input: models.User{
				ID:       testID,
				Email:    "updated@example.com",
				GithubID: "github-updated",
				Role:     "writer",
			},
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.User{
					ID:        testID,
					Email:     "original@example.com",
					CreatedAt: mockTime,
				}, nil).Once()
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(user models.User) bool {
					return user.ID == testID && user.Email == "updated@example.com"
				})).Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name: "success - updates only email",
			input: models.User{
				ID:    testID,
				Email: "newemail@example.com",
			},
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.User{ID: testID}, nil).Once()
				repo.EXPECT().Update(ctx, mock.Anything).Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name: "success - updates role to admin",
			input: models.User{
				ID:   testID,
				Role: "admin",
			},
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.User{ID: testID, Role: "reader"}, nil).Once()
				repo.EXPECT().Update(ctx, mock.Anything).Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name: "error - user not found",
			input: models.User{
				ID:    "non-existent-id",
				Email: "updated@example.com",
			},
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, "non-existent-id").Return(models.User{}, errors.New("user not found")).Once()
			},
			expectedErrMsg: "user not found",
		},
		{
			name: "error - repository update fails",
			input: models.User{
				ID:    testID,
				Email: "updated@example.com",
			},
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.User{ID: testID}, nil).Once()
				repo.EXPECT().Update(ctx, mock.Anything).Return(errors.New("update failed")).Once()
			},
			expectedErrMsg: "update failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewUserRepository(t)
			tc.setupMocks(repo)

			svc := users.NewService(repo, nil, nil)

			// Act
			err := svc.UpdateUser(ctx, tc.input)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		userID         string
		setupMocks     func(repo *automock.UserRepository)
		expectedErrMsg string
	}{
		{
			name:   "success - deletes existing user",
			userID: "test-uuid-123",
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Delete(ctx, "test-uuid-123").Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - user not found",
			userID: "non-existent-id",
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Delete(ctx, "non-existent-id").Return(errors.New("user not found")).Once()
			},
			expectedErrMsg: "user not found",
		},
		{
			name:   "error - database error",
			userID: "test-uuid-123",
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Delete(ctx, "test-uuid-123").Return(errors.New("database error")).Once()
			},
			expectedErrMsg: "database error",
		},
		{
			name:   "error - empty id",
			userID: "",
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Delete(ctx, "").Return(errors.New("invalid id")).Once()
			},
			expectedErrMsg: "invalid id",
		},
		{
			name:   "error - foreign key constraint",
			userID: "user-with-refs",
			setupMocks: func(repo *automock.UserRepository) {
				repo.EXPECT().Delete(ctx, "user-with-refs").Return(errors.New("foreign key constraint violation")).Once()
			},
			expectedErrMsg: "foreign key constraint violation",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewUserRepository(t)
			tc.setupMocks(repo)

			svc := users.NewService(repo, nil, nil)

			// Act
			err := svc.DeleteUser(ctx, tc.userID)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
