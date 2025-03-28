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
	Date              string    `json:"date"`
	WorkoutToday      string    `json:"workout_today"`
	Streak            int       `json:"streak"`
	WeightToday       float64   `json:"weight_today"`
	Quote             string    `json:"quote"`
	WeightChartData   []float64 `json:"weight_chart_data"`
	DeficitChartData  []float64 `json:"deficit_chart_data"`
	ProteinChartData  []float64 `json:"protein_chart_data"`
	CaloriesChartData []float64 `json:"calories_chart_data"`
	Deficit           float64   `json:"deficit"`
	Protein           float64   `json:"protein"`
	Calories          float64   `json:"calories"`
}
