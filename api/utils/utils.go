package utils

import (
	"encoding/gob"
	"plefi/api/models"
)

// Register concrete types stored in session (for securecookie/gob)
func init() {
	gob.Register(&models.PlexAuth{})
	gob.Register(&models.UserInfo{})
}
