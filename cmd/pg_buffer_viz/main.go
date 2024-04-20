package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/bonnefoa/pg_buffer_viz/internal/util"

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

func pgBufferVizFun(cmd *cobra.Command, args []string) {
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "pg_buffer_viz",
		Run: pgBufferVizFun,
	}
	rootFlags := rootCmd.PersistentFlags()
	err := viper.BindPFlags(rootFlags)
	util.FatalIf(err)

	versionCmd := &cobra.Command{
		Use:   "version",
		Run:   versionFun,
		Short: "Print command version",
	}
	rootCmd.AddCommand(versionCmd)

	defer pprof.StopCPUProfile()
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Root command failed: %v", err)
	}
}
