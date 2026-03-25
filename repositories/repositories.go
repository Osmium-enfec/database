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
// User Repository
// ============================================================================

// UserRepository defines user data access operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, perPage int) ([]*models.User, int, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = true
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users WHERE email = $1 AND is_active = true
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET name = $1, email = $2, password_hash = $3, role = $4, is_active = $5, updated_at = $6
		WHERE id = $7
	`
	user.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return err
}

func (r *userRepository) List(ctx context.Context, page, perPage int) ([]*models.User, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM users WHERE is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users WHERE is_active = true
		ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, rows.Err()
}

// ============================================================================
// Content Repository
// ============================================================================

// ContentRepository defines content data access operations
type ContentRepository interface {
	Create(ctx context.Context, content *models.Content) error
	GetByID(ctx context.Context, id string) (*models.Content, error)
	Update(ctx context.Context, content *models.Content) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, perPage int, filters map[string]string) ([]*models.Content, int, error)
	GetByCreator(ctx context.Context, creatorID string, page, perPage int) ([]*models.Content, int, error)
}

type contentRepository struct {
	db *sql.DB
}

// NewContentRepository creates a new content repository
func NewContentRepository(db *sql.DB) ContentRepository {
	return &contentRepository{db: db}
}

func (r *contentRepository) Create(ctx context.Context, content *models.Content) error {
	query := `
		INSERT INTO contents (id, type, program_id, topic_id, subtopic_id, difficulty, estimated_time_minutes, 
			status, created_by, current_version_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		content.ID, content.Type, content.ProgramID, content.TopicID, content.SubtopicID,
		content.Difficulty, content.EstimatedTimeMinutes, content.Status, content.CreatedBy,
		content.CurrentVersionID, content.IsActive, content.CreatedAt, content.UpdatedAt)
	return err
}

func (r *contentRepository) GetByID(ctx context.Context, id string) (*models.Content, error) {
	query := `
		SELECT id, type, program_id, topic_id, subtopic_id, difficulty, estimated_time_minutes,
			status, created_by, current_version_id, is_active, created_at, updated_at
		FROM contents WHERE id = $1 AND is_active = true
	`
	content := &models.Content{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&content.ID, &content.Type, &content.ProgramID, &content.TopicID, &content.SubtopicID,
		&content.Difficulty, &content.EstimatedTimeMinutes, &content.Status, &content.CreatedBy,
		&content.CurrentVersionID, &content.IsActive, &content.CreatedAt, &content.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content not found")
		}
		return nil, err
	}
	return content, nil
}

func (r *contentRepository) Update(ctx context.Context, content *models.Content) error {
	query := `
		UPDATE contents SET type = $1, program_id = $2, topic_id = $3, subtopic_id = $4,
			difficulty = $5, estimated_time_minutes = $6, status = $7, current_version_id = $8,
			is_active = $9, updated_at = $10
		WHERE id = $11
	`
	content.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		content.Type, content.ProgramID, content.TopicID, content.SubtopicID,
		content.Difficulty, content.EstimatedTimeMinutes, content.Status, content.CurrentVersionID,
		content.IsActive, content.UpdatedAt, content.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("content not found")
	}
	return err
}

func (r *contentRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE contents SET is_active = false, updated_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("content not found")
	}
	return err
}

func (r *contentRepository) List(ctx context.Context, page, perPage int, filters map[string]string) ([]*models.Content, int, error) {
	baseQuery := `SELECT id, type, program_id, topic_id, subtopic_id, difficulty, estimated_time_minutes,
		status, created_by, current_version_id, is_active, created_at, updated_at
		FROM contents WHERE is_active = true`

	// Build filter conditions
	args := []interface{}{}
	argCount := 1

	if filters["type"] != "" {
		baseQuery += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, filters["type"])
		argCount++
	}
	if filters["program_id"] != "" {
		baseQuery += fmt.Sprintf(" AND program_id = $%d", argCount)
		args = append(args, filters["program_id"])
		argCount++
	}
	if filters["topic_id"] != "" {
		baseQuery += fmt.Sprintf(" AND topic_id = $%d", argCount)
		args = append(args, filters["topic_id"])
		argCount++
	}
	if filters["difficulty"] != "" {
		baseQuery += fmt.Sprintf(" AND difficulty = $%d", argCount)
		args = append(args, filters["difficulty"])
		argCount++
	}
	if filters["status"] != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters["status"])
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

	var contents []*models.Content
	for rows.Next() {
		content := &models.Content{}
		if err := rows.Scan(&content.ID, &content.Type, &content.ProgramID, &content.TopicID,
			&content.SubtopicID, &content.Difficulty, &content.EstimatedTimeMinutes, &content.Status,
			&content.CreatedBy, &content.CurrentVersionID, &content.IsActive, &content.CreatedAt, &content.UpdatedAt); err != nil {
			return nil, 0, err
		}
		contents = append(contents, content)
	}

	return contents, total, rows.Err()
}

