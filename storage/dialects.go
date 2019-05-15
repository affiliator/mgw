package storage

import (
	"bytes"
	"fmt"
	"github.com/affiliator/mgw/config"
	"reflect"
	"strings"
	"text/template"
)

const (
	SQLiteFormat DSNFormat = "{{ .Path}}"
	MySQLFormat  DSNFormat = "{{ .Username}}:{{ .Password}}@/{{ .Database}}?charset={{ .Charset}}&parseTime={{ .ParseTime}}"
)

var (
	availableDialects []string
	dialects          Dialects
)

func init() {
	dialects = NewDialectCollection()
}

func NewDialectCollection() Dialects {
	sqlite := Dialect{SQLiteFormat}
	mysql := Dialect{MySQLFormat}

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
	Format DSNFormat
}

func (d Dialect) Name() string {
	name := reflect.TypeOf(d).Name()
	return strings.ToLower(name)
}

func (d Dialect) Config() (config.Connection, error) {
	return config.Ptr().Storage.Connections.FindByDialect(d.Name())
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

func (d Dialects) Valid(dialect string) bool {
	if availableDialects == nil {
		availableDialects = d.getAvailable()
	}

	for _, v := range availableDialects {
		if strings.EqualFold(v, dialect) {
			return true
		}
	}

	return false
}

func (d Dialects) getAvailable() []string {
	t := reflect.TypeOf(d)
	dt := reflect.TypeOf(Dialect{})
	available := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type == dt {
			name := t.Field(i).Name
			available[i] = name
		}
	}

	return available
}

func (d Dialects) Configured() (Dialect, error) {
	var dialect Dialect
	configured := config.Ptr().Storage.Default

	switch configured {
	case "mysql":
		return d.MySQL, nil
	case "sqlite":
		return d.SQLite, nil
	}

	return dialect, fmt.Errorf("configured dialect `%s` is not implemented", configured)
}
