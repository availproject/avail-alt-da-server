package main

import (
	"context"
	"errors"
	"fmt"

	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	availService "avail-alt-da-server/service"

	cli "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	RPC     string
	SEED    string
	APPID   int
	TIMEOUT time.Duration
)

func Check() error {
	if RPC == "" {
		return errors.New("no rpc url provided")
	}
	if APPID == 0 {
		return errors.New("no app id provided")
	}
	if SEED == "" {
		return errors.New("seedphrase not provided")
	}
	return nil
}

func TestAvailDAClientService(t *testing.T) {
	logger := testlog.Logger(t, log.LevelDebug)
	if err := godotenv.Load(); err != nil {
		logger.Crit("Error loading .env file: ", err)
	}
	RPC = os.Getenv("AVAIL_RPC")
	SEED = os.Getenv("AVAIL_SEED")

	appID, err := strconv.ParseInt(os.Getenv("AVAIL_APPID"), 10, 64)
	if err != nil {
		log.Crit("Error parsing APP_ID: ", err)
	}
	APPID = int(appID)

	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Info("Error parsing TIMEOUT: ", err)
		timeout = 100
	}
	TIMEOUT = time.Duration(timeout) * time.Second

	err = Check()
	if err != nil {
		panic(err)
	}

	store := availService.NewAvailService(RPC, SEED, APPID, TIMEOUT)

	ctx := context.Background()

	server := NewAvailDAServer("127.0.0.1", 0, store, logger, true)

	require.NoError(t, server.Start())

	cfg := cli.CLIConfig{
		Enabled:      true,
		DAServerURL:  fmt.Sprintf("http://%s", server.Endpoint()),
		VerifyOnRead: false,
		GenericDA:    true,
	}
	require.NoError(t, cfg.Check())

	client := cfg.NewDAClient()

	rng := rand.New(rand.NewSource(1234))

	input := RandomData(rng, 2000)

	comm, err := client.SetInput(ctx, input)
	println("comm", comm)
	require.NoError(t, err)

	stored, err := client.GetInput(ctx, comm)
	require.NoError(t, err)
	require.Equal(t, stored, input)
}

func RandomData(rng *rand.Rand, size int) []byte {
	out := make([]byte, size)
	rng.Read(out)
	return out
}
