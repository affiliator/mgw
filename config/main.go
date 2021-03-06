package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var config Configuration
var initialized = false

type Configuration struct {
	Paths     Paths
	Defaults  DefaultValues
	Providers Providers `json:"provider,omitempty"`
}

type Paths struct {
	Config      File
	Pid         File
	Credentials File
}

type File struct {
	Name     string
	Relative bool
}

func (f File) GetPath() string {
	if f.Relative == true {
		return config.getWorkingDirectory() + f.Name
	}

	return f.Name
}

func (f File) Exists() (bool, error) {
	_, err := os.Stat(f.GetPath())
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func (f File) ToString() string {
	return f.GetPath()
}

func (f File) Read() ([]byte, error) {
	file, err := os.Open(f.GetPath())
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}

func (f File) ReadTo(v interface{}) error {
	if result, _ := f.Exists(); result == false {
		return errors.New("File " + f.GetPath() + " does not exist.")
	}
	data, err := f.Read()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (Configuration) getWorkingDirectory() string {
	Path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if strings.HasSuffix(Path, "/") {
		return Path
	}

	return Path + "/"
}

type DefaultValues struct{}

const (
	defaultPidFile         string = "/var/run/mgw.pid"
	defaultCredentialsFile string = "credentials.json"
	defaultConfigFile      string = "config.json"
)

func (d DefaultValues) getPidFile() File {
	return File{defaultPidFile, false}
}

func (d DefaultValues) getCredentialsFile() File {
	return File{defaultCredentialsFile, true}
}

func (d DefaultValues) getConfigFile() File {
	return File{defaultConfigFile, true}
}

func Initialize() *Configuration {
	config = Configuration{}
	initialized = true
	return &config
}

func Ptr() *Configuration {
	if initialized == false {
		Initialize().FillWithDefaults()
	}

	return &config
}

func (c *Configuration) Read() ([]byte, error) {
	return c.Paths.Config.Read()
}

func (c *Configuration) ReadTo(v interface{}) error {
	return c.Paths.Config.ReadTo(v)
}

func (c *Configuration) Reload() error {
	return c.Paths.Config.ReadTo(c)
}

func (c *Configuration) FillWithDefaults() *Configuration {
	c.Paths.Config = c.Defaults.getConfigFile()
	c.Paths.Pid = c.Defaults.getPidFile()
	c.Paths.Credentials = c.Defaults.getCredentialsFile()

	return c
}
