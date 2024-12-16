package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
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

type JobApplication struct {
	ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`                                                           // Unique ID for the application
	JobID        primitive.ObjectID  `json:"jobID" bson:"jobID" validate:"required"`                                            // Reference to the Job ID
	UserDetails  ApplicantDetails    `json:"userDetails" bson:"userDetails" validate:"required"`                                // User's personal and contact details
	Resume       string              `json:"resume" bson:"resume" validate:"required,url"`                                      // URL to the resume
	CoverLetter  string              `json:"coverLetter,omitempty" bson:"coverLetter,omitempty"`                                // Optional cover letter
	Education    []EducationDetails  `json:"education" bson:"education" validate:"required,min=1,dive"`                         // Education history
	Experience   []ExperienceDetails `json:"experience" bson:"experience,omitempty" validate:"omitempty,dive"`                  // Work experience
	AppliedDate  time.Time           `json:"appliedDate" bson:"appliedDate"`                                                    // Date of application submission
	Status       string              `json:"status" bson:"status" validate:"required,oneof=Pending Shortlisted Rejected Hired"` // Application status
	Notes        string              `json:"notes,omitempty" bson:"notes,omitempty" validate:"omitempty,max=500"`               // Admin notes
	ReferralCode string              `json:"referralCode,omitempty" bson:"referralCode,omitempty" validate:"omitempty,max=50"`  // Optional referral code
}

// ApplicantDetails stores the user's personal and contact details
type ApplicantDetails struct {
	Name          string    `json:"name" bson:"name" validate:"required,max=100"`                                                                      // Applicant's name
	Email         string    `json:"email" bson:"email" validate:"required,email"`                                                                      // Applicant's email
	Phone         string    `json:"phone,omitempty" bson:"phone,omitempty" validate:"omitempty,e164"`                                                  // Optional phone number
	DateOfBirth   time.Time `json:"dateOfBirth" bson:"dateOfBirth" validate:"required"`                                                                // Date of birth
	MaritalStatus string    `json:"maritalStatus,omitempty" bson:"maritalStatus,omitempty" validate:"omitempty,oneof=Single Married Divorced Widowed"` // Marital status
	Address       string    `json:"address,omitempty" bson:"address,omitempty" validate:"omitempty,max=200"`
	LinkedIn      string    `json:"linkedIn,omitempty" bson:"linkedIn,omitempty" validate:"omitempty,url"` // LinkedIn profile URL                                          // Address
}

// EducationDetails stores details of a user's educational qualifications
type EducationDetails struct {
	Institution string     `json:"institution" bson:"institution" validate:"required,max=100"`                         // Institution name
	Degree      string     `json:"degree" bson:"degree" validate:"required,max=100"`                                   // Degree name (e.g., B.Tech, MBA)
	Field       string     `json:"field" bson:"field" validate:"required,max=100"`                                     // Field of study (e.g., Computer Science)
	StartDate   CustomDate `json:"startDate,omitempty" bson:"startDate,omitempty" validate:"required"`                 // Start date
	EndDate     CustomDate `json:"endDate,omitempty" bson:"endDate,omitempty" validate:"omitempty,gtefield=StartDate"` // End date
}

// ExperienceDetails stores details of a user's work experience
type ExperienceDetails struct {
	CompanyName       string     `json:"companyName" bson:"companyName" validate:"required,max=100"`                         // Name of the company
	IsCurrentEmployer bool       `json:"isCurrentEmployer" bson:"isCurrentEmployer" validate:"required,bool"`                // Is the current employer
	Role              string     `json:"role" bson:"role" validate:"required,max=100"`                                       // Role (e.g., Software Engineer)
	Responsibilities  string     `json:"responsibilities" bson:"responsibilities" validate:"required,max=100"`               // Job responsibilities
	StartDate         CustomDate `json:"startDate,omitempty" bson:"startDate,omitempty" validate:"required"`                 // Start date
	EndDate           CustomDate `json:"endDate,omitempty" bson:"endDate,omitempty" validate:"omitempty,gtefield=StartDate"` // End date
	Description       string     `json:"description,omitempty" bson:"description,omitempty" validate:"omitempty,max=500"`    // Optional job description
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
