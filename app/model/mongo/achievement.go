package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Achievement represents a student's achievement document stored in MongoDB.
type Achievement struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Required references
	StudentID string `bson:"studentId" json:"studentId"`

	// Generic dynamic achievement information
	Title    string                 `bson:"title" json:"title"`
	Type     string                 `bson:"type" json:"type"`           // academic / non-academic / etc.
	Category string                 `bson:"category" json:"category"`   // lomba, sertifikasi, karya tulis, dll.
	Level    string                 `bson:"level" json:"level"`         // lokal, nasional, internasional
	Details  map[string]interface{} `bson:"details" json:"details"`     // flexible dynamic fields
	Tags     []string               `bson:"tags,omitempty" json:"tags"` // optional custom tags

	// Attachments (URLs, file names, metadata)
	Attachments []Attachment `bson:"attachments,omitempty" json:"attachments"`

	// Metadata
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// Attachment represents a single file stored externally (e.g. Cloud Storage / server folder).
type Attachment struct {
	FileName string `bson:"fileName" json:"fileName"`
	URL      string `bson:"url" json:"url"`
	MimeType string `bson:"mimeType" json:"mimeType"`
	Size     int64  `bson:"size" json:"size"` // bytes
}
