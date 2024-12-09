package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	Title               string             `json:"title" bson:"title" validate:"required,max=100"`
	UUID                string             `json:"-" bson:"uuid"`
	JobID               string             `json:"jobID" bson:"jobID"`
	OpenPosition        int64              `json:"openPosition" bson:"openPosition" validate:"required"`
	Description         string             `json:"description" bson:"description" validate:"required"`
	KeyResponsibilities []string           `json:"keyResponsibilities" bson:"keyResponsibilities" validate:"required,min=1,max=30,dive"`
	Requirements        []string           `json:"requirements" bson:"requirements" validate:"required,min=1,max=30,dive"`
	Department          string             `json:"department,omitempty" bson:"department,omitempty"`                                         // Optional department
	JobType             string             `json:"jobType" bson:"jobType" validate:"required,oneof=Full-Time Part-Time Internship Contract"` // Type of job
	Location            []string           `json:"location" bson:"location" validate:"required"`                                             // Job location
	Remote              bool               `json:"remote,omitempty" bson:"remote,omitempty"`                                                 // Whether the job can be done remotely
	Qualifications      []string           `json:"qualifications" bson:"qualifications" validate:"required,min=1,dive"`                      // Required qualifications
	Experience          string             `json:"experience,omitempty" bson:"experience,omitempty"`                                         // Optional experience details
	MustHaveSkills      []string           `json:"mustHaveSkills,omitempty" bson:"mustHaveSkills" validate:"min=1,max=20,dive"`
	GoodToHaveSkills    []string           `json:"goodToHaveSkills,omitempty" bson:"goodToHaveSkills" validate:"min=1,max=20,dive"`
	Salary              SalaryRange        `json:"salary,omitempty" bson:"salary,omitempty"`                                        // Optional salary range
	PostedBy            string             `json:"postedBy,omitempty" bson:"postedBy,omitempty"`                                    // Name or ID of the person posting
	PostingDate         time.Time          `json:"postingDate,omitempty" bson:"postingDate,omitempty"`                              // Auto-generated posting date
	ApplicationDeadline CustomDate         `json:"applicationDeadline,omitempty" bson:"applicationDeadline,omitempty"`              // Optional application deadline
	Status              string             `json:"status,omitempty" bson:"status,omitempty" validate:"omitempty,oneof=Open Closed"` // Job status
}

// SalaryRange represents the minimum and maximum salary for a job posting
type SalaryRange struct {
	Min      float64 `json:"min" bson:"min" validate:"required,gt=0"`
	Max      float64 `json:"max" bson:"max" validate:"required,gt=0,gtefield=Min"`
	Currency string  `json:"currency" bson:"currency" validate:"oneof=INR USD"`
}

// CustomDate is a type for handling flexible date formats
type CustomDate struct {
	time.Time
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	str := string(data)
	str = str[1 : len(str)-1]

	// Try parsing as a full ISO 8601 datetime
	t, err := time.Parse(time.RFC3339, str)
	if err == nil {
		cd.Time = t
		return nil
	}

	// Try parsing as a date only (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", str)
	if err == nil {
		cd.Time = t
		return nil
	}

	// Return error if no formats matched
	return errors.New("invalid date format")
}
