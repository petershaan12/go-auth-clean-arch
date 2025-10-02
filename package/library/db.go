package library

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database struct {
	DB *gorm.DB
}

func openDB(dsn string) (*gorm.DB, error) {
	dialect := mysql.Open(dsn)
	db, err := gorm.Open(dialect, &gorm.Config{
		// singular table
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(MaxIdleConns())
	conn.SetMaxOpenConns(MaxOpenConns())
	conn.SetConnMaxLifetime(ConnMaxLifeTime())
	conn.SetConnMaxIdleTime(ConnMaxIdleTime())

	return db, nil
}

func GetSqlDB() (*sql.DB, error) {
	db, err := openDB(DBDSN())

	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	return sqlDB, nil
}

func GetDatabase() (Database, error) {
	gormDB, err := openDB(DBDSN())
	if err != nil {
		return Database{}, err
	}
	return Database{DB: gormDB}, nil
}
