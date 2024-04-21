package util

import (
	"flag"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func SetCommonCliFlags(fs *pflag.FlagSet, defaultLogLevel string) {
	fs.String("log-level", defaultLogLevel, "Log level to use")
	fs.String("cpu-profile", "", "Destination file for cpu profiling")
	fs.String("mem-profile", "", "Destination file for memory profiling")
	fs.Duration("timeout", 5*time.Second, "Timeout")
}

func CommonInitialization() {
	configureLog()
	cpuProfile := viper.GetString("cpu-profile")
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		FatalIf(err)
		pprof.StartCPUProfile(f)
	}
}

func ConfigureViper() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	err := viper.BindPFlags(pflag.CommandLine)
	FatalIf(err)

	viper.SetEnvPrefix("PG_BUFFER_VIZ")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetConfigName(".pg_buffer_viz")
	viper.AddConfigPath("$HOME")
	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		FatalIf(err)
	}
}
