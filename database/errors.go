package database

import (
	"github.com/whitesmith/brand-digital-box"
)

// database errors.
const (
	ErrTransaction       = box.Error("failed to start transaction")
	ErrCreateCollection  = box.Error("failed to create collection")
	ErrRecordNotFound    = box.Error("record does not exist")
	ErrCreateRecord      = box.Error("failed to insert record")
	ErrDeleteRecord      = box.Error("failed to delete record")
	ErrIterateCollection = box.Error("failed to iterate over collection")
)
