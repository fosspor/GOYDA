package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	Interests    []string  `json:"interests"`
	PasswordHash string    `json:"-"`
}

type Location struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Seasons     []string  `json:"seasons"`
	Lat         *float64  `json:"lat,omitempty"`
	Lng         *float64  `json:"lng,omitempty"`
	MediaURLs   []string  `json:"media_urls"`
}

type Route struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Title     string          `json:"title"`
	Season    string          `json:"season"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

type Store struct {
	Pool *pgxpool.Pool
}

func (s *Store) CreateUser(ctx context.Context, email, hash, display string, interests []string) (User, error) {
	id := uuid.New()
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, display_name, interests)
		VALUES ($1,$2,$3,$4,$5)
	`, id, strings.ToLower(strings.TrimSpace(email)), hash, display, interests)
	if err != nil {
		return User{}, err
	}
	em := strings.ToLower(strings.TrimSpace(email))
	return User{ID: id, Email: em, DisplayName: display, Interests: interests, PasswordHash: hash}, nil
}

func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.Pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, COALESCE(interests, '{}')
		FROM users WHERE lower(email)=lower($1)
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.Interests)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (s *Store) UserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var u User
	err := s.Pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, COALESCE(interests, '{}')
		FROM users WHERE id=$1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.Interests)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (s *Store) UpdateUserInterests(ctx context.Context, id uuid.UUID, interests []string) error {
	tag, err := s.Pool.Exec(ctx, `UPDATE users SET interests=$2 WHERE id=$1`, id, interests)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) ListLocations(ctx context.Context, search string) ([]Location, error) {
	search = strings.TrimSpace(search)
	var rows pgx.Rows
	var err error
	if search == "" {
		rows, err = s.Pool.Query(ctx, `
			SELECT id, name, description, category, COALESCE(seasons,'{}'),
			       ST_Y(geom::geometry) as lat, ST_X(geom::geometry) as lng,
			       COALESCE(media_urls,'{}')
			FROM locations ORDER BY name
		`)
	} else {
		p := "%" + strings.ToLower(search) + "%"
		rows, err = s.Pool.Query(ctx, `
			SELECT id, name, description, category, COALESCE(seasons,'{}'),
			       ST_Y(geom::geometry) as lat, ST_X(geom::geometry) as lng,
			       COALESCE(media_urls,'{}')
			FROM locations
			WHERE lower(name) LIKE $1 OR lower(description) LIKE $1 OR lower(category) LIKE $1
			ORDER BY name
		`, p)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Description, &loc.Category, &loc.Seasons, &loc.Lat, &loc.Lng, &loc.MediaURLs); err != nil {
			return nil, err
		}
		out = append(out, loc)
	}
	return out, rows.Err()
}

