module github.com/Victor-Uzunov/devops-project/graphqlServer

go 1.26.1

require (
	github.com/99designs/gqlgen v0.17.49
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/stretchr/testify v1.9.0
	github.com/vektah/gqlparser/v2 v2.5.16
)

require github.com/Victor-Uzunov/devops-project/todoservice v0.0.0-20241222102949-251bf1f433bb

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/cors v1.11.1 //direct
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/sys v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
