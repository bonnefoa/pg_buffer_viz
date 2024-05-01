package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/bonnefoa/pg_buffer_viz/pkg/bufferviz"
	"github.com/bonnefoa/pg_buffer_viz/pkg/db"
	"github.com/bonnefoa/pg_buffer_viz/pkg/httpserver"
	"github.com/bonnefoa/pg_buffer_viz/pkg/render"
	"github.com/bonnefoa/pg_buffer_viz/pkg/util"
	"github.com/rotisserie/eris"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version   = "dev"
	gitCommit = "none"
	gitBranch = "unknown"
	goVersion = "unknown"
	buildDate = "unknown"
)

func versionFun(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Git hash: %s\n", gitCommit)
	fmt.Printf("Git branch: %s\n", gitBranch)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Go Version: %s\n", goVersion)
	os.Exit(0)
}

func generateFun(cmd *cobra.Command, args []string) {
	dbConfig := db.GetDbConfigCli()

	timeout := viper.GetDuration("timeout")
	output := viper.GetString("output")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	d, err := db.NewDbPool(ctx, dbConfig.ConnectUrl)
	if err != nil {
		logrus.Fatalf("Error connecting to PostgreSQL: %s", eris.ToString(err, true))
	}
	table, err := d.FetchTable(ctx, dbConfig.Relation)
	if err != nil {
		logrus.Fatalf("Error when fetching table information: %s", eris.ToString(err, true))
	}

	canvas := render.NewFileCanvas(output)
	b := bufferviz.NewBufferViz(canvas.SVG, 30, 20)
	b.DrawTable(table)

	os.Exit(0)
}

func handleSignals(cancel context.CancelFunc) {
	sigIn := make(chan os.Signal, 100)
	signal.Notify(sigIn)
	for sig := range sigIn {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			logrus.Errorf("Caught signal '%s' (%d); terminating.", sig, sig)
			cancel()
		}
	}
}

func serveFun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	httpServerConfCli := httpserver.GetHttpServerConfigCli()
	_, err := httpserver.StartHttpServer(ctx, &httpServerConfCli)
	if err != nil {
		logrus.Fatalf("Error starting http server: %s", eris.ToString(err, true))
	}

	go func() {
		logrus.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	<-ctx.Done()
}

func pgBufferVizFun(cmd *cobra.Command, args []string) {
	os.Exit(0)
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "pg_buffer_viz",
		Run: pgBufferVizFun,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Run:   versionFun,
		Short: "Print command version",
	}
	rootCmd.AddCommand(versionCmd)

	generate := &cobra.Command{
		Use:   "generate",
		Run:   generateFun,
		Short: "Generate visualisation for a given relation",
	}
	serve := &cobra.Command{
		Use:   "serve",
		Run:   serveFun,
		Short: "Start the http server",
	}
	rootCmd.AddCommand(generate)
	rootCmd.AddCommand(serve)

	// Setup Flags
	rootFlags := rootCmd.PersistentFlags()
	util.SetCommonCliFlags(rootFlags, "info")
	db.SetDbConfigFlags(rootFlags)
	rootFlags.String("output", "output.svg", "Output filename")
	err := viper.BindPFlags(rootFlags)
	util.FatalIf(err)

	serveFlags := serve.Flags()
	httpserver.SetHttpServerConfigFlags(serveFlags)
	err = viper.BindPFlags(serveFlags)
	util.FatalIf(err)

	util.ConfigureViper()
	cobra.OnInitialize(util.CommonInitialization)

	defer pprof.StopCPUProfile()
	defer util.DoMemoryProfile()
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Root command failed: %v", eris.ToString(err, true))
	}
}
