package lists_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateList(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testID := "test-uuid-123"

	tests := []struct {
		name           string
		input          models.List
		setupMocks     func(repo *automock.ListRepository, uuid *automock.UUIDService, timeService *automock.TimeService)
		expectedID     string
		expectedErrMsg string
	}{
		{
			name: "success - creates list with all fields",
			input: models.List{
				Name:        "Test List",
				Description: "Test description",
				OwnerID:     "owner-1",
				SharedWith:  []string{"user-1", "user-2"},
			},
			setupMocks: func(repo *automock.ListRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.MatchedBy(func(list models.List) bool {
					return list.ID == testID &&
						list.Name == "Test List" &&
						list.OwnerID == "owner-1" &&
						list.CreatedAt == mockTime
				})).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "success - creates list with minimal fields",
			input: models.List{
				Name:    "Minimal List",
				OwnerID: "owner-1",
			},
			setupMocks: func(repo *automock.ListRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "error - repository create fails",
			input: models.List{
				Name:    "Test List",
				OwnerID: "owner-1",
			},
			setupMocks: func(repo *automock.ListRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return("", errors.New("database error")).Once()
			},
			expectedID:     "",
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			uuidService := automock.NewUUIDService(t)
			timeService := automock.NewTimeService(t)

			tc.setupMocks(repo, uuidService, timeService)

			svc := lists.NewService(repo, uuidService, timeService)

			id, err := svc.CreateList(ctx, tc.input)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}

func TestService_GetList(t *testing.T) {
	ctx := context.Background()
	testID := "test-uuid-123"
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		listID         string
		setupMocks     func(repo *automock.ListRepository)
		expectedList   models.List
		expectedErrMsg string
	}{
		{
			name:   "success - returns list by id",
			listID: testID,
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.List{
					ID:          testID,
					Name:        "Test List",
					Description: "Test description",
					OwnerID:     "owner-1",
					CreatedAt:   mockTime,
					UpdatedAt:   mockTime,
				}, nil).Once()
			},
			expectedList: models.List{
				ID:          testID,
				Name:        "Test List",
				Description: "Test description",
				OwnerID:     "owner-1",
				CreatedAt:   mockTime,
				UpdatedAt:   mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - list not found",
			listID: "non-existent",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().Get(ctx, "non-existent").Return(models.List{}, errors.New("list not found")).Once()
			},
			expectedList:   models.List{},
			expectedErrMsg: "list not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			list, err := svc.GetList(ctx, tc.listID)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedList, list)
			}
		})
	}
}

func TestService_UpdateList(t *testing.T) {
	ctx := context.Background()
	testID := "test-uuid-123"
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          models.List
		setupMocks     func(repo *automock.ListRepository)
		expectedErrMsg string
	}{
		{
			name: "success - updates list",
			input: models.List{
				ID:          testID,
				Name:        "Updated List",
				Description: "Updated description",
				OwnerID:     "owner-1",
			},
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.List{ID: testID, CreatedAt: mockTime}, nil).Once()
				repo.EXPECT().Update(ctx, mock.Anything).Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name: "error - list not found",
			input: models.List{
				ID:   "non-existent",
				Name: "Updated List",
			},
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().Get(ctx, "non-existent").Return(models.List{}, errors.New("list not found")).Once()
			},
			expectedErrMsg: "list not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)

			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			err := svc.UpdateList(ctx, tc.input)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_DeleteList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		listID         string
		setupMocks     func(repo *automock.ListRepository)
		expectedErrMsg string
	}{
		{
			name:   "success - deletes existing list",
			listID: "test-uuid-123",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().Delete(ctx, "test-uuid-123").Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - list not found",
			listID: "non-existent",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().Delete(ctx, "non-existent").Return(errors.New("list not found")).Once()
			},
			expectedErrMsg: "list not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			err := svc.DeleteList(ctx, tc.listID)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetAllLists(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		setupMocks     func(repo *automock.ListRepository)
		expectedLists  []models.List
		expectedErrMsg string
	}{
		{
			name: "success - returns all lists",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetAll(ctx).Return([]models.List{
					{ID: "list-1", Name: "List 1", OwnerID: "owner-1", CreatedAt: mockTime},
					{ID: "list-2", Name: "List 2", OwnerID: "owner-2", CreatedAt: mockTime},
				}, nil).Once()
			},
			expectedLists: []models.List{
				{ID: "list-1", Name: "List 1", OwnerID: "owner-1", CreatedAt: mockTime},
				{ID: "list-2", Name: "List 2", OwnerID: "owner-2", CreatedAt: mockTime},
			},
			expectedErrMsg: "",
		},
		{
			name: "success - returns empty list",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetAll(ctx).Return([]models.List{}, nil).Once()
			},
			expectedLists:  []models.List{},
			expectedErrMsg: "",
		},
		{
			name: "error - database error",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetAll(ctx).Return(nil, errors.New("database error")).Once()
			},
			expectedLists:  nil,
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			result, err := svc.GetAllLists(ctx)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedLists, result)
			}
		})
	}
}

