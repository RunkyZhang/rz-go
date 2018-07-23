package managements

import (
	"time"
	"rz/middleware/notifycenter/models"
)

type managementBase struct {
}

func (myself *managementBase) setPoBase(poBase *models.PoBase) {
	var mixTime time.Time
	if mixTime == poBase.CreatedTime {
		now := time.Now()
		poBase.CreatedTime = now
	}
	poBase.UpdatedTime = poBase.CreatedTime
	poBase.Deleted = false
}
