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
}

func NewProfessional(id, userID, licenseNumber string, category Category, yearsOfExperience int, resume string, metadata []byte) (*Professional, error) {
	if userID == "" {
		return nil, ErrEmptyUserID
	}
	if licenseNumber == "" {
		return nil, ErrEmptyLicenseNumber
	}
	if !category.IsValid() {
		return nil, ErrInvalidCategory
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
	if licenseNumber != nil {
		if *licenseNumber == "" {
			return ErrEmptyLicenseNumber
		}
		p.LicenseNumber = *licenseNumber
	}
	if category != nil {
		if !category.IsValid() {
			return ErrInvalidCategory
		}
		p.Category = *category
	}
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

