package warehouse_helper

import (
	"encoding/json"
	"log"
	"os"
)

type WarehouseConfig struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

const filePath = "config/warehouse.json"

//goland:noinspection GoUnusedExportedFunction
func ParseWarehouseConfig() (map[string]WarehouseConfig, error) {
	if data, err := os.ReadFile(filePath); err == nil {
		var warehouses []WarehouseConfig
		if err = json.Unmarshal(data, &warehouses); err != nil {
			log.Printf("[warehouse_helper] Error parsing JSON: %v", err)
			return nil, err
		}

		result := make(map[string]WarehouseConfig)

		for _, warehouse := range warehouses {
			log.Printf("[warehouse_helper] Warehouse Name: %s, IP: %s, Port: %s\n", warehouse.Name, warehouse.IP, warehouse.Port)
			result[warehouse.Name] = warehouse
		}

		return result, nil
	} else {
		log.Printf("[warehouse_helper] Error reading JSON file: %v", err)
		return nil, err
	}
}
