package managements

import (
	"rz/middleware/notifycenter/global"
)

var (
	IncreasingManagement = increasingManagement{}
)

type increasingManagement struct {
}

// To sharding
func (*increasingManagement) Increase() (int64, error) {
	key := global.RedisKeyMessage + "increasing"

	return global.GetRedisClient().Increment(key, 1)
}
