package testutil

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

// MockPaymentRepository is a mock implementation of repository.PaymentRepository.
type MockPaymentRepository struct {
	Payments map[string]*entity.Payment // keyed by ID

	FindByIDFunc              func(ctx context.Context, id string) (*entity.Payment, error)
	FindByEnrollmentIDFunc    func(ctx context.Context, enrollmentID string) ([]entity.Payment, error)
	FindByGatewayPaymentIDFunc func(ctx context.Context, gw, gatewayPaymentID string) (*entity.Payment, error)
	FindAllFunc               func(ctx context.Context, filters repository.PaymentFilters) ([]entity.Payment, int, error)
}

func NewMockPaymentRepository() *MockPaymentRepository {
	return &MockPaymentRepository{
		Payments: make(map[string]*entity.Payment),
	}
}

func (m *MockPaymentRepository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	p, ok := m.Payments[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *MockPaymentRepository) FindByEnrollmentID(ctx context.Context, enrollmentID string) ([]entity.Payment, error) {
	if m.FindByEnrollmentIDFunc != nil {
		return m.FindByEnrollmentIDFunc(ctx, enrollmentID)
	}
	var result []entity.Payment
	for _, p := range m.Payments {
		if p.EnrollmentID == enrollmentID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *MockPaymentRepository) FindByGatewayPaymentID(ctx context.Context, gw, gatewayPaymentID string) (*entity.Payment, error) {
	if m.FindByGatewayPaymentIDFunc != nil {
		return m.FindByGatewayPaymentIDFunc(ctx, gw, gatewayPaymentID)
	}
	for _, p := range m.Payments {
		if p.Gateway == gw && p.GatewayPaymentID != nil && *p.GatewayPaymentID == gatewayPaymentID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockPaymentRepository) FindAll(ctx context.Context, filters repository.PaymentFilters) ([]entity.Payment, int, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, filters)
	}
	var result []entity.Payment
	for _, p := range m.Payments {
		result = append(result, *p)
	}
	return result, len(result), nil
}

func (m *MockPaymentRepository) Create(ctx context.Context, p *entity.Payment) error {
	m.Payments[p.ID] = p
	return nil
}

func (m *MockPaymentRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, p *entity.Payment) error {
	m.Payments[p.ID] = p
	return nil
}

func (m *MockPaymentRepository) Update(ctx context.Context, p *entity.Payment) error {
	m.Payments[p.ID] = p
	return nil
}

func (m *MockPaymentRepository) UpdateWithTx(ctx context.Context, tx *sqlx.Tx, p *entity.Payment) error {
	m.Payments[p.ID] = p
	return nil
}

func (m *MockPaymentRepository) UpdateStatus(ctx context.Context, id, status string) error {
	if p, ok := m.Payments[id]; ok {
		p.Status = status
	}
	return nil
}

func (m *MockPaymentRepository) UpdateStatusWithTx(ctx context.Context, tx *sqlx.Tx, id, status string) error {
	return m.UpdateStatus(ctx, id, status)
}

// MockCouponRepository is a mock implementation of repository.CouponRepository.
type MockCouponRepository struct {
	Coupons     map[string]*entity.Coupon // keyed by ID
	CouponCodes map[string]string         // code -> ID
	Usages      map[string]int            // "couponID:userID" -> count
}

func NewMockCouponRepository() *MockCouponRepository {
	return &MockCouponRepository{
		Coupons:     make(map[string]*entity.Coupon),
		CouponCodes: make(map[string]string),
		Usages:      make(map[string]int),
	}
}

func (m *MockCouponRepository) FindByID(ctx context.Context, id string) (*entity.Coupon, error) {
	c, ok := m.Coupons[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *MockCouponRepository) FindByCode(ctx context.Context, code string) (*entity.Coupon, error) {
	id, ok := m.CouponCodes[code]
	if !ok {
		return nil, nil
	}
	return m.FindByID(ctx, id)
}

func (m *MockCouponRepository) FindAll(ctx context.Context, activeOnly bool, page, perPage int) ([]entity.Coupon, int, error) {
	var result []entity.Coupon
	for _, c := range m.Coupons {
		if activeOnly && !c.IsActive {
			continue
		}
		result = append(result, *c)
	}
	return result, len(result), nil
}

func (m *MockCouponRepository) Create(ctx context.Context, c *entity.Coupon) error {
	m.Coupons[c.ID] = c
	m.CouponCodes[c.Code] = c.ID
	return nil
}

func (m *MockCouponRepository) Update(ctx context.Context, c *entity.Coupon) error {
	m.Coupons[c.ID] = c
	m.CouponCodes[c.Code] = c.ID
	return nil
}

func (m *MockCouponRepository) Delete(ctx context.Context, id string) error {
	if c, ok := m.Coupons[id]; ok {
		delete(m.CouponCodes, c.Code)
		delete(m.Coupons, id)
	}
	return nil
}

func (m *MockCouponRepository) IncrementUsage(ctx context.Context, id string) error {
	if c, ok := m.Coupons[id]; ok {
		c.CurrentUses++
	}
	return nil
}

func (m *MockCouponRepository) IncrementUsageWithTx(ctx context.Context, tx *sqlx.Tx, id string) error {
	return m.IncrementUsage(ctx, id)
}

func (m *MockCouponRepository) CreateUsage(ctx context.Context, u *entity.CouponUsage) error {
	key := u.CouponID
	if u.UserID != nil {
		key += ":" + *u.UserID
	}
	m.Usages[key]++
	return nil
}

func (m *MockCouponRepository) CreateUsageWithTx(ctx context.Context, tx *sqlx.Tx, u *entity.CouponUsage) error {
	return m.CreateUsage(ctx, u)
}

func (m *MockCouponRepository) CountUsageByUser(ctx context.Context, couponID, userID string) (int, error) {
	key := couponID + ":" + userID
	return m.Usages[key], nil
}

// MockPaymentTransactionRepository is a mock implementation.
type MockPaymentTransactionRepository struct {
	Transactions []*entity.PaymentTransaction
}

func NewMockPaymentTransactionRepository() *MockPaymentTransactionRepository {
	return &MockPaymentTransactionRepository{}
}

func (m *MockPaymentTransactionRepository) FindByPaymentID(ctx context.Context, paymentID string) ([]entity.PaymentTransaction, error) {
	var result []entity.PaymentTransaction
	for _, t := range m.Transactions {
		if t.PaymentID == paymentID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (m *MockPaymentTransactionRepository) Create(ctx context.Context, txLog *entity.PaymentTransaction) error {
	m.Transactions = append(m.Transactions, txLog)
	return nil
}

func (m *MockPaymentTransactionRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, txLog *entity.PaymentTransaction) error {
	return m.Create(ctx, txLog)
}
