package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConsentType string

const (
	ConsentPrivacyPolicy ConsentType = "privacy_policy"
	ConsentDataCampaign  ConsentType = "data_campaign"
)

// User Consent Status solo guarda el último estado.
type UserConsentStatus struct {
	UserID                 uuid.UUID
	ConsentedPrivacyPolicy bool
	ConsentedDataCampaign  bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type UserConsent struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      ConsentType
	Accepted  bool
	Version   string
	IPAddress string
	CreatedAt time.Time
}

type NewUserConsentParams struct {
	UserID    uuid.UUID
	Type      ConsentType
	Accepted  bool
	Version   string
	IPAddress string
}

func NewUserConsentStatus(userId uuid.UUID, consentPrivacy, consentDataCampaign bool) *UserConsentStatus {
	now := time.Now()
	return &UserConsentStatus{
		UserID:                 userId,
		ConsentedPrivacyPolicy: consentPrivacy,
		ConsentedDataCampaign:  consentDataCampaign,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
}

func NewUserConsent(uc NewUserConsentParams) *UserConsent {
	now := time.Now()
	return &UserConsent{
		ID:        uuid.New(),
		UserID:    uc.UserID,
		Type:      uc.Type,
		Accepted:  uc.Accepted,
		Version:   uc.Version,
		IPAddress: uc.IPAddress,
		CreatedAt: now,
	}
}
