package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
)

const (
	ListenAddrFlagName = "addr"
	PortFlagName       = "port"
	AvailRPCUrl        = "avail.rpc"
	Seed               = "avail.seed"
	AppID              = "avail.appid"
	Timeout            = "avail.timeout"
	TurboDAEnabled     = "avail.turboda"
	TurboDAURL         = "avail.turboda.url"
	TurboDAKey         = "avail.turboda.key"
)

const EnvVarPrefix = "OP_PLASMA_AVAIL_DA_SERVER"

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

var (
	ListenAddrFlag = &cli.StringFlag{
		Name:    ListenAddrFlagName,
		Usage:   "server listening address",
		Value:   "127.0.0.1",
		EnvVars: prefixEnvVars("ADDR"),
	}
	PortFlag = &cli.IntFlag{
		Name:    PortFlagName,
		Usage:   "server listening port",
		Value:   3100,
		EnvVars: prefixEnvVars("PORT"),
	}
	AvailRPCFlag = &cli.StringFlag{
		Name:    AvailRPCUrl,
		Usage:   "rpc url for avail node",
		EnvVars: prefixEnvVars("AVAIL_RPC"),
	}
	SeedFlag = &cli.StringFlag{
		Name:    Seed,
		Usage:   "avail seed phrase",
		EnvVars: prefixEnvVars("AVAIL_SEED"),
	}
	AppIDFlag = &cli.StringFlag{
		Name:    AppID,
		Usage:   "avail app id for the rollup",
		EnvVars: prefixEnvVars("AVAIL_APPID"),
	}
	TimeoutFlag = &cli.DurationFlag{
		Name:    Timeout,
		Usage:   "timeout parameter for request to avail",
		EnvVars: prefixEnvVars("AVAIL_TIMEOUT"),
		Value:   100 * time.Second,
	}
	TurboDAEnabledFlag = &cli.BoolFlag{
		Name:    TurboDAEnabled,
		Usage:   "enable turbo da",
		EnvVars: prefixEnvVars("AVAIL_TURBODA"),
		Value:   false,
	}
	TurboDAURLFlag = &cli.StringFlag{
		Name:    TurboDAURL,
		Usage:   "turbo da url",
		EnvVars: prefixEnvVars("AVAIL_TURBODA_URL"),
	}
	TurboDAKeyFlag = &cli.StringFlag{
		Name:    TurboDAKey,
		Usage:   "turbo da key",
		EnvVars: prefixEnvVars("AVAIL_TURBODA_KEY"),
	}
)

var requiredFlags = []cli.Flag{
	ListenAddrFlag,
	PortFlag,
	AppIDFlag,
}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	AvailRPCFlag,
	SeedFlag,
	TurboDAEnabledFlag,
	TurboDAURLFlag,
	TurboDAKeyFlag,
}

func init() {
	optionalFlags = append(optionalFlags, oplog.CLIFlags(EnvVarPrefix)...)
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

type CLIConfig struct {
	RPC        string
	Seed       string
	AppId      int
	Timeout    time.Duration
	TurboDA    bool
	TurboDAURL string
	TurboDAKey string
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		RPC:        ctx.String(AvailRPCUrl),
		Seed:       ctx.String(Seed),
		AppId:      ctx.Int(AppID),
		Timeout:    ctx.Duration(Timeout),
		TurboDA:    ctx.Bool(TurboDAEnabled),
		TurboDAURL: ctx.String(TurboDAURL),
		TurboDAKey: ctx.String(TurboDAKey),
	}
}

func (c CLIConfig) Check() error {
	if !c.TurboDA {
		if c.RPC == "" {
			return errors.New("no rpc url provided")
		}
		if c.Seed == "" {
			return errors.New("seedphrase not provided")
		}
	}
	if c.AppId == 0 {
		return errors.New("no app id provided")
	}
	if c.TurboDA {
		if c.TurboDAURL == "" {
			return errors.New("turbo da url not provided")
		}
		if c.TurboDAKey == "" {
			return errors.New("turbo da key not provided")
		}
	}
	return nil
}

func CheckRequired(ctx *cli.Context) error {

	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
