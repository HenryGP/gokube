package logging

import (
	"fmt"

	"go.uber.org/zap"
)

func InitLogger(preset string) (*zap.SugaredLogger, error) {
	var zapLogger *zap.Logger
	var err error

	switch preset {
	case "PROD":
		zapLogger, err = zap.NewProduction()
	case "DEV":
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		fmt.Println("Failed to create logger, will use the default one")
		fmt.Println(err)
	}

	zap.ReplaceGlobals(zapLogger)
	var logger *zap.SugaredLogger
	logger = zap.S()
	return logger, nil
}
