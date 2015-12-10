package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/chbmuc/cec"
    "time"
    "fmt"
    "log"
    "os"
)

func changedState(state string) {
    fmt.Print(state)
}

func setupLocalDatabase() {
    err := os.MkdirAll("/opt/eyes",0755)
    db, err := sql.Open("sqlite3", "/opt/eyes/eyes.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS state (metricsID INTEGER PRIMARY KEY, timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, state BOOLEAN NOT NULL)")
    if err != nil {
        log.Printf("%q\n", err)
        return
    }
}

func main() {

    cec.Open("", "cec.go")
    var previousState string = ""

    // Initial Startup tasks
    setupLocalDatabase()

    // Check the TV state every second
    for {
        state := cec.GetDevicePowerStatus(0)
        if state != previousState{
            go changedState(state)
        }
        previousState = state
        time.Sleep(1 * time.Second)    
    }
}
