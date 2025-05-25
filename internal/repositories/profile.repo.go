package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"minitask1.go/internal/models"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (p *ProfileRepository) UpdateUserProfile(ctx context.Context, userID int, profile models.ProfileUserForm, filepath string) (models.Profile, error) {

	query := `UPDATE profiles SET `
	values := []any{}
	if profile.FirstName != "" {
		query += fmt.Sprintf(`first_name = $%d,`, len(values)+1)
		values = append(values, profile.FirstName)
	}

	if profile.LastName != "" {
		query += fmt.Sprintf(`last_name = $%d,`, len(values)+1)
		values = append(values, profile.LastName)
	}

	if profile.PhoneNumber != "" {
		query += fmt.Sprintf(`phone_number = $%d,`, len(values)+1)
		values = append(values, profile.PhoneNumber)
	}

	if filepath != "" {
		// if len(values) > 0 {
		// 	query += `, `
		// }
		query += fmt.Sprintf(`profile_picture = $%d`, len(values)+1)
		values = append(values, filepath)
	}
	query += fmt.Sprintf(" WHERE user_id = $%d", len(values)+1)
	values = append(values, userID)

	query += ` RETURNING first_name, last_name, profile_picture, phone_number`
	var result models.Profile
	if err := p.db.QueryRow(ctx, query, values...).Scan(&result.FirstName, &result.LastName, &result.ProfilePicture, &result.PhoneNumber); err != nil {
		return models.Profile{}, err
	}

	return result, nil
}

func (p *ProfileRepository) GetUserProfile(ctx context.Context, userID int) (*models.ProfileUser, error) {
	query := `
		SELECT 
    p.first_name,
    p.last_name,
    p.points,
    p.phone_number,
    t.name AS status
FROM profiles p
LEFT JOIN tiers t ON p.status = t.id
WHERE p.user_id = $1;
	`

	var profile models.ProfileUser
	err := p.db.QueryRow(ctx, query, userID).Scan(
		&profile.FirstName,
		&profile.LastName,
		&profile.Points,
		&profile.PhoneNumber,
		// &profile.ProfilePicture,
		&profile.Status,
	)

	if err != nil {
		return nil, err
	}
	return &profile, nil
}
