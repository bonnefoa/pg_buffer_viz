package db

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DbConfigCli struct {
	ConnectUrl string
	Relation   string
}

func SetDbConfigFlags(fs *pflag.FlagSet) {
	fs.String("connect-url", "", "Connection url to PostgreSQL db")
	fs.String("relation", "", "Target relation")
}

func GetDbConfigCli() DbConfigCli {
	d := DbConfigCli{}
	d.ConnectUrl = viper.GetString("connect-url")
	d.Relation = viper.GetString("relation")
	return d
}
