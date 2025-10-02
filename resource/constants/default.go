package constants

import "time"

const (
	DefaultServiceTimeOut       time.Duration = 10 * time.Second
	DefaultConnMaxLifeTime      time.Duration = 1 * time.Hour
	DefaultConnMaxIdleTime      time.Duration = 15 * time.Minute
	DefaultAccessTokenDuration  time.Duration = 1 * time.Hour
	DefaultRefreshTokenDuration time.Duration = 24 * time.Hour * 7 // 7 days
	DefaultDBPingInterval       time.Duration = 1 * time.Second
	DefaultDBRetryAttempts      int           = 3

	DefaultWorkerNamespace     string = "default"
	DefaultWorkerConcurrency   int    = 10
	DefaultWorkerRetryAttempts int    = 3

	// status
	InternalServerError string = "Internal Server Error"
	BadRequest          string = "Bad Request"
	LayoutDateTime      string = "2006-01-02 15:04:05"
	LayoutDate          string = "2006-01-02"
	TimeZone            string = "Asia/Jakarta"
	DBTransaction       string = "db_trx"
	DbContext           string = "ctx_db"
	Unauthorized        string = "Unauthorized"
)
