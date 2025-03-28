package services

import (
	"cutpanionKiosk/internal/models"
	"testing"
)

func TestDayKeepsStreak(t *testing.T) {
	tests := []struct {
		name string
		data models.AggregatedData
		want bool
	}{
		{
			name: "valid streak day",
			data: models.AggregatedData{
				WorkoutToday: "Push Day",
				WeightToday:  175.5,
				Deficit:      -250,
				Protein:      130,
				Calories:     2200,
			},
			want: true,
		},
		{
			name: "missing workout",
			data: models.AggregatedData{
				WorkoutToday: "",
				WeightToday:  175.5,
				Deficit:      -250,
				Protein:      130,
				Calories:     2200,
			},
			want: false,
		},
		{
			name: "large deficit breaks streak",
			data: models.AggregatedData{
				WorkoutToday: "Pull Day",
				WeightToday:  175.5,
				Deficit:      -500,
				Protein:      130,
				Calories:     2200,
			},
			want: false,
		},
		{
			name: "low protein breaks streak",
			data: models.AggregatedData{
				WorkoutToday: "Push Day",
				WeightToday:  175.5,
				Deficit:      -250,
				Protein:      70,
				Calories:     2200,
			},
			want: false,
		},
		{
			name: "zero calories breaks streak",
			data: models.AggregatedData{
				WorkoutToday: "Push Day",
				WeightToday:  175.5,
				Deficit:      -250,
				Protein:      130,
				Calories:     0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dayKeepsStreak(tt.data)
			if got != tt.want {
				t.Errorf("dayKeepsStreak() = %v, want %v", got, tt.want)
			}
		})
	}
}
