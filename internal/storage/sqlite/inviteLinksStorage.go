package sqlite

import (
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"fmt"
)

type InviteLinksStorage struct {
	db Database
}

func (s *InviteLinksStorage) MustInit() {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS invite_links(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		code TEXT NOT NULL UNIQUE,
		status TEXT NOT NULL DEFAULT %v,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`, models.InviteLinkStatusCreated)

	_, err := s.db.Exec(query)

	if err != nil {
		panic(fmt.Errorf("failed to init invite_links table: %w", err))
	}
}

func (s *InviteLinksStorage) GetInviteLinkByCode(code string) (*models.InviteLink, error) {
	query := `SELECT 
		id,
		user_id,
		code,
		status,
		created_at,
		updated_at 
		FROM invite_links WHERE code = ?`

	var inviteLink models.InviteLink

	err := s.db.Get(&inviteLink, query, code)

	if err != nil {
		return nil, coreErrors.ErrorNotFound
	}

	return &inviteLink, nil
}

func (s *InviteLinksStorage) CreateInviteLink(payload *models.InviteLinkCreatePayload) (*models.InviteLink, error) {
	query := `INSERT INTO invite_links (user_id, code, status) VALUES (:user_id, :code, :status) RETURNING *`

	rows, err := s.db.NamedQuery(
		query,
		payload,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create invite link: %w", err)
	}

	defer rows.Close()

	var inviteLink models.InviteLink

	if !rows.Next() {
		return nil, fmt.Errorf("not rows next")
	}

	err = rows.StructScan(&inviteLink)

	if err != nil {
		return nil, fmt.Errorf("failed to scan inserted invite link: %w", err)
	}

	return &inviteLink, nil
}
