package repositories

import (
	"rz/core/common"
	"rz/middleware/notifycenter/global"
)

type repositoryBase struct {
	common.RepositoryBase
}

func init() {
	common.SetConnectionStrings(global.GetConfig().Databases)
}