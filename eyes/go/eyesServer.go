package main

import (
    "net"
    "fmt"
)

func printMessages(msgchan <-chan string) {
    for {
        msg := <-msgchan
        fmt.Printf("Database ID: %s\n", msg)
    }
}

func handleConnection(c net.Conn, msgchan chan<- string) {
    fmt.Print("Connection")
    buf := make([]byte, 4096)

    for {
        n, err := c.Read(buf)
        if err != nil || n == 0 {
            c.Close()
            break
        }
        msgchan <- string(buf[0:n])
        n, err = c.Write(buf[0:n])
        if err != nil {
            c.Close()
            break
        }
    }
    fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
}

func main() {
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        // handle error
    }
    
    msgchan := make(chan string)
    go printMessages(msgchan)
    
    for {
        conn, err := ln.Accept()
        if err != nil {
            // handle error
        }
        go handleConnection(conn, msgchan)
    }
}
