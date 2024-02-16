package clients

import "{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/clients/centrifugo"

type Clients struct {
	Centrifugo centrifugo.Client
}

type NewClientsParams struct {
	CentrifugoURL    string
	CentrifugoAPIKey string
}

func New(params NewClientsParams) Clients {
	return Clients{
		Centrifugo: centrifugo.New(params.CentrifugoURL, params.CentrifugoAPIKey),
	}
}
