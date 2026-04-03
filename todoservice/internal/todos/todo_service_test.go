package todos_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateTodo(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testID := "test-uuid-123"
	testListID := "list-uuid-456"

	tests := []struct {
		name           string
		input          models.Todo
		setupMocks     func(repo *automock.TodoRepository, uuid *automock.UUIDService, timeService *automock.TimeService)
		expectedID     string
		expectedErrMsg string
	}{
		{
			name: "success - creates todo with all fields",
			input: models.Todo{
				Title:       "Test Todo",
				Description: "Test description",
				ListID:      testListID,
				Priority:    constants.PriorityHigh,
			},
			setupMocks: func(repo *automock.TodoRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.MatchedBy(func(todo models.Todo) bool {
					return todo.ID == testID &&
						todo.Title == "Test Todo" &&
						todo.Description == "Test description" &&
						todo.ListID == testListID &&
						todo.Priority == constants.PriorityHigh &&
						todo.CreatedAt == mockTime &&
						todo.UpdatedAt == mockTime
				})).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "success - creates todo with minimum fields",
			input: models.Todo{
				Title:  "Minimal Todo",
				ListID: testListID,
			},
			setupMocks: func(repo *automock.TodoRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "success - creates todo with low priority",
			input: models.Todo{
				Title:    "Low Priority Todo",
				ListID:   testListID,
				Priority: constants.PriorityLow,
			},
			setupMocks: func(repo *automock.TodoRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return(testID, nil).Once()
			},
			expectedID:     testID,
			expectedErrMsg: "",
		},
		{
			name: "error - repository create fails",
			input: models.Todo{
				Title:       "Test Todo",
				Description: "Test description",
				ListID:      testListID,
				Priority:    constants.PriorityLow,
			},
			setupMocks: func(repo *automock.TodoRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
				uuid.EXPECT().Generate().Return(testID).Once()
				timeService.EXPECT().Now().Return(mockTime).Times(2)
				repo.EXPECT().Create(ctx, mock.Anything).Return("", errors.New("database connection failed")).Once()
			},
			expectedID:     "",
			expectedErrMsg: "database connection failed",
		},
		{
			name: "error - repository returns constraint violation",
			input: models.Todo{
				Title:  "Duplicate Todo",
				ListID: testListID,
			},
			setupMocks: func(repo *automock.TodoRepository, uuid *automock.UUIDService, timeService *automock.TimeService) {
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
			repo := automock.NewTodoRepository(t)
			uuidService := automock.NewUUIDService(t)
			timeService := automock.NewTimeService(t)

			tc.setupMocks(repo, uuidService, timeService)

			svc := todos.NewService(repo, uuidService, timeService)

			// Act
			id, err := svc.CreateTodo(ctx, tc.input)

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

func TestService_GetTodo(t *testing.T) {
	ctx := context.Background()
	testID := "test-uuid-123"
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		todoID         string
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodo   models.Todo
		expectedErrMsg string
	}{
		{
			name:   "success - returns todo by id",
			todoID: testID,
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.Todo{
					ID:          testID,
					Title:       "Test Todo",
					Description: "Test description",
					ListID:      "list-123",
					Priority:    constants.PriorityMedium,
					Completed:   false,
					CreatedAt:   mockTime,
					UpdatedAt:   mockTime,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:          testID,
				Title:       "Test Todo",
				Description: "Test description",
				ListID:      "list-123",
				Priority:    constants.PriorityMedium,
				Completed:   false,
				CreatedAt:   mockTime,
				UpdatedAt:   mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:   "success - returns completed todo",
			todoID: testID,
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.Todo{
					ID:        testID,
					Title:     "Completed Todo",
					Completed: true,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:        testID,
				Title:     "Completed Todo",
				Completed: true,
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - todo not found",
			todoID: "non-existent-id",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Get(ctx, "non-existent-id").Return(models.Todo{}, errors.New("todo not found")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "todo not found",
		},
		{
			name:   "error - database connection error",
			todoID: testID,
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Get(ctx, testID).Return(models.Todo{}, errors.New("database connection refused")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "database connection refused",
		},
		{
			name:   "error - empty id",
			todoID: "",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Get(ctx, "").Return(models.Todo{}, errors.New("invalid id")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "invalid id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todo, err := svc.GetTodo(ctx, tc.todoID)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, todo)
			}
		})
	}
}

func TestService_UpdateTodo(t *testing.T) {
	ctx := context.Background()
	testID := "test-uuid-123"
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	newMockTime := time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		input          models.Todo
		setupMocks     func(repo *automock.TodoRepository, timeService *automock.TimeService)
		expectedErrMsg string
	}{
		{
			name: "success - updates all fields",
			input: models.Todo{
				ID:          testID,
				Title:       "Updated Title",
				Description: "Updated Description",
				ListID:      "list-123",
				Priority:    constants.PriorityHigh,
			},
			setupMocks: func(repo *automock.TodoRepository, timeService *automock.TimeService) {
				repo.EXPECT().Get(ctx, testID).Return(models.Todo{
					ID:        testID,
					Title:     "Original Title",
					CreatedAt: mockTime,
				}, nil).Once()
				timeService.EXPECT().Now().Return(newMockTime).Once()
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(todo models.Todo) bool {
					return todo.ID == testID &&
						todo.Title == "Updated Title" &&
						todo.Description == "Updated Description" &&
						todo.UpdatedAt == newMockTime
				})).Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name: "success - updates only title",
			input: models.Todo{
				ID:    testID,
				Title: "New Title Only",
			},
			setupMocks: func(repo *automock.TodoRepository, timeService *automock.TimeService) {
				repo.EXPECT().Get(ctx, testID).Return(models.Todo{ID: testID}, nil).Once()
				timeService.EXPECT().Now().Return(newMockTime).Once()
				repo.EXPECT().Update(ctx, mock.Anything).Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name: "error - todo not found",
			input: models.Todo{
				ID:    "non-existent-id",
				Title: "Updated Title",
			},
			setupMocks: func(repo *automock.TodoRepository, timeService *automock.TimeService) {
				repo.EXPECT().Get(ctx, "non-existent-id").Return(models.Todo{}, errors.New("todo not found")).Once()
			},
			expectedErrMsg: "todo not found",
		},
		{
			name: "error - repository update fails",
			input: models.Todo{
				ID:    testID,
				Title: "Updated Title",
			},
			setupMocks: func(repo *automock.TodoRepository, timeService *automock.TimeService) {
				repo.EXPECT().Get(ctx, testID).Return(models.Todo{ID: testID}, nil).Once()
				timeService.EXPECT().Now().Return(newMockTime).Once()
				repo.EXPECT().Update(ctx, mock.Anything).Return(errors.New("update failed")).Once()
			},
			expectedErrMsg: "update failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			timeService := automock.NewTimeService(t)

			tc.setupMocks(repo, timeService)

			svc := todos.NewService(repo, nil, timeService)

			// Act
			err := svc.UpdateTodo(ctx, tc.input)

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

func TestService_DeleteTodo(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		todoID         string
		setupMocks     func(repo *automock.TodoRepository)
		expectedErrMsg string
	}{
		{
			name:   "success - deletes existing todo",
			todoID: "test-uuid-123",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Delete(ctx, "test-uuid-123").Return(nil).Once()
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - todo not found",
			todoID: "non-existent-id",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Delete(ctx, "non-existent-id").Return(errors.New("todo not found")).Once()
			},
			expectedErrMsg: "todo not found",
		},
		{
			name:   "error - database error",
			todoID: "test-uuid-123",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Delete(ctx, "test-uuid-123").Return(errors.New("database error")).Once()
			},
			expectedErrMsg: "database error",
		},
		{
			name:   "error - empty id",
			todoID: "",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().Delete(ctx, "").Return(errors.New("invalid id")).Once()
			},
			expectedErrMsg: "invalid id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			err := svc.DeleteTodo(ctx, tc.todoID)

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

func TestService_ListTodosByListID(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		listID         string
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodos  []models.Todo
		expectedErrMsg string
	}{
		{
			name:   "success - returns todos for list",
			listID: "list-123",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAllByListID(ctx, "list-123").Return([]models.Todo{
					{ID: "todo-1", Title: "Todo 1", ListID: "list-123", CreatedAt: mockTime},
					{ID: "todo-2", Title: "Todo 2", ListID: "list-123", CreatedAt: mockTime},
				}, nil).Once()
			},
			expectedTodos: []models.Todo{
				{ID: "todo-1", Title: "Todo 1", ListID: "list-123", CreatedAt: mockTime},
				{ID: "todo-2", Title: "Todo 2", ListID: "list-123", CreatedAt: mockTime},
			},
			expectedErrMsg: "",
		},
		{
			name:   "success - returns empty list",
			listID: "empty-list",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAllByListID(ctx, "empty-list").Return([]models.Todo{}, nil).Once()
			},
			expectedTodos:  []models.Todo{},
			expectedErrMsg: "",
		},
		{
			name:   "success - returns single todo",
			listID: "list-single",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAllByListID(ctx, "list-single").Return([]models.Todo{
					{ID: "todo-1", Title: "Single Todo", ListID: "list-single"},
				}, nil).Once()
			},
			expectedTodos: []models.Todo{
				{ID: "todo-1", Title: "Single Todo", ListID: "list-single"},
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - list not found",
			listID: "non-existent-list",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAllByListID(ctx, "non-existent-list").Return([]models.Todo{}, errors.New("list not found")).Once()
			},
			expectedTodos:  []models.Todo{},
			expectedErrMsg: "list not found",
		},
		{
			name:   "error - database error",
			listID: "list-123",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAllByListID(ctx, "list-123").Return([]models.Todo{}, errors.New("database error")).Once()
			},
			expectedTodos:  []models.Todo{},
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todos, err := svc.ListTodosByListID(ctx, tc.listID)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodos, todos)
			}
		})
	}
}

