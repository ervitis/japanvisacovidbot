package model

import "time"

type (
	JapanCovidData struct {
		Date             string
		DateCovid        time.Time
		Pcr              int
		Positive         int
		Symptom          int
		Symptomless      int
		SymtomConfirming int
		Hospitalize      int
		Mild             int
		Severe           int
		Confirming       int
		Waiting          int
		Discharge        int
		Death            int
	}
)
