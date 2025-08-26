package cronApp

import (
	"log"
	"os"
	"privat-unmei/internal/db"

	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
)

type CronUtil struct {
	cr *cron.Cron
}

func NewCronUtil() *CronUtil {
	cr := cron.New()
	return &CronUtil{cr}
}

func (rcu *CronUtil) Start() {
	rcu.cr.Start()
}

func (rcu *CronUtil) Stop() {
	rcu.cr.Stop()
}

func (rcu *CronUtil) AddJob(spec string, callback func()) error {
	_, err := rcu.cr.AddFunc(spec, callback)
	return err
}

func Run() {
	rcu := NewCronUtil()
	driver, err := db.ConnectDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	crc := NewCourseRequestCron(driver)
	if err := rcu.AddJob(os.Getenv("CRON_SPEC"), crc.UpdateExpiredRequest); err != nil {
		log.Fatalln(err.Error())
	}
	rcu.Start()
	select {}
}
