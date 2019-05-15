package config

import (
	"fmt"
)

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

type Storage struct {
	Default     string      `json:"default"`
	Connections Connections `json:"connections"`
}

func (s Storage) Configured() (Connection, error) {
	return s.Connections.FindByDialect(s.Default)
}

type Connections struct {
	MySQL  MySQL  `json:"mysql"`
	SQLite SQLite `json:"sqlite"`
}

func (c Connections) FindByDialect(dialect string) (Connection, error) {
	switch dialect {
	case "mysql":
		return c.MySQL, nil
	case "sqlite":
		return c.SQLite, nil
	}

	return nil, fmt.Errorf("configured dialect `%s` is not implemented", dialect)
}

type Connection interface {
	getDSN() string
}

type MySQL struct {
	Hostname  string `json:"hostname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	Charset   string `json:"charset"`
	ParseTime bool   `json:"parseTime"`
}

func (m MySQL) getDSN() string {
	return fmt.Sprintf(
		"%s:%s@/%s?charset=%s&parseTime=%t",
		m.Username, m.Password, m.Database, m.Charset, m.ParseTime)
}

type SQLite struct {
	Path string `json:"path"`
}

func (c SQLite) getDSN() string {
	return c.Path
}
