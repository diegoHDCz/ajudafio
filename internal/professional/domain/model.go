package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyUserID          = errors.New("user ID cannot be empty")
	ErrEmptyLicenseNumber   = errors.New("license number cannot be empty")
	ErrInvalidCategory      = errors.New("invalid professional category")
	ErrNegativeYearsOfExp   = errors.New("years of experience cannot be negative")
)

type Category string

const (
	HospitalCompanion Category = "HOSPITAL_COMPANION"
	ElderlyCaregiver  Category = "ELDERLY_CAREGIVER"
	Nurse             Category = "NURSE"
	Physiotherapist   Category = "PHYSIOTHERAPIST"
)

func (c Category) IsValid() bool {
	switch c {
	case HospitalCompanion, ElderlyCaregiver, Nurse, Physiotherapist:
		return true
	}
	return false
}

func (c Category) RequiresLicenseNumber() bool {
	return c == Nurse || c == Physiotherapist
}

type Professional struct {
	ID                string
	UserID            string
	LicenseNumber     string
	Category          Category
	YearsOfExperience int
	Verified          bool
	Resume            string
	Metadata          []byte
	CreatedAt         time.Time
	UpdatedAt         time.Time
	UserName          *string
	UserAvatarURL     *string
	UserEmail         *string
	UserRole          *string
}

func NewProfessional(id, userID, licenseNumber string, category Category, yearsOfExperience int, resume string, metadata []byte) (*Professional, error) {
	if userID == "" {
		return nil, ErrEmptyUserID
	}
	if !category.IsValid() {
		return nil, ErrInvalidCategory
	}
	if licenseNumber == "" && category.RequiresLicenseNumber() {
		return nil, ErrEmptyLicenseNumber
	}
	if yearsOfExperience < 0 {
		return nil, ErrNegativeYearsOfExp
	}
	return &Professional{
		ID:                id,
		UserID:            userID,
		LicenseNumber:     licenseNumber,
		Category:          category,
		YearsOfExperience: yearsOfExperience,
		Resume:            resume,
		Metadata:          metadata,
	}, nil
}

func (p *Professional) ApplyUpdate(licenseNumber *string, category *Category, yearsOfExperience *int, verified *bool, resume *string, metadata []byte) error {
	effectiveCategory := p.Category
	if category != nil {
		if !category.IsValid() {
			return ErrInvalidCategory
		}
		effectiveCategory = *category
	}
	if licenseNumber != nil {
		if *licenseNumber == "" && effectiveCategory.RequiresLicenseNumber() {
			return ErrEmptyLicenseNumber
		}
		p.LicenseNumber = *licenseNumber
	}
	p.Category = effectiveCategory
	if yearsOfExperience != nil {
		if *yearsOfExperience < 0 {
			return ErrNegativeYearsOfExp
		}
		p.YearsOfExperience = *yearsOfExperience
	}
	if verified != nil {
		p.Verified = *verified
	}
	if resume != nil {
		p.Resume = *resume
	}
	if metadata != nil {
		p.Metadata = metadata
	}
	return nil
}

