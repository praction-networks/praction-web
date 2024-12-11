package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInterest struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name" validate:"required,max=100"` // User's full name
	UUID          string             `json:"uuid" bson:"uuid"`                             // Auto-generated UUID, excluded from JSON
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`                   // Auto-generated creation timestamp, excluded from JSON
	VerifiedAt    time.Time          `json:"verifiedAt" bson:"verifiedAt"`                 // Auto-generated verification timestamp, excluded from JSON
	Email         string             `json:"email" bson:"email" validate:"required,email"`
	PinCode       int64              `json:"pincode" bson:"pincode" validate:"required,pincode"` // Valid email address
	Mobile        string             `json:"mobile" bson:"mobile" validate:"required,len=10"`    // 10-digit mobile number
	OTP           int64              `json:"-" bson:"otp"`
	OTPExpireTime time.Time          `json:"-" bson:"otpExpireTime"`
	Address       string             `json:"address" bson:"address" validate:"required,max=255"` // User's address
	IsVerified    bool               `json:"isVerified" bson:"isVerified"`
}
