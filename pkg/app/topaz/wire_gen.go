// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package topaz

import (
	"github.com/aserto-dev/logger"
	"github.com/aserto-dev/topaz/pkg/app"
	"github.com/aserto-dev/topaz/pkg/app/auth"
	"github.com/aserto-dev/topaz/pkg/app/impl"
	"github.com/aserto-dev/topaz/pkg/app/server"
	"github.com/aserto-dev/topaz/pkg/cc"
	"github.com/aserto-dev/topaz/pkg/cc/config"
	"github.com/aserto-dev/topaz/resolvers"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus"
)

// Injectors from wire.go:

func BuildApp(logOutput logger.Writer, errOutput logger.ErrWriter, configPath config.Path, overrides config.Overrider) (*app.Authorizer, func(), error) {
	ccCC, cleanup, err := cc.NewCC(logOutput, errOutput, configPath, overrides)
	if err != nil {
		return nil, nil, err
	}
	context := ccCC.Context
	zerologLogger := ccCC.Log
	configConfig := ccCC.Config
	common := &configConfig.Common
	group := ccCC.ErrGroup
	resolversResolvers := resolvers.New()
	authorizerServer := impl.NewAuthorizerServer(zerologLogger, common, resolversResolvers)
	grpcRegistrations, err := GRPCServerRegistrations(context, zerologLogger, configConfig, authorizerServer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	handlerRegistrations := GatewayServerRegistrations()
	serveMux := server.GatewayMux()
	registerer := _wireRegistererValue
	httpServer, err := server.NewGatewayServer(zerologLogger, common, serveMux, registerer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	serverServer, cleanup2, err := server.NewServer(context, zerologLogger, common, group, grpcRegistrations, handlerRegistrations, httpServer, serveMux)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	authorizer := &app.Authorizer{
		Context:       context,
		Logger:        zerologLogger,
		Configuration: configConfig,
		Server:        serverServer,
		Resolver:      resolversResolvers,
	}
	return authorizer, func() {
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireRegistererValue = prometheus.DefaultRegisterer
)

func BuildTestApp(logOutput logger.Writer, errOutput logger.ErrWriter, configPath config.Path, overrides config.Overrider) (*app.Authorizer, func(), error) {
	ccCC, cleanup, err := cc.NewTestCC(logOutput, errOutput, configPath, overrides)
	if err != nil {
		return nil, nil, err
	}
	context := ccCC.Context
	zerologLogger := ccCC.Log
	configConfig := ccCC.Config
	common := &configConfig.Common
	group := ccCC.ErrGroup
	resolversResolvers := resolvers.New()
	authorizerServer := impl.NewAuthorizerServer(zerologLogger, common, resolversResolvers)
	grpcRegistrations, err := GRPCServerRegistrations(context, zerologLogger, configConfig, authorizerServer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	handlerRegistrations := GatewayServerRegistrations()
	serveMux := server.GatewayMux()
	registry := prometheus.NewRegistry()
	httpServer, err := server.NewGatewayServer(zerologLogger, common, serveMux, registry)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	serverServer, cleanup2, err := server.NewServer(context, zerologLogger, common, group, grpcRegistrations, handlerRegistrations, httpServer, serveMux)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	authorizer := &app.Authorizer{
		Context:       context,
		Logger:        zerologLogger,
		Configuration: configConfig,
		Server:        serverServer,
		Resolver:      resolversResolvers,
	}
	return authorizer, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var (
	commonSet = wire.NewSet(server.NewServer, server.NewGatewayServer, server.GatewayMux, resolvers.New, impl.NewAuthorizerServer, GRPCServerRegistrations,
		GatewayServerRegistrations, auth.NewAPIKeyAuthMiddleware, wire.FieldsOf(new(*cc.CC), "Config", "Log", "Context", "ErrGroup"), wire.FieldsOf(new(*config.Config), "Common", "DecisionLogger"), wire.Struct(new(app.Authorizer), "*"),
	)

	appTestSet = wire.NewSet(
		commonSet, cc.NewTestCC, prometheus.NewRegistry, wire.Bind(new(prometheus.Registerer), new(*prometheus.Registry)),
	)

	appSet = wire.NewSet(
		commonSet, cc.NewCC, wire.InterfaceValue(new(prometheus.Registerer), prometheus.DefaultRegisterer),
	)
)