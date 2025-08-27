package cronApp

import (
	"context"
	"log"
)

func (crc *CourseRequestCron) UpdateCompletedRequest() {
	ctx := context.Background()
	if err := crc.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crc.csr.CompleteSchedule(ctx); err != nil {
			return err
		}
		if err := crc.crr.CompleteRequest(ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Println(err.Error())
	}
}
