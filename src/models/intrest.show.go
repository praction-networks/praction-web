package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInterest struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name" validate:"required,max=100"`     // User's full name
	UUID          string             `json:"uuid" bson:"uuid"`                                 // Auto-generated UUID
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`                       // Auto-generated creation timestamp
	VerifiedAt    time.Time          `json:"verifiedAt,omitempty" bson:"verifiedAt,omitempty"` // Auto-generated verification timestamp
	IsVerified    bool               `json:"isVerified,omitempty" bson:"isVerified,omitempty"`
	Email         string             `json:"email" bson:"email" validate:"required,email"`
	PinCode       int64              `json:"pincode" bson:"pincode" validate:"required,pincode"`                      // Valid PIN code
	Mobile        string             `json:"mobile" bson:"mobile" validate:"required,len=10"`                         // 10-digit mobile number
	Address       string             `json:"address,omitempty" bson:"address,omitempty" validate:"omitempty,max=255"` // User's address
	InterestStage string             `json:"interestStage,omitempty" bson:"interestStage,omitempty"`
	OTP           int64              `json:"-" bson:"otp"`           // OTP for verification
	OTPExpireTime time.Time          `json:"-" bson:"otpExpireTime"` // OTP expiration time
}

type UserInterestUpdate struct {
	Comments                  []string `json:"comments,omitempty" bson:"comments,omitempty"`                                                                              // Additional user comments                                                                             // Status of verification
	InterestStage             string   `json:"interestStage,omitempty" bson:"interestStage,omitempty" validate:"oneof=Verified FollowUp InstallationScheduled Completed"` // (New, Verified, Follow-up, Installation Scheduled, Completed)
	FollowUpDate              string   `json:"followUpDate,omitempty" bson:"followUpDate,omitempty" validate:"omitempty,datetime=2006-01-02"`                             // Date for next follow-up
	PreferredInstallationDate string   `json:"preferredInstallationDate,omitempty" bson:"preferredInstallationDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	InstallationDate          string   `json:"installationDate,omitempty" bson:"installationDate,omitempty" validate:"omitempty,datetime=2006-01-02"`
	IsInstallationAgreed      bool     `json:"isInstallationAgreed,omitempty" bson:"isInstallationAgreed,omitempty"`                                                     // User agreement to install                                                         // Confirmed installation date
	InstallationStatus        string   `json:"installationStatus,omitempty" bson:"installationStatus,omitempty" validate:"oneof=Pending InProgress Completed Cancelled"` // (Pending, InProgress, Completed, Cancelled)
	InstallationNotes         []string `json:"installationNotes,omitempty" bson:"installationNotes,omitempty"`                                                           // Additional info about installation
	SelectedPlan              string   `json:"selectedPlan,omitempty" bson:"selectedPlan,omitempty"`                                                                     // Plan chosen by user
}
