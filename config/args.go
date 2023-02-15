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
	grpcServerAddr := flag.String("grpcAddr", "", "grpc server listen addr")
	httpServerAddr := flag.String("httpAddr", "", "http server listen addr")
	maxIdPoolSize := flag.Uint64("maxPoolSize", 50000, "OceanId max id pool size")
	minIdPoolSize := flag.Uint64("minPoolSize", 5000, "OceanId min pool size")
	idMdPath := flag.String("metadata", "./id_data", "OceanId metadata path")
	flag.Parse()
	*ctx = context.WithValue(*ctx, ArgsContextKey, Args{
		"GRPC_SERVER_ADDR": *grpcServerAddr,
		"HTTP_SERVER_ADDR": *httpServerAddr,
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
