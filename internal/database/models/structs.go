package models

import "time"

type DefaultBranch struct {
	Name          string `bson:"name" json:"name"`
	CommittedDate string `bson:"committed_date" json:"committed_date"`
	Author        string `bson:"author" json:"author"`
}

type PrimaryLanguage struct {
	Name  string `bson:"name" json:"name"`
	Color string `bson:"color" json:"color"`
}

type RepositoryMetadata struct {
	SyncedAt time.Time `bson:"synced_at" json:"synced_at"`
}
