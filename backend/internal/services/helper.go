package services

import (
	"privat-unmei/internal/customerrors"
	"time"
)

func CalculateEndTime(startTime string, duration int) (string, error) {
	start, err := time.Parse("15:04:05", startTime)
	if err != nil {
		return "", customerrors.NewError(
			"failed to convert time",
			err,
			customerrors.CommonErr,
		)
	}
	end := start.Add(time.Duration(duration) * time.Minute)
	return end.Format("15:04:05"), nil
}
