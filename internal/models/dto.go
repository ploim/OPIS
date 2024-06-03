package models

import "time"

type ProgramsDTO struct {
	Name          string        `json:"name"`
	ContentType   string        `json:"content_type"`
	ContractStart time.Time     `json:"contract_start"`
	ContractEnd   time.Time     `json:"contract_end"`
	AirDuration   time.Duration `json:"air_duration"`
	AirFrequency  string        `json:"air_frequency"`
	TimeTypes     []string      `json:"time_types"`
}
