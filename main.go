package main

import (
	"os"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/servers"
	"github.com/NatthawutSK/ri-shop/pkg/databases"
)


func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg	:= config.LoadConfig(envPath())
	
	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	// fmt.Println(db)

	servers.NewSever(cfg, db).Start()
}