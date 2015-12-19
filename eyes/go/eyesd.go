package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/tbenz9/cec"
    "time"
    "fmt"
    "os"
    "math/rand"
    "net"
    "strconv"
    "encoding/json"
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
    if debug {fmt.Printf("Database returned ID %v\n", id)}
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

func setupWebsocket(c chan string) {
    if debug{fmt.Printf("Remote Server = %v\n",remoteServerAddress)}
    conn, err := net.Dial("tcp", "192.168.0.106:8080")
    checkErr(err)

    for {
        fmt.Fprintf(conn, <-c)
    }
}


/////////////////////////////////////////////////////////////
//
// Other Functions and Variables
//
/////////////////////////////////////////////////////////////

var debug bool = true
var emulate bool = true
var sleepTime int = 1
func changedState(state int, c chan string) {

    type Device struct {
        Identifier string
        currentState int
    }
    var currentDevice = Device{"dev1", 0}

    id := strconv.FormatInt(insertStateIntoDatabase(state),10)
    if debug {fmt.Printf("ID is: %v\n", id)}
    m, _ := json.Marshal(currentDevice)
    c <- (string(m))
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

    if !emulate {cec.Open("", "cec.go")}
    previousState := 5
    state := 0
    c := make(chan string)

    // Initial Startup tasks
    go setupWebsocket(c)
    go setupLocalDatabase()

    // Check the TV state every second
    for {
        if !emulate {state = cec.GetDevicePowerStatus(0)}
        if emulate {state = rand.Intn(2)}
        if debug {fmt.Printf("The TV is %v\n", state)}
        if state != previousState {
            go changedState(state, c)
        }
        previousState = state
        time.Sleep(1 * time.Second)
        if emulate {time.Sleep(time.Duration(sleepTime) * time.Second)}
    }
}
