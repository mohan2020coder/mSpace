package models

import "gorm.io/gorm"

// Community groups collections
type Community struct {
	gorm.Model
	Name        string       `json:"name" gorm:"type:text;uniqueIndex"`
	Description string       `json:"description" gorm:"type:text"`
	Collections []Collection `json:"collections"`
}

// Collection groups items
type Collection struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:text"`
	Description string `json:"description" gorm:"type:text"`
	CommunityID uint   `json:"community_id" gorm:"index"`
	Items       []Item `json:"items"`
}

// Item with versioning and workflow
type Item struct {
	gorm.Model
	Title        string     `json:"title" gorm:"type:text"`
	Author       string     `json:"author" gorm:"type:text"`
	Abstract     string     `json:"abstract" gorm:"type:text"`
	Status       string     `json:"status" gorm:"index"` // DRAFT/SUBMITTED/PUBLISHED/REJECTED
	FileURL      string     `json:"file_url" gorm:"type:text"`
	Version      int        `json:"version"`
	CollectionID uint       `json:"collection_id" gorm:"index"`
	Metadata     []Metadata `json:"metadata"`
	Visibility   string     `json:"visibility" gorm:"index"` // PUBLIC/PRIVATE
	FullText     string     `json:"full_text" gorm:"type:text"`

	LegalJSON string `gorm:"type:json"`
}

// Metadata for arbitrary fields
type Metadata struct {
	gorm.Model
	ItemID uint   `json:"item_id" gorm:"index"`
	Key    string `json:"key" gorm:"index"`
	Value  string `json:"value" gorm:"type:text"`
}
