package cmd

import (
	"fmt"
	faktory "github.com/contribsys/faktory_worker_go"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/spf13/cobra"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/db/postgres"
	"github.com/tyrm/supreme-robot/kv/redis"
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
			"POSTGRES_DSN",
			"REDIS_DNS_ADDRESS",
		}
		c, err := config.CollectConfig(requiredVars)
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
		rc, err := redis.NewClient(c.RedisDNSAddress, c.RedisDNSDB, c.RedisDNSPassword)
		if err != nil {
			logger.Errorf("new redis client: %s", err.Error())
			return
		}

		// create db client
		dc, err := postgres.NewClient(c)
		if err != nil {
			logger.Errorf("new models client: %s", err.Error())
			return
		}

		// create faktory manager
		manager := faktory.NewManager()

		// create worker
		wkr, err := worker.NewWorker(rc, manager, dc)
		if err != nil {
			logger.Errorf("new worker: %s", err.Error())
			return
		}

		// ** start application **
		errChan := make(chan error)

		go func() {
			err := wkr.Run()
			if err != nil {
				errChan <- fmt.Errorf("worker: %s", err.Error())
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
