package models

import "time"

type Workout struct {
	Date  time.Time
	Title string
}

type WeightEntry struct {
	Date   time.Time
	Weight float64
}

type AggregatedData struct {
	Date         string
	Time         string
	WorkoutToday string
	Streak       int
	WeightToday  float64
	Quotes       string
	ChartData    []float64
	// etc...
}
