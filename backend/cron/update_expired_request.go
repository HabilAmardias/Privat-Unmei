package cronApp

import (
	"context"
	"database/sql"
	"log"
	"privat-unmei/internal/repositories"
)

type CourseRequestCron struct {
	crr *repositories.CourseRequestRepositoryImpl
	csr *repositories.CourseScheduleRepositoryImpl
	cr  *repositories.CourseRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func NewCourseRequestCron(db *sql.DB) *CourseRequestCron {
	crr := repositories.CreateCourseRequestRepository(db)
	csr := repositories.CreateCourseScheduleRepository(db)
	cr := repositories.CreateCourseRepository(db)
	tmr := repositories.CreateTransactionManager(db)
	return &CourseRequestCron{crr, csr, cr, tmr}
}

func (crc *CourseRequestCron) UpdateExpiredRequest() {
	ctx := context.Background()
	ids := new([]int)
	if err := crc.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crc.cr.UpdateCourseTransactionCount(ctx); err != nil {
			return err
		}
		if err := crc.crr.CancelExpiredRequest(ctx, ids); err != nil {
			return err
		}
		if len(*ids) > 0 {
			log.Println(*ids)
			if err := crc.csr.CancelExpiredSchedule(ctx, *ids); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Println(err.Error())
	}
}
