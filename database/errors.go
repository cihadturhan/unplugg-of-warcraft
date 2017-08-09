package database

import "github.com/whitesmith/unplugg-of-warcraft"

// database errors.
const (
	errDatabaseFailed  = warcraft.Error("failed to start database client")
	errDatabaseIndex   = warcraft.Error("failed to create database index")
	errDatabaseInsert  = warcraft.Error("failed to insert data to the database")
	errDatabaseQuery   = warcraft.Error("failed to query data from the database")
	errDatabaseDelete  = warcraft.Error("failed to delete data from the database")
	errDatabaseUpdate  = warcraft.Error("failed to update data from the database")
	errDatabaseHash    = warcraft.Error("failed to hash user password")
	errDatabaseCounter = warcraft.Error("failed to get a new ID counter")
	errMerge           = warcraft.Error("failed to merge data")
)
