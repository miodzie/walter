package seras

import "gorm.io/gorm"

// Modules need some type of interface for storing data.
// Pass them a separated database? E.g. each module gets their own
// sqlite database.
// I don't want to tie it directly to Gorm.db, sql.DB is better, but
// having an interface is the best way to not couple seras to an
// implementation. (What if, I want my storage to be in ElasticSearch,
// or just a text file?)

// I could have a simple key-store setup...

type HasDatabase interface {
	setDB(*gorm.DB)
	DB() *gorm.DB
}
