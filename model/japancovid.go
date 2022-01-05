package model

type (
	JapanCovidResponse struct {
		Date             string `json:"date"`
		Pcr              int    `json:"pcr"`
		Positive         int    `json:"positive"`
		Symptom          int    `json:"symptom"`
		Symptomless      int    `json:"symptomless"`
		SymtomConfirming int    `json:"symtomConfirming"`
		Hospitalize      int    `json:"hospitalize"`
		Mild             int    `json:"mild"`
		Severe           int    `json:"severe"`
		Confirming       int    `json:"confirming"`
		Waiting          int    `json:"waiting"`
		Discharge        int    `json:"discharge"`
		Death            int    `json:"death"`
	}
)
