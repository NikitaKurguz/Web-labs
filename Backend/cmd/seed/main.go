package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"primarte/internal/app/apiserver"
	"primarte/internal/app/model"
	"primarte/internal/app/store"
)

type SeedItem struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Color      string          `json:"color"`
	Category   json.RawMessage `json:"category"`
	Price      string          `json:"price"`
	PriceLabel string          `json:"priceLabel"`
	Img        string          `json:"img"`
}

func main() {
	raw, err := os.ReadFile("../L1-site/data/data.json")
	if err != nil {
		log.Fatalf("read json: %v", err)
	}

	var items []SeedItem
	if err := json.Unmarshal(raw, &items); err != nil {
		log.Fatalf("parse json: %v", err)
	}
	log.Printf("прочитано %d товаров", len(items))

	srvCfg := apiserver.NewConfig()
	if err := srvCfg.Load("configs/config.toml"); err != nil {
		log.Fatalf("server config: %v", err)
	}
	dbCfg := store.NewConfig()
	if err := dbCfg.Load("configs/config.toml"); err != nil {
		log.Fatalf("db config: %v", err)
	}
	st, err := store.New(dbCfg.DataBase.ConnectionString())
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer st.Close()

	for _, item := range items {
		p := &model.Product{
			Name:     item.Name,
			Img:      item.Img,
			Type:     item.Type,
			Color:    item.Color,
			Category: parseCategory(item.Category),
			Price:    parsePrice(item.PriceLabel),
		}
		if err := st.CreateProduct(p); err != nil {
			log.Fatalf("create %s: %v", item.Name, err)
		}
	}
	log.Println("готово")
}

func parseCategory(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		return strings.Join(arr, ",")
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return ""
}

func parsePrice(s string) int {
	s = strings.ReplaceAll(s, "₽", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "Бесплатно") {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
