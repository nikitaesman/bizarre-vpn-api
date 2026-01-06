package services

import (
	"bizarre-vpn-api/internal/models"
	"bizarre-vpn-api/internal/services/interfaces"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"bizarre-vpn-api/internal/shared/random"
	"errors"
	"fmt"
	"log/slog"
)

const (
	LinkUserCodeLength = 6
)

type InviteLinksService struct {
	log                    *slog.Logger
	inviteLinksStorage     interfaces.InviteLinksStorage
	lnkUserProviderStorage interfaces.LnkUserProviderStorage
}

func NewInviteLinksService(
	log *slog.Logger,
	inviteLinksStorage interfaces.InviteLinksStorage,
	lnkUserProviderStorage interfaces.LnkUserProviderStorage,
) *InviteLinksService {
	return &InviteLinksService{
		inviteLinksStorage:     inviteLinksStorage,
		log:                    log,
		lnkUserProviderStorage: lnkUserProviderStorage,
	}
}

func (s *InviteLinksService) LinkUserWithTgProviderByCode(code string, externalUserId string) (userId int64, Err error) {
	op := "internal.services.InviteLinksService.LinkUserWithTgProviderByCode"

	inviteLink, err := s.inviteLinksStorage.GetInviteLinkByCode(code)

	if err != nil {
		if errors.Is(err, coreErrors.ErrorNotFound) {
			return 0, err
		}

		return 0, fmt.Errorf("%v: %w", op, err)
	}

	userProviders, err := s.lnkUserProviderStorage.GetListByUserId(inviteLink.UserId)

	if err != nil {
		return 0, fmt.Errorf("%v: %w", op, err)
	}

	if len(*userProviders) != 0 {
		return 0, coreErrors.ErrorUserAlreadyLinked
	}

	createLnkUserProviderPayload := models.CreateLnkUserProviderPayload{
		ProviderType:   models.TelegramProviderName,
		ExternalUserId: externalUserId,
		UserId:         inviteLink.UserId,
	}

	createdLnkUserProvider, err := s.lnkUserProviderStorage.CreateLnkUserProvider(&createLnkUserProviderPayload, nil)

	if err != nil {
		return 0, fmt.Errorf("%v: create LnkUserProvider error: %w", op, err)
	}

	return createdLnkUserProvider.UserId, nil
}

func (s *InviteLinksService) CreateItem(userId int64) (*models.InviteLink, error) {
	op := "internal.services.InviteLinksService.CreateItem"

	userProviders, err := s.lnkUserProviderStorage.GetListByUserId(userId)

	if err != nil {
		return nil, fmt.Errorf("%v: %w", op, err)
	}

	if len(*userProviders) != 0 {
		return nil, coreErrors.ErrorUserAlreadyLinked
	}

	randomCode, err := random.GetRandomString(LinkUserCodeLength)

	if err != nil {
		return nil, fmt.Errorf("%v: generate random code error: %w", op, err)
	}

	payload := &models.InviteLinkCreatePayload{
		UserId: userId,
		Code:   randomCode,
		Status: models.InviteLinkStatusCreated,
	}

	inviteLink, err := s.inviteLinksStorage.CreateInviteLink(payload)

	if err != nil {
		return nil, fmt.Errorf("%v: create inviteLink error: %w", op, err)
	}

	return inviteLink, nil
}
