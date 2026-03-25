package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"project/models"
)

// ============================================================================
// Program Repository
// ============================================================================

// ProgramRepository defines program data access operations
type ProgramRepository interface {
	Create(ctx context.Context, program *models.Program) error
	GetByID(ctx context.Context, id string) (*models.Program, error)
	Update(ctx context.Context, program *models.Program) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, perPage int) ([]*models.Program, int, error)
}

type programRepository struct {
	db *sql.DB
}

// NewProgramRepository creates a new program repository
func NewProgramRepository(db *sql.DB) ProgramRepository {
	return &programRepository{db: db}
}

func (r *programRepository) Create(ctx context.Context, program *models.Program) error {
	query := `
		INSERT INTO programs (id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		program.ID, program.Name, program.Description, program.IsActive, program.CreatedAt, program.UpdatedAt)
	return err
}

func (r *programRepository) GetByID(ctx context.Context, id string) (*models.Program, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM programs WHERE id = $1 AND is_active = true
	`
	program := &models.Program{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&program.ID, &program.Name, &program.Description, &program.IsActive, &program.CreatedAt, &program.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("program not found")
		}
		return nil, err
	}
	return program, nil
}

func (r *programRepository) Update(ctx context.Context, program *models.Program) error {
	query := `
		UPDATE programs SET name = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`
	program.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		program.Name, program.Description, program.IsActive, program.UpdatedAt, program.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("program not found")
	}
	return err
}

func (r *programRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE programs SET is_active = false, updated_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("program not found")
	}
	return err
}

func (r *programRepository) List(ctx context.Context, page, perPage int) ([]*models.Program, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM programs WHERE is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM programs WHERE is_active = true
		ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var programs []*models.Program
	for rows.Next() {
		program := &models.Program{}
		if err := rows.Scan(&program.ID, &program.Name, &program.Description, &program.IsActive, &program.CreatedAt, &program.UpdatedAt); err != nil {
			return nil, 0, err
		}
		programs = append(programs, program)
	}

	return programs, total, rows.Err()
}

// ============================================================================
// Topic Repository
// ============================================================================

// TopicRepository defines topic data access operations
type TopicRepository interface {
	Create(ctx context.Context, topic *models.Topic) error
	GetByID(ctx context.Context, id string) (*models.Topic, error)
	Update(ctx context.Context, topic *models.Topic) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, programID string, page, perPage int) ([]*models.Topic, int, error)
}

type topicRepository struct {
	db *sql.DB
}

// NewTopicRepository creates a new topic repository
func NewTopicRepository(db *sql.DB) TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) Create(ctx context.Context, topic *models.Topic) error {
	query := `
		INSERT INTO topics (id, program_id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		topic.ID, topic.ProgramID, topic.Name, topic.Description, topic.IsActive, topic.CreatedAt, topic.UpdatedAt)
	return err
}

func (r *topicRepository) GetByID(ctx context.Context, id string) (*models.Topic, error) {
	query := `
		SELECT id, program_id, name, description, is_active, created_at, updated_at
		FROM topics WHERE id = $1 AND is_active = true
	`
	topic := &models.Topic{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&topic.ID, &topic.ProgramID, &topic.Name, &topic.Description, &topic.IsActive, &topic.CreatedAt, &topic.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, err
	}
	return topic, nil
}

func (r *topicRepository) Update(ctx context.Context, topic *models.Topic) error {
	query := `
		UPDATE topics SET program_id = $1, name = $2, description = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	topic.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		topic.ProgramID, topic.Name, topic.Description, topic.IsActive, topic.UpdatedAt, topic.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("topic not found")
	}
	return err
}

func (r *topicRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE topics SET is_active = false, updated_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("topic not found")
	}
	return err
}

func (r *topicRepository) List(ctx context.Context, programID string, page, perPage int) ([]*models.Topic, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM topics WHERE program_id = $1 AND is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, programID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, program_id, name, description, is_active, created_at, updated_at
		FROM topics WHERE program_id = $1 AND is_active = true
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, programID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var topics []*models.Topic
	for rows.Next() {
		topic := &models.Topic{}
		if err := rows.Scan(&topic.ID, &topic.ProgramID, &topic.Name, &topic.Description, &topic.IsActive, &topic.CreatedAt, &topic.UpdatedAt); err != nil {
			return nil, 0, err
		}
		topics = append(topics, topic)
	}

	return topics, total, rows.Err()
}

// ============================================================================
// Subtopic Repository
// ============================================================================

// SubtopicRepository defines subtopic data access operations
type SubtopicRepository interface {
	Create(ctx context.Context, subtopic *models.Subtopic) error
	GetByID(ctx context.Context, id string) (*models.Subtopic, error)
	Update(ctx context.Context, subtopic *models.Subtopic) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, topicID string, page, perPage int) ([]*models.Subtopic, int, error)
}

type subtopicRepository struct {
	db *sql.DB
}

// NewSubtopicRepository creates a new subtopic repository
func NewSubtopicRepository(db *sql.DB) SubtopicRepository {
	return &subtopicRepository{db: db}
}

