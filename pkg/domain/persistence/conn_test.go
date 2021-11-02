package persistence

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/billing/credits/internal/conf"
	"testing"
)

func TestConn(t *testing.T) {
	var cfg conf.Database
	require.NoError(t, cfg.Parse())

	db, err := OpenConn(cfg)
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	assert.NoError(t, sqlDB.Ping())
}
