package main

import (
	"crypto/md5"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//go:embed embed/help.txt
var embedHelp string

//go:embed embed/ddl.sql
var embedDDL string

//go:embed embed/views.sql
var embedViews string

//go:embed embed/triggers.sql
var embedTriggers string

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func (conf *Config) PrepareDB(db *sql.DB) (err error) {

	fmt.Println("Creating tables")
	if _, err = db.Exec(embedDDL); err != nil {
		return err
	}
	fmt.Println("Creating views")
	if _, err = db.Exec(embedViews); err != nil {
		return err
	}
	fmt.Println("Creating triggers")
	if _, err = db.Exec(embedTriggers); err != nil {
		return err
	}
	fmt.Println("Inserting default Project")
	query := "insert into Project (id, project_name) values (1, ?);"
	if _, err = db.Exec(query, conf.DefaultProject); err != nil {
		return err
	}

	// Calculate ddl hash
	csDB := fmt.Sprintf("%x", md5.Sum([]byte(embedDDL)))
	fmt.Printf("Updating SysVersion table: %s\n", csDB)

	query = `
        insert into SysVersion (id, cs_db) values (0, ?)
        on conflict do update set cs_db = ?;
        `
	if _, err = db.Exec(query, csDB, csDB); err != nil {
		return err
	}
	fmt.Println("Done")
	return err
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
			dbfile := filepath.Join(conf.DBLocation, conf.DBName)

			prep := false
			if !Exists(dbfile) {
				prep = true
			}

			db, err = sql.Open("sqlite3", dbfile)
			if err != nil {
				return nil, err
			}
			if prep {
				if err = conf.PrepareDB(db); err != nil {
					return nil, err
				}
			}

			if err = CheckDB(db); err != nil {
				return nil, err
			}
			return db, err
		}
	}
	return nil, ErrNoConfig
}