func TestService_ListAllByUserID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		userID         string
		setupMocks     func(repo *automock.ListRepository)
		expectedAccess []models.Access
		expectedErrMsg string
	}{
		{
			name:   "success - returns access for user",
			userID: "user-1",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().ListAllByUserID(ctx, "user-1").Return([]models.Access{
					{ListID: "list-1", UserID: "user-1", Role: "admin"},
					{ListID: "list-2", UserID: "user-1", Role: "reader"},
				}, nil).Once()
			},
			expectedAccess: []models.Access{
				{ListID: "list-1", UserID: "user-1", Role: "admin"},
				{ListID: "list-2", UserID: "user-1", Role: "reader"},
			},
			expectedErrMsg: "",
		},
		{
			name:   "success - returns empty list",
			userID: "user-2",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().ListAllByUserID(ctx, "user-2").Return([]models.Access{}, nil).Once()
			},
			expectedAccess: []models.Access{},
			expectedErrMsg: "",
		},
		{
			name:   "error - database error",
			userID: "user-1",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().ListAllByUserID(ctx, "user-1").Return(nil, errors.New("database error")).Once()
			},
			expectedAccess: nil,
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			result, err := svc.ListAllByUserID(ctx, tc.userID)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAccess, result)
			}
		})
	}
}

func TestService_GetUsersByListID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		listID         string
		setupMocks     func(repo *automock.ListRepository)
		expectedAccess []models.Access
		expectedErrMsg string
	}{
		{
			name:   "success - returns users for list",
			listID: "list-1",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetUsersByListID(ctx, "list-1").Return([]models.Access{
					{ListID: "list-1", UserID: "user-1", Role: "admin"},
					{ListID: "list-1", UserID: "user-2", Role: "writer"},
					{ListID: "list-1", UserID: "user-3", Role: "reader"},
				}, nil).Once()
			},
			expectedAccess: []models.Access{
				{ListID: "list-1", UserID: "user-1", Role: "admin"},
				{ListID: "list-1", UserID: "user-2", Role: "writer"},
				{ListID: "list-1", UserID: "user-3", Role: "reader"},
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - list not found",
			listID: "non-existent",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetUsersByListID(ctx, "non-existent").Return(nil, errors.New("list not found")).Once()
			},
			expectedAccess: nil,
			expectedErrMsg: "list not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			result, err := svc.GetUsersByListID(ctx, tc.listID)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAccess, result)
			}
		})
	}
}

func TestService_GetListOwnerID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		listID         string
		setupMocks     func(repo *automock.ListRepository)
		expectedOwner  string
		expectedErrMsg string
	}{
		{
			name:   "success - returns owner id",
			listID: "list-1",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetListOwnerID(ctx, "list-1").Return("owner-1", nil).Once()
			},
			expectedOwner:  "owner-1",
			expectedErrMsg: "",
		},
		{
			name:   "error - list not found",
			listID: "non-existent",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetListOwnerID(ctx, "non-existent").Return("", errors.New("list not found")).Once()
			},
			expectedOwner:  "",
			expectedErrMsg: "list not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			result, err := svc.GetListOwnerID(ctx, tc.listID)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOwner, result)
			}
		})
	}
}

func TestService_CreateAccess(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		input          models.Access
		setupMocks     func(repo *automock.ListRepository)
		expectedAccess models.Access
		expectedErrMsg string
	}{
		{
			name: "success - creates access",
			input: models.Access{
				ListID: "list-1",
				UserID: "user-1",
				Role:   "reader",
			},
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().CreateAccess(ctx, models.Access{
					ListID: "list-1",
					UserID: "user-1",
					Role:   "reader",
				}).Return(models.Access{
					ListID: "list-1",
					UserID: "user-1",
					Role:   "reader",
				}, nil).Once()
			},
			expectedAccess: models.Access{
				ListID: "list-1",
				UserID: "user-1",
				Role:   "reader",
			},
			expectedErrMsg: "",
		},
		{
			name: "error - database error",
			input: models.Access{
				ListID: "list-1",
				UserID: "user-1",
				Role:   "reader",
			},
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().CreateAccess(ctx, mock.Anything).Return(models.Access{}, errors.New("database error")).Once()
			},
			expectedAccess: models.Access{},
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			result, err := svc.CreateAccess(ctx, tc.input)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAccess, result)
			}
		})
	}
}

func TestService_GetAccess(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		listID         string
		userID         string
		setupMocks     func(repo *automock.ListRepository)
		expectedAccess models.Access
		expectedErrMsg string
	}{
		{
			name:   "success - returns access",
			listID: "list-1",
			userID: "user-1",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetAccess(ctx, "list-1", "user-1").Return(models.Access{
					ListID: "list-1",
					UserID: "user-1",
					Role:   "reader",
				}, nil).Once()
			},
			expectedAccess: models.Access{
				ListID: "list-1",
				UserID: "user-1",
				Role:   "reader",
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - access not found",
			listID: "list-1",
			userID: "user-2",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().GetAccess(ctx, "list-1", "user-2").Return(models.Access{}, errors.New("access not found")).Once()
			},
			expectedAccess: models.Access{},
			expectedErrMsg: "access not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			result, err := svc.GetAccess(ctx, tc.listID, tc.userID)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAccess, result)
			}
		})
	}
}

func TestService_DeleteAccess(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		listID         string
		userID         string
		setupMocks     func(repo *automock.ListRepository)
		expectedErrMsg string
	}{
		{
			name:   "success - deletes access",
			listID: "list-1",
			userID: "user-1",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().DeleteAccess(ctx, "list-1", "user-1").Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - access not found",
			listID: "list-1",
			userID: "user-2",
			setupMocks: func(repo *automock.ListRepository) {
				repo.EXPECT().DeleteAccess(ctx, "list-1", "user-2").Return(errors.New("access not found")).Once()
			},
			expectedErrMsg: "access not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := automock.NewListRepository(t)
			tc.setupMocks(repo)

			svc := lists.NewService(repo, nil, nil)

			err := svc.DeleteAccess(ctx, tc.listID, tc.userID)

			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