func TestService_GetAllTodos(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodos  []models.Todo
		expectedErrMsg string
	}{
		{
			name: "success - returns all todos",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAll(ctx).Return([]models.Todo{
					{ID: "todo-1", Title: "Todo 1", ListID: "list-1", CreatedAt: mockTime},
					{ID: "todo-2", Title: "Todo 2", ListID: "list-2", CreatedAt: mockTime},
					{ID: "todo-3", Title: "Todo 3", ListID: "list-1", CreatedAt: mockTime},
				}, nil).Once()
			},
			expectedTodos: []models.Todo{
				{ID: "todo-1", Title: "Todo 1", ListID: "list-1", CreatedAt: mockTime},
				{ID: "todo-2", Title: "Todo 2", ListID: "list-2", CreatedAt: mockTime},
				{ID: "todo-3", Title: "Todo 3", ListID: "list-1", CreatedAt: mockTime},
			},
			expectedErrMsg: "",
		},
		{
			name: "success - returns empty list when no todos",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAll(ctx).Return([]models.Todo{}, nil).Once()
			},
			expectedTodos:  []models.Todo{},
			expectedErrMsg: "",
		},
		{
			name: "error - database error",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().GetAll(ctx).Return(nil, errors.New("database connection failed")).Once()
			},
			expectedTodos:  nil,
			expectedErrMsg: "database connection failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todos, err := svc.GetAllTodos(ctx)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodos, todos)
			}
		})
	}
}

