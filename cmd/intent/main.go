package main

import (
    "log"
    "github.com/intentregistry/intent-cli/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        log.Fatal(err)
    }
}