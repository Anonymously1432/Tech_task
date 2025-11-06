package wallets_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"tech_task/internal/handler/wallets"
	mockuc "tech_task/internal/usecase/wallets/mocks"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap/zaptest"
)

func TestHandler_GetWalletAmount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Get("/wallet/:id", handler.GetWalletAmount)

	walletID := uuid.New()
	expectedAmount := float32(123.45)

	mockUc.EXPECT().
		GetWalletAmount(gomock.Any(), walletID).
		Return(expectedAmount, nil)

	req := httptest.NewRequest(http.MethodGet, "/wallet/"+walletID.String(), nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestHandler_GetWalletAmount_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Get("/wallet/:id", handler.GetWalletAmount)

	req := httptest.NewRequest(http.MethodGet, "/wallet/invalid-uuid", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestHandler_GetWalletAmount_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Get("/wallet/:id", handler.GetWalletAmount)

	walletID := uuid.New()

	mockUc.EXPECT().
		GetWalletAmount(gomock.Any(), walletID).
		Return(float32(0), pgx.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/wallet/"+walletID.String(), nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestHandler_GetWalletAmount_UcError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUc := mockuc.NewMockIUseCase(ctrl)
	handler := wallets.NewHandler(zaptest.NewLogger(t), mockUc)

	app := fiber.New()
	app.Get("/wallet/:id", handler.GetWalletAmount)

	walletID := uuid.New()
	mockUc.EXPECT().
		GetWalletAmount(gomock.Any(), walletID).
		Return(float32(0), errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/wallet/"+walletID.String(), nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}
