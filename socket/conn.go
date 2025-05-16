
package socket

import (
    "sync"
    "github.com/gorilla/websocket"
)

var (
    mu   sync.Mutex
    conn *websocket.Conn
)

// Save conn ให้เรียกตอน client websocket connect เข้ามา
func SetConn(c *websocket.Conn) {
    mu.Lock()
    conn = c
    mu.Unlock()
}

// ดึง conn ไปใช้ในที่อื่น
func GetConn() *websocket.Conn {
    mu.Lock()
    defer mu.Unlock()
    return conn
}
