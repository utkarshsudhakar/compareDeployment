package main

import (
    "fmt"
    
    "log"
    "io/ioutil"
)

func main() {
    
    files, err := ioutil.ReadDir("C:\\HAWK\\Repo\\ccgf-qastaging-hawk-config")
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        fmt.Println(file.Name())
    }
}