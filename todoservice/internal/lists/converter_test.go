package lists_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverter_ConvertListToModel(t *testing.T) {
	converter := lists.NewConverter()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		entity   lists.Entity
		expected models.List
	}{
		{
			name: "converts complete entity with all fields",
			entity: lists.Entity{
				ID:          "list-1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-1",
				SharedWith:  []string{"user-1", "user-2"},
				Tags:        pkg.NewValidNullableString(`["work","personal"]`),
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Visibility:  constants.VisibilityShared,
			},
			expected: models.List{
				ID:          "list-1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-1",
				SharedWith:  []string{"user-1", "user-2"},
				Tags:        json.RawMessage(`["work","personal"]`),
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Visibility:  constants.VisibilityShared,
			},
		},
		{
			name: "converts entity with minimal fields",
			entity: lists.Entity{
				ID:      "list-2",
				Name:    "Minimal List",
				OwnerID: "owner-1",
			},
			expected: models.List{
				ID:      "list-2",
				Name:    "Minimal List",
				OwnerID: "owner-1",
				Tags:    nil,
			},
		},
		{
			name: "converts private list",
			entity: lists.Entity{
				ID:         "list-3",
				Name:       "Private List",
				OwnerID:    "owner-1",
				Visibility: constants.VisibilityPrivate,
			},
			expected: models.List{
				ID:         "list-3",
				Name:       "Private List",
				OwnerID:    "owner-1",
				Visibility: constants.VisibilityPrivate,
				Tags:       nil,
			},
		},
		{
			name: "converts entity with empty shared list",
			entity: lists.Entity{
				ID:         "list-4",
				Name:       "Unshared List",
				OwnerID:    "owner-1",
				SharedWith: []string{},
			},
			expected: models.List{
				ID:         "list-4",
				Name:       "Unshared List",
				OwnerID:    "owner-1",
				SharedWith: []string{},
				Tags:       nil,
			},
		},
		{
			name: "converts entity with nil shared list",
			entity: lists.Entity{
				ID:         "list-5",
				Name:       "Nil Shared List",
				OwnerID:    "owner-1",
				SharedWith: nil,
			},
			expected: models.List{
				ID:         "list-5",
				Name:       "Nil Shared List",
				OwnerID:    "owner-1",
				SharedWith: nil,
				Tags:       nil,
			},
		},
		{
			name: "converts entity with many shared users",
			entity: lists.Entity{
				ID:         "list-6",
				Name:       "Popular List",
				OwnerID:    "owner-1",
				SharedWith: []string{"user-1", "user-2", "user-3", "user-4", "user-5"},
			},
			expected: models.List{
				ID:         "list-6",
				Name:       "Popular List",
				OwnerID:    "owner-1",
				SharedWith: []string{"user-1", "user-2", "user-3", "user-4", "user-5"},
				Tags:       nil,
			},
		},
		{
			name: "converts entity with empty tags",
			entity: lists.Entity{
				ID:      "list-7",
				Name:    "No Tags List",
				OwnerID: "owner-1",
				Tags:    pkg.NewValidNullableString(""),
			},
			expected: models.List{
				ID:      "list-7",
				Name:    "No Tags List",
				OwnerID: "owner-1",
				Tags:    json.RawMessage(""),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := converter.ConvertListToModel(tc.entity)

			// Assert
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.Description, result.Description)
			assert.Equal(t, tc.expected.OwnerID, result.OwnerID)
			assert.Equal(t, tc.expected.Visibility, result.Visibility)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)

			if tc.expected.SharedWith != nil {
				require.NotNil(t, result.SharedWith)
				assert.Equal(t, tc.expected.SharedWith, result.SharedWith)
			} else {
				assert.Nil(t, result.SharedWith)
			}
		})
	}
}

func TestConverter_ConvertListToEntity(t *testing.T) {
	converter := lists.NewConverter()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		model    models.List
		expected lists.Entity
	}{
		{
			name: "converts complete model with all fields",
			model: models.List{
				ID:          "list-1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-1",
				SharedWith:  []string{"user-1", "user-2"},
				Tags:        json.RawMessage(`["work","personal"]`),
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Visibility:  constants.VisibilityShared,
			},
			expected: lists.Entity{
				ID:          "list-1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-1",
				SharedWith:  []string{"user-1", "user-2"},
				Tags:        pkg.NewValidNullableString(`["work","personal"]`),
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Visibility:  constants.VisibilityShared,
			},
		},
		{
			name: "converts model with minimal fields",
			model: models.List{
				ID:      "list-2",
				Name:    "Minimal List",
				OwnerID: "owner-1",
			},
			expected: lists.Entity{
				ID:      "list-2",
				Name:    "Minimal List",
				OwnerID: "owner-1",
			},
		},
		{
			name: "converts private model",
			model: models.List{
				ID:         "list-3",
				Name:       "Private List",
				OwnerID:    "owner-1",
				Visibility: constants.VisibilityPrivate,
			},
			expected: lists.Entity{
				ID:         "list-3",
				Name:       "Private List",
				OwnerID:    "owner-1",
				Visibility: constants.VisibilityPrivate,
			},
		},
		{
			name: "converts model with empty description",
			model: models.List{
				ID:          "list-4",
				Name:        "List without description",
				Description: "",
				OwnerID:     "owner-1",
			},
			expected: lists.Entity{
				ID:          "list-4",
				Name:        "List without description",
				Description: "",
				OwnerID:     "owner-1",
			},
		},
		{
			name: "converts model with long description",
			model: models.List{
				ID:          "list-5",
				Name:        "Detailed List",
				Description: "This is a very long description that contains many characters and explains the purpose of this list in great detail.",
				OwnerID:     "owner-1",
			},
			expected: lists.Entity{
				ID:          "list-5",
				Name:        "Detailed List",
				Description: "This is a very long description that contains many characters and explains the purpose of this list in great detail.",
				OwnerID:     "owner-1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := converter.ConvertListToEntity(tc.model)

			// Assert
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.Description, result.Description)
			assert.Equal(t, tc.expected.OwnerID, result.OwnerID)
			assert.Equal(t, tc.expected.Visibility, result.Visibility)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)

			if tc.expected.SharedWith != nil {
				require.NotNil(t, result.SharedWith)
				assert.Equal(t, tc.expected.SharedWith, result.SharedWith)
			}
		})
	}
}

func TestNewConverter(t *testing.T) {
	// Act
	converter := lists.NewConverter()

	// Assert
	assert.NotNil(t, converter)
}
