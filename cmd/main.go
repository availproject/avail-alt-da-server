package main

import (
	"avail-alt-da-server/avail"
	flags "avail-alt-da-server/flags"
	server "avail-alt-da-server/server"
	"context"
	"fmt"
	"os"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

var Version = "v0.0.3"

func main() {

	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, "", "", "")
	app.Name = "avail-alt-da-da-server"
	app.Usage = "Alt DA Avail Service"
	app.Description = "Service for interacting with Avail DA"
	app.Action = StartDAServer

	ctx := opio.WithInterruptBlocker(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

func StartDAServer(cliCtx *cli.Context) error {
	if err := flags.CheckRequired(cliCtx); err != nil {
		return err
	}

	cfg := flags.ReadCLIConfig(cliCtx)

	if err := cfg.Check(); err != nil {
		return err
	}

	logCfg := oplog.ReadCLIConfig(cliCtx)

	l := oplog.NewLogger(oplog.AppOut(cliCtx), logCfg)
	oplog.SetGlobalLogHandler(l.Handler())

	l.Info("Initializing Alt DA server...")

	daprovider, err := avail.NewDAProvider(cfg, l)
	if err != nil {
		return fmt.Errorf("failed to create da service: %w", err)
	}

	server := server.NewDAServer(cfg.Addr, cfg.Port, daprovider, l, true)

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start the DA server")
	} else {
		l.Info("Started DA Server")
	}

	defer func() {
		if err := server.Stop(); err != nil {
			l.Error("failed to stop DA server", "err", err)
		}
	}()

	opio.BlockOnInterrupts()

	return nil
}
