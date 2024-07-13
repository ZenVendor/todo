package main

import (
	"fmt"

	f "github.com/ZenVendor/todo/internal/functions"
)



func main() {
    var conf f.Config
    conf.ReadConfig()
    db, err := conf.OpenDB() 
    if err != nil {
        fmt.Print("[!]")
        return
    }
    defer db.Close()

    op, od, _ := f.CountPrompt(db)
    prompt := fmt.Sprintf("[%d:%d]", op, od)
    if op + od == -2 {
        fmt.Print("[!]")
    }
    if op == -1 {
        prompt = fmt.Sprintf("[!:%d]", od)
    }
    if od == -1 {
        prompt = fmt.Sprintf("[%d:!]", op)
    }
    fmt.Print(prompt)
}
