package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/tbenz9/cec"
    "time"
    "fmt"
    "os"
)



/////////////////////////////////////////////////////////////
//
// Database Functions and Variables
//
/////////////////////////////////////////////////////////////

var localDatabasePath string = "/tmp/eyes/eyesd.db"

func insertStateIntoDatabase(state int) int64 {
    db, err := sql.Open("sqlite3", localDatabasePath)
    checkErr(err)
 
    stmt, err := db.Prepare("INSERT INTO state (state) values(?)")
    checkErr(err)
    res, err := stmt.Exec(state)
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)
    if debug {fmt.Printf("Database returned ID %v", id)}
    db.Close()
    return id
}

func setupLocalDatabase() {
    err := os.MkdirAll("/tmp/eyes",0755)
    db, err := sql.Open("sqlite3", localDatabasePath)
    checkErr(err)

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS state (metricsID INTEGER PRIMARY KEY, timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, state INTEGER NOT NULL)")
    checkErr(err)
    db.Close()
}

/////////////////////////////////////////////////////////////
//
// Websocket Functions
//
/////////////////////////////////////////////////////////////

var remoteServerAddress string = "192.168.0.106"

func setupWebsocket() {
    if debug{fmt.Println("Attempting to establish a websocket")}
    if debug{fmt.Printf("Remote Server = %v",remoteServerAddress)}
}

/////////////////////////////////////////////////////////////
//
// Other Functions and Variables
//
/////////////////////////////////////////////////////////////

var debug bool = false

func changedState(state int) {
    insertStateIntoDatabase(state)
    if debug {fmt.Printf("TV Changed State to %v\n", state)}
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

/////////////////////////////////////////////////////////////
//
// Main
//
/////////////////////////////////////////////////////////////

func main() {

    cec.Open("", "cec.go")
    var previousState int = 5

    // Initial Startup tasks
    setupWebsocket()
    setupLocalDatabase()

    // Check the TV state every second
    for {
        state := cec.GetDevicePowerStatus(0)
        if debug {fmt.Print(state)}
        if state != previousState {
            go changedState(state)
        }
        previousState = state
        time.Sleep(1 * time.Second)    
    }
}
