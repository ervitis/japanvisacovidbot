package model

import "time"

type (
	JapanCovidResponse struct {
		Date             int `json:"date"`
		Pcr              int `json:"pcr"`
		Positive         int `json:"positive"`
		Symptom          int `json:"symptom"`
		Symptomless      int `json:"symptomless"`
		SymtomConfirming int `json:"symtomConfirming"`
		Hospitalize      int `json:"hospitalize"`
		Mild             int `json:"mild"`
		Severe           int `json:"severe"`
		Confirming       int `json:"confirming"`
		Waiting          int `json:"waiting"`
		Discharge        int `json:"discharge"`
		Death            int `json:"death"`
	}

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
