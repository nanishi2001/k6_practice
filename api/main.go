package main

import (
	"log"

	"k6-practice/api/config"
	"k6-practice/api/server"
)

func main() {
	// 設定読み込み
	cfg := config.Load()

	// JWT_SECRETの警告
	if !config.IsJWTSecretSet() {
		log.Println("[WARNING] JWT_SECRET is not set. Using default secret. Set JWT_SECRET in production!")
	}

	// サーバー起動
	srv := server.New(cfg)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
