package cache

import "cutpanionKiosk/internal/models"

var currentData models.AggregatedData

func SetLatest(data models.AggregatedData) {
	currentData = data
}

func GetLatest() models.AggregatedData {
	return currentData
}
