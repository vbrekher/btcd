package main

import (
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/v2"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ffldb"
	"github.com/btcsuite/btcd/wire/v2"
	"github.com/stretchr/testify/require"
)

// TestHandleGetBlockNullVerbosity ensures an explicit JSON null verbosity is
// treated like the default verbosity instead of causing a nil dereference.
func TestHandleGetBlockNullVerbosity(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "blocks")
	db, err := database.Create("ffldb", dbPath, wire.MainNet)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	chain, err := blockchain.New(&blockchain.Config{
		DB:          db,
		ChainParams: &chaincfg.MainNetParams,
		TimeSource:  blockchain.NewMedianTime(),
	})
	require.NoError(t, err)

	cmd := &btcjson.GetBlockCmd{
		Hash:      chaincfg.MainNetParams.GenesisHash.String(),
		Verbosity: nil,
	}
	result, err := handleGetBlock(
		&rpcServer{cfg: rpcserverConfig{
			DB:          db,
			Chain:       chain,
			ChainParams: &chaincfg.MainNetParams,
		}},
		cmd, make(chan struct{}),
	)
	require.NoError(t, err)
	reply, ok := result.(btcjson.GetBlockVerboseResult)
	require.True(t, ok)
	require.Equal(t, chaincfg.MainNetParams.GenesisHash.String(), reply.Hash)
	require.NotEmpty(t, reply.Tx)
	require.Empty(t, reply.RawTx)
}
