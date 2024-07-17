package main

import "time"

type TrendingGo struct {
	Name        string
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
