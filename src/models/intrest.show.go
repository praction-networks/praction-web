package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInterest struct {
	ID                        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name                      string             `json:"name" bson:"name" validate:"required,max=100"`
	UUID                      string             `json:"uuid" bson:"uuid"`
	CreatedAt                 time.Time          `json:"createdAt" bson:"createdAt"`
	VerifiedAt                *time.Time         `json:"verifiedAt,omitempty" bson:"verifiedAt,omitempty"`
	IsVerified                bool               `json:"isVerified,omitempty" bson:"isVerified,omitempty"`
	Email                     string             `json:"email" bson:"email" validate:"required,email"`
	PinCode                   int64              `json:"pincode" bson:"pincode" validate:"required,len=6,numeric"`
	Mobile                    string             `json:"mobile" bson:"mobile" validate:"required,len=10"`
	Address                   string             `json:"address,omitempty" bson:"address,omitempty" validate:"omitempty,max=255"`
	OTP                       int64              `json:"-" bson:"otp"`
	OTPExpireTime             time.Time          `json:"-" bson:"otpExpireTime"`
	InterestStage             string             `json:"interestStage,omitempty" bson:"interestStage,omitempty"`
	IsInstallationAgreed      bool               `json:"isInstallationAgreed,omitempty" bson:"isInstallationAgreed,omitempty"`
	Comments                  []string           `json:"comments,omitempty" bson:"comments,omitempty"`
	PreferredInstallationDate *string            `json:"preferredInstallationDate,omitempty" bson:"preferredInstallationDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	FollowUpDate              *string            `json:"followUpDate,omitempty" bson:"followUpDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	InstallationDate          *string            `json:"installationDate,omitempty" bson:"installationDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	InstallationStatus        string             `json:"installationStatus,omitempty" bson:"installationStatus,omitempty"`
	InstallationNotes         []string           `json:"installationNotes,omitempty" bson:"installationNotes,omitempty"`
	SelectedPlan              string             `json:"selectedPlan,omitempty" bson:"selectedPlan,omitempty"`
}
