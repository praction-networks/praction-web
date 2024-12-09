package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRefrence struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	ReferedBy     UserType           `json:"referredBy" bson:"referredBy" validate:"required"` // Ensure key matches 'referredBy' in JSON
	Referels      []UserType         `json:"referrels" bson:"referrels" validate:"required,min=1,max=5,uniqueEmailsAndMobiles"`
	RefrelCoupon  string             `json:"refrelCoupon" bson:"refrelCoupon"`
	UUID          string             `json:"uuid" bson:"uuid" validate:"omitempty"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	VerifiedAt    time.Time          `json:"verifiedAt" bson:"verifiedAt"`
	OTP           int64              `json:"_" bson:"otp"`
	OTPExpireTime time.Time          `json:"-" bson:"otpExpireTime"`
	IsVerified    bool               `json:"isVerified" bson:"isVerified"`
}

type UserType struct {
	Name    string `json:"name" bson:"name" validate:"required,max=30"`
	Mobile  string `json:"mobile" bson:"mobile" validate:"required,len=10"`
	Email   string `json:"email" bson:"email" validate:"required,email"`
	PinCode int64  `json:"pincode" bson:"pincode" validate:"required,pincode"`
}
