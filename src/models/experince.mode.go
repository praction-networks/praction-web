package models

type Experince struct {
	Name               string `json:"name" bson:"name" validate:"required,oneof=admin user"`
	CandidateExperince string
}
