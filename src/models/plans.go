package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plan struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Category   string             `json:"category" bson:"category" validate:"required,oneof=Internet Landline Leased_Line Business_Broadband Smart_Business_Broadband Broadband Internet+OTT Internet+IPTV  Internet+OTT+IPTV Internet_Stream+"`
	UUID       string             `json:"-" bson:"uuid"`
	PlanDetail []PlanSpecific     `json:"planDetails" bson:"planDetails" validate:"required,min=1,dive"`
}

type PlanSpecific struct {
	PlanID       string               `json:"planID" bson:"planID"`
	Name         string               `json:"name" bson:"name" validate:"required,max=30"`
	PlanCategory string               `json:"planCategory" bson:"planCategory" validate:"max=20"`
	Category     string               `json:"category" bson:"category" validate:"required,oneof=Internet Landline Leased_Line Business_Broadband Smart_Business_Broadband Broadband Internet+OTT Internet+IPTV Internet+OTT+IPTV Internet_Stream+"`
	Speed        float64              `json:"speed" bson:"speed" validate:"required,gt=0"`
	SpeedUnit    string               `json:"speedUnit" bson:"speedUnit" validate:"required,oneof=Mbps Gbps"`
	Price        float64              `json:"price" bson:"price" validate:"required,gt=0"`
	Period       int                  `json:"period" bson:"period" validate:"required,gt=0"`
	PeriodUnit   string               `json:"periodUnit" bson:"periodUnit" validate:"required,oneof=Month Year Day Days Months Years"`
	Offering     map[string]bool      `json:"offering" bson:"offering" validate:"required"`
	OTTs         []primitive.ObjectID `json:"otts,omitempty" bson:"otts,omitempty"`
	IPTVs        []primitive.ObjectID `json:"iptvs,omitempty" bson:"iptvs,omitempty"`
	OttDetails   []Image              `json:"ottDetails,omitempty" bson:"-"`
	IPTVDetails  []Image              `json:"iptvDetails,omitempty" bson:"iptvDetails,omitempty"`
}

type AvailableUserRequest struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UUID          string             `json:"uuid" bson:"uuid"`
	PlanID        string             `json:"planID" bson:"planID" validate:"required"`
	FirstName     string             `json:"firstName" bson:"firstName" validate:"required,oneword"`
	MiddleName    string             `json:"middleName,omitempty" bson:"middleName,omitempty" validate:"omitempty,oneword"`
	LastName      string             `json:"lastName" bson:"lastName" validate:"required,oneword"`
	Email         string             `json:"email" bson:"email" validate:"required,email"`
	Mobile        string             `json:"mobile" bson:"mobile" validate:"required,mobile"`
	HearAboutUs   string             `json:"hearAboutUs" bson:"hearAboutUs" validate:"required"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsVerified    bool               `json:"isVerified" bson:"isVerified"`
	VerifiedAt    time.Time          `json:"verifiedAt" bson:"verifiedAt"`
	OTP           int64              `json:"-" bson:"otp"`
	OTPExpireTime time.Time          `json:"-" bson:"otpExpireTime"`
}

type UnAvailableArea struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UUID         string             `json:"-" bson:"uuid"`
	FirstName    string             `json:"firstName" bson:"firstName" validate:"required,oneword"`
	MiddleName   string             `json:"middleName,omitempty" bson:"middleName,omitempty" validate:"omitempty,oneword"`
	LastName     string             `json:"lastName" bson:"lastName" validate:"required,oneword"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Mobile       string             `json:"mobile" bson:"mobile" validate:"required,mobile"`
	PropertyType string             `json:"propertyType" bson:"propertyType" validate:"required"`
	HearAboutUs  string             `json:"hearAboutUs" bson:"hearAboutUs" validate:"required"`
	Area         Area               `json:"area" bson:"area" validate:"required"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Area struct {
	Coordinates      UnPointRequest `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
	FormattedAddress string         `bson:"formatted_address,omitempty" json:"formatted_address,omitempty"`
	PlaceID          string         `bson:"place_id,omitempty" json:"place_id,omitempty"`
	Pincode          string         `bson:"pincode,omitempty" json:"pincode,omitempty" validate:"len=6"`
}

type UnPointRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type UserOTPPlanVerify struct {
	Email  string `json:"email" validate:"email"`
	Mobile string `json:"mobile" validate:"len=10"`
	OTP    int64  `json:"otp" validate:"required,otp_len"`
	PlanID string `json:"planID" validate:"required"`
}

type UserOTPPlanResend struct {
	Email  string `json:"email" validate:"email"`
	Mobile string `json:"mobile" validate:"len=10"`
	Resend string `json:"resend" validate:"required,oneof=text voice"` // Corrected closing double-quote
	PlanID string `json:"planID" validate:"required"`
}
