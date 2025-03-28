package services

import (
	"cutpanionKiosk/internal/cache"
	"cutpanionKiosk/internal/models"
	"log"
	"strconv"
	"time"

	"google.golang.org/api/sheets/v4"
)

func FetchSheetData() {
	srv := GetSheetsService()
	spreadsheetId := "16YW9JyuJtI91NY4Kttgo_WbkhvoDG9tM7DG5vDX8fzk"
	readRange := "DailyLog!A2:F" // Adjust based on your layout

	resp, err := fetchSheetResponse(srv, spreadsheetId, readRange)
	if err != nil {
		log.Fatalf("Unable to retrieve data: %v", err)
	}

	today := time.Now()
	todayStr := today.Format("2006-01-02")
	streak := 0

	// 1) Extract historical maps for each metric:
	historicalWeights := extractHistoricalValues(resp.Values, 0, 1)  // date col=0, weight col=1
	historicalDeficit := extractHistoricalValues(resp.Values, 0, 3)  // deficit col=3
	historicalProtein := extractHistoricalValues(resp.Values, 0, 4)  // protein col=4
	historicalCalories := extractHistoricalValues(resp.Values, 0, 5) // calories col=5

	// 2) Build chart data arrays
	weightChartData := buildChartData(historicalWeights, today)
	deficitChartData := buildChartData(historicalDeficit, today)
	proteinChartData := buildChartData(historicalProtein, today)
	caloriesChartData := buildChartData(historicalCalories, today)

	for _, row := range resp.Values {
		// Skip rows that do not have at least a date and one more column.
		if len(row) < 2 {
			continue
		}

		data := parseRow(row)
		// If this row is not for today, update the streak
		if data.Date != todayStr {
			streak = updateStreak(streak, data)
			continue
		}

		// For today's row, update the data with the streak and a random quote
		data.Streak = streak
		data.Quote = getRandomQuote()
		data.WeightChartData = weightChartData
		data.DeficitChartData = deficitChartData
		data.ProteinChartData = proteinChartData
		data.CaloriesChartData = caloriesChartData
		cache.SetLatest(data)
		log.Println("✅ Cached today's data from Google Sheets")
		return
	}

	log.Println("⚠️ No matching row found for today's date in sheet")
}

// fetchSheetResponse wraps the call to the Sheets API.
func fetchSheetResponse(srv *sheets.Service, spreadsheetId, readRange string) (*sheets.ValueRange, error) {
	return srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
}

// parseRow converts a raw row (slice of interface{}) to an AggregatedData struct.
func parseRow(row []interface{}) models.AggregatedData {
	return models.AggregatedData{
		Date:         parseString(row, 0),
		WeightToday:  parseFloat(row, 1),
		WorkoutToday: parseString(row, 2),
		Deficit:      parseFloat(row, 3),
		Protein:      parseFloat(row, 4),
		Calories:     parseFloat(row, 5),
	}
}

// updateStreak increments or resets the streak based on whether the day's data qualifies.
func updateStreak(currentStreak int, data models.AggregatedData) int {
	// For streak checking, we need to use today's date in the aggregated data.
	todayData := models.AggregatedData{
		Date:         time.Now().Format("2006-01-02"),
		WeightToday:  data.WeightToday,
		WorkoutToday: data.WorkoutToday,
		Deficit:      data.Deficit,
		Protein:      data.Protein,
		Calories:     data.Calories,
	}

	if dayKeepsStreak(todayData) {
		return currentStreak + 1
	}
	return 0
}

// parseString safely extracts a string value from a row at the given index.
func parseString(row []interface{}, index int) string {
	if len(row) > index {
		if s, ok := row[index].(string); ok {
			return s
		}
	}
	return ""
}

// parseFloat safely extracts a float64 from a row at the given index.
func parseFloat(row []interface{}, index int) float64 {
	if len(row) > index {
		if s, ok := row[index].(string); ok {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return f
			}
		}
	}
	return 0
}

func dayKeepsStreak(row models.AggregatedData) bool {
	return row.WorkoutToday != "" && row.WeightToday != 0 && row.Deficit > -300 && row.Protein > 90 && row.Calories > 0
}

// extractHistoricalValues iterates over rows and returns a map of date -> the float value
// at the given valueIndex (e.g., weight is column 1, deficit is column 3, etc.).
func extractHistoricalValues(rows [][]interface{}, dateIndex, valueIndex int) map[string]float64 {
	valuesMap := make(map[string]float64)
	for _, row := range rows {
		if len(row) <= valueIndex {
			continue
		}
		dateStr := parseString(row, dateIndex)
		val := parseFloat(row, valueIndex)
		if dateStr != "" && val != 0 {
			valuesMap[dateStr] = val
		}
	}
	return valuesMap
}

// extractHistoricalWeights iterates over all rows and returns a map of date -> weight.
func extractHistoricalWeights(rows [][]interface{}) map[string]float64 {
	weights := make(map[string]float64)
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		dateStr := parseString(row, 0)
		weight := parseFloat(row, 1)
		if dateStr != "" && weight != 0 {
			weights[dateStr] = weight
		}
	}
	return weights
}

// buildChartData creates a slice of 7 numbers (previous 7 days' data)
// interpolating missing values and backfilling initial gaps.
func buildChartData(historical map[string]float64, today time.Time) []float64 {
	const days = 7
	chartData := make([]float64, days)
	dates := make([]time.Time, days)
	// Build the target dates: from today-6 to today.
	for i := 0; i < days; i++ {
		dates[i] = today.AddDate(0, 0, -days+1+i)
	}

	// For each target day, check for an existing value.
	for i, date := range dates {
		dateStr := date.Format("2006-01-02")
		if val, ok := historical[dateStr]; ok {
			chartData[i] = val
		} else {
			// Missing data: attempt to interpolate using previous/next.
			prevIdx, nextIdx := -1, -1
			for j := i - 1; j >= 0; j-- {
				dStr := dates[j].Format("2006-01-02")
				if _, ok := historical[dStr]; ok {
					prevIdx = j
					break
				}
			}
			for j := i + 1; j < days; j++ {
				dStr := dates[j].Format("2006-01-02")
				if _, ok := historical[dStr]; ok {
					nextIdx = j
					break
				}
			}

			switch {
			case prevIdx != -1 && nextIdx != -1:
				// Both previous and next available: linear interpolation.
				prevVal := historical[dates[prevIdx].Format("2006-01-02")]
				nextVal := historical[dates[nextIdx].Format("2006-01-02")]
				factor := float64(i-prevIdx) / float64(nextIdx-prevIdx)
				chartData[i] = prevVal + (nextVal-prevVal)*factor
			case prevIdx != -1:
				// Backfill with previous.
				chartData[i] = historical[dates[prevIdx].Format("2006-01-02")]
			case nextIdx != -1:
				// Use next available.
				chartData[i] = historical[dates[nextIdx].Format("2006-01-02")]
			default:
				// No data at all: default to 0.
				chartData[i] = 0
			}
		}
	}

	return chartData
}

// buildWeightChartData creates a slice of 7 numbers (previous 7 days' weights)
// interpolating missing values and backfilling initial gapsweightC
func buildWeightChartData(historicalWeights map[string]float64, today time.Time) []float64 {
	const days = 7
	chartData := make([]float64, days)
	dates := make([]time.Time, days)
	// Build the target dates: from today-6 to today.
	for i := range days {
		dates[i] = today.AddDate(0, 0, -days+1+i)
	}

	// For each target day, check for an existing weight.
	for i, date := range dates {
		dateStr := date.Format("2006-01-02")
		if weight, ok := historicalWeights[dateStr]; ok {
			chartData[i] = weight
		} else {
			// Missing weight: attempt to interpolate using previous and next available values.
			prevIdx, nextIdx := -1, -1
			// Look backwards for the previous available weight.
			for j := i - 1; j >= 0; j-- {
				dStr := dates[j].Format("2006-01-02")
				if _, ok := historicalWeights[dStr]; ok {
					prevIdx = j
					break
				}
			}
			// Look forward for the next available weight.
			for j := i + 1; j < days; j++ {
				dStr := dates[j].Format("2006-01-02")
				if _, ok := historicalWeights[dStr]; ok {
					nextIdx = j
					break
				}
			}

			if prevIdx != -1 && nextIdx != -1 {
				// Both previous and next available: perform linear interpolation.
				prevDateStr := dates[prevIdx].Format("2006-01-02")
				nextDateStr := dates[nextIdx].Format("2006-01-02")
				prevWeight := historicalWeights[prevDateStr]
				nextWeight := historicalWeights[nextDateStr]
				factor := float64(i-prevIdx) / float64(nextIdx-prevIdx)
				chartData[i] = prevWeight + (nextWeight-prevWeight)*factor
			} else if prevIdx != -1 {
				// No next value; backfill with the previous available weight.
				prevDateStr := dates[prevIdx].Format("2006-01-02")
				chartData[i] = historicalWeights[prevDateStr]
			} else if nextIdx != -1 {
				// No previous value; use the next available weight.
				nextDateStr := dates[nextIdx].Format("2006-01-02")
				chartData[i] = historicalWeights[nextDateStr]
			} else {
				// No data at all: default to 0 (or you could choose another default).
				chartData[i] = 0
			}
		}
	}

	return chartData
}
