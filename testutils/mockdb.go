package testutils

import (
	"testing"

	"booking-system/config"

	"github.com/pashagolub/pgxmock/v2"
)

func NewMockDB(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	config.DB = mock

	t.Cleanup(func() {
		mock.Close()
	})

	return mock
}
