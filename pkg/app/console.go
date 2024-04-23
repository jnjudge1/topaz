package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	builder "github.com/aserto-dev/service-host"
	"github.com/aserto-dev/topaz/pkg/app/handlers"
	"github.com/aserto-dev/topaz/pkg/cc/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ConsoleService struct{}

const (
	consoleService = "console"
)

func NewConsole() ServiceTypes {
	return &ConsoleService{}
}

func (e *ConsoleService) AvailableServices() []string {
	return []string{"console"}
}

func (e *ConsoleService) GetGRPCRegistrations(services ...string) builder.GRPCRegistrations {
	return func(server *grpc.Server) {
	}
}

func (e *ConsoleService) GetGatewayRegistration(services ...string) builder.HandlerRegistrations {
	return func(ctx context.Context, mux *runtime.ServeMux, grpcEndpoint string, opts []grpc.DialOption) error {
		return mux.HandlePath("GET", "/", runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			http.Redirect(w, r, "/ui/directory/model", http.StatusSeeOther)
		}))
	}
}

func (e *ConsoleService) Cleanups() []func() {
	return nil
}

func (e *ConsoleService) PrepareConfig(cfg *config.Config) *handlers.TopazCfg {
	authorizerURL := ""
	if serviceConfig, ok := cfg.APIConfig.Services[authorizerService]; ok {
		authorizerURL = getGatewayAddress(serviceConfig)
	}
	readerURL := ""
	if serviceConfig, ok := cfg.APIConfig.Services[readerService]; ok {
		readerURL = getGatewayAddress(serviceConfig)
	}
	writerURL := ""
	if serviceConfig, ok := cfg.APIConfig.Services[writerService]; ok {
		writerURL = getGatewayAddress(serviceConfig)
	}
	importerURL := ""
	if serviceConfig, ok := cfg.APIConfig.Services[importerService]; ok {
		importerURL = getGatewayAddress(serviceConfig)
	}

	exporterURL := ""
	if serviceConfig, ok := cfg.APIConfig.Services[exporterService]; ok {
		exporterURL = getGatewayAddress(serviceConfig)
	}

	modelURL := ""
	if serviceConfig, ok := cfg.APIConfig.Services[modelService]; ok {
		modelURL = getGatewayAddress(serviceConfig)
	}

	consoleURL := ""
	if serviceConfig, ok := cfg.APIConfig.Services[consoleService]; ok {
		consoleURL = getGatewayAddress(serviceConfig)
	}

	authorizerAPIKey := ""
	if _, ok := cfg.APIConfig.Services[authorizerService]; ok {
		for key := range cfg.Auth.APIKeys {
			// we only need a key
			authorizerAPIKey = key
			break
		}
	}

	directoryAPIKey := ""
	if _, ok := cfg.APIConfig.Services[readerService]; ok {
		for key := range cfg.Auth.APIKeys {
			// we only need a key
			directoryAPIKey = key
			break
		}
	}

	return &handlers.TopazCfg{
		AuthorizerServiceURL:        authorizerURL,
		AuthorizerAPIKey:            authorizerAPIKey,
		DirectoryServiceURL:         readerURL,
		DirectoryAPIKey:             directoryAPIKey,
		DirectoryTenantID:           cfg.DirectoryResolver.TenantID,
		DirectoryReaderServiceURL:   readerURL,
		DirectoryWriterServiceURL:   writerURL,
		DirectoryImporterServiceURL: importerURL,
		DirectoryExporterServiceURL: exporterURL,
		DirectoryModelServiceURL:    modelURL,
		ConsoleURL:                  consoleURL,
	}
}

func getGatewayAddress(serviceConfig *builder.API) string {
	if serviceConfig.Gateway.FQDN != "" {
		return serviceConfig.Gateway.FQDN
	}
	addr := serviceAddress(serviceConfig.Gateway.ListenAddress)

	if serviceConfig.Gateway.HTTP {
		return fmt.Sprintf("http://%s", addr)
	} else {
		return fmt.Sprintf("https://%s", addr)
	}
}

func serviceAddress(listenAddress string) string {
	return strings.Replace(listenAddress, "0.0.0.0", "localhost", 1)
}