// ListLocationsPage возвращает страницу локаций и общее число строк (с тем же фильтром search).
func (s *Store) ListLocationsPage(ctx context.Context, search string, limit, offset int) ([]Location, int, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}
	search = strings.TrimSpace(search)
	var total int
	var err error
	if search == "" {
		err = s.Pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM locations`).Scan(&total)
	} else {
		p := "%" + strings.ToLower(search) + "%"
		err = s.Pool.QueryRow(ctx, `
			SELECT COUNT(*)::int FROM locations
			WHERE lower(name) LIKE $1 OR lower(description) LIKE $1 OR lower(category) LIKE $1
		`, p).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}
	var rows pgx.Rows
	if search == "" {
		rows, err = s.Pool.Query(ctx, `
			SELECT id, name, description, category, COALESCE(seasons,'{}'),
			       ST_Y(geom::geometry) as lat, ST_X(geom::geometry) as lng,
			       COALESCE(media_urls,'{}')
			FROM locations ORDER BY name
			LIMIT $1 OFFSET $2
		`, limit, offset)
	} else {
		p := "%" + strings.ToLower(search) + "%"
		rows, err = s.Pool.Query(ctx, `
			SELECT id, name, description, category, COALESCE(seasons,'{}'),
			       ST_Y(geom::geometry) as lat, ST_X(geom::geometry) as lng,
			       COALESCE(media_urls,'{}')
			FROM locations
			WHERE lower(name) LIKE $1 OR lower(description) LIKE $1 OR lower(category) LIKE $1
			ORDER BY name
			LIMIT $2 OFFSET $3
		`, p, limit, offset)
	}
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()
	var out []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Description, &loc.Category, &loc.Seasons, &loc.Lat, &loc.Lng, &loc.MediaURLs); err != nil {
			return nil, total, err
		}
		out = append(out, loc)
	}
	if out == nil {
		out = []Location{}
	}
	return out, total, rows.Err()
}

func (s *Store) GetLocation(ctx context.Context, id uuid.UUID) (Location, error) {
	var loc Location
	err := s.Pool.QueryRow(ctx, `
		SELECT id, name, description, category, COALESCE(seasons,'{}'),
		       ST_Y(geom::geometry) as lat, ST_X(geom::geometry) as lng,
		       COALESCE(media_urls,'{}')
		FROM locations WHERE id=$1
	`, id).Scan(&loc.ID, &loc.Name, &loc.Description, &loc.Category, &loc.Seasons, &loc.Lat, &loc.Lng, &loc.MediaURLs)
	if errors.Is(err, pgx.ErrNoRows) {
		return Location{}, ErrNotFound
	}
	return loc, err
}

func (s *Store) CreateLocation(ctx context.Context, loc Location) (Location, error) {
	if loc.ID == uuid.Nil {
		loc.ID = uuid.New()
	}
	var err error
	if loc.Lat != nil && loc.Lng != nil {
		_, err = s.Pool.Exec(ctx, `
			INSERT INTO locations (id, name, description, category, seasons, geom, media_urls)
			VALUES ($1,$2,$3,$4,$5, ST_SetSRID(ST_MakePoint($7,$6),4326)::geography, $8)
		`, loc.ID, loc.Name, loc.Description, loc.Category, loc.Seasons, *loc.Lat, *loc.Lng, loc.MediaURLs)
	} else {
		_, err = s.Pool.Exec(ctx, `
			INSERT INTO locations (id, name, description, category, seasons, media_urls)
			VALUES ($1,$2,$3,$4,$5,$6)
		`, loc.ID, loc.Name, loc.Description, loc.Category, loc.Seasons, loc.MediaURLs)
	}
	if err != nil {
		return Location{}, err
	}
	return s.GetLocation(ctx, loc.ID)
}

// UpdateLocation обновляет поля; если lat/lng заданы — обновляет geom, иначе geom не трогает.
func (s *Store) UpdateLocation(ctx context.Context, loc Location) error {
	if loc.Lat != nil && loc.Lng != nil {
		tag, err := s.Pool.Exec(ctx, `
			UPDATE locations SET
				name=$2, description=$3, category=$4, seasons=$5,
				geom=ST_SetSRID(ST_MakePoint($7,$6),4326)::geography,
				media_urls=$8
			WHERE id=$1
		`, loc.ID, loc.Name, loc.Description, loc.Category, loc.Seasons, *loc.Lat, *loc.Lng, loc.MediaURLs)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return nil
	}
	tag, err := s.Pool.Exec(ctx, `
		UPDATE locations SET
			name=$2, description=$3, category=$4, seasons=$5, media_urls=$6
		WHERE id=$1
	`, loc.ID, loc.Name, loc.Description, loc.Category, loc.Seasons, loc.MediaURLs)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	tag, err := s.Pool.Exec(ctx, `DELETE FROM locations WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) ListRoutesForUser(ctx context.Context, userID uuid.UUID) ([]Route, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT id, user_id, title, season, COALESCE(payload, '{}'::jsonb), created_at
		FROM routes WHERE user_id=$1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Route
	for rows.Next() {
		var r Route
		var payload []byte
		if err := rows.Scan(&r.ID, &r.UserID, &r.Title, &r.Season, &payload, &r.CreatedAt); err != nil {
			return nil, err
		}
		r.Payload = json.RawMessage(payload)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) CreateRoute(ctx context.Context, userID uuid.UUID, title, season string, payload json.RawMessage) (Route, error) {
	id := uuid.New()
	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO routes (id, user_id, title, season, payload) VALUES ($1,$2,$3,$4,$5)
	`, id, userID, title, season, payload)
	if err != nil {
		return Route{}, err
	}
	return s.GetRoute(ctx, id)
}

func (s *Store) GetRoute(ctx context.Context, id uuid.UUID) (Route, error) {
	var r Route
	var payload []byte
	err := s.Pool.QueryRow(ctx, `
		SELECT id, user_id, title, season, COALESCE(payload, '{}'::jsonb), created_at
		FROM routes WHERE id=$1
	`, id).Scan(&r.ID, &r.UserID, &r.Title, &r.Season, &payload, &r.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Route{}, ErrNotFound
	}
	r.Payload = json.RawMessage(payload)
	return r, err
}

func (s *Store) RouteBelongsTo(ctx context.Context, routeID, userID uuid.UUID) (Route, error) {
	r, err := s.GetRoute(ctx, routeID)
	if err != nil {
		return Route{}, err
	}
	if r.UserID != userID {
		return Route{}, ErrNotFound
	}
	return r, nil
}

func (s *Store) UpdateRoute(ctx context.Context, routeID, userID uuid.UUID, title, season string, payload json.RawMessage) error {
	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}
	tag, err := s.Pool.Exec(ctx, `
		UPDATE routes SET title=$3, season=$4, payload=$5
		WHERE id=$1 AND user_id=$2
	`, routeID, userID, title, season, payload)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) DeleteRouteForUser(ctx context.Context, routeID, userID uuid.UUID) error {
	tag, err := s.Pool.Exec(ctx, `DELETE FROM routes WHERE id=$1 AND user_id=$2`, routeID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) SeedDemo(ctx context.Context) error {
	var cnt int
	if err := s.Pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM locations`).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}
	demos := []Location{
		{Name: "Станица Тамань", Description: "История, море, локальная кухня.", Category: "culture", Seasons: []string{"spring", "summer", "autumn"}, Lat: fp(45.2129), Lng: fp(36.7184), MediaURLs: []string{}},
		{Name: "Винодельня в долине", Description: "Дегустации и экскурсии по виноградникам.", Category: "wine", Seasons: []string{"summer", "autumn"}, Lat: fp(45.0), Lng: fp(38.9), MediaURLs: []string{}},
		{Name: "Фермерское хозяйство", Description: "Агротуризм, дегустация сыров.", Category: "agro", Seasons: []string{"spring", "summer"}, Lat: fp(45.2), Lng: fp(39.0), MediaURLs: []string{}},
	}
	for _, d := range demos {
		if _, err := s.CreateLocation(ctx, d); err != nil {
			return fmt.Errorf("seed %s: %w", d.Name, err)
		}
	}
	return nil
}

func fp(f float64) *float64 { return &f }