func TestService_CompleteTodo(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		todoID         string
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodo   models.Todo
		expectedErrMsg string
	}{
		{
			name:   "success - marks todo as complete",
			todoID: "todo-123",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().CompleteTodo(ctx, "todo-123").Return(models.Todo{
					ID:        "todo-123",
					Title:     "Test Todo",
					Completed: true,
					UpdatedAt: mockTime,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:        "todo-123",
				Title:     "Test Todo",
				Completed: true,
				UpdatedAt: mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - todo not found",
			todoID: "non-existent",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().CompleteTodo(ctx, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "todo not found",
		},
		{
			name:   "error - database error",
			todoID: "todo-123",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().CompleteTodo(ctx, "todo-123").Return(models.Todo{}, errors.New("database error")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todo, err := svc.CompleteTodo(ctx, tc.todoID)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, todo)
			}
		})
	}
}

func TestService_UpdateTodoTitle(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		todoID         string
		newTitle       string
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodo   models.Todo
		expectedErrMsg string
	}{
		{
			name:     "success - updates title",
			todoID:   "todo-123",
			newTitle: "Updated Title",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoTitle(ctx, "todo-123", "Updated Title").Return(models.Todo{
					ID:        "todo-123",
					Title:     "Updated Title",
					UpdatedAt: mockTime,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:        "todo-123",
				Title:     "Updated Title",
				UpdatedAt: mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:     "success - updates to short title",
			todoID:   "todo-123",
			newTitle: "A",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoTitle(ctx, "todo-123", "A").Return(models.Todo{
					ID:    "todo-123",
					Title: "A",
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:    "todo-123",
				Title: "A",
			},
			expectedErrMsg: "",
		},
		{
			name:     "error - todo not found",
			todoID:   "non-existent",
			newTitle: "New Title",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoTitle(ctx, "non-existent", "New Title").Return(models.Todo{}, errors.New("todo not found")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "todo not found",
		},
		{
			name:     "error - database error",
			todoID:   "todo-123",
			newTitle: "New Title",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoTitle(ctx, "todo-123", "New Title").Return(models.Todo{}, errors.New("database error")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todo, err := svc.UpdateTodoTitle(ctx, tc.todoID, tc.newTitle)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, todo)
			}
		})
	}
}

func TestService_UpdateTodoDescription(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		todoID         string
		newDescription string
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodo   models.Todo
		expectedErrMsg string
	}{
		{
			name:           "success - updates description",
			todoID:         "todo-123",
			newDescription: "Updated Description",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoDescription(ctx, "todo-123", "Updated Description").Return(models.Todo{
					ID:          "todo-123",
					Title:       "Test Todo",
					Description: "Updated Description",
					UpdatedAt:   mockTime,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:          "todo-123",
				Title:       "Test Todo",
				Description: "Updated Description",
				UpdatedAt:   mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:           "success - updates to long description",
			todoID:         "todo-123",
			newDescription: "This is a very long description that contains many words and sentences to test the system's ability to handle large text inputs.",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoDescription(ctx, "todo-123", mock.Anything).Return(models.Todo{
					ID:          "todo-123",
					Description: "This is a very long description that contains many words and sentences to test the system's ability to handle large text inputs.",
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:          "todo-123",
				Description: "This is a very long description that contains many words and sentences to test the system's ability to handle large text inputs.",
			},
			expectedErrMsg: "",
		},
		{
			name:           "error - todo not found",
			todoID:         "non-existent",
			newDescription: "New Description",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoDescription(ctx, "non-existent", "New Description").Return(models.Todo{}, errors.New("todo not found")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "todo not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todo, err := svc.UpdateTodoDescription(ctx, tc.todoID, tc.newDescription)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, todo)
			}
		})
	}
}

func TestService_UpdateTodoPriority(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		todoID         string
		newPriority    constants.PriorityLevel
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodo   models.Todo
		expectedErrMsg string
	}{
		{
			name:        "success - updates to high priority",
			todoID:      "todo-123",
			newPriority: constants.PriorityHigh,
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoPriority(ctx, "todo-123", constants.PriorityHigh).Return(models.Todo{
					ID:        "todo-123",
					Title:     "Test Todo",
					Priority:  constants.PriorityHigh,
					UpdatedAt: mockTime,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:        "todo-123",
				Title:     "Test Todo",
				Priority:  constants.PriorityHigh,
				UpdatedAt: mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:        "success - updates to medium priority",
			todoID:      "todo-123",
			newPriority: constants.PriorityMedium,
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoPriority(ctx, "todo-123", constants.PriorityMedium).Return(models.Todo{
					ID:       "todo-123",
					Priority: constants.PriorityMedium,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:       "todo-123",
				Priority: constants.PriorityMedium,
			},
			expectedErrMsg: "",
		},
		{
			name:        "success - updates to low priority",
			todoID:      "todo-123",
			newPriority: constants.PriorityLow,
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoPriority(ctx, "todo-123", constants.PriorityLow).Return(models.Todo{
					ID:       "todo-123",
					Priority: constants.PriorityLow,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:       "todo-123",
				Priority: constants.PriorityLow,
			},
			expectedErrMsg: "",
		},
		{
			name:        "error - todo not found",
			todoID:      "non-existent",
			newPriority: constants.PriorityHigh,
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateTodoPriority(ctx, "non-existent", constants.PriorityHigh).Return(models.Todo{}, errors.New("todo not found")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "todo not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todo, err := svc.UpdateTodoPriority(ctx, tc.todoID, tc.newPriority)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, todo)
			}
		})
	}
}

func TestService_UpdateAssignedTo(t *testing.T) {
	ctx := context.Background()
	mockTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	assignedUser := "user-456"

	tests := []struct {
		name           string
		todoID         string
		userID         string
		setupMocks     func(repo *automock.TodoRepository)
		expectedTodo   models.Todo
		expectedErrMsg string
	}{
		{
			name:   "success - assigns todo to user",
			todoID: "todo-123",
			userID: "user-456",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateAssignedTo(ctx, "todo-123", "user-456").Return(models.Todo{
					ID:         "todo-123",
					Title:      "Test Todo",
					AssignedTo: &assignedUser,
					UpdatedAt:  mockTime,
				}, nil).Once()
			},
			expectedTodo: models.Todo{
				ID:         "todo-123",
				Title:      "Test Todo",
				AssignedTo: &assignedUser,
				UpdatedAt:  mockTime,
			},
			expectedErrMsg: "",
		},
		{
			name:   "error - todo not found",
			todoID: "non-existent",
			userID: "user-456",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateAssignedTo(ctx, "non-existent", "user-456").Return(models.Todo{}, errors.New("todo not found")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "todo not found",
		},
		{
			name:   "error - user not found",
			todoID: "todo-123",
			userID: "invalid-user",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateAssignedTo(ctx, "todo-123", "invalid-user").Return(models.Todo{}, errors.New("user not found")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "user not found",
		},
		{
			name:   "error - database error",
			todoID: "todo-123",
			userID: "user-456",
			setupMocks: func(repo *automock.TodoRepository) {
				repo.EXPECT().UpdateAssignedTo(ctx, "todo-123", "user-456").Return(models.Todo{}, errors.New("database error")).Once()
			},
			expectedTodo:   models.Todo{},
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo := automock.NewTodoRepository(t)
			tc.setupMocks(repo)

			svc := todos.NewService(repo, nil, nil)

			// Act
			todo, err := svc.UpdateAssignedTo(ctx, tc.todoID, tc.userID)

			// Assert
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, todo)
			}
		})
	}
}
