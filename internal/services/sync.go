package services

import "log"

func SyncAll() {
	log.Println("Syncing calendar and sheets...")

	FetchSheetData()
	FetchCalendarEvents()

	log.Println("Done syncing Syncing calendar and sheets.")

}
