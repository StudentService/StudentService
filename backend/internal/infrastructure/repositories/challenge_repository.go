package repositories

import (
	"context"
	"log"

	"backend/internal/domain/challenge"
	"backend/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ChallengeRepository struct{}

// GetByUserID получает все вызовы пользователя
func (r *ChallengeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*challenge.Challenge, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, user_id, title, description, goal, 
		       start_date, end_date, status, progress,
		       created_at, updated_at
		FROM challenges
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		log.Printf("Database error in GetByUserID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var challenges []*challenge.Challenge
	for rows.Next() {
		c := &challenge.Challenge{}
		err := rows.Scan(
			&c.ID, &c.UserID, &c.Title, &c.Description, &c.Goal,
			&c.StartDate, &c.EndDate, &c.Status, &c.Progress,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning challenge row: %v", err)
			continue
		}
		challenges = append(challenges, c)
	}

	return challenges, nil
}

// GetByID получает вызов по ID
func (r *ChallengeRepository) GetByID(ctx context.Context, id uuid.UUID) (*challenge.Challenge, error) {
	c := &challenge.Challenge{}
	err := db.Pool.QueryRow(ctx, `
		SELECT id, user_id, title, description, goal, 
		       start_date, end_date, status, progress,
		       created_at, updated_at
		FROM challenges
		WHERE id = $1
	`, id).Scan(
		&c.ID, &c.UserID, &c.Title, &c.Description, &c.Goal,
		&c.StartDate, &c.EndDate, &c.Status, &c.Progress,
		&c.CreatedAt, &c.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetByID: %v", err)
		return nil, err
	}
	return c, nil
}

// Create создает новый вызов
func (r *ChallengeRepository) Create(ctx context.Context, c *challenge.Challenge) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO challenges (user_id, title, description, goal, start_date, end_date, status, progress)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, c.UserID, c.Title, c.Description, c.Goal, c.StartDate, c.EndDate, c.Status, c.Progress,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		log.Printf("Database error in Create: %v", err)
		return err
	}
	return nil
}

// Update обновляет вызов
func (r *ChallengeRepository) Update(ctx context.Context, c *challenge.Challenge) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE challenges
		SET title = $1, description = $2, goal = $3,
		    start_date = $4, end_date = $5,
		    status = $6, progress = $7,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`, c.Title, c.Description, c.Goal, c.StartDate, c.EndDate,
		c.Status, c.Progress, c.ID)

	if err != nil {
		log.Printf("Database error in Update: %v", err)
		return err
	}
	return nil
}

// Delete удаляет вызов
func (r *ChallengeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM challenges WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in Delete: %v", err)
		return err
	}
	return nil
}

// GetCheckpoints получает чекпоинты вызова
func (r *ChallengeRepository) GetCheckpoints(ctx context.Context, challengeID uuid.UUID) ([]*challenge.Checkpoint, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, challenge_id, title, description, due_date, is_completed, completed_at, order_num, created_at, updated_at
		FROM checkpoints
		WHERE challenge_id = $1
		ORDER BY order_num
	`, challengeID)

	if err != nil {
		log.Printf("Database error in GetCheckpoints: %v", err)
		return nil, err
	}
	defer rows.Close()

	var checkpoints []*challenge.Checkpoint
	for rows.Next() {
		cp := &challenge.Checkpoint{}
		err := rows.Scan(
			&cp.ID, &cp.ChallengeID, &cp.Title, &cp.Description, &cp.DueDate,
			&cp.IsCompleted, &cp.CompletedAt, &cp.OrderNum,
			&cp.CreatedAt, &cp.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning checkpoint: %v", err)
			continue
		}
		checkpoints = append(checkpoints, cp)
	}

	return checkpoints, nil
}

// CreateCheckpoint создает чекпоинт
func (r *ChallengeRepository) CreateCheckpoint(ctx context.Context, cp *challenge.Checkpoint) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO checkpoints (challenge_id, title, description, due_date, order_num)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, cp.ChallengeID, cp.Title, cp.Description, cp.DueDate, cp.OrderNum,
	).Scan(&cp.ID, &cp.CreatedAt, &cp.UpdatedAt)

	if err != nil {
		log.Printf("Database error in CreateCheckpoint: %v", err)
		return err
	}
	return nil
}

