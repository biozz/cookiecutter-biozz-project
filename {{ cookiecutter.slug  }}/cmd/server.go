package cmd

import (
	"context"
	"fmt"

	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/clients"
	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/server"
	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/storage/mock"
	"github.com/urfave/cli/v2"
)

func ServerCmd() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "start a {{ cookiecutter.slug }} server",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			addr := cCtx.String("addr")
			db := cCtx.String("db")
			s := mock.New(db)
			// s := sqlite.New(db)
			err := s.Init(context.TODO())
			if err != nil {
				fmt.Println(err)
				return err
			}
			defer s.Close(context.TODO())

			w := server.New(server.NewServerParams{
				Storage:          s,
				OIDCProviderURL:  cCtx.String("oidc-provider-url"),
				OIDCClientID:     cCtx.String("oidc-client-id"),
				OIDCClientSecret: cCtx.String("oidc-client-secret"),
				JWTSecret:        cCtx.String("jwt-secret"),
				BaseURL:          cCtx.String("base-url"),
				Dev:              cCtx.String("env") == "dev",
				Clients: clients.New(clients.NewClientsParams{
					CentrifugoURL:    cCtx.String("centrifugo-url"),
					CentrifugoAPIKey: cCtx.String("centrifugo-api-key"),
				}),
			})
			w.Start(addr)
			return nil
		},
	}
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "addr",
		Value:   ":3002",
		EnvVars: []string{"SERVER_ADDR"},
	},
	&cli.StringFlag{
		Name:    "db",
		EnvVars: []string{"SERVER_DB"},
	},
	&cli.StringFlag{
		Name:    "env",
		EnvVars: []string{"SERVER_ENV"},
	},
	&cli.StringFlag{
		Name:    "jwt-secret",
		EnvVars: []string{"SERVER_JWT_SECRET"},
	},
	&cli.StringFlag{
		Name:    "base-url",
		EnvVars: []string{"SERVER_BASE_URL"},
	},
	&cli.StringFlag{
		Name:    "oidc-provider-url",
		EnvVars: []string{"SERVER_OIDC_PROVIDER_URL"},
        
	},
	&cli.StringFlag{
		Name:    "oidc-client-id",
		EnvVars: []string{"SERVER_OIDC_CLIENT_ID"},
	},
	&cli.StringFlag{
		Name:    "oidc-client-secret",
		EnvVars: []string{"SERVER_OIDC_CLIENT_SECRET"},
	},
	&cli.StringFlag{
		Name:    "centrifugo-url",
		EnvVars: []string{"SERVER_CENTRIFUGO_URL"},
	},
	&cli.StringFlag{
		Name:    "centrifugo-api-key",
		EnvVars: []string{"SERVER_CENTRIFUGO_API_KEY"},
	},
}