func (r *contentRepository) GetByCreator(ctx context.Context, creatorID string, page, perPage int) ([]*models.Content, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM contents WHERE created_by = $1 AND is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, creatorID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, type, program_id, topic_id, subtopic_id, difficulty, estimated_time_minutes,
			status, created_by, current_version_id, is_active, created_at, updated_at
		FROM contents WHERE created_by = $1 AND is_active = true
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, creatorID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contents []*models.Content
	for rows.Next() {
		content := &models.Content{}
		if err := rows.Scan(&content.ID, &content.Type, &content.ProgramID, &content.TopicID,
			&content.SubtopicID, &content.Difficulty, &content.EstimatedTimeMinutes, &content.Status,
			&content.CreatedBy, &content.CurrentVersionID, &content.IsActive, &content.CreatedAt, &content.UpdatedAt); err != nil {
			return nil, 0, err
		}
		contents = append(contents, content)
	}

	return contents, total, rows.Err()
}

// ============================================================================
// Content Version Repository
// ============================================================================

// ContentVersionRepository defines content version data access operations
type ContentVersionRepository interface {
	Create(ctx context.Context, version *models.ContentVersion) error
	GetByID(ctx context.Context, id string) (*models.ContentVersion, error)
	GetLatestByContentID(ctx context.Context, contentID string) (*models.ContentVersion, error)
	GetByContentID(ctx context.Context, contentID string, page, perPage int) ([]*models.ContentVersion, int, error)
	GetPending(ctx context.Context, page, perPage int) ([]*models.ContentVersion, int, error)
	Update(ctx context.Context, version *models.ContentVersion) error
}

type contentVersionRepository struct {
	db *sql.DB
}

// NewContentVersionRepository creates a new content version repository
func NewContentVersionRepository(db *sql.DB) ContentVersionRepository {
	return &contentVersionRepository{db: db}
}

