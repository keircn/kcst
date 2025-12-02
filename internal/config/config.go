package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server    ServerConfig
	Storage   StorageConfig
	Retention RetentionConfig
}

type ServerConfig struct {
	Address string
}

type StorageConfig struct {
	UploadDir string
	DBPath    string
}

type RetentionConfig struct {
	MinTTL          time.Duration
	MaxTTL          time.Duration
	MaxFileSize     int64
	CleanupInterval time.Duration
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Address: ":8080",
		},
		Storage: StorageConfig{
			UploadDir: "./uploads",
			DBPath:    "./data/kcst.db",
		},
		Retention: RetentionConfig{
			MinTTL:          1 * time.Hour,
			MaxTTL:          28 * 24 * time.Hour,
			MaxFileSize:     100 * 1024 * 1024,
			CleanupInterval: 1 * time.Hour,
		},
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	defer file.Close()

	section := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")

		switch section {
		case "server":
			switch key {
			case "address":
				cfg.Server.Address = value
			}
		case "storage":
			switch key {
			case "upload_dir":
				cfg.Storage.UploadDir = value
			case "db_path":
				cfg.Storage.DBPath = value
			}
		case "retention":
			switch key {
			case "min_ttl":
				if d, err := parseDuration(value); err == nil {
					cfg.Retention.MinTTL = d
				}
			case "max_ttl":
				if d, err := parseDuration(value); err == nil {
					cfg.Retention.MaxTTL = d
				}
			case "max_file_size":
				if size, err := parseSize(value); err == nil {
					cfg.Retention.MaxFileSize = size
				}
			case "cleanup_interval":
				if d, err := parseDuration(value); err == nil {
					cfg.Retention.CleanupInterval = d
				}
			}
		}
	}

	return cfg, scanner.Err()
}

func parseDuration(s string) (time.Duration, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	if strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	return time.ParseDuration(s)
}

func parseSize(s string) (int64, error) {
	s = strings.ToUpper(strings.TrimSpace(s))

	multiplier := int64(1)
	if strings.HasSuffix(s, "GIB") {
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "GIB")
	} else if strings.HasSuffix(s, "MIB") {
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "MIB")
	} else if strings.HasSuffix(s, "KIB") {
		multiplier = 1024
		s = strings.TrimSuffix(s, "KIB")
	} else if strings.HasSuffix(s, "GB") {
		multiplier = 1000 * 1000 * 1000
		s = strings.TrimSuffix(s, "GB")
	} else if strings.HasSuffix(s, "MB") {
		multiplier = 1000 * 1000
		s = strings.TrimSuffix(s, "MB")
	} else if strings.HasSuffix(s, "KB") {
		multiplier = 1000
		s = strings.TrimSuffix(s, "KB")
	}

	val, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, err
	}

	return val * multiplier, nil
}
