package todos_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_ConvertTodoToModel(t *testing.T) {
	converter := todos.NewConverter()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	assignedUser := "user-123"

	tests := []struct {
		name     string
		entity   todos.Entity
		expected models.Todo
	}{
		{
			name: "converts complete entity with all fields",
			entity: todos.Entity{
				ID:          "todo-1",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "list-1",
				Tags:        pkg.NewValidNullableString(`["work","urgent"]`),
				Completed:   false,
				DueDate:     sql.NullTime{Time: fixedTime, Valid: true},
				StartDate:   sql.NullTime{Time: fixedTime, Valid: true},
				Priority:    constants.PriorityHigh,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				AssignedTo:  &assignedUser,
			},
			expected: models.Todo{
				ID:          "todo-1",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "list-1",
				Tags:        json.RawMessage(`["work","urgent"]`),
				Completed:   false,
				DueDate:     &fixedTime,
				StartDate:   &fixedTime,
				Priority:    constants.PriorityHigh,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				AssignedTo:  &assignedUser,
			},
		},
		{
			name: "converts entity with minimal fields",
			entity: todos.Entity{
				ID:     "todo-2",
				Title:  "Minimal Todo",
				ListID: "list-1",
			},
			expected: models.Todo{
				ID:     "todo-2",
				Title:  "Minimal Todo",
				ListID: "list-1",
				Tags:   nil,
			},
		},
		{
			name: "converts completed todo",
			entity: todos.Entity{
				ID:        "todo-3",
				Title:     "Completed Todo",
				ListID:    "list-1",
				Completed: true,
				Priority:  constants.PriorityLow,
			},
			expected: models.Todo{
				ID:        "todo-3",
				Title:     "Completed Todo",
				ListID:    "list-1",
				Completed: true,
				Priority:  constants.PriorityLow,
				Tags:      nil,
			},
		},
		{
			name: "converts entity with null dates",
			entity: todos.Entity{
				ID:        "todo-4",
				Title:     "Todo without dates",
				ListID:    "list-1",
				DueDate:   sql.NullTime{Valid: false},
				StartDate: sql.NullTime{Valid: false},
			},
			expected: models.Todo{
				ID:        "todo-4",
				Title:     "Todo without dates",
				ListID:    "list-1",
				DueDate:   nil,
				StartDate: nil,
				Tags:      nil,
			},
		},
		{
			name: "converts entity with empty tags",
			entity: todos.Entity{
				ID:     "todo-5",
				Title:  "Todo with empty tags",
				ListID: "list-1",
				Tags:   pkg.NewValidNullableString(""),
			},
			expected: models.Todo{
				ID:     "todo-5",
				Title:  "Todo with empty tags",
				ListID: "list-1",
				Tags:   json.RawMessage(""),
			},
		},
		{
			name: "converts entity with medium priority",
			entity: todos.Entity{
				ID:       "todo-6",
				Title:    "Medium Priority Todo",
				ListID:   "list-1",
				Priority: constants.PriorityMedium,
			},
			expected: models.Todo{
				ID:       "todo-6",
				Title:    "Medium Priority Todo",
				ListID:   "list-1",
				Priority: constants.PriorityMedium,
				Tags:     nil,
			},
		},
		{
			name: "converts entity without assigned user",
			entity: todos.Entity{
				ID:         "todo-7",
				Title:      "Unassigned Todo",
				ListID:     "list-1",
				AssignedTo: nil,
			},
			expected: models.Todo{
				ID:         "todo-7",
				Title:      "Unassigned Todo",
				ListID:     "list-1",
				AssignedTo: nil,
				Tags:       nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := converter.ConvertTodoToModel(tc.entity)

			// Assert
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Title, result.Title)
			assert.Equal(t, tc.expected.Description, result.Description)
			assert.Equal(t, tc.expected.ListID, result.ListID)
			assert.Equal(t, tc.expected.Completed, result.Completed)
			assert.Equal(t, tc.expected.Priority, result.Priority)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)

			if tc.expected.DueDate != nil {
				require.NotNil(t, result.DueDate)
				assert.Equal(t, *tc.expected.DueDate, *result.DueDate)
			} else {
				assert.Nil(t, result.DueDate)
			}

			if tc.expected.StartDate != nil {
				require.NotNil(t, result.StartDate)
				assert.Equal(t, *tc.expected.StartDate, *result.StartDate)
			} else {
				assert.Nil(t, result.StartDate)
			}

			if tc.expected.AssignedTo != nil {
				require.NotNil(t, result.AssignedTo)
				assert.Equal(t, *tc.expected.AssignedTo, *result.AssignedTo)
			} else {
				assert.Nil(t, result.AssignedTo)
			}
		})
	}
}

