package cronApp

import (
	"flag"
	"log"
	"os"
	"privat-unmei/internal/db"
	"privat-unmei/internal/logger"

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
	var isProd bool
	flag.BoolVar(&isProd, "release", false, "Run production environemnt")
	flag.Parse()
	rcu := NewCronUtil()
	logger, err := logger.CreateNewLogger(isProd)
	if err != nil {
		log.Fatalln(err.Error())
	}
	driver, err := db.ConnectDB(logger)
	if err != nil {
		log.Fatalln(err.Error())
	}
	crc := NewCourseRequestCron(driver, logger)
	if err := rcu.AddJob(os.Getenv("EXPIRED_CRON_SPEC"), crc.UpdateExpiredRequest); err != nil {
		log.Fatalln(err.Error())
	}
	if err := rcu.AddJob(os.Getenv("COMPLETED_CRON_SPEC"), crc.UpdateCompletedRequest); err != nil {
		log.Fatalln(err.Error())
	}
	rcu.Start()
	select {}
}
