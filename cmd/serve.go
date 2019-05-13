package cmd

import (
	"errors"
	"fmt"
	"github.com/flashmob/go-guerrilla"
	"github.com/spf13/cobra"
	"mailgun-mgw/config"
	mailgun_processor "mailgun-mgw/processor"
	"os"
	"os/signal"
	"syscall"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start smtp daemon.",
	Run:   serve,
}

var (
	d guerrilla.Daemon
	cfg *config.Configuration
	signalChannel = make(chan os.Signal, 1)
)

func init() {
	rootCmd.AddCommand(serveCmd)

	cfg = config.Ptr()

	serveCmd.PersistentFlags().StringVarP(&cfg.Paths.Config.Name, "config", "c",
		"", "Path to the configuration file")

	serveCmd.PersistentFlags().StringVarP(&cfg.Paths.Pid.Name, "pid", "p",
		"", "Path to the pid file")

	serveCmd.PersistentFlags().StringVarP(&cfg.Paths.Credentials.Name, "credentials", "a",
		"", "Path to the file containing mailgun api credentials")
}

func serve(cmd *cobra.Command, args []string) {
	d = guerrilla.Daemon{}
	d.AddProcessor("Mailgun", mailgun_processor.MailgunProcessor)

	err := readConfig(config.Ptr())
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
func readConfig(c *config.Configuration) error {
	if err := c.Reload(); err != nil {
		return err
	}

	if _, err := d.LoadConfig(c.Paths.Config.GetPath()); err != nil {
		return err
	}

	if d.Config.PidFile == "" {
		d.Config.PidFile = c.Paths.Pid.GetPath()
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
			_ = d.ReloadConfigFile(config.Ptr().Paths.Config.GetPath())
		} else if sig == syscall.SIGUSR1 {
			_ = d.ReopenLogs()
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
