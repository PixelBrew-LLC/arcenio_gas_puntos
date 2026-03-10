package usecase

import (
	"context"
	"strconv"
	"time"

	"PixelBrew-LLC/arcenio_gas_puntos_backend/internal/domain"

	"github.com/google/uuid"
)

type transactionUsecase struct {
	ledgerRepo   domain.PointsLedgerRepository
	clientRepo   domain.ClientRepository
	settingsRepo domain.SettingsRepository
}

func NewTransactionUsecase(
	lr domain.PointsLedgerRepository,
	cr domain.ClientRepository,
	sr domain.SettingsRepository,
) domain.TransactionUsecase {
	return &transactionUsecase{
		ledgerRepo:   lr,
		clientRepo:   cr,
		settingsRepo: sr,
	}
}

func (u *transactionUsecase) EarnPoints(ctx context.Context, clientID string, gallons float64, processedByUserID string) (*domain.EarnResult, error) {
	if gallons <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	// 1. Verificar que el cliente existe
	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// 2. Leer configuración
	minGallonsSetting, err := u.settingsRepo.Get(ctx, domain.SettingMinGallons)
	if err != nil {
		return nil, err
	}
	minGallons, _ := strconv.ParseFloat(minGallonsSetting.Value, 64)

	// 3. Validar mínimo de galones
	if gallons < minGallons {
		return nil, domain.ErrBelowMinGallons
	}

	// 4. Calcular puntos
	ppgSetting, err := u.settingsRepo.Get(ctx, domain.SettingPointsPerGallon)
	if err != nil {
		return nil, err
	}
	pointsPerGallon, _ := strconv.ParseFloat(ppgSetting.Value, 64)
	pointsEarned := gallons * pointsPerGallon

	// 5. Calcular fecha de expiración
	expiryMonthsSetting, err := u.settingsRepo.Get(ctx, domain.SettingPointsExpiryMonths)
	if err != nil {
		return nil, err
	}
	expiryMonths, _ := strconv.Atoi(expiryMonthsSetting.Value)
	expiresAt := time.Now().AddDate(0, expiryMonths, 0)

	// 6. Crear entry en el ledger
	processedBy, _ := uuid.Parse(processedByUserID)
	entry := &domain.PointsLedger{
		ID:                uuid.New(),
		ClientID:          client.ID,
		Points:            pointsEarned,
		TransactionType:   domain.TransactionTypeEarn,
		GallonsAmount:     gallons,
		ProcessedByUserID: processedBy,
		CreatedAt:         time.Now(),
		ExpiresAt:         &expiresAt,
	}

	if err := u.ledgerRepo.Create(ctx, entry); err != nil {
		return nil, err
	}

	// 7. Obtener nuevo balance
	newBalance, err := u.ledgerRepo.GetBalance(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return &domain.EarnResult{
		PointsEarned: pointsEarned,
		NewBalance:   newBalance,
		Transaction:  entry,
		Client:       client,
	}, nil
}

func (u *transactionUsecase) RedeemPoints(ctx context.Context, clientID string, processedByUserID string) (*domain.RedeemResult, error) {
	// 1. Verificar que el cliente existe
	client, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// 2. Obtener balance actual
	currentBalance, err := u.ledgerRepo.GetBalance(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// 3. Verificar mínimo para canje
	minRedeemSetting, err := u.settingsRepo.Get(ctx, domain.SettingMinRedeemPoints)
	if err != nil {
		return nil, err
	}
	minRedeem, _ := strconv.ParseFloat(minRedeemSetting.Value, 64)

	if currentBalance < minRedeem {
		return nil, domain.ErrBelowMinRedeem
	}

	// 4. Canjear TODOS los puntos
	processedBy, _ := uuid.Parse(processedByUserID)
	entry := &domain.PointsLedger{
		ID:                uuid.New(),
		ClientID:          client.ID,
		Points:            -currentBalance, // Negativo: canjear todo el saldo
		TransactionType:   domain.TransactionTypeRedeem,
		GallonsAmount:     0,
		ProcessedByUserID: processedBy,
		CreatedAt:         time.Now(),
		ExpiresAt:         nil, // Los canjes no expiran
	}

	if err := u.ledgerRepo.Create(ctx, entry); err != nil {
		return nil, err
	}

	// 5. Obtener nuevo balance (debería ser 0)
	newBalance, err := u.ledgerRepo.GetBalance(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return &domain.RedeemResult{
		PointsRedeemed: currentBalance,
		NewBalance:     newBalance,
		Transaction:    entry,
		Client:         client,
	}, nil
}

func (u *transactionUsecase) GetClientBalance(ctx context.Context, clientID string) (float64, error) {
	// Verificar que el cliente existe
	_, err := u.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return 0, err
	}
	return u.ledgerRepo.GetBalance(ctx, clientID)
}

func (u *transactionUsecase) GetMinRedeemPoints(ctx context.Context) (float64, error) {
	setting, err := u.settingsRepo.Get(ctx, domain.SettingMinRedeemPoints)
	if err != nil {
		return 0, err
	}
	minRedeem, _ := strconv.ParseFloat(setting.Value, 64)
	return minRedeem, nil
}

func (u *transactionUsecase) GetTransactionHistory(ctx context.Context, filter domain.TransactionFilter) ([]*domain.PointsLedger, error) {
	return u.ledgerRepo.ListFiltered(ctx, filter)
}

func (u *transactionUsecase) GetDashboardStats(ctx context.Context, month, year int) (*domain.DashboardStats, error) {
	return u.ledgerRepo.GetDashboardStats(ctx, month, year)
}
