package logger

import "go.uber.org/zap"

type CustomLogger interface {
	Info(args ...interface{})
	Infoln(args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
}

func CreateNewLogger(isProd bool) (*zap.SugaredLogger, error) {
	if isProd {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		return logger.Sugar(), nil
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
