package main

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DBLocation string `yaml:"dblocation"`
	DBName     string `yaml:"dbname"`
	DateFormat string `yaml:"dateformat"`
	GroupName  string `yaml:"defaultgroup"`
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func OpenLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

func (conf *Config) Prepare() (db *sql.DB, err error) {
	configDirs := []string{filepath.Dir("")}

	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		configDirs = append(configDirs, filepath.Join(dir, "", "todo"))
	}
	if dir, ok := os.LookupEnv("HOME"); ok {
		configDirs = append(configDirs, filepath.Join(dir, ".config", "todo"))
	}

	for _, cd := range configDirs {
		configFile := filepath.Join(cd, CONFIG_FILE)
		if Exists(configFile) {
			f, err := os.ReadFile(configFile)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(f, &conf)
			if err != nil {
				return nil, err
			}
			if conf.DBLocation == "" {
				conf.DBLocation = filepath.Dir("")
			}

			db, err = sql.Open("sqlite3", filepath.Join(conf.DBLocation, conf.DBName))
			if err != nil {
				return nil, err
			}
			if err = CheckDB(db); err != nil {
				return nil, err
			}
			return db, err
		}
	}
	return nil, ErrNoConfig
}