func (r *subtopicRepository) Create(ctx context.Context, subtopic *models.Subtopic) error {
	query := `
		INSERT INTO subtopics (id, topic_id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		subtopic.ID, subtopic.TopicID, subtopic.Name, subtopic.Description, subtopic.IsActive, subtopic.CreatedAt, subtopic.UpdatedAt)
	return err
}

func (r *subtopicRepository) GetByID(ctx context.Context, id string) (*models.Subtopic, error) {
	query := `
		SELECT id, topic_id, name, description, is_active, created_at, updated_at
		FROM subtopics WHERE id = $1 AND is_active = true
	`
	subtopic := &models.Subtopic{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subtopic.ID, &subtopic.TopicID, &subtopic.Name, &subtopic.Description, &subtopic.IsActive, &subtopic.CreatedAt, &subtopic.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subtopic not found")
		}
		return nil, err
	}
	return subtopic, nil
}

func (r *subtopicRepository) Update(ctx context.Context, subtopic *models.Subtopic) error {
	query := `
		UPDATE subtopics SET topic_id = $1, name = $2, description = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	subtopic.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		subtopic.TopicID, subtopic.Name, subtopic.Description, subtopic.IsActive, subtopic.UpdatedAt, subtopic.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("subtopic not found")
	}
	return err
}

func (r *subtopicRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE subtopics SET is_active = false, updated_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("subtopic not found")
	}
	return err
}

func (r *subtopicRepository) List(ctx context.Context, topicID string, page, perPage int) ([]*models.Subtopic, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM subtopics WHERE topic_id = $1 AND is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, topicID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, topic_id, name, description, is_active, created_at, updated_at
		FROM subtopics WHERE topic_id = $1 AND is_active = true
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, topicID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subtopics []*models.Subtopic
	for rows.Next() {
		subtopic := &models.Subtopic{}
		if err := rows.Scan(&subtopic.ID, &subtopic.TopicID, &subtopic.Name, &subtopic.Description, &subtopic.IsActive, &subtopic.CreatedAt, &subtopic.UpdatedAt); err != nil {
			return nil, 0, err
		}
		subtopics = append(subtopics, subtopic)
	}

	return subtopics, total, rows.Err()
}

// ============================================================================
// Audit Log Repository
// ============================================================================

// AuditLogRepository defines audit log data access operations
type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	GetByID(ctx context.Context, id string) (*models.AuditLog, error)
	List(ctx context.Context, page, perPage int, filters map[string]string) ([]*models.AuditLog, int, error)
	GetByResource(ctx context.Context, resourceType, resourceID string, page, perPage int) ([]*models.AuditLog, int, error)
}

type auditLogRepository struct {
	db *sql.DB
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	changesJSON, err := json.Marshal(log.Changes)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO audit_logs (id, user_id, action, resource_type, resource_id, changes, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = r.db.ExecContext(ctx, query,
		log.ID, log.UserID, log.Action, log.ResourceType, log.ResourceID,
		changesJSON, log.IPAddress, log.UserAgent, log.CreatedAt)
	return err
}

func (r *auditLogRepository) GetByID(ctx context.Context, id string) (*models.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource_type, resource_id, changes, ip_address, user_agent, created_at
		FROM audit_logs WHERE id = $1
	`
	log := &models.AuditLog{}
	var changesJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
		&changesJSON, &log.IPAddress, &log.UserAgent, &log.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("audit log not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(changesJSON, &log.Changes); err != nil {
		return nil, err
	}

	return log, nil
}

func (r *auditLogRepository) List(ctx context.Context, page, perPage int, filters map[string]string) ([]*models.AuditLog, int, error) {
	baseQuery := `SELECT id, user_id, action, resource_type, resource_id, changes, ip_address, user_agent, created_at
		FROM audit_logs WHERE 1=1`

	// Build filter conditions
	args := []interface{}{}
	argCount := 1

	if filters["action"] != "" {
		baseQuery += fmt.Sprintf(" AND action = $%d", argCount)
		args = append(args, filters["action"])
		argCount++
	}
	if filters["resource_type"] != "" {
		baseQuery += fmt.Sprintf(" AND resource_type = $%d", argCount)
		args = append(args, filters["resource_type"])
		argCount++
	}
	if filters["user_id"] != "" {
		baseQuery += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, filters["user_id"])
		argCount++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") AS cnt"
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := baseQuery + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, perPage, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		var changesJSON []byte
		if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&changesJSON, &log.IPAddress, &log.UserAgent, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal(changesJSON, &log.Changes); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, rows.Err()
}

func (r *auditLogRepository) GetByResource(ctx context.Context, resourceType, resourceID string, page, perPage int) ([]*models.AuditLog, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE resource_type = $1 AND resource_id = $2`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, resourceType, resourceID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, user_id, action, resource_type, resource_id, changes, ip_address, user_agent, created_at
		FROM audit_logs WHERE resource_type = $1 AND resource_id = $2
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, query, resourceType, resourceID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		var changesJSON []byte
		if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&changesJSON, &log.IPAddress, &log.UserAgent, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal(changesJSON, &log.Changes); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, rows.Err()
}
