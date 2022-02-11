package config_parser

import (
	"srm/lib/logger"

	"github.com/Terry-Mao/goconf"
)

var config *goconf.Config

var configPath = "/etc/srm/srm.conf"
var err error
var configParam *goconf.Section

func Init() error {

	config = goconf.New()
	if err = config.Parse(configPath); err != nil {
		logger.Error(err)
	}

	if parseParam() != nil {
		return err
	}

	return nil
}

func parseParam() error {

	configParam = config.Get("setting")
	if configParam == nil {
		logger.Fatal("Fail to parse '/etc/srm/srm.conf' file.")
		return err
	}
	var isError int = 0
	var tint int64

	Setting = setting{}

	tint, err = configParam.Int("mape")
	Setting.Mape = int(tint)
	if err != nil {
		logger.Error(err)
		isError = 1
	}

	tint, err = configParam.Int("clean")
	Setting.Clean = int(tint)
	if err != nil {
		logger.Error(err)
		isError = 1
	}

	tint, err = configParam.Int("threshold")
	Setting.Threshold = int(tint)
	if err != nil {
		logger.Error(err)
		isError = 1
	}

	Setting.Minimum, err = configParam.String("minimum")
	if err != nil {
		logger.Error(err)
		isError = 1
	}

	if isError == 1 {
		return err
	}
	logger.Info("Success to parse /etc/srm/srm.conf")

	return nil
}