// UpdateCheckpoint обновляет чекпоинт
func (r *ChallengeRepository) UpdateCheckpoint(ctx context.Context, cp *challenge.Checkpoint) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE checkpoints
		SET title = $1, description = $2, due_date = $3,
		    is_completed = $4, completed_at = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`, cp.Title, cp.Description, cp.DueDate, cp.IsCompleted, cp.CompletedAt, cp.ID)

	if err != nil {
		log.Printf("Database error in UpdateCheckpoint: %v", err)
		return err
	}
	return nil
}

// DeleteCheckpoint удаляет чекпоинт
func (r *ChallengeRepository) DeleteCheckpoint(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM checkpoints WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in DeleteCheckpoint: %v", err)
		return err
	}
	return nil
}

// GetArtifacts получает артефакты
func (r *ChallengeRepository) GetArtifacts(ctx context.Context, challengeID uuid.UUID) ([]*challenge.Artifact, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, challenge_id, type, name, url, created_at
		FROM artifacts
		WHERE challenge_id = $1
	`, challengeID)

	if err != nil {
		log.Printf("Database error in GetArtifacts: %v", err)
		return nil, err
	}
	defer rows.Close()

	var artifacts []*challenge.Artifact
	for rows.Next() {
		a := &challenge.Artifact{}
		err := rows.Scan(&a.ID, &a.ChallengeID, &a.Type, &a.Name, &a.URL, &a.CreatedAt)
		if err != nil {
			log.Printf("Error scanning artifact: %v", err)
			continue
		}
		artifacts = append(artifacts, a)
	}

	return artifacts, nil
}

// CreateArtifact создает артефакт
func (r *ChallengeRepository) CreateArtifact(ctx context.Context, a *challenge.Artifact) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO artifacts (challenge_id, type, name, url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, a.ChallengeID, a.Type, a.Name, a.URL,
	).Scan(&a.ID, &a.CreatedAt)

	if err != nil {
		log.Printf("Database error in CreateArtifact: %v", err)
		return err
	}
	return nil
}

// DeleteArtifact удаляет артефакт
func (r *ChallengeRepository) DeleteArtifact(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM artifacts WHERE id = $1`, id)
	if err != nil {
		log.Printf("Database error in DeleteArtifact: %v", err)
		return err
	}
	return nil
}

// GetSelfAssessment получает самооценку
func (r *ChallengeRepository) GetSelfAssessment(ctx context.Context, challengeID uuid.UUID) (*challenge.SelfAssessment, error) {
	sa := &challenge.SelfAssessment{}
	err := db.Pool.QueryRow(ctx, `
		SELECT id, challenge_id, rating, comment, created_at
		FROM self_assessments
		WHERE challenge_id = $1
	`, challengeID).Scan(
		&sa.ID, &sa.ChallengeID, &sa.Rating, &sa.Comment, &sa.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error in GetSelfAssessment: %v", err)
		return nil, err
	}
	return sa, nil
}

// CreateSelfAssessment создает самооценку
func (r *ChallengeRepository) CreateSelfAssessment(ctx context.Context, sa *challenge.SelfAssessment) error {
	err := db.Pool.QueryRow(ctx, `
		INSERT INTO self_assessments (challenge_id, rating, comment)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, sa.ChallengeID, sa.Rating, sa.Comment,
	).Scan(&sa.ID, &sa.CreatedAt)

	if err != nil {
		log.Printf("Database error in CreateSelfAssessment: %v", err)
		return err
	}
	return nil
}

// UpdateSelfAssessment обновляет самооценку
func (r *ChallengeRepository) UpdateSelfAssessment(ctx context.Context, sa *challenge.SelfAssessment) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE self_assessments
		SET rating = $1, comment = $2
		WHERE id = $3
	`, sa.Rating, sa.Comment, sa.ID)

	if err != nil {
		log.Printf("Database error in UpdateSelfAssessment: %v", err)
		return err
	}
	return nil
}
