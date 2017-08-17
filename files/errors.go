package files

import "github.com/whitesmith/unplugg-of-warcraft"

const (
	errFailedRead         warcraft.Error = "Failed to read raw file"
	errFailedUnmarshal    warcraft.Error = "Failed to unmarshal binary data"
	errFailedStringToInt  warcraft.Error = "Failed to convert filename to int"
	errFailedRemoveFile   warcraft.Error = "Failed to remove file"
	errFailedDatabaseSave warcraft.Error = "Failed to save file dump to database"
	errFailedDumpFilter   warcraft.Error = "Failed to filter dump file"
)
