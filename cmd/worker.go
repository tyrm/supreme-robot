package cmd

import (
	"errors"
	"fmt"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/spf13/cobra"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/redis"
	"github.com/tyrm/supreme-robot/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run the worker",
	//TODO Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		requiredVars := []string{
			"EXT_HOSTNAME",
			"POSTGRES_DSN",
			"REDIS_DNS_ADDRESS",
			"SECRET",
		}
		c, err := config.CollectStartupConfig(requiredVars)
		if err != nil {
			log.Fatalf("error gathering configuration: %s", err.Error())
			return
		}

		// start logger
		err = loggo.ConfigureLoggers(c.LoggerConfig)
		if err != nil {
			log.Fatalf("error configuring logger: %s", err.Error())
			return
		}
		_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
		if err != nil {
			log.Fatalf("error configuring color logger: %s", err.Error())
			return
		}

		logger := loggo.GetLogger("main")
		logger.Infof("starting worker process")

		// create redis client
		rc, err := redis.NewClient(c.RedisDnsAddress, c.RedisDnsDB, c.RedisDnsPassword)
		if err != nil {
			logger.Errorf("new redis client: %s", err.Error())
			return
		}

		// create models client
		dc, err := models.NewClient(c)
		if err != nil {
			logger.Errorf("new models client: %s", err.Error())
			return
		}

		// create worker
		wkr, err := worker.NewWorker(rc, dc)
		if err != nil {
			logger.Errorf("new worker: %s", err.Error())
			return
		}

		// ** start application **
		errChan := make(chan error)

		go func() {
			err := wkr.Run()
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("worker: %s", err.Error()))
			}
		}()

		// Wait for SIGINT and SIGTERM (HIT CTRL-C)
		nch := make(chan os.Signal)
		signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-nch:
			logger.Infof("got sig: %s", sig)
		case err := <-errChan:
			logger.Criticalf(err.Error())
		}

		logger.Infof("worker process done")
	},
}
