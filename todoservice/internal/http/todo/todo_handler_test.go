package todo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Victor-Uzunov/devops-project/todoservice/internal/http/todo"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestHandler_CreateTodo(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "success - creates todo",
			requestBody: models.Todo{
				ID:          "1",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "list-1",
				Tags:        json.RawMessage(`null`),
			},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().CreateTodo(mock.Anything, mock.MatchedBy(func(t models.Todo) bool {
					return t.Title == "Test Todo" && t.Description == "Test Description"
				})).Return("1", nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "1",
		},
		{
			name: "success - creates todo with minimal fields",
			requestBody: models.Todo{
				Title:  "Minimal Todo",
				ListID: "list-1",
				Tags:   json.RawMessage(`null`),
			},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().CreateTodo(mock.Anything, mock.Anything).Return("new-id", nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "new-id",
		},
		{
			name: "success - creates todo with priority",
			requestBody: models.Todo{
				Title:    "High Priority Todo",
				ListID:   "list-1",
				Priority: constants.PriorityHigh,
				Tags:     json.RawMessage(`null`),
			},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().CreateTodo(mock.Anything, mock.MatchedBy(func(t models.Todo) bool {
					return t.Priority == constants.PriorityHigh
				})).Return("priority-id", nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "priority-id",
		},
		{
			name: "error - service fails",
			requestBody: models.Todo{
				Title:  "Test Todo",
				ListID: "list-1",
				Tags:   json.RawMessage(`null`),
			},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().CreateTodo(mock.Anything, mock.Anything).Return("", errors.New("service error")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "error - invalid json body",
			requestBody: "invalid json",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				// No DB expectations - handler fails before starting transaction for JSON decode
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "error - empty body",
			requestBody: nil,
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				// No DB expectations as handler should fail before starting transaction
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			var body []byte
			if tc.requestBody != nil {
				if str, ok := tc.requestBody.(string); ok {
					body = []byte(str)
				} else {
					body, _ = json.Marshal(tc.requestBody)
				}
			}

			tc.setupMocks(svc, mockDB)

			req := httptest.NewRequest(http.MethodPost, "/todos/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()

			// Act
			handler.CreateTodo(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedBody != "" && tc.expectedStatusCode == http.StatusCreated {
				var respID string
				err := json.NewDecoder(resp.Body).Decode(&respID)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedBody, respID)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_GetTodo(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name:   "success - returns todo",
			todoID: "1",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{
					ID:          "1",
					Title:       "Test Todo",
					Description: "Test Description",
					ListID:      "list-1",
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.Equal(t, "1", respTodo.ID)
				assert.Equal(t, "Test Todo", respTodo.Title)
			},
		},
		{
			name:   "success - returns completed todo",
			todoID: "completed-1",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "completed-1").Return(models.Todo{
					ID:        "completed-1",
					Title:     "Completed Todo",
					Completed: true,
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.Equal(t, "completed-1", respTodo.ID)
				assert.True(t, respTodo.Completed)
			},
		},
		{
			name:   "error - todo not found",
			todoID: "non-existent",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
		{
			name:   "error - empty id",
			todoID: "",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "").Return(models.Todo{}, errors.New("invalid id")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			req := httptest.NewRequest(http.MethodGet, "/todos/"+tc.todoID, nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.GetTodo(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_UpdateTodo(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		requestBody        interface{}
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
	}{
		{
			name:   "success - updates todo",
			todoID: "1",
			requestBody: models.Todo{
				Title:       "Updated Title",
				Description: "Updated Description",
				ListID:      "list-1",
				Tags:        json.RawMessage(`null`),
			},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().UpdateTodo(mock.Anything, mock.MatchedBy(func(t models.Todo) bool {
					return t.ID == "1" && t.Title == "Updated Title"
				})).Return(nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "success - updates only title",
			todoID: "1",
			requestBody: models.Todo{
				Title: "New Title Only",
				Tags:  json.RawMessage(`null`),
			},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().UpdateTodo(mock.Anything, mock.Anything).Return(nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "error - todo not found",
			todoID: "non-existent",
			requestBody: models.Todo{
				Title: "Updated Title",
				Tags:  json.RawMessage(`null`),
			},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:        "error - invalid json body",
			todoID:      "1",
			requestBody: "invalid json",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				// No DB expectations - handler fails before starting transaction for JSON decode
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			var body []byte
			if str, ok := tc.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tc.requestBody)
			}

			tc.setupMocks(svc, mockDB)

			req := httptest.NewRequest(http.MethodPut, "/todos/update/"+tc.todoID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.UpdateTodo(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_DeleteTodo(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
	}{
		{
			name:   "success - deletes todo",
			todoID: "1",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().DeleteTodo(mock.Anything, "1").Return(nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:   "error - todo not found on get",
			todoID: "non-existent",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:   "error - delete fails",
			todoID: "1",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().DeleteTodo(mock.Anything, "1").Return(errors.New("delete failed")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			req := httptest.NewRequest(http.MethodDelete, "/todos/delete/"+tc.todoID, nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.DeleteTodo(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_ListTodosByListID(t *testing.T) {
	tests := []struct {
		name               string
		listID             string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name:   "success - returns todos for list",
			listID: "list-1",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().ListTodosByListID(mock.Anything, "list-1").Return([]models.Todo{
					{ID: "1", Title: "Todo 1", ListID: "list-1"},
					{ID: "2", Title: "Todo 2", ListID: "list-1"},
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodos []models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodos)
				require.NoError(t, err)
				assert.Len(t, respTodos, 2)
			},
		},
		{
			name:   "success - returns empty list",
			listID: "empty-list",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().ListTodosByListID(mock.Anything, "empty-list").Return([]models.Todo{}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodos []models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodos)
				require.NoError(t, err)
				assert.Len(t, respTodos, 0)
			},
		},
		{
			name:   "error - list not found",
			listID: "non-existent",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().ListTodosByListID(mock.Anything, "non-existent").Return([]models.Todo{}, errors.New("list not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			req := httptest.NewRequest(http.MethodGet, "/lists/"+tc.listID+"/todos", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"list_id": tc.listID})
			w := httptest.NewRecorder()

			// Act
			handler.ListTodosByListID(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_GetAllTodos(t *testing.T) {
	tests := []struct {
		name               string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name: "success - returns all todos",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetAllTodos(mock.Anything).Return([]models.Todo{
					{ID: "1", Title: "Todo 1", ListID: "list-1"},
					{ID: "2", Title: "Todo 2", ListID: "list-2"},
					{ID: "3", Title: "Todo 3", ListID: "list-1"},
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodos []models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodos)
				require.NoError(t, err)
				assert.Len(t, respTodos, 3)
			},
		},
		{
			name: "success - returns empty list when no todos",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetAllTodos(mock.Anything).Return([]models.Todo{}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodos []models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodos)
				require.NoError(t, err)
				assert.Len(t, respTodos, 0)
			},
		},
		{
			name: "error - service fails",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetAllTodos(mock.Anything).Return(nil, errors.New("database error")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()

			// Act
			handler.GetAllTodos(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_CompleteTodo(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name:   "success - marks todo as complete",
			todoID: "1",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1", Completed: false}, nil).Once()
				svc.EXPECT().CompleteTodo(mock.Anything, "1").Return(models.Todo{
					ID:        "1",
					Title:     "Test Todo",
					Completed: true,
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.True(t, respTodo.Completed)
			},
		},
		{
			name:   "error - todo not found",
			todoID: "non-existent",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
		{
			name:   "error - complete fails",
			todoID: "1",
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().CompleteTodo(mock.Anything, "1").Return(models.Todo{}, errors.New("complete failed")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusInternalServerError,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			req := httptest.NewRequest(http.MethodPatch, "/todos/"+tc.todoID+"/complete", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.CompleteTodo(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_UpdateTodoTitle(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		requestBody        map[string]string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name:        "success - updates title",
			todoID:      "1",
			requestBody: map[string]string{"title": "Updated Title"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().UpdateTodoTitle(mock.Anything, "1", "Updated Title").Return(models.Todo{
					ID:    "1",
					Title: "Updated Title",
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.Equal(t, "Updated Title", respTodo.Title)
			},
		},
		{
			name:        "error - empty title",
			todoID:      "1",
			requestBody: map[string]string{"title": ""},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				// No DB expectations - handler fails before starting transaction for empty title check
			},
			expectedStatusCode: http.StatusBadRequest,
			validateResponse:   nil,
		},
		{
			name:        "error - todo not found",
			todoID:      "non-existent",
			requestBody: map[string]string{"title": "New Title"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/todos/"+tc.todoID+"/title", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.UpdateTodoTitle(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_UpdateTodoDescription(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		requestBody        map[string]string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name:        "success - updates description",
			todoID:      "1",
			requestBody: map[string]string{"description": "Updated Description"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().UpdateTodoDescription(mock.Anything, "1", "Updated Description").Return(models.Todo{
					ID:          "1",
					Description: "Updated Description",
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.Equal(t, "Updated Description", respTodo.Description)
			},
		},
		{
			name:        "error - empty description",
			todoID:      "1",
			requestBody: map[string]string{"description": ""},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				// No DB expectations - handler fails before starting transaction for empty description check
			},
			expectedStatusCode: http.StatusBadRequest,
			validateResponse:   nil,
		},
		{
			name:        "error - todo not found",
			todoID:      "non-existent",
			requestBody: map[string]string{"description": "New Description"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/todos/"+tc.todoID+"/description", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.UpdateTodoDescription(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_UpdateTodoPriority(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		requestBody        map[string]string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name:        "success - updates to high priority",
			todoID:      "1",
			requestBody: map[string]string{"priority": "high"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().UpdateTodoPriority(mock.Anything, "1", constants.PriorityHigh).Return(models.Todo{
					ID:       "1",
					Priority: constants.PriorityHigh,
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.Equal(t, constants.PriorityHigh, respTodo.Priority)
			},
		},
		{
			name:        "success - updates to medium priority",
			todoID:      "1",
			requestBody: map[string]string{"priority": "medium"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().UpdateTodoPriority(mock.Anything, "1", constants.PriorityMedium).Return(models.Todo{
					ID:       "1",
					Priority: constants.PriorityMedium,
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.Equal(t, constants.PriorityMedium, respTodo.Priority)
			},
		},
		{
			name:        "success - updates to low priority",
			todoID:      "1",
			requestBody: map[string]string{"priority": "low"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				svc.EXPECT().UpdateTodoPriority(mock.Anything, "1", constants.PriorityLow).Return(models.Todo{
					ID:       "1",
					Priority: constants.PriorityLow,
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.Equal(t, constants.PriorityLow, respTodo.Priority)
			},
		},
		{
			name:        "error - invalid priority",
			todoID:      "1",
			requestBody: map[string]string{"priority": "invalid"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				// No DB expectations - handler fails before starting transaction for invalid priority
			},
			expectedStatusCode: http.StatusBadRequest,
			validateResponse:   nil,
		},
		{
			name:        "error - todo not found",
			todoID:      "non-existent",
			requestBody: map[string]string{"priority": "high"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/todos/"+tc.todoID+"/priority", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.UpdateTodoPriority(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestHandler_UpdateAssignedTo(t *testing.T) {
	tests := []struct {
		name               string
		todoID             string
		requestBody        map[string]string
		setupMocks         func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock)
		expectedStatusCode int
		validateResponse   func(t *testing.T, resp *http.Response)
	}{
		{
			name:        "success - assigns todo to user",
			todoID:      "1",
			requestBody: map[string]string{"user_id": "550e8400-e29b-41d4-a716-446655440000"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{ID: "1"}, nil).Once()
				assignedUser := "550e8400-e29b-41d4-a716-446655440000"
				svc.EXPECT().UpdateAssignedTo(mock.Anything, "1", "550e8400-e29b-41d4-a716-446655440000").Return(models.Todo{
					ID:         "1",
					Title:      "Test Todo",
					AssignedTo: &assignedUser,
				}, nil).Once()
				mockDB.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, resp *http.Response) {
				var respTodo models.Todo
				err := json.NewDecoder(resp.Body).Decode(&respTodo)
				require.NoError(t, err)
				assert.NotNil(t, respTodo.AssignedTo)
			},
		},
		{
			name:        "error - invalid uuid",
			todoID:      "1",
			requestBody: map[string]string{"user_id": "invalid-uuid"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				// No DB expectations - handler fails before starting transaction for invalid UUID
			},
			expectedStatusCode: http.StatusBadRequest,
			validateResponse:   nil,
		},
		{
			name:        "error - todo not found",
			todoID:      "non-existent",
			requestBody: map[string]string{"user_id": "550e8400-e29b-41d4-a716-446655440000"},
			setupMocks: func(svc *automock.TodoService, mockDB sqlxmock.Sqlmock) {
				mockDB.ExpectBegin()
				svc.EXPECT().GetTodo(mock.Anything, "non-existent").Return(models.Todo{}, errors.New("todo not found")).Once()
				mockDB.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			validateResponse:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mockDB, err := sqlxmock.Newx()
			require.NoError(t, err)
			defer db.Close()

			svc := automock.NewTodoService(t)
			handler := todo.NewHandler(svc, db)

			tc.setupMocks(svc, mockDB)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/todos/"+tc.todoID+"/assign", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, map[string]string{"id": tc.todoID})
			w := httptest.NewRecorder()

			// Act
			handler.UpdateAssignedTo(w, req)

			// Assert
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.validateResponse != nil {
				tc.validateResponse(t, resp)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
