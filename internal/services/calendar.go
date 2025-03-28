package services

import (
	"log"
	"time"
)

func FetchCalendarEvents() {
	srv := GetCalendarService()
	calendarId := "primary"

	tNow := time.Now().Format(time.RFC3339)
	tTomorrow := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	events, err := srv.Events.List(calendarId).
		TimeMin(tNow).TimeMax(tTomorrow).
		SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve calendar events: %v", err)
	}

	for _, item := range events.Items {
		log.Printf("Workout Event: %s (%s)", item.Summary, item.Start.DateTime)
	}
}
