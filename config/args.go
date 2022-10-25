package config

import (
	"OceanID/utils/values"
	"context"
	"flag"
	"fmt"
)

// Args map readonly
type Args map[string]any

func (a Args) Get(key string) any {
	value, ok := a[key]
	if !ok {
		return ""
	}
	return value
}

const ArgsContextKey string = "ARGS"

func parseArgs(ctx *context.Context) {
	// configPath := flag.String("config", "./config.toml", "config file path")
	serverAddr := flag.String("addr", "127.0.0.1:7890", "grpc server listen addr")
	maxIdPoolSize := flag.Uint64("mps", 50000, "OceanId max id pool size")
	minIdPoolSize := flag.Uint64("nps", 5000, "OceanId min pool size")
	idMdPath := flag.String("imp", "id_data", "OceanId metadata path")
	flag.Parse()
	*ctx = context.WithValue(*ctx, ArgsContextKey, Args{
		"SERVER_ADDR":      *serverAddr,
		"MAX_ID_POOL_SIZE": *maxIdPoolSize,
		"MIN_ID_POOL_SIZE": *minIdPoolSize,
		"ID_METADATA_PATH": *idMdPath,
	})
}

func AssertArgs(ctx context.Context) (Args, error) {
	return values.ContextAssertion[Args](ctx, ArgsContextKey)
}

func PrintContextArgs(ctx context.Context) {
	value, err := AssertArgs(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(value)
}
