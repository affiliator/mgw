package config

type Providers struct {
	Mailgun Mailgun `json:"mailgun,omitempty"`
}

type Credentials struct {
	Domain string `json:"domain"`
	ApiKey string `json:"api_key"`
}

type Mailgun struct {
	ApiBase     string      `json:"api_base"`
	Credentials Credentials `json:"credentials"`
}