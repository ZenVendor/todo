package main

import (
	"fmt"
	"os"

	"github.com/ZenVendor/todo/internal/functions"
)



func main() {
    
    home, ok := os.LookupEnv("HOME")
    if !ok {
        fmt.Print("[!]")
        return
    }
    configFile := fmt.Sprintf("%s/.config/todo/todo.yml", home)
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        fmt.Print("[!]")
        return
    }
    var conf functions.Config
    if err := conf.ReadConfig(configFile); err != nil {
        fmt.Print("[!]")
        return
    }
    db, err := conf.OpenDB() 
    if err != nil {
        fmt.Print("[!]")
        return
    }
    defer db.Close()

    op, od, _ := functions.CountPrompt(db)
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
