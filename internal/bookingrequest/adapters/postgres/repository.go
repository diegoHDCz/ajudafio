package bookingrequestpostgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/domain"
	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.BookingRequestRepository {
	return &repository{db: db, queries: New(db)}
}

func (r *repository) Create(ctx context.Context, req *domain.BookingRequest) (*domain.BookingRequest, error) {
	clientID, err := shared.ParseUUID(req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.Create: invalid client_id: %w", err)
	}
	professionalID, err := shared.ParseUUID(req.ProfessionalID)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.Create: invalid professional_id: %w", err)
	}
	addressID, err := shared.ParseUUID(req.AddressID)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.Create: invalid address_id: %w", err)
	}
	scheduleJSON, err := json.Marshal(req.ScheduleDetails)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.Create: marshal schedule_details: %w", err)
	}
	row, err := r.queries.CreateBookingRequest(ctx, CreateBookingRequestParams{
		ClientID:        clientID,
		ProfessionalID:  professionalID,
		AddressID:       addressID,
		ProposedValue:   float64ToNumeric(req.ProposedValue),
		ScheduleDetails: scheduleJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.Create: %w", err)
	}
	return rowToDomain(row)
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.BookingRequest, error) {
	uid, err := shared.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.GetByID: invalid id: %w", err)
	}
	row, err := r.queries.GetBookingRequestByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("booking request not found: %w", err)
		}
		return nil, fmt.Errorf("bookingrequestpostgres.GetByID: %w", err)
	}
	return rowToDomain(row)
}

func (r *repository) GetByClientID(ctx context.Context, clientID string) ([]*domain.BookingRequest, error) {
	uid, err := shared.ParseUUID(clientID)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.GetByClientID: invalid client_id: %w", err)
	}
	rows, err := r.queries.GetBookingRequestsByClientID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.GetByClientID: %w", err)
	}
	return rowsToDomain(rows)
}

func (r *repository) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.BookingRequest, error) {
	uid, err := shared.ParseUUID(professionalID)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.GetByProfessionalID: invalid professional_id: %w", err)
	}
	rows, err := r.queries.GetBookingRequestsByProfessionalID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.GetByProfessionalID: %w", err)
	}
	return rowsToDomain(rows)
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string, rejectionReason *string) (*domain.BookingRequest, error) {
	uid, err := shared.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres.UpdateStatus: invalid id: %w", err)
	}
	row, err := r.queries.UpdateBookingRequestStatus(ctx, UpdateBookingRequestStatusParams{
		ID:              uid,
		Status:          status,
		RejectionReason: rejectionReason,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("booking request not found: %w", err)
		}
		return nil, fmt.Errorf("bookingrequestpostgres.UpdateStatus: %w", err)
	}
	return rowToDomain(row)
}

func rowToDomain(row BookingRequest) (*domain.BookingRequest, error) {
	var schedule domain.ScheduleDetails
	if err := json.Unmarshal(row.ScheduleDetails, &schedule); err != nil {
		return nil, fmt.Errorf("bookingrequestpostgres: unmarshal schedule_details: %w", err)
	}
	return &domain.BookingRequest{
		ID:              uuid.UUID(row.ID.Bytes).String(),
		ClientID:        uuid.UUID(row.ClientID.Bytes).String(),
		ProfessionalID:  uuid.UUID(row.ProfessionalID.Bytes).String(),
		AddressID:       uuid.UUID(row.AddressID.Bytes).String(),
		ProposedValue:   numericToFloat64(row.ProposedValue),
		ScheduleDetails: schedule,
		Status:          domain.BookingRequestStatus(row.Status),
		RejectionReason: row.RejectionReason,
		CreatedAt:       row.CreatedAt.Time,
		RespondedAt:     pgTimestamptzToTimePtr(row.RespondedAt),
	}, nil
}

func rowsToDomain(rows []BookingRequest) ([]*domain.BookingRequest, error) {
	result := make([]*domain.BookingRequest, 0, len(rows))
	for _, row := range rows {
		d, err := rowToDomain(row)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(f, 'f', 2, 64))
	return n
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid || n.Int == nil {
		return 0
	}
	val := new(big.Float).SetInt(n.Int)
	if n.Exp != 0 {
		factor := new(big.Float).SetFloat64(math.Pow(10, float64(n.Exp)))
		val.Mul(val, factor)
	}
	f, _ := val.Float64()
	return f
}

func pgTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
