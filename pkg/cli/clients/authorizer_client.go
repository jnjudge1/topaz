package clients

import (
	"context"
	"fmt"

	azc "github.com/aserto-dev/go-aserto/client"
	"github.com/fullstorydev/grpcurl"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/aserto-dev/go-authorizer/aserto/authorizer/v2"
	"github.com/aserto-dev/topaz/pkg/cli/cc"
)

type AuthorizerConfig struct {
	Host     string `flag:"host" short:"H" default:"${authorizer_svc}" env:"TOPAZ_AUTHORIZER_SVC" help:"authorizer service address"`
	APIKey   string `flag:"api-key" short:"k" default:"${authorizer_key}" env:"TOPAZ_AUTHORIZER_KEY" help:"authorizer API key"`
	Token    string `flag:"token" default:"${authorizer_token}" env:"TOPAZ_AUTHORIZER_TOKEN" help:"authorizer OAuth2.0 token" hidden:""`
	Insecure bool   `flag:"insecure" short:"i" default:"${insecure}" env:"TOPAZ_INSECURE" help:"skip TLS verification"`
	TenantID string `flag:"tenant-id" help:"" default:"${tenant_id}" env:"ASERTO_TENANT_ID" `
}

func NewAuthorizerClient(c *cc.CommonCtx, cfg *AuthorizerConfig) (authorizer.AuthorizerClient, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("no host specified")
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	opts := []azc.ConnectionOption{
		azc.WithAddr(cfg.Host),
		azc.WithInsecure(cfg.Insecure),
	}

	if cfg.APIKey != "" {
		opts = append(opts, azc.WithAPIKeyAuth(cfg.APIKey))
	}

	if cfg.Token != "" {
		opts = append(opts, azc.WithTokenAuth(cfg.Token))
	}

	if cfg.TenantID != "" {
		opts = append(opts, azc.WithTenantID(cfg.TenantID))
	}

	conn, err := azc.NewConnection(c.Context, opts...)
	if err != nil {
		return nil, err
	}

	return authorizer.NewAuthorizerClient(conn), nil
}

func (cfg *AuthorizerConfig) validate() error {
	ctx := context.Background()

	tlsConf, err := grpcurl.ClientTLSConfig(cfg.Insecure, "", "", "")
	if err != nil {
		return errors.Wrap(err, "failed to create TLS config")
	}

	creds := credentials.NewTLS(tlsConf)

	opts := []grpc.DialOption{
		grpc.WithUserAgent("topaz/dev-build (no version set)"), // TODO: verify user-agent value
	}

	if cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if _, err := grpcurl.BlockingDial(ctx, "tcp", cfg.Host, creds, opts...); err != nil {
		return err
	}

	return nil
}
