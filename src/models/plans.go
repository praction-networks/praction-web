package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Plan struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Category   string             `json:"category" bson:"category" validate:"required,oneof=Internet Landline Leased_Line Business_Broadband Smart_Business_Broadband Broadband Internet+OTT Internet+OTT+IPTV"`
	UUID       string             `json:"-" bson:"uuid"`
	PlanDetail []PlanSpecific     `json:"planDetails" bson:"planDetails" validate:"required,min=1,dive"`
}

type PlanSpecific struct {
	Name       string          `json:"name" bson:"name" validate:"required,max=30"`
	Speed      float64         `json:"speed" bson:"speed" validate:"required,gt=0"`
	SpeedUnit  string          `json:"speedUnit" bson:"speedUnit" validate:"required,oneof=Mbps Gbps"`
	Price      float64         `json:"price" bson:"price" validate:"required,gt=0"`
	Period     int             `json:"period" bson:"period" validate:"required,gt=0"`
	PeriodUnit string          `json:"periodUnit" bson:"periodUnit" validate:"required,oneof=Month Year Day Days Months Years"`
	Offering   map[string]bool `json:"offering" bson:"offering" validate:"required"`
}
