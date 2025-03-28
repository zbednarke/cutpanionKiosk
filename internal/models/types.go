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
	WorkoutToday string
	Streak       int
	WeightToday  float64
	Quote        string
	ChartData    []float64
	Deficit      float64
	Protein      float64
	Calories     float64
	// etc...
}