func TestConverter_ConvertTodoToEntity(t *testing.T) {
	converter := todos.NewConverter()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	zeroTime := time.Time{}
	assignedUser := "user-123"

	tests := []struct {
		name     string
		model    models.Todo
		expected todos.Entity
	}{
		{
			name: "converts complete model with all fields",
			model: models.Todo{
				ID:          "todo-1",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "list-1",
				Tags:        json.RawMessage(`["work","urgent"]`),
				Completed:   false,
				DueDate:     &fixedTime,
				StartDate:   &fixedTime,
				Priority:    constants.PriorityHigh,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				AssignedTo:  &assignedUser,
			},
			expected: todos.Entity{
				ID:          "todo-1",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "list-1",
				Tags:        pkg.NewValidNullableString(`["work","urgent"]`),
				Completed:   false,
				DueDate:     sql.NullTime{Time: fixedTime, Valid: true},
				StartDate:   sql.NullTime{Time: fixedTime, Valid: true},
				Priority:    constants.PriorityHigh,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				AssignedTo:  &assignedUser,
			},
		},
		{
			name: "converts model with zero time dates",
			model: models.Todo{
				ID:        "todo-2",
				Title:     "Todo with zero time",
				ListID:    "list-1",
				DueDate:   &zeroTime,
				StartDate: &zeroTime,
			},
			expected: todos.Entity{
				ID:        "todo-2",
				Title:     "Todo with zero time",
				ListID:    "list-1",
				DueDate:   sql.NullTime{Valid: false},
				StartDate: sql.NullTime{Valid: false},
			},
		},
		{
			name: "converts completed model",
			model: models.Todo{
				ID:        "todo-3",
				Title:     "Completed Todo",
				ListID:    "list-1",
				Completed: true,
				DueDate:   &fixedTime,
				StartDate: &fixedTime,
			},
			expected: todos.Entity{
				ID:        "todo-3",
				Title:     "Completed Todo",
				ListID:    "list-1",
				Completed: true,
				DueDate:   sql.NullTime{Time: fixedTime, Valid: true},
				StartDate: sql.NullTime{Time: fixedTime, Valid: true},
			},
		},
		{
			name: "converts model with all priority levels - low",
			model: models.Todo{
				ID:        "todo-4",
				Title:     "Low Priority",
				ListID:    "list-1",
				Priority:  constants.PriorityLow,
				DueDate:   &fixedTime,
				StartDate: &fixedTime,
			},
			expected: todos.Entity{
				ID:        "todo-4",
				Title:     "Low Priority",
				ListID:    "list-1",
				Priority:  constants.PriorityLow,
				DueDate:   sql.NullTime{Time: fixedTime, Valid: true},
				StartDate: sql.NullTime{Time: fixedTime, Valid: true},
			},
		},
		{
			name: "converts model with all priority levels - medium",
			model: models.Todo{
				ID:        "todo-5",
				Title:     "Medium Priority",
				ListID:    "list-1",
				Priority:  constants.PriorityMedium,
				DueDate:   &fixedTime,
				StartDate: &fixedTime,
			},
			expected: todos.Entity{
				ID:        "todo-5",
				Title:     "Medium Priority",
				ListID:    "list-1",
				Priority:  constants.PriorityMedium,
				DueDate:   sql.NullTime{Time: fixedTime, Valid: true},
				StartDate: sql.NullTime{Time: fixedTime, Valid: true},
			},
		},
		{
			name: "converts model with empty description",
			model: models.Todo{
				ID:          "todo-6",
				Title:       "Todo without description",
				Description: "",
				ListID:      "list-1",
				DueDate:     &fixedTime,
				StartDate:   &fixedTime,
			},
			expected: todos.Entity{
				ID:          "todo-6",
				Title:       "Todo without description",
				Description: "",
				ListID:      "list-1",
				DueDate:     sql.NullTime{Time: fixedTime, Valid: true},
				StartDate:   sql.NullTime{Time: fixedTime, Valid: true},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := converter.ConvertTodoToEntity(tc.model)

			// Assert
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Title, result.Title)
			assert.Equal(t, tc.expected.Description, result.Description)
			assert.Equal(t, tc.expected.ListID, result.ListID)
			assert.Equal(t, tc.expected.Completed, result.Completed)
			assert.Equal(t, tc.expected.Priority, result.Priority)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
			assert.Equal(t, tc.expected.DueDate.Valid, result.DueDate.Valid)
			assert.Equal(t, tc.expected.StartDate.Valid, result.StartDate.Valid)

			if tc.expected.AssignedTo != nil {
				require.NotNil(t, result.AssignedTo)
				assert.Equal(t, *tc.expected.AssignedTo, *result.AssignedTo)
			}
		})
	}
}

func TestNewConverter(t *testing.T) {
	// Act
	converter := todos.NewConverter()

	// Assert
	assert.NotNil(t, converter)
}
