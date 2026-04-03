# TodoList Project

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Technology Stack](#technology-stack)
4. [System Components](#system-components)
5. [Installation and Setup](#installation-and-setup)
6. [API Documentation](#api-documentation)
7. [Database](#database)
8. [Security and Authentication](#security-and-authentication)
9. [Deployment](#deployment)
10. [CI/CD Pipeline](#cicd-pipeline)
11. [Usage](#usage)
12. [Troubleshooting](#troubleshooting)

---

## Overview

**TodoList Project** is a modern microservices application for task management (Todo Lists), inspired by Todoist. The project implements a full-featured task management system with emphasis on security, scalability, and modern DevOps practices.

### Key Features

- **Secure authentication** via GitHub OAuth 2.0
- **JWT-based authorization** with refresh token mechanism
- **Role-based access control (RBAC)** - Reader, Writer, Admin roles
- **REST API** for core CRUD operations
- **GraphQL API** as a facade for more flexible queries
- **PostgreSQL** database with migrations
- **SAP UI5** based user interface
- **Kubernetes deployment** with Istio service mesh
- **CI/CD pipeline** with GitHub Actions
- **Docker containerization** for all components

### Business Features

1. **Todo List Management**
    - Create, edit, and delete lists
    - Share lists with other users
    - Three visibility levels: Private, Shared, Public
    - Add tags to lists

2. **Task Management (Todos)**
    - Create, edit, and delete tasks
    - Mark tasks as completed
    - Priorities: Low, Medium, High
    - Due dates and start dates
    - Assign tasks to users
    - Tags for organization

3. **Collaboration**
    - Share lists with colleagues
    - Three access levels: Reader, Writer, Admin
    - Invitation and approval system
    - Collaborator management

---

## Architecture

### Microservices Architecture

The project follows a microservices architectural pattern with clear separation of responsibilities:

```
┌─────────────────────────────────────────────────────────────┐
│                      User Browser                            │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    UI Service (SAP UI5)                      │
│                    Port: 80 (Nginx)                          │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              GraphQL Server (Go + gqlgen)                    │
│                    Port: 8000                                │
│              - Schema validation                             │
│              - Request aggregation                           │
│              - REST API facade                               │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              REST API Server (Go + Gin)                      │
│                    Port: 5000                                │
│              - Business logic                                │
│              - Authentication/Authorization                  │
│              - Data validation                               │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL Database                             │
│                    Port: 5432                                │
│              - Data persistence                              │
│              - Relational integrity                          │
└─────────────────────────────────────────────────────────────┘
```

### Kubernetes Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Istio Ingress Gateway                     │
│                    (Port 8000 exposed)                       │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
    ┌────────┐    ┌──────────┐    ┌──────────┐
    │   UI   │    │ GraphQL  │    │   REST   │
    │  Pod   │    │   Pod    │    │   Pod    │
    └────────┘    └──────────┘    └─────┬────┘
                                         │
                                         ▼
                                  ┌──────────┐
                                  │ Postgres │
                                  │   Pod    │
                                  └──────────┘
```

### Layered Architecture (REST Service)

```
┌─────────────────────────────────────────┐
│         HTTP Layer (Routers)            │
│  - Route definitions                    │
│  - Request/Response handling            │
│  - Middleware integration               │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Service Layer                   │
│  - Business logic                       │
│  - Data validation                      │
│  - Authorization checks                 │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Repository Layer                │
│  - Database queries                     │
│  - Data mapping                         │
│  - Transaction management               │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Database Layer                  │
│  - PostgreSQL                           │
│  - Schema management                    │
└─────────────────────────────────────────┘
```

---

## Technology Stack

### Backend

#### REST API Service
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database Driver**: pgx (PostgreSQL)
- **Authentication**: OAuth2, JWT
- **Configuration**: envconfig, godotenv
- **Logging**: Structured logging

#### GraphQL Service
- **Language**: Go 1.21+
- **GraphQL Framework**: gqlgen
- **HTTP Client**: Standard Go http client
- **Schema**: GraphQL SDL (Schema Definition Language)

### Frontend

- **Framework**: SAP UI5
- **Web Server**: Nginx
- **UI Components**: SAP Fiori design system
- **Routing**: SAP UI5 Router

### Database

- **RDBMS**: PostgreSQL 14+
- **Migrations**: golang-migrate
- **Connection Pooling**: pgx pool

### Infrastructure

- **Containerization**: Docker
- **Orchestration**: Kubernetes (k3d)
- **Service Mesh**: Istio 1.20+
- **Package Manager**: Helm 3
- **Local Cluster**: k3d (k3s in Docker)

### DevOps & CI/CD

- **CI/CD**: GitHub Actions
- **Code Quality**: SonarCloud
- **Security Scanning**: Snyk
- **Linting**: golangci-lint
- **Container Registry**: Docker Hub

---

## System Components

### 1. REST API Service (`todoservice`)

**Location**: `/todoservice`

#### Structure

```
todoservice/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── db/                     # Database connection
│   ├── http/                   # HTTP server & routers
│   ├── lists/                  # Lists domain
│   │   ├── repository.go       # DB operations
│   │   ├── service.go          # Business logic
│   │   └── handlers.go         # HTTP handlers
│   ├── todos/                  # Todos domain
│   ├── users/                  # Users domain
│   └── oauth2/                 # OAuth2 handlers
├── pkg/
│   ├── jwt/                    # JWT utilities
│   ├── log/                    # Logging
│   ├── models/                 # Data models
│   └── constants/              # Constants
├── Dockerfile
└── go.mod
```

#### Core Functions

**Authentication & Authorization**
- GitHub OAuth2 login flow
- JWT token generation
- Refresh token mechanism
- Role-based middleware
- Tenant isolation

**Lists Management**
- CRUD operations for lists
- Visibility control (private/shared/public)
- Owner management
- Tag support

**Todos Management**
- CRUD operations for tasks
- Priority management
- Due dates and start dates
- Task assignment
- Completion tracking

**Users Management**
- User profile management
- Role assignment
- List access management

#### API Endpoints

**Authentication**
```
GET  /login/github               # Initiate OAuth flow
GET  /login/github/callback      # OAuth callback
POST /auth/refresh               # Refresh access token
GET  /logout                     # Logout user
```

**Lists**
```
GET    /lists                    # Get all accessible lists
GET    /lists/:id                # Get specific list
POST   /lists                    # Create new list
PUT    /lists/:id                # Update list
DELETE /lists/:id                # Delete list
POST   /lists/:id/access         # Grant access to user
DELETE /lists/:id/access/:userId # Remove access
```

**Todos**
```
GET    /todos                    # Get all todos
GET    /todos/:id                # Get specific todo
POST   /todos                    # Create new todo
PUT    /todos/:id                # Update todo
DELETE /todos/:id                # Delete todo
PATCH  /todos/:id/complete       # Mark as completed
```

**Users**
```
GET    /users                    # Get all users (admin)
GET    /users/:id                # Get user by ID
GET    /users/me                 # Get current user
PUT    /users/:id                # Update user
```

### 2. GraphQL Service (`graphqlServer`)

**Location**: `/graphqlServer`

#### Structure

```
graphqlServer/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── graph/
│   │   └── schema.graphqls     # GraphQL schema
│   ├── resolvers/              # Query/Mutation resolvers
│   │   ├── root_resolver.go
│   │   ├── query_resolver.go
│   │   ├── mutation_resolver.go
│   │   └── directives.go       # Custom directives
│   ├── client/
│   │   └── todoservice.go      # REST API client
│   ├── converters/             # Data converters
│   └── server/
│       └── server.go           # HTTP server setup
├── generated/                  # Generated code (gqlgen)
├── gqlgen.yml                  # gqlgen configuration
└── Dockerfile
```

#### GraphQL Schema

**Types**
```graphql
type User {
  id: ID!
  email: String!
  githubID: String!
  role: UserRole!
  createdAt: String!
  updatedAt: String!
}

type List {
  id: ID!
  name: String!
  description: String
  owner: User!
  visibility: Visibility!
  tags: [String!]
  createdAt: String!
  updatedAt: String!
  todos: [Todo!]!
  collaborators: [ListAccess!]!
}

type Todo {
  id: ID!
  list: List!
  title: String!
  description: String
  completed: Boolean!
  dueDate: String
  startDate: String
  priority: Priority
  tags: [String!]
  assignedTo: User
  createdAt: String!
  updatedAt: String!
}
```

**Queries**
```graphql
type Query {
  # Users
  users: [User!]!
  user(id: ID!): User
  userByEmail: User
  usersByList(id: ID!): [User!]!

  # Lists
  lists: [List!]!
  list(id: ID!): List
  listsGlobal: [List!]!
  listsPending: [List!]!
  listsAccepted: [List!]!

  # Todos
  todos: [Todo!]!
  todo(id: ID!): Todo
  todosGlobal: [Todo!]!
  todosByList(id: ID!): [Todo!]!

  # Access
  getListAccesses(listId: ID!): [ListAccess!]!
}
```

**Mutations**
```graphql
type Mutation {
  # User mutations
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): User!

  # List mutations
  createList(input: CreateListInput!): List!
  updateList(id: ID!, input: UpdateListInput!): List!
  deleteList(id: ID!): List!

  # Todo mutations
  createTodo(input: CreateTodoInput!): Todo!
  updateTodo(id: ID!, input: UpdateTodoInput!): Todo!
  completeTodo(id: ID!): Todo!
  deleteTodo(id: ID!): Todo!

  # Access mutations
  addListAccess(input: GrantListAccessInput!): ListAccess!
  removeListAccess(listId: ID!): ListAccess!
  acceptList(listId: ID!): Boolean
  removeCollaborator(listId: ID!, userId: ID!): ListAccess!
}
```

#### Custom Directives

**@validate** - Input data validation
```graphql
directive @validate(type: String!) on INPUT_FIELD_DEFINITION

input CreateUserInput {
  email: String! @validate(type: "email")
  githubId: String!
  role: UserRole!
}
```

### 3. UI Service (`ui`)

**Location**: `/ui`

#### Structure

```
ui/
├── webapp/
│   ├── controller/             # UI Controllers
│   │   ├── App.controller.js
│   │   ├── Login.controller.js
│   │   ├── Lists.controller.js
│   │   ├── TodoDetails.controller.js
│   │   └── Settings.controller.js
│   ├── view/                   # XML Views
│   │   ├── App.view.xml
│   │   ├── Login.view.xml
│   │   ├── Lists.view.xml
│   │   └── TodoDetails.view.xml
│   ├── css/                    # Stylesheets
│   ├── util/                   # Utilities
│   ├── formatter/              # Data formatters
│   ├── Component.js            # UI5 Component
│   ├── manifest.json           # App descriptor
│   └── index.html              # Entry point
├── nginx.conf                  # Nginx configuration
├── Dockerfile
└── package.json
```

#### Features

**Login Page**
- GitHub OAuth initiation
- Error handling
- Redirect after login

**Lists View**
- List of all lists
- Filter by visibility
- Create new list
- Edit list
- Delete list
- Invite collaborators

**Todo Details View**
- List of tasks in selected list
- Create new task
- Edit task
- Mark as completed
- Delete task
- Assign task

**Settings View**
- User profile information
- Pending invitations
- Accept/Reject invitations

### 4. Database (`migrations`)

**Location**: `/migrations`

#### Schema Overview

**Users Table**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    github_id VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Lists Table**
```sql
CREATE TABLE lists (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID REFERENCES users(id),
    visibility VARCHAR(50) NOT NULL,
    tags JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Todos Table**
```sql
CREATE TABLE todos (
    id UUID PRIMARY KEY,
    list_id UUID REFERENCES lists(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    due_date TIMESTAMP,
    start_date TIMESTAMP,
    priority VARCHAR(50),
    tags JSONB,
    assigned_to UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**List Access Table**
```sql
CREATE TABLE list_access (
    list_id UUID REFERENCES lists(id),
    user_id UUID REFERENCES users(id),
    access_level VARCHAR(50) NOT NULL,
    status VARCHAR(50),
    PRIMARY KEY (list_id, user_id)
);
```

**Refresh Tokens Table**
```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### Migrations

The project uses golang-migrate for database migration management:

```
migrations/
├── 20240819110300_initial_schema.up.sql
├── 20240819110300_initial_schema.down.sql
├── 20240819121433_update_timestamp_trigger.up.sql
├── 20240819121433_update_timestamp_trigger.down.sql
├── 20241006064212_add_status_column_to_list_access.up.sql
├── 20241006064212_add_status_column_to_list_access.down.sql
├── 20241007083941_allow_null_assigned_to.up.sql
├── 20241007083941_allow_null_assigned_to.down.sql
├── 20241009095830_add_start_date_column.up.sql
├── 20241009095830_add_start_date_column.down.sql
├── 20241014115800_modify_start_and_due_date_columns.up.sql
├── 20241014115800_modify_start_and_due_date_columns.down.sql
├── 20241017111456_refresh_token.up.sql
└── 20241017111456_refresh_token.down.sql
```

---

## Installation and Setup

### Prerequisites

Before starting, ensure you have the following installed:

- **Docker** (v20.10+)
- **k3d** (v5.0+)
- **kubectl** (v1.20+)
- **Helm** (v3.0+)
- **istioctl** (v1.20+)

### Step 1: Clone the Repository

```bash
git clone <your-repo-url>
cd Todolist
```

### Step 2: Create Configuration File

```bash
cp config.example.yaml config.yaml
```

### Step 3: Set Up GitHub OAuth Application

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click **"New OAuth App"**
3. Fill in the details:
    - **Application name**: TodoList App
    - **Homepage URL**: `http://localhost:8000`
    - **Authorization callback URL**: `http://localhost:8000/login/github/callback`
4. Copy the **Client ID** and **Client Secret**

### Step 4: Generate Security Keys

Use the script:

```bash
chmod +x generate-keys.sh
./generate-keys.sh
```

Or manually:

```bash
# OAuth2 State
openssl rand -base64 32

# JWT Key
openssl rand -base64 32
```

### Step 5: Configure config.yaml

Edit `config.yaml`:

```yaml
github:
  oauth:
    clientId: "YOUR_GITHUB_CLIENT_ID"
    clientSecret: "YOUR_GITHUB_CLIENT_SECRET"
    redirectUrl: "http://localhost:8000/login/github/callback"
    scopes: "read:org,user"
  organization: "YOUR_GITHUB_ORG"

security:
  oauth2State: "GENERATED_STATE"
  jwtKey: "GENERATED_JWT_KEY"

database:
  password: "YOUR_SECURE_PASSWORD"
```

### Step 6: Run the Setup Script

```bash
chmod +x setup.sh
./setup.sh
```

The script will:
- Create a k3d Kubernetes cluster
- Install Istio service mesh
- Create the necessary namespaces
- Deploy all components
- Run database migrations

### Step 7: Verify Deployment

```bash
# Check pods
kubectl get pods -n todoapp-system

# Check services
kubectl get svc -n todoapp-system

# Check Istio gateway
kubectl get gateway -n todoapp-system
```

### Step 8: Access the Application

Open your browser at:
```
http://localhost:8000
```

---

## API Documentation

### REST API

**Base URL**: `http://localhost:5000`

#### Authentication Flow

**1. Initiate OAuth Login**
```http
GET /login/github
```

Redirects to GitHub for authentication.

**2. OAuth Callback**
```http
GET /login/github/callback?code=<code>&state=<state>
```

Response:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 3600,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "role": "writer"
  }
}
```

**3. Refresh Token**
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

Response:
```json
{
  "access_token": "eyJhbGc...",
  "expires_in": 3600
}
```

#### Lists API

**Get All Lists**
```http
GET /lists
Authorization: Bearer <token>
```

Response:
```json
[
  {
    "id": "uuid",
    "name": "Work Tasks",
    "description": "Tasks for work projects",
    "owner_id": "uuid",
    "visibility": "private",
    "tags": ["work", "important"],
    "created_at": "2026-01-01T10:00:00Z",
    "updated_at": "2026-01-01T10:00:00Z"
  }
]
```

**Create List**
```http
POST /lists
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Shopping List",
  "description": "Weekly shopping items",
  "visibility": "private",
  "tags": ["personal", "shopping"]
}
```

**Update List**
```http
PUT /lists/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Name",
  "description": "Updated description",
  "visibility": "shared"
}
```

**Delete List**
```http
DELETE /lists/:id
Authorization: Bearer <token>
```

**Grant List Access**
```http
POST /lists/:id/access
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": "uuid",
  "access_level": "writer"
}
```

#### Todos API

**Get All Todos**
```http
GET /todos
Authorization: Bearer <token>
```

**Get Todos by List**
```http
GET /todos?list_id=<uuid>
Authorization: Bearer <token>
```

**Create Todo**
```http
POST /todos
Authorization: Bearer <token>
Content-Type: application/json

{
  "list_id": "uuid",
  "title": "Complete documentation",
  "description": "Write comprehensive docs",
  "priority": "high",
  "due_date": "2026-01-15T17:00:00Z",
  "tags": ["documentation", "urgent"]
}
```

**Update Todo**
```http
PUT /todos/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated title",
  "completed": true,
  "priority": "medium"
}
```

**Complete Todo**
```http
PATCH /todos/:id/complete
Authorization: Bearer <token>
```

**Delete Todo**
```http
DELETE /todos/:id
Authorization: Bearer <token>
```

### GraphQL API

**Endpoint**: `http://localhost:8000/graphql`

#### Example Queries

**Get All Lists with Todos**
```graphql
query {
  lists {
    id
    name
    description
    visibility
    owner {
      email
    }
    todos {
      id
      title
      completed
      priority
    }
  }
}
```

**Get Specific List**
```graphql
query GetList($id: ID!) {
  list(id: $id) {
    id
    name
    description
    todos {
      id
      title
      completed
      dueDate
      assignedTo {
        email
      }
    }
    collaborators {
      user {
        email
      }
      accessLevel
    }
  }
}
```

**Get Current User Info**
```graphql
query {
  userByEmail {
    id
    email
    role
    createdAt
  }
}
```

#### Example Mutations

**Create List**
```graphql
mutation CreateList($input: CreateListInput!) {
  createList(input: $input) {
    id
    name
    description
    visibility
  }
}

# Variables
{
  "input": {
    "name": "My New List",
    "description": "Description here",
    "visibility": "PRIVATE",
    "tags": ["work"]
  }
}
```

**Create Todo**
```graphql
mutation CreateTodo($input: CreateTodoInput!) {
  createTodo(input: $input) {
    id
    title
    priority
    dueDate
  }
}

# Variables
{
  "input": {
    "listId": "uuid",
    "title": "New Task",
    "description": "Task description",
    "priority": "HIGH",
    "dueDate": "2026-01-15T17:00:00Z"
  }
}
```

**Update Todo**
```graphql
mutation UpdateTodo($id: ID!, $input: UpdateTodoInput!) {
  updateTodo(id: $id, input: $input) {
    id
    title
    completed
  }
}

# Variables
{
  "id": "uuid",
  "input": {
    "completed": true,
    "priority": "MEDIUM"
  }
}
```

**Grant List Access**
```graphql
mutation AddListAccess($input: GrantListAccessInput!) {
  addListAccess(input: $input) {
    list {
      name
    }
    user {
      email
    }
    accessLevel
  }
}

# Variables
{
  "input": {
    "listId": "uuid",
    "userId": "uuid",
    "accessLevel": "WRITER"
  }
}
```

---

## Database

### Entity Relationship Diagram

```
┌─────────────────┐
│     USERS       │
├─────────────────┤
│ id (PK)         │───┐
│ email           │   │
│ github_id       │   │
│ role            │   │
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
                      │
                      │ owner_id
                      │
┌─────────────────┐   │
│     LISTS       │◄──┘
├─────────────────┤
│ id (PK)         │───┐
│ name            │   │
│ description     │   │
│ owner_id (FK)   │   │ list_id
│ visibility      │   │
│ tags            │   │
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
        │             │
        │             │
        │ list_id     │
        │             │
        ▼             ▼
┌─────────────────┐ ┌─────────────────┐
│     TODOS       │ │  LIST_ACCESS    │
├─────────────────┤ ├─────────────────┤
│ id (PK)         │ │ list_id (PK,FK) │
│ list_id (FK)    │ │ user_id (PK,FK) │
│ title           │ │ access_level    │
│ description     │ │ status          │
│ completed       │ └─────────────────┘
│ due_date        │
│ start_date      │
│ priority        │
│ tags            │
│ assigned_to(FK) │───┐
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
                      │
        ┌─────────────┘
        │ user_id
        ▼
┌─────────────────────┐
│  REFRESH_TOKENS     │
├─────────────────────┤
│ id (PK)             │
│ user_id (FK)        │
│ token               │
│ expires_at          │
│ created_at          │
└─────────────────────┘
```

### Indexes

```sql
-- Lists
CREATE INDEX idx_lists_owner_id ON lists(owner_id);
CREATE INDEX idx_lists_visibility ON lists(visibility);

-- Todos
CREATE INDEX idx_todos_list_id ON todos(list_id);
CREATE INDEX idx_todos_assigned_to ON todos(assigned_to);
CREATE INDEX idx_todos_completed ON todos(completed);

-- List Access
CREATE INDEX idx_list_access_user_id ON list_access(user_id);
CREATE INDEX idx_list_access_status ON list_access(status);

-- Refresh Tokens
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

### Triggers

**Automatic Updated_at Timestamp**
```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lists_updated_at 
    BEFORE UPDATE ON lists 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_todos_updated_at 
    BEFORE UPDATE ON todos 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

---

## Security and Authentication

### OAuth 2.0 Flow

```
┌─────────┐                                  ┌─────────┐
│         │                                  │         │
│  User   │                                  │ GitHub  │
│         │                                  │         │
└────┬────┘                                  └────┬────┘
     │                                            │
     │ 1. Click "Login with GitHub"              │
     │───────────────────────────────────────────►
     │                                            │
     │ 2. Redirect to GitHub OAuth               │
     │◄───────────────────────────────────────────
     │                                            │
     │ 3. GitHub Login & Authorization           │
     │───────────────────────────────────────────►
     │                                            │
     │ 4. Redirect with code                     │
     │◄───────────────────────────────────────────
     │                                            │
     │ 5. Exchange code for tokens               │
     │───────────────────────────────────────────►
     │                                            │
     │ 6. Return access_token & user info        │
     │◄───────────────────────────────────────────
     │                                            │
     │ 7. Store JWT & Refresh Token              │
     │                                            │
```

### JWT Token Structure

**Access Token** (valid 1 hour):
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "user_id": "uuid",
    "email": "user@example.com",
    "role": "writer",
    "github_id": "12345",
    "exp": 1704988800,
    "iat": 1704985200
  }
}
```

**Refresh Token** (valid 7 days):
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "user_id": "uuid",
    "token_id": "uuid",
    "exp": 1705590000,
    "iat": 1704985200
  }
}
```

### Role-Based Access Control (RBAC)

#### Roles and Permissions

**Reader**
- Can read lists they have access to
- Can read tasks in lists
- Cannot modify lists or tasks
- Cannot create new lists

**Writer**
- All Reader permissions
- Can create new lists
- Can modify and delete tasks in lists they have access to
- Can add users to lists they created
- Cannot modify lists they don't own

**Admin**
- All Writer permissions
- Can read and modify all lists
- Can manage all users
- Can delete any resources

#### List-Level Access Control

Each list can have additional collaborators:

- **List Owner** - full control
- **List Reader** - read only
- **List Writer** - read and write
- **List Admin** - full control like owner

### Security Best Practices

1. **Environment Variables** - Sensitive data is stored as environment variables
2. **HTTPS** - Production deployment uses TLS/SSL
3. **JWT Expiration** - Short-lived access tokens with refresh mechanism
4. **Password Hashing** - N/A (using OAuth)
5. **Input Validation** - Validation of all input data
6. **SQL Injection Prevention** - Using prepared statements
7. **CORS Configuration** - Properly configured CORS
8. **Rate Limiting** - Protection against brute force attacks (in production)

---

## Deployment

### Kubernetes Resources

#### Namespaces

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: todoapp-system
  labels:
    istio-injection: enabled
```

#### Deployments

**REST Service**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todoapp-rest
  namespace: todoapp-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: todoapp-rest
  template:
    metadata:
      labels:
        app: todoapp-rest
    spec:
      containers:
      - name: rest
        image: victoruzunov/todoapp-rest:project
        ports:
        - containerPort: 5000
        env:
        - name: DB_HOST
          value: "todoapp-postgres"
        - name: DB_PORT
          value: "5432"
        resources:
          limits:
            cpu: "1"
            memory: "256Mi"
          requests:
            cpu: "512m"
            memory: "128Mi"
```

**GraphQL Service**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todoapp-graphql
  namespace: todoapp-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: todoapp-graphql
  template:
    metadata:
      labels:
        app: todoapp-graphql
    spec:
      containers:
      - name: graphql
        image: victoruzunov/todoapp-graphql:project
        ports:
        - containerPort: 8000
```

**PostgreSQL**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todoapp-postgres
  namespace: todoapp-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: todoapp-postgres
  template:
    metadata:
      labels:
        app: todoapp-postgres
    spec:
      containers:
      - name: postgres
        image: postgres:14
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
```

#### Services

```yaml
apiVersion: v1
kind: Service
metadata:
  name: todoapp-rest
  namespace: todoapp-system
spec:
  selector:
    app: todoapp-rest
  ports:
  - port: 5000
    targetPort: 5000

---
apiVersion: v1
kind: Service
metadata:
  name: todoapp-graphql
  namespace: todoapp-system
spec:
  selector:
    app: todoapp-graphql
  ports:
  - port: 8000
    targetPort: 8000

---
apiVersion: v1
kind: Service
metadata:
  name: todoapp-postgres
  namespace: todoapp-system
spec:
  selector:
    app: todoapp-postgres
  ports:
  - port: 5432
    targetPort: 5432
```

### Istio Configuration

**Gateway**
```yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: todoapp-gateway
  namespace: todoapp-system
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
```

**VirtualService**
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: todoapp-routes
  namespace: todoapp-system
spec:
  hosts:
  - "*"
  gateways:
  - todoapp-gateway
  http:
  - match:
    - uri:
        prefix: "/graphql"
    route:
    - destination:
        host: todoapp-graphql
        port:
          number: 8000
  - match:
    - uri:
        prefix: "/api"
    route:
    - destination:
        host: todoapp-rest
        port:
          number: 5000
  - route:
    - destination:
        host: todoapp-ui
        port:
          number: 80
```

**mTLS Policy**
```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: todoapp-system
spec:
  mtls:
    mode: PERMISSIVE
```

### Helm Chart Structure

```
charts/todoapp/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── ingress-gateway.yaml
│   ├── mTLS.yaml
│   └── _helpers.tpl
└── charts/
    ├── rest/
    │   ├── Chart.yaml
    │   ├── values.yaml
    │   └── templates/
    │       ├── deployment.yaml
    │       ├── service.yaml
    │       ├── configmap.yaml
    │       └── secret.yaml
    ├── graphql/
    │   ├── Chart.yaml
    │   ├── values.yaml
    │   └── templates/
    │       ├── deployment.yaml
    │       └── service.yaml
    ├── postgres/
    │   ├── Chart.yaml
    │   ├── values.yaml
    │   └── templates/
    │       ├── deployment.yaml
    │       ├── service.yaml
    │       ├── pvc.yaml
    │       └── migration-job.yaml
    └── ui/
        ├── Chart.yaml
        ├── values.yaml
        └── templates/
            ├── deployment.yaml
            └── service.yaml
```

---

## CI/CD Pipeline

### GitHub Actions Workflows

The project uses GitHub Actions for CI/CD with the following jobs:

#### 1. Markdown Lint
```yaml
name: Markdown Check
on: [push]
jobs:
  markdown-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install markdownlint
        run: npm install -g markdownlint-cli
      - name: Lint README
        run: markdownlint README.md
```

#### 2. Lint and Format Check
- Checks code formatting with `gofmt`
- Checks code quality with `golangci-lint`
- Separate jobs for REST and GraphQL services

#### 3. Build
- Compiles Go code
- Checks for compilation errors
- Caches dependencies

#### 4. SonarCloud Analysis
- Static code analysis
- Code quality metrics
- Security vulnerability detection
- Technical debt tracking

#### 5. Snyk Security Scan
- Dependency vulnerability scanning
- License compliance check
- Automatic PR creation for fixes

#### 6. Docker Build and Push
- Builds Docker images
- Tags with commit SHA and branch name
- Pushes to Docker Hub registry

#### 7. Deploy to Kubernetes
- Automatic deployment on successful build
- Rolling update strategy
- Health checks

### Pipeline Flow

```
┌──────────────┐
│  Git Push    │
└──────┬───────┘
       │
       ▼
┌──────────────────────┐
│  Markdown Lint       │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│  Lint & Format       │
│  (REST + GraphQL)    │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│  Build               │
│  (REST + GraphQL)    │
└──────┬───────────────┘
       │
       ├─────────────────┐
       │                 │
       ▼                 ▼
┌──────────────┐  ┌──────────────┐
│  SonarCloud  │  │    Snyk      │
│   Analysis   │  │   Security   │
└──────┬───────┘  └──────┬───────┘
       │                 │
       └────────┬────────┘
                │
                ▼
       ┌──────────────────┐
       │  Docker Build    │
       │  & Push          │
       └────────┬─────────┘
                │
                ▼
       ┌──────────────────┐
       │  Deploy to K8s   │
       └──────────────────┘
```

---

## Usage

### Starting the Application

**Method 1: Automatic Setup**
```bash
./setup.sh
```

**Method 2: Manual Setup**

1. Create cluster:
```bash
k3d cluster create todoapp-cluster \
  --port 8000:80@loadbalancer \
  --wait
```

2. Install Istio:
```bash
istioctl install --set profile=demo -y
```

3. Create namespace:
```bash
kubectl create namespace todoapp-system
kubectl label namespace todoapp-system istio-injection=enabled
```

4. Deploy the application:
```bash
helm install todoapp ./charts/todoapp \
  --namespace todoapp-system \
  --set-file config=config.yaml
```

### Using the UI

1. **Login**
    - Open `http://localhost:8000`
    - Click "Login with GitHub"
    - Authorize the application

2. **Create a List**
    - Click the "+" button
    - Enter name and description
    - Select visibility
    - Add tags (optional)

3. **Add Todos**
    - Click on a list
    - Click "Add Todo"
    - Fill in the details
    - Set priority and due date

4. **Share a List**
    - Open a list
    - Click "Share"
    - Enter user email
    - Select access level

5. **Managing Tasks**
    - Check a task to mark it as completed
    - Edit to change details
    - Assign to another user
    - Delete to remove

### Using GraphQL Playground

1. Open `http://localhost:8000/graphql`
2. View the schema in the documentation
3. Write queries and mutations
4. Test different scenarios

Example:
```graphql
query MyLists {
  lists {
    id
    name
    todos {
      title
      completed
    }
  }
}
```

### Using the REST API

**With curl:**

```bash
# Login
curl http://localhost:8000/login/github

# Get Lists (after login)
curl -H "Authorization: Bearer <token>" \
     http://localhost:5000/lists

# Create List
curl -X POST \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"name":"My List","visibility":"private"}' \
     http://localhost:5000/lists

# Create Todo
curl -X POST \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"list_id":"<uuid>","title":"Task 1","priority":"high"}' \
     http://localhost:5000/todos
```

**With Postman:**
1. Import collection
2. Set environment variables (token, base_url)
3. Execute requests

---

## Troubleshooting

### Common Issues

#### 1. Pods Not Starting

**Problem**: Pods stuck in Pending status

**Solution**:
```bash
# Check pod status
kubectl describe pod <pod-name> -n todoapp-system

# Check resources
kubectl top nodes

# Check events
kubectl get events -n todoapp-system --sort-by='.lastTimestamp'
```

#### 2. Database Connection Failed

**Problem**: REST API cannot connect to PostgreSQL

**Solution**:
```bash
# Check Postgres pod
kubectl logs -n todoapp-system -l app=todoapp-postgres

# Check service
kubectl get svc -n todoapp-system todoapp-postgres

# Test connection
kubectl run -it --rm --restart=Never postgres-test \
  --image=postgres:14 \
  --namespace=todoapp-system \
  -- psql -h todoapp-postgres -U postgres
```

#### 3. OAuth Redirect Issues

**Problem**: GitHub OAuth redirect fails

**Solution**:
1. Check GitHub OAuth app settings
2. Verify redirect URL matches exactly
3. Check config.yaml settings
4. Restart REST pod

```bash
kubectl rollout restart deployment/todoapp-rest -n todoapp-system
```

#### 4. JWT Token Expired

**Problem**: API returns 401 Unauthorized

**Solution**:
- Use refresh token endpoint
- Login again
- Check token expiration time

```bash
# Refresh token
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<token>"}' \
  http://localhost:5000/auth/refresh
```

#### 5. Istio Gateway Not Working

**Problem**: Cannot access application on port 8000

**Solution**:
```bash
# Check Istio installation
kubectl get pods -n istio-system

# Check gateway
kubectl get gateway -n todoapp-system

# Check virtual service
kubectl get virtualservice -n todoapp-system

# Port forward as fallback
kubectl port-forward -n todoapp-system \
  svc/todoapp-ui 8000:80
```

### Debugging Commands

**View Logs**
```bash
# REST API logs
kubectl logs -n todoapp-system -l app=todoapp-rest -f

# GraphQL logs
kubectl logs -n todoapp-system -l app=todoapp-graphql -f

# UI logs
kubectl logs -n todoapp-system -l app=todoapp-ui -f

# Postgres logs
kubectl logs -n todoapp-system -l app=todoapp-postgres -f
```

**Execute Commands in Pods**
```bash
# Access REST API container
kubectl exec -it -n todoapp-system \
  deployment/todoapp-rest -- /bin/sh

# Access Database
kubectl exec -it -n todoapp-system \
  deployment/todoapp-postgres -- psql -U postgres
```

**Check Resources**
```bash
# CPU and Memory usage
kubectl top pods -n todoapp-system

# Describe deployment
kubectl describe deployment todoapp-rest -n todoapp-system

# Get all resources
kubectl get all -n todoapp-system
```

### Clean Up

**Remove Application**
```bash
# Uninstall Helm release
helm uninstall todoapp -n todoapp-system

# Delete namespace
kubectl delete namespace todoapp-system
```

**Delete Cluster**
```bash
# Delete k3d cluster
k3d cluster delete todoapp-cluster
```

**Complete Cleanup**
```bash
# Stop and remove everything
k3d cluster delete todoapp-cluster
docker system prune -a
```

---

## Additional Resources

### Documentation

- [Go Documentation](https://golang.org/doc/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Istio Documentation](https://istio.io/latest/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [GraphQL Documentation](https://graphql.org/learn/)
- [SAP UI5 Documentation](https://ui5.sap.com/)

### Related Files

- `SETUP_GUIDE.md` - Detailed setup instructions
- `PIPELINE.md` - CI/CD pipeline documentation
- `config.example.yaml` - Configuration template
- `README.md` - Project overview

### Contact and Support

For questions and issues:
- Open an Issue in the GitHub repository
- Check existing issues for solutions
- Consult the documentation

---

## Conclusion

The TodoList Project is a full-featured microservices application that demonstrates modern best practices in:

- **Software Architecture** - Microservices architecture with clear separation
- **Security** - OAuth 2.0, JWT, RBAC
- **DevOps** - Containerization, Kubernetes, CI/CD
- **API Design** - RESTful and GraphQL APIs
- **Database Design** - Normalization, indexes, migrations
- **Frontend Development** - Modern UI framework

The project is ready for production deployment with minor modifications for the specific production environment (SSL certificates, external database, monitoring, logging, etc.).

---

**Version**: 1.0  
**Last Updated**: January 2026  
**Author**: TodoList Project Team
