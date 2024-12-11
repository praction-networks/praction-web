package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogCategory represents the structure for blog categories
type BlogCategory struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required,min=2,max=100"`
	UUID        string             `json:"uuid" bson:"uuid"`
	Slug        string             `json:"slug" bson:"slug" validate:"required,slug"`
	Parent      string             `json:"parent,omitempty" bson:"parent,omitempty"` // Optional
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsActive    bool               `json:"-" bson:"isActive"`
	IsDeleted   bool               `json:"-" bson:"isDeleted"`
}

// BlogTag represents the structure for blog tags
type BlogTag struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required,min=2,max=50,oneWord"`
	UUID        string             `json:"uuid" bson:"uuid"`
	Slug        string             `json:"slug" bson:"slug" validate:"required,slug"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsActive    bool               `json:"-" bson:"isActive"`
	IsDeleted   bool               `json:"-" bson:"isDeleted"`
}

// Post represents the structure for a blog post
type Blog struct {
	ID                  primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	BlogTitle           string               `json:"blogTitle" bson:"blogTitle" validate:"required,min=2,max=100"`
	UUID                string               `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Slug                string               `json:"slug" bson:"slug" validate:"required,slug"`
	BlogImage           string               `json:"blogImage,omitempty" bson:"blogImage,omitempty" validate:"omitempty,uuid4"`
	BlogDescription     string               `json:"blogDescription" bson:"blogDescription" validate:"required,min=10,max=7500"`
	BlogDescriptionType string               `json:"blogDescriptionType" bson:"blogDescriptionType" validate:"required,oneof=html text"`
	BlogAuthor          string               `json:"blogAuthor" bson:"blogAuthor" validate:"required,min=2,max=100"`
	Tag                 []string             `json:"tag" bson:"tag" validate:"required,dive,alphanumunicode"`
	Category            []string             `json:"category" bson:"category" validate:"required,dive"`
	MetaDescription     string               `json:"metaDescription,omitempty" bson:"metaDescription,omitempty" validate:"omitempty,max=160"`
	MetaKeywords        []string             `json:"metaKeywords,omitempty" bson:"metaKeywords,omitempty" validate:"omitempty,dive,max=30"`
	EmbeddedMedia       []string             `json:"embeddedMedia,omitempty" bson:"embeddedMedia,omitempty" validate:"omitempty,dive,http_url"`
	Summary             string               `json:"summary,omitempty" bson:"summary,omitempty" validate:"omitempty,max=250"`
	FeatureImage        string               `json:"featureImage,omitempty" bson:"featureImage,omitempty" validate:"omitempty,uuid4"`
	Status              string               `json:"status" bson:"status" validate:"omitempty,oneof=draft published"`
	CreatedAt           time.Time            `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt           time.Time            `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	View                int64                `json:"view" bson:"view" validate:"gte=0"`
	CommentsCount       int64                `json:"commentsCount" bson:"commentsCount" validate:"gte=0"`
	Shares              int64                `json:"shares" bson:"shares" validate:"gte=0"`
	IsApproved          bool                 `json:"isApproved" bson:"isApproved"`
	IsActive            bool                 `json:"-" bson:"isActive"`
	IsDeleted           bool                 `json:"-" bson:"isDeleted"`
	CommentsList        []primitive.ObjectID `json:"commentsList" bson:"commentsList"`
	Comments            []Comments           `json:"comments" bson:"comments"`
}

// Comments represents the structure for comments on a post
type Comments struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid"`
	Name        string             `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Mobile      string             `json:"mobile" bson:"mobile" validate:"required,numeric,len=10"`
	Email       string             `json:"email" bson:"email" validate:"required,email"`
	Description string             `json:"description" bson:"description" validate:"required,min=1,max=250"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	IsActive    bool               `json:"-" bson:"IsActive"`
	IsDeleted   bool               `json:"-" bson:"isDeleted"`
}

type Image struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID      string             `json:"uuid" bson:"uuid"`
	FileID    string             `json:"fileID" bson:"fileID"`
	MimeType  string             `json:"mimeType" bson:"mimeType"`
	FileName  string             `json:"fileName" bson:"FileName"`
	ImageURL  string             `json:"imageURL" bson:"imageURL"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"-" bson:"updated_at"`
	IsActive  bool               `json:"-" bson:"IsActive" validate:"required"`
	IsDeleted bool               `json:"-" bson:"isDeleted"`
}
