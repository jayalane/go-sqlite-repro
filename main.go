package main

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	globals "github.com/jayalane/go-globals"
	"os"
	"time"
)

const perfTune = `
pragma journal_mode = WAL;
pragma synchronous = normal;
pragma temp_store = memory;
pragma mmap_size = 30000000000;`

var dbNameTemplate = "persist_%d_%s.sqlite"
var dbPath = "./dbs/"

func main() {

	g := globals.NewGlobal("docker = true\nlogStdout = true\ndebugLevel = network\n", true)

	err := os.Mkdir(dbPath, 0770)

	if err != nil && !os.IsExist(err) {
		s := fmt.Sprintf("Can't make state DB ./dbs/ for SQLITE %s", err)
		panic(s)
	}
	dbName := fmt.Sprintf(dbNameTemplate, 1, "a")

	db, err := sql.Open("sqlite", dbPath+dbName)
	if err != nil {
		g.Ml.La("Error opening", dbPath+dbName, err)
		db = nil
		return
	}
	g.Ml.Ls("Opened DB", dbPath+dbName)
	_, err = db.Exec(perfTune)
	if err != nil {
		g.Ml.La("Error tuning", err)
		return
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS ips (
  key TEXT,
  value INT,
  ttl INT)`)
	if err != nil {
		g.Ml.La("Error creating if not", err)
		return
	}
	g.Ml.Ls("Opened, tuned and made table DB", dbPath+dbName)
	defer db.Close()
	time.Sleep(3 * time.Second)
}
