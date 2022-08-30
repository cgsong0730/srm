package main

import (
	"os"
	config "srm/lib/config_parser"
	"srm/lib/logger"
	"srm/module/mape"
	"strings"
)

func init() {
	err := logger.Init()
	if err != nil {
		logger.Fatal("Fail to initialize logger.")
		end()
	}

	err = config.Init()
	if err != nil {
		logger.Fatal("Fail to initialize config_parser")
	}

	logger.Info("#####[srm start]#####")
}

func main() {

	if len(os.Args) < 2 {
		logger.Fatal("Check the command format - srm < run | visual >")
	}

	firstArg := os.Args[1]

	if strings.Compare(firstArg, "run") == 0 {
		err := mape.Run()
		if err != nil {
			logger.Fatal("Fail to run mape module.")
		}
	} else if strings.Compare(firstArg, "visual") == 0 {
		err := mape.Visual()
		if err != nil {
			logger.Fatal("Fail to visualize workflow.")
		}
	} else {
		logger.Fatal("Check the command format - srm < run | visual >")
	}

	end()
}

func end() {
	logger.Info("#####[srm end]#####")
}
