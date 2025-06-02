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
	Kafka struct {
		Broker       string `json:"broker"`
		UserTopic    string `json:"userTopic"`
		ViewTopic    string `json:"viewTopic"`
		LikeTopic    string `json:"likeTopic"`
		CommentTopic string `json:"commentTopic"`
	} `json:"kafka"`
	StatsService struct {
        Port int `json:"port"`
    } `json:"statsService"`
}

type DBConfig struct {
	Port int `json:"port"`
}

func LoadConfig(filename string) Config {
	possiblePaths := []string{
		filename,
		"common/config/" + filename,
		"../common/config/" + filename,
		"../../common/config/" + filename,
	}

	for _, path := range possiblePaths {
		config, err := tryLoadConfig(path)
		if err == nil {
			log.Printf("Successfully loaded config from %s", path)
			return config
		}
	}

	log.Printf("Warning: Failed to load config file, using default configuration")
	return getDefaultConfig()
}

func tryLoadConfig(path string) (Config, error) {
	config := Config{}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func getDefaultConfig() Config {
	return Config{
		APIService: struct {
			Port int `json:"port"`
		}{
			Port: 8080,
		},
		UserService: struct {
			Port     int      `json:"port"`
			DBConfig DBConfig `json:"dbConfig"`
		}{
			Port:     8081,
			DBConfig: DBConfig{Port: 5432},
		},
		PostService: struct {
			Port     int      `json:"port"`
			DBConfig DBConfig `json:"dbConfig"`
		}{
			Port:     8082,
			DBConfig: DBConfig{Port: 5432},
		},
		Kafka: struct {
			Broker       string `json:"broker"`
			UserTopic    string `json:"userTopic"`
			ViewTopic    string `json:"viewTopic"`
			LikeTopic    string `json:"likeTopic"`
			CommentTopic string `json:"commentTopic"`
		}{
			Broker:       "localhost:9092",
			UserTopic:    "user_registrations",
			ViewTopic:    "post_views",
			LikeTopic:    "post_likes",
			CommentTopic: "post_comments",
		},
	}
}
