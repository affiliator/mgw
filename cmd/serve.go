package cmd

import (
	"errors"
	"fmt"
	"github.com/flashmob/go-guerrilla"
	"github.com/spf13/cobra"
	mailgun_processor "mailgun-mgw/processor"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start smtp daemon.",
	Run:   serve,
}

const (
	defaultPidFile         string = "/var/run/mgw.pid"
	defaultCredentialsFile string = "credentials.json"
	defaultConfigFile      string = "config.example.json"
)

type Configuration struct {
	Path        string
	Credentials string
	PIDFile     string
}

func (c Configuration) getCurrentDirectory() string {
	Path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if strings.HasSuffix(Path, "/") {
		return Path
	}

	return Path + "/"
}

func (c Configuration) getDirectoryTo(file string) string {
	return c.getCurrentDirectory() + file
}

func (c Configuration) getConfiguration() string {
	return c.getDirectoryTo(c.Path)
}

func (c Configuration) getPidFile() string {
	return c.PIDFile
}

func (c Configuration) getCredentials() string {
	return c.getDirectoryTo(c.Credentials)
}

var (
	d             guerrilla.Daemon
	config        Configuration
	signalChannel = make(chan os.Signal, 1)
)

func init() {
	rootCmd.AddCommand(serveCmd)

	config = Configuration{defaultConfigFile, defaultCredentialsFile, defaultPidFile}

	serveCmd.PersistentFlags().StringVarP(&config.Path, "config", "c",
		defaultConfigFile, "Path to the configuration file")

	serveCmd.PersistentFlags().StringVarP(&config.PIDFile, "pid", "p",
		defaultPidFile, "Path to the pid file")

	serveCmd.PersistentFlags().StringVarP(&config.Credentials, "credentials", "a",
		defaultCredentialsFile, "Path to the file containing mailgun api credentials")
}

func initDaemon() guerrilla.Daemon {
	d = guerrilla.Daemon{}

	err := mailgun_processor.InitCredentials(config.getCredentials())
	if err != nil {
		panic(err)
	}

	d.AddProcessor("Mailgun", mailgun_processor.MailgunProcessor)

	return d
}

func serve(cmd *cobra.Command, args []string) {
	d = initDaemon()

	err := readConfig(config)
	if err != nil {
		panic(err)
	}

	if d.Start() == nil {
		fmt.Println("Server Started")
	}

	sigHandler()
}

// Superset of `guerrilla.AppConfig` containing options specific
// the the command line interface.
type DaemonConfig struct {
	guerrilla.AppConfig
}

func (c *DaemonConfig) emitChangeEvents(oldConfig *DaemonConfig, app guerrilla.Guerrilla) {
	// if your CmdConfig has any extra fields, you can emit events here
	// ...

	// call other emitChangeEvents
	c.AppConfig.EmitChangeEvents(&oldConfig.AppConfig, app)
}

// ReadConfig which should be called at startup
func readConfig(configuration Configuration) error {
	if _, err := d.LoadConfig(configuration.getConfiguration()); err != nil {
		return err
	}

	if d.Config.PidFile == "" {
		d.Config.PidFile = configuration.getPidFile()
	}

	if len(d.Config.AllowedHosts) == 0 {
		return errors.New("`allowed_hosts` should not be empty")
	}

	return nil
}

func sigHandler() {
	// handle SIGHUP for reloading the configuration while running
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGUSR1,
	)
	// Keep the daemon busy by waiting for signals to come
	for sig := range signalChannel {
		if sig == syscall.SIGHUP {
			d.ReloadConfigFile(config.Path)
		} else if sig == syscall.SIGUSR1 {
			d.ReopenLogs()
		} else if sig == syscall.SIGTERM || sig == syscall.SIGQUIT || sig == syscall.SIGINT {
			fmt.Println("Shutdown signal caught")
			d.Shutdown()
			fmt.Println("Shutdown completed, exiting.")
			return
		} else {
			fmt.Println("Shutdown, unknown signal caught")
			return
		}
	}
}
