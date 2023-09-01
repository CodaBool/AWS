package main

import (
	basicLog "log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"
)

type TrendingGo struct {
	Name        string
	Href        string
	FullName    string
	Stars       int64
	Description string
	UpdatedAt   time.Time
}

type TrendingGithub struct {
	Name        string
	FullName    string
	Stars       int64
	Description string
	UpdatedAt   time.Time
}

type UpcomingMovie struct {
	Title     string
	Release   time.Time
	UpdatedAt time.Time
}

type TrendingMovie struct {
	Title     string
	Rank      int
	Rating    string
	Velocity  string
	UpdatedAt time.Time
}

type TrendingTV struct {
	Title     string
	Rank      int
	Velocity  string
	Rating    string
	UpdatedAt time.Time
}

type TrendingGame struct {
	Title     string
	Price     string
	MSRP      string
	UpdatedAt time.Time
}

type TrendingJS struct {
	Subject     string
	Page        int
	Rank        int
	Description string
	Title       string
	UpdatedAt   time.Time
}

type TrendingPY struct {
	Downloads   int64
	Name        string
	Description string
	UpdatedAt   time.Time
}

var db *gorm.DB

func dbInit(migrate bool) {
	var err error
	newLogger := gLog.New(
		basicLog.New(os.Stdout, "\r\n", basicLog.LstdFlags), // io writer
		gLog.Config{
			SlowThreshold: 5 * time.Second, // Slow SQL threshold
			Colorful:      true,            // Disable color
		},
	)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("PG_URI"),
		PreferSimpleProtocol: true, // necessary when pg pooling
	}), &gorm.Config{Logger: newLogger})
	check(err)

	slog.Info("Migrating")
	if migrate {
		db.AutoMigrate(&TrendingGo{}, &TrendingGithub{}, &TrendingTV{}, &UpcomingMovie{}, &TrendingGame{}, &TrendingJS{}, &TrendingPY{}, &TrendingMovie{})
	}
}

func upload(table string, data []any) {
	slog.Info("Clearing previous trending go data")
	db.Exec("DELETE FROM " + table)

	slog.Info("Inserting data")
	result := db.Create(data)
	check(result.Error)
}
