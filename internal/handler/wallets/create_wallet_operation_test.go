package wallets_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"tech_task/internal/domain"
	"tech_task/internal/handler/wallets"
	"testing"

	mockuc "tech_task/internal/usecase/wallets/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/zap/zaptest"
)

func TestHandler_CreateWalletOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Post("/wallet/operation", handler.CreateWalletOperation)

	walletID := uuid.New()
	amount := float32(100)
	operationType := "DEPOSIT"

	reqBody := domain.WalletOperationRequest{
		WalletId:      walletID,
		OperationType: operationType,
		Amount:        amount,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	mockUc.EXPECT().
		CreateWalletOperation(gomock.Any(), walletID, operationType, amount).
		Return(float32(500), nil)

	req := httptest.NewRequest(http.MethodPost, "/wallet/operation", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestHandler_CreateWalletOperation_BodyParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Post("/wallet/operation", handler.CreateWalletOperation)

	req := httptest.NewRequest(http.MethodPost, "/wallet/operation", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestHandler_CreateWalletOperation_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Post("/wallet/operation", handler.CreateWalletOperation)

	reqBody := domain.WalletOperationRequest{
		WalletId:      uuid.New(),
		OperationType: "DEPOSIT",
		Amount:        0,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	resp, _ := app.Test(httptest.NewRequest(http.MethodPost, "/wallet/operation", bytes.NewReader(bodyBytes)))
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestHandler_CreateWalletOperation_UcError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Post("/wallet/operation", handler.CreateWalletOperation)

	walletID := uuid.New()
	amount := float32(100)
	operationType := "DEPOSIT"

	reqBody := domain.WalletOperationRequest{
		WalletId:      walletID,
		OperationType: operationType,
		Amount:        amount,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockUc.EXPECT().
		CreateWalletOperation(gomock.Any(), walletID, operationType, amount).
		Return(float32(0), errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/wallet/operation", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}
