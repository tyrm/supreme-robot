package cmd

import (
	"errors"
	"fmt"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/spf13/cobra"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/graphql"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/redis"
	"github.com/tyrm/supreme-robot/scheduler"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(graphqlCmd)
}

var graphqlCmd = &cobra.Command{
	Use:   "graphql",
	Short: "Run the graphql server",
	//TODO Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		requiredVars := []string{
			"ACCESS_SECRET",
			"EXT_HOSTNAME",
			"POSTGRES_DSN",
			"PRIMARY_NS",
			"REDIS_WEBAPP_ADDRESS",
			"REFRESH_SECRET",
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
		logger.Infof("starting main process")

		// create scheduler client
		sc, err := scheduler.NewClient()
		if err != nil {
			logger.Errorf("new scheduler client: %s", err.Error())
			return
		}

		// create models client
		dc, err := models.NewClient(c)
		if err != nil {
			logger.Errorf("new models client: %s", err.Error())
			return
		}

		// create redis client
		rc, err := redis.NewClient(c.RedisWebappAddress, c.RedisWebappDB, c.RedisWebappPassword)
		if err != nil {
			logger.Errorf("new redis client: %s", err.Error())
			return
		}

		// create web server
		ws, err := graphql.NewServer(c, sc, dc, rc)
		if err != nil {
			logger.Errorf("new webapp server: %s", err.Error())
			return
		}

		// ** start application **
		errChan := make(chan error)

		// start web server
		logger.Infof("starting web app")
		go func(errChan chan error) {
			err := ws.ListenAndServe()
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("webapp: %s", err.Error()))
			}
		}(errChan)
		defer ws.Close()

		// Wait for SIGINT and SIGTERM (HIT CTRL-C)
		nch := make(chan os.Signal)
		signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-nch:
			logger.Infof("got sig: %s", sig)
		case err := <-errChan:
			logger.Criticalf(err.Error())
		}

		logger.Infof("main process done")
	},
}
