package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error) {

	query := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comment_count
		FROM posts p
			LEFT JOIN comments c ON c.post_id = p.id
			LEFT JOIN users u ON p.user_id = u.id
			JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE f.user_id = $2 OR p.user_id = $3
		GROUP BY p.id, u.username
		ORDER BY p.id DESC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, userID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata

	for rows.Next() {

		var post PostWithMetadata

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, post)
	}

	return feed, nil
}

// ...................................................................

func (s *PostStore) Create(ctx context.Context, post *Post) error {

	query := `
		INSERT INTO posts
			(content, title, user_id, tags)
		Values ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	args := []any{post.Content, post.Title, post.UserID, pq.Array(post.Tags)}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {

	query := `
		SELECT 
			id, user_id, title, content, created_at, updated_at, tags, version
		FROM posts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}

	return &post, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {

	query := `
		UPDATE posts
		SET title = $1, content =$2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	args := []any{post.Title, post.Content, post.ID, post.Version}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {

	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
