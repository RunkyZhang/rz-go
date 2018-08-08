package healths

import (
	"rz/core/common"
	"rz/middleware/notifycenter/global"
)

type RuntimeHealthIndicator struct {
}

func (myself *RuntimeHealthIndicator) Indicate() (*common.HealthReport) {
	healthReport := &common.HealthReport{
		Ok:     true,
		Name:   "Runtime",
		Type:   "Runtime",
		Level:  1,
		Detail: make(map[string]interface{}),
	}

	healthReport.Detail["Version"] = global.Version

	return healthReport
}
