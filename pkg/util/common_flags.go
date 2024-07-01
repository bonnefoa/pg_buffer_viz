package util

import (
	"flag"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/bonnefoa/pg_buffer_viz/pkg/model"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func SetCommonCliFlags(fs *pflag.FlagSet, defaultLogLevel string) {
	fs.String("log-level", defaultLogLevel, "Log level to use")
	fs.String("cpu-profile", "", "Destination file for cpu profiling")
	fs.String("mem-profile", "", "Destination file for memory profiling")
	fs.Int("block-width", 10, "Width of a block")
	fs.Int("block-height", 10, "Height of a block")
	fs.Int("margin-width", 3, "Width margin in block between elements")
	fs.Int("margin-height", 3, "Height margin in block between elements")
	fs.Duration("timeout", 5*time.Second, "Timeout")
}

func GetBlockSize() model.Size {
	res := model.Size{
		Width:  viper.GetInt("block-width"),
		Height: viper.GetInt("block-height"),
	}
	return res
}

func GetMarginSize() model.Size {
	res := model.Size{
		Width:  viper.GetInt("margin-width"),
		Height: viper.GetInt("margin-height"),
	}
	return res
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
