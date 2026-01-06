package interfaces

import "bizarre-vpn-api/internal/models"

type InviteLinksStorage interface {
	GetInviteLinkByCode(code string) (*models.InviteLink, error)
	CreateInviteLink(payload *models.InviteLinkCreatePayload) (*models.InviteLink, error)
}
