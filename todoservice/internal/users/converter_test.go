package users_test

import (
	"testing"
	"time"

	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestConverter_ConvertUserToModel(t *testing.T) {
	converter := users.NewConverter()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		entity   users.Entity
		expected models.User
	}{
		{
			name: "converts complete entity with all fields",
			entity: users.Entity{
				ID:        "user-1",
				Email:     "test@example.com",
				GithubID:  "github-123",
				Role:      "admin",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: models.User{
				ID:        "user-1",
				Email:     "test@example.com",
				GithubID:  "github-123",
				Role:      "admin",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
		},
		{
			name: "converts entity with minimal fields",
			entity: users.Entity{
				ID:    "user-2",
				Email: "minimal@example.com",
			},
			expected: models.User{
				ID:    "user-2",
				Email: "minimal@example.com",
			},
		},
		{
			name: "converts entity with reader role",
			entity: users.Entity{
				ID:       "user-3",
				Email:    "reader@example.com",
				GithubID: "github-456",
				Role:     "reader",
			},
			expected: models.User{
				ID:       "user-3",
				Email:    "reader@example.com",
				GithubID: "github-456",
				Role:     "reader",
			},
		},
		{
			name: "converts entity with writer role",
			entity: users.Entity{
				ID:       "user-4",
				Email:    "writer@example.com",
				GithubID: "github-789",
				Role:     "writer",
			},
			expected: models.User{
				ID:       "user-4",
				Email:    "writer@example.com",
				GithubID: "github-789",
				Role:     "writer",
			},
		},
		{
			name: "converts entity with zero time",
			entity: users.Entity{
				ID:        "user-5",
				Email:     "notime@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			expected: models.User{
				ID:        "user-5",
				Email:     "notime@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		{
			name: "converts entity with empty github id",
			entity: users.Entity{
				ID:       "user-6",
				Email:    "nogithub@example.com",
				GithubID: "",
				Role:     "reader",
			},
			expected: models.User{
				ID:       "user-6",
				Email:    "nogithub@example.com",
				GithubID: "",
				Role:     "reader",
			},
		},
		{
			name: "converts entity with long email",
			entity: users.Entity{
				ID:       "user-7",
				Email:    "very.long.email.address.that.is.quite.lengthy@example.company.com",
				GithubID: "github-long",
				Role:     "admin",
			},
			expected: models.User{
				ID:       "user-7",
				Email:    "very.long.email.address.that.is.quite.lengthy@example.company.com",
				GithubID: "github-long",
				Role:     "admin",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := converter.ConvertUserToModel(tc.entity)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConverter_ConvertUserToEntity(t *testing.T) {
	converter := users.NewConverter()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		model    models.User
		expected users.Entity
	}{
		{
			name: "converts complete model with all fields",
			model: models.User{
				ID:        "user-1",
				Email:     "test@example.com",
				GithubID:  "github-123",
				Role:      "admin",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: users.Entity{
				ID:        "user-1",
				Email:     "test@example.com",
				GithubID:  "github-123",
				Role:      "admin",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
		},
		{
			name: "converts model with minimal fields",
			model: models.User{
				ID:    "user-2",
				Email: "minimal@example.com",
			},
			expected: users.Entity{
				ID:    "user-2",
				Email: "minimal@example.com",
			},
		},
		{
			name: "converts model with reader role",
			model: models.User{
				ID:       "user-3",
				Email:    "reader@example.com",
				GithubID: "github-456",
				Role:     "reader",
			},
			expected: users.Entity{
				ID:       "user-3",
				Email:    "reader@example.com",
				GithubID: "github-456",
				Role:     "reader",
			},
		},
		{
			name: "converts model with writer role",
			model: models.User{
				ID:       "user-4",
				Email:    "writer@example.com",
				GithubID: "github-789",
				Role:     "writer",
			},
			expected: users.Entity{
				ID:       "user-4",
				Email:    "writer@example.com",
				GithubID: "github-789",
				Role:     "writer",
			},
		},
		{
			name: "converts model with zero time",
			model: models.User{
				ID:        "user-5",
				Email:     "notime@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			expected: users.Entity{
				ID:        "user-5",
				Email:     "notime@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		{
			name: "converts model with different created and updated times",
			model: models.User{
				ID:        "user-6",
				Email:     "times@example.com",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime.Add(24 * time.Hour),
			},
			expected: users.Entity{
				ID:        "user-6",
				Email:     "times@example.com",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime.Add(24 * time.Hour),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := converter.ConvertUserToEntity(tc.model)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNewConverter(t *testing.T) {
	// Act
	converter := users.NewConverter()

	// Assert
	assert.NotNil(t, converter)
}