func (r *contentVersionRepository) Create(ctx context.Context, version *models.ContentVersion) error {
	dataJSON, err := json.Marshal(version.Data)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO content_versions (id, content_id, version_number, data, created_by, review_status, review_comment, reviewed_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = r.db.ExecContext(ctx, query,
		version.ID, version.ContentID, version.VersionNumber, dataJSON, version.CreatedBy,
		version.ReviewStatus, version.ReviewComment, version.ReviewedBy, version.CreatedAt, version.UpdatedAt)
	return err
}

func (r *contentVersionRepository) GetByID(ctx context.Context, id string) (*models.ContentVersion, error) {
	query := `
		SELECT id, content_id, version_number, data, created_by, review_status, review_comment, reviewed_by, created_at, updated_at
		FROM content_versions WHERE id = $1
	`
	version := &models.ContentVersion{}
	var dataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&version.ID, &version.ContentID, &version.VersionNumber, &dataJSON,
		&version.CreatedBy, &version.ReviewStatus, &version.ReviewComment, &version.ReviewedBy,
		&version.CreatedAt, &version.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("version not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(dataJSON, &version.Data); err != nil {
		return nil, err
	}

	return version, nil
}

func (r *contentVersionRepository) GetLatestByContentID(ctx context.Context, contentID string) (*models.ContentVersion, error) {
	query := `
		SELECT id, content_id, version_number, data, created_by, review_status, review_comment, reviewed_by, created_at, updated_at
		FROM content_versions WHERE content_id = $1
		ORDER BY version_number DESC LIMIT 1
	`
	version := &models.ContentVersion{}
	var dataJSON []byte
	err := r.db.QueryRowContext(ctx, query, contentID).Scan(
		&version.ID, &version.ContentID, &version.VersionNumber, &dataJSON,
		&version.CreatedBy, &version.ReviewStatus, &version.ReviewComment, &version.ReviewedBy,
		&version.CreatedAt, &version.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("version not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(dataJSON, &version.Data); err != nil {
		return nil, err
	}

	return version, nil
}

func (r *contentVersionRepository) GetByContentID(ctx context.Context, contentID string, page, perPage int) ([]*models.ContentVersion, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM content_versions WHERE content_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, contentID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, content_id, version_number, data, created_by, review_status, review_comment, reviewed_by, created_at, updated_at
		FROM content_versions WHERE content_id = $1
		ORDER BY version_number DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, contentID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var versions []*models.ContentVersion
	for rows.Next() {
		version := &models.ContentVersion{}
		var dataJSON []byte
		if err := rows.Scan(&version.ID, &version.ContentID, &version.VersionNumber, &dataJSON,
			&version.CreatedBy, &version.ReviewStatus, &version.ReviewComment, &version.ReviewedBy,
			&version.CreatedAt, &version.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal(dataJSON, &version.Data); err != nil {
			return nil, 0, err
		}
		versions = append(versions, version)
	}

	return versions, total, rows.Err()
}

func (r *contentVersionRepository) GetPending(ctx context.Context, page, perPage int) ([]*models.ContentVersion, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM content_versions WHERE review_status = 'pending'`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	query := `
		SELECT id, content_id, version_number, data, created_by, review_status, review_comment, reviewed_by, created_at, updated_at
		FROM content_versions WHERE review_status = 'pending'
		ORDER BY created_at ASC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var versions []*models.ContentVersion
	for rows.Next() {
		version := &models.ContentVersion{}
		var dataJSON []byte
		if err := rows.Scan(&version.ID, &version.ContentID, &version.VersionNumber, &dataJSON,
			&version.CreatedBy, &version.ReviewStatus, &version.ReviewComment, &version.ReviewedBy,
			&version.CreatedAt, &version.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal(dataJSON, &version.Data); err != nil {
			return nil, 0, err
		}
		versions = append(versions, version)
	}

	return versions, total, rows.Err()
}

func (r *contentVersionRepository) Update(ctx context.Context, version *models.ContentVersion) error {
	dataJSON, err := json.Marshal(version.Data)
	if err != nil {
		return err
	}

	query := `
		UPDATE content_versions SET version_number = $1, data = $2, created_by = $3,
			review_status = $4, review_comment = $5, reviewed_by = $6, updated_at = $7
		WHERE id = $8
	`
	version.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		version.VersionNumber, dataJSON, version.CreatedBy,
		version.ReviewStatus, version.ReviewComment, version.ReviewedBy, version.UpdatedAt, version.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("version not found")
	}
	return err
}

// ============================================================================
// Tag Repository
// ============================================================================

// TagRepository defines tag data access operations
type TagRepository interface {
	Create(ctx context.Context, tag *models.Tag) error
	GetByID(ctx context.Context, id string) (*models.Tag, error)
	GetByName(ctx context.Context, name string) (*models.Tag, error)
	Update(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*models.Tag, error)
	GetTagsByContentID(ctx context.Context, contentID string) ([]string, error)
	AssignTag(ctx context.Context, contentID, tagID string) error
	RemoveTag(ctx context.Context, contentID, tagID string) error
}

type tagRepository struct {
	db *sql.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *models.Tag) error {
	query := `
		INSERT INTO tags (id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		tag.ID, tag.Name, tag.Description, tag.IsActive, tag.CreatedAt, tag.UpdatedAt)
	return err
}

func (r *tagRepository) GetByID(ctx context.Context, id string) (*models.Tag, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM tags WHERE id = $1 AND is_active = true
	`
	tag := &models.Tag{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tag.ID, &tag.Name, &tag.Description, &tag.IsActive, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, err
	}
	return tag, nil
}

func (r *tagRepository) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM tags WHERE name = $1 AND is_active = true
	`
	tag := &models.Tag{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&tag.ID, &tag.Name, &tag.Description, &tag.IsActive, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, err
	}
	return tag, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *models.Tag) error {
	query := `
		UPDATE tags SET name = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`
	tag.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		tag.Name, tag.Description, tag.IsActive, tag.UpdatedAt, tag.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tag not found")
	}
	return err
}

func (r *tagRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tags SET is_active = false, updated_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tag not found")
	}
	return err
}

func (r *tagRepository) List(ctx context.Context) ([]*models.Tag, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM tags WHERE is_active = true
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.IsActive, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func (r *tagRepository) GetTagsByContentID(ctx context.Context, contentID string) ([]string, error) {
	query := `
		SELECT t.name FROM tags t
		INNER JOIN content_tags ct ON t.id = ct.tag_id
		WHERE ct.content_id = $1 AND t.is_active = true
		ORDER BY t.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tags = append(tags, name)
	}

	return tags, rows.Err()
}

func (r *tagRepository) AssignTag(ctx context.Context, contentID, tagID string) error {
	query := `
		INSERT INTO content_tags (id, content_id, tag_id, created_at)
		VALUES (gen_random_uuid(), $1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (content_id, tag_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, contentID, tagID)
	return err
}

func (r *tagRepository) RemoveTag(ctx context.Context, contentID, tagID string) error {
	query := `DELETE FROM content_tags WHERE content_id = $1 AND tag_id = $2`
	_, err := r.db.ExecContext(ctx, query, contentID, tagID)
	return err
}
