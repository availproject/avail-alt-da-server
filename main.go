package main

import (
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

var Version = "v0.0.2"

func main() {
	fmt.Println("ADDR:", os.Getenv("ADDR"))
	fmt.Println("PORT:", os.Getenv("PORT"))
	fmt.Println("AVAIL_RPC:", os.Getenv("AVAIL_RPC"))
	fmt.Println("AVAIL_SEED:", os.Getenv("AVAIL_SEED"))
	fmt.Println("AVAIL_APPID:", os.Getenv("AVAIL_APPID"))
	fmt.Println("AVAIL_TIMEOUT:", os.Getenv("AVAIL_TIMEOUT"))
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(Flags)
	app.Version = opservice.FormatVersion(Version, "", "", "")
	app.Name = "avail-da-server"
	app.Usage = "Plasma Avail DA Service"
	app.Description = "Service for interacting with Avail DA"
	app.Action = StartDAServer

	ctx := opio.WithInterruptBlocker(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}
