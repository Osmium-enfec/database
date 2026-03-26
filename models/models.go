package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// User represents a system user
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // creator, reviewer, admin
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Program represents an educational program
type Program struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Topic represents a topic within a program
type Topic struct {
	ID          string    `json:"id"`
	ProgramID   string    `json:"program_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Subtopic represents a subtopic within a topic
type Subtopic struct {
	ID          string    `json:"id"`
	TopicID     string    `json:"topic_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag represents a tag for content
type Tag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Content represents a piece of educational content
type Content struct {
	ID                   string    `json:"id"`
	Type                 string    `json:"type"` // question, code_problem, documentation
	ProgramID            string    `json:"program_id"`
	TopicID              string    `json:"topic_id"`
	SubtopicID           string    `json:"subtopic_id"`
	Difficulty           string    `json:"difficulty"` // easy, medium, hard
	EstimatedTimeMinutes int       `json:"estimated_time_minutes"`
	Status               string    `json:"status"` // draft, pending_review, approved, rejected
	CreatedBy            string    `json:"created_by"`
	CurrentVersionID     *string   `json:"current_version_id"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Relationships (populated when needed)
	Tags    []Tag           `json:"tags,omitempty"`
	Creator *User           `json:"creator,omitempty"`
	Version *ContentVersion `json:"version,omitempty"`
}

// ContentData represents the versioned content structure
type ContentData struct {
	Type              string             `json:"type"`
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	QuestionType      string             `json:"question_type,omitempty"` // mcq, msq, fill, short
	QuestionText      string             `json:"question_text,omitempty"`
	Hints             []string           `json:"hints,omitempty"` // At least 3 hints for questions
	Options           []string           `json:"options,omitempty"`
	CorrectOptions    []int              `json:"correct_options,omitempty"`
	CodeProblemData   *CodeProblemData   `json:"code_problem_data,omitempty"`
	DocumentationData *DocumentationData `json:"documentation_data,omitempty"`
}

// CodeProblemData represents code problem specific data
type CodeProblemData struct {
	StarterCode       string         `json:"starter_code"`
	SolutionCode      string         `json:"solution_code"`
	ExecutionTemplate string         `json:"execution_template"`
	TestCases         []CodeTestCase `json:"test_cases"`
}

// CodeTestCase represents a test case for code problems
type CodeTestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsHidden       bool   `json:"is_hidden"`
}

// DocumentationData represents documentation specific data
type DocumentationData struct {
	MarkdownContent string `json:"markdown_content"`
}

// ContentVersion represents a specific version of content
type ContentVersion struct {
	ID            string      `json:"id"`
	ContentID     string      `json:"content_id"`
	VersionNumber int         `json:"version_number"`
	Data          ContentData `json:"data"`
	CreatedBy     string      `json:"created_by"`
	ReviewStatus  string      `json:"review_status"` // pending, approved, rejected
	ReviewComment string      `json:"review_comment"`
	ReviewedBy    *string     `json:"reviewed_by"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`

	// Relationships
	Creator  *User `json:"creator,omitempty"`
	Reviewer *User `json:"reviewer,omitempty"`
}

// Scan implements sql.Scanner interface for ContentData
func (cd *ContentData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &cd)
}

// Value implements driver.Valuer interface for ContentData
func (cd ContentData) Value() (driver.Value, error) {
	return json.Marshal(cd)
}

// RefreshToken represents a refresh token for JWT auth
type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditLog represents an action audit log
type AuditLog struct {
	ID           string                 `json:"id"`
	UserID       *string                `json:"user_id"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Changes      map[string]interface{} `json:"changes"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	CreatedAt    time.Time              `json:"created_at"`
}
