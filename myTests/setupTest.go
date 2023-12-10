package myTests

import (
	"encoding/json"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/servers"
	"github.com/NatthawutSK/ri-shop/pkg/databases"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := databases.DbConnect(cfg.Db())

	s := servers.NewSever(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

// CompressToJSON is a function that compresses any object to JSON string. for testing purpose.
func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
