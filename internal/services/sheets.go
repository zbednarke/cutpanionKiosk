package services

import (
	"log"
)

func FetchSheetData() {
	srv := GetSheetsService()

	spreadsheetId := "16YW9JyuJtI91NY4Kttgo_WbkhvoDG9tM7DG5vDX8fzk"
	readRange := "Sheet1!A2:F" // Adjust based on your layout

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data: %v", err)
	}

	for _, row := range resp.Values {
		log.Printf("Date: %s, Weight: %s, Workout: %s\n", row[0], row[1], row[2])
		// Parse and cache rows as needed
	}
}
