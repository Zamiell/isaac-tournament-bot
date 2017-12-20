package main

import (
    "fmt"
    "github.com/aodin/date"
)

func main() {
    a, err := date.Parse("tuesday")
    if err != nil {
        fmt.Println("error:", err)
    } else {
        fmt.Println(a)
    }

}
