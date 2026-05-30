package services

import (
	"testing"

	"booking-system/models"
	"booking-system/testutils"

	"github.com/pashagolub/pgxmock/v2"
)

func TestPromoteToAdmin_NoOpWhenAlreadyAdmin(t *testing.T) {
	_ = testutils.NewMockDB(t)

	user := models.User{Email: "a@example.com", Role: true}
	if err := PromoteToAdmin(user); err != nil {
		t.Fatalf("PromoteToAdmin error: %v", err)
	}
}

func TestPromoteToAdmin_UpdatesUser(t *testing.T) {
	mock := testutils.NewMockDB(t)

	user := models.User{Email: "a@example.com", Role: false}
	mock.ExpectExec("UPDATE users SET admin=true").
		WithArgs(user.Email).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	if err := PromoteToAdmin(user); err != nil {
		t.Fatalf("PromoteToAdmin error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
