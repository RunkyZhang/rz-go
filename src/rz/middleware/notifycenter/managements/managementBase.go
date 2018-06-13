package managements

import (
	"time"
	"rz/middleware/notifycenter/models"
)

type managementBase struct {
}

func (myself *managementBase) setPoBase(poBase *models.PoBase) {
	now := time.Now()
	poBase.CreatedTime = now
	poBase.UpdatedTime = now
	poBase.Deleted = false
}