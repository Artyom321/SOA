package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	APIService struct {
		Port int `json:"port"`
	} `json:"apiService"`
	UserService struct {
		Port     int      `json:"port"`
		DBConfig DBConfig `json:"dbConfig"`
	} `json:"userService"`
	PostService struct {
		Port     int      `json:"port"`
		DBConfig DBConfig `json:"dbConfig"`
	} `json:"postService"`
}

type DBConfig struct {
	Port int `json:"port"`
}

func LoadConfig(filename string) Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open config: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return config
}
