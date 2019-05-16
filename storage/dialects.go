package storage

import (
	"bytes"
	"fmt"
	"github.com/affiliator/mgw/config"
	"reflect"
	"text/template"
)

const (
	SQLiteFormat DSNFormat = "{{ .Path}}"
	MySQLFormat  DSNFormat = "{{ .Username}}:{{ .Password}}@/{{ .Database}}?charset={{ .Charset}}&parseTime={{ .ParseTime}}"

	SQLite string = "sqlite"
	MySQL  string = "mysql"
)

var (
	dialects Dialects
)

func init() {
	dialects = NewDialectCollection()
}

func NewDialectCollection() Dialects {
	sqlite := Dialect{SQLite, SQLiteFormat}
	mysql := Dialect{MySQL, MySQLFormat}

	return Dialects{
		MySQL:  mysql,
		SQLite: sqlite,
	}
}

func GetDialects() *Dialects {
	return &dialects
}

type DSNFormat string

func (f DSNFormat) String() string {
	return reflect.ValueOf(f).String()
}

type IDialect interface {
	Name() string
	Config() (config.Connection, error)
	DSN() string
}

type Dialect struct {
	Name   string
	Format DSNFormat
}

func (d Dialect) Config() (config.Connection, error) {
	return config.Ptr().Storage.Connections.ByDialect(d.Name)
}

func (d Dialect) DSN() string {
	c, err := d.Config()
	if err != nil {
		// Todo: Log.
		return ""
	}

	parse, err := template.New("dsn").Parse(d.Format.String())
	if err != nil {
		// Todo: Log
		return ""
	}

	var buff bytes.Buffer
	err = parse.Execute(&buff, c)
	if err != nil {
		// Todo: Log
		return ""
	}

	return buff.String()
}

type Dialects struct {
	MySQL  Dialect
	SQLite Dialect
}

func (d Dialects) Mapping() *map[string]DSNFormat {
	return &map[string]DSNFormat{
		SQLite: SQLiteFormat,
		MySQL:  MySQLFormat,
	}
}

func (d Dialects) Valid(dialect string) bool {
	if _, exists := (*d.Mapping())[dialect]; exists {
		return true
	}

	return false
}

func (d Dialects) Configured() (Dialect, error) {
	storage := config.Ptr().Storage.Default

	switch storage {
	case "mysql":
		return d.MySQL, nil
	case "sqlite":
		return d.SQLite, nil
	}

	return Dialect{}, fmt.Errorf("configured dialect `%s` is not implemented", storage)
}
