package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const CONFIG_FN = "todo_config.yml"

type Config struct {
	DBLocation   string `yaml:"dblocation"`
	DBName       string `yaml:"dbname"`
	DateFormat   string `yaml:"dateformat"`
	UseNerdFonts bool   `yaml:"usenerdfonts"`
}

func OpenLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

func (conf *Config) Prepare(local, reset bool) {
	if local {
		log.Printf("LOCAL: Using current dir.\n")
	}
	configDir := filepath.Dir("")

	if !local {
		if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
			configDir = filepath.Join(dir, "todo")
		}
		if dir, ok := os.LookupEnv("HOME"); ok {
			configDir = filepath.Join(dir, ".config", "todo")
		}
	}
	log.Printf("Config dir: %s\n", configDir)

	// default config
	conf.DBLocation = configDir
	conf.DBName = "todo.db"
	conf.DateFormat = "2006-01-02"
	conf.UseNerdFonts = false

	configFile := filepath.Join(configDir, CONFIG_FN)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Printf("Config dir does not exist. Creating directory.\n")
		if err = os.MkdirAll(configDir, 0700); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Config dir exists.")
	}
	if reset {
		log.Printf("RESET: Removing config and db files\n")
		if err := os.Remove(configFile); err != nil {
			log.Println(err)
		}
		if err := os.Remove(filepath.Join(conf.DBLocation, conf.DBName)); err != nil {
			log.Println(err)
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("Config file does not exist. Creating default file.\n")
		writeConf := fmt.Sprintf("dblocation: %s\ndbname: %s\ndateformat: %s", conf.DBLocation, conf.DBName, conf.DateFormat)
		if err := os.WriteFile(configFile, []byte(writeConf), 0700); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Config file exists.")
	}

	if _, err := os.Stat(filepath.Join(conf.DBLocation, conf.DBName)); !os.IsNotExist(err) {
		log.Println("Database file exists.")
	}
	log.Printf("Opening/Creating database file.")
	db, err := conf.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Printf("Checking tables.\n")
	if !TableExists(db) {
		err = CreateTable(db)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Config file and database are ready.\nConfig dir: %s\n", configDir)
}

func (conf *Config) ReadConfig() {
	configDirs := []string{filepath.Dir("")}

	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		configDirs = append(configDirs, filepath.Join(dir, "", "todo"))
	}
	if dir, ok := os.LookupEnv("HOME"); ok {
		configDirs = append(configDirs, filepath.Join(dir, ".config", "todo"))
	}

	for _, cd := range configDirs {
		configFile := filepath.Join(cd, CONFIG_FN)
		if _, err := os.Stat(configFile); err == nil {
			//log.Printf("Found config file in %s, reading config.\n", cd)
			f, err := os.ReadFile(configFile)
			if err != nil {
				log.Fatal(err)
			}
			err = yaml.Unmarshal(f, &conf)
			if err != nil {
				log.Fatal(err)
			}
			if conf.DBLocation == "" {
				conf.DBLocation = filepath.Dir("")
			}

			//log.Printf("Opening database file: %s/%s", conf.DBLocation, conf.DBName)
			db, err := conf.OpenDB()
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			//log.Printf("Checking tables.")
			if !TableExists(db) {
				fmt.Printf("Tables do not exist. Run todo prepare.\n")
				//log.Fatal("Tables do not exist.\n")
			}
			break
		}
	}
}
