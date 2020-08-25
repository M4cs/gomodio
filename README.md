# gomodio
### A wrapper for mod.io in Golang

#### (This is still a work in progress. See below for completion amount)

## Installation

```
go get -u github.com/M4cs/gomodio
```

## Usage

### Basic Usage w/ API Key (Read Only Access)

```go
package main

import (
    "fmt"

    "github.com/M4cs/gomodio"
)

func main() {
    // Create a User
    user := gomodio.NewUser("YOUR_API_KEY", "YOUR_EMAIL")

    // Search for games
    games, err := user.GetGames(map[string]string{"q": "Skater XL"})
    if err != nil {
        fmt.Println(err.Error())
    }
    for _, g := range games.Data {
        fmt.Println("Name:", g.Name)
        fmt.Println("Summary:", g.Summary)
    }
    // Grab Game Object
    game, err := user.GetGame(games.Data[0].ID, nil)
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println(game.ID)
}
```

## Completion

- [X] Authentication (No Third-Party)
- [X] Games
- [X] Mods
- [X] Files
- [X] Subscribe
- [X] Comments
- [X] Media
- [X] Events
- [X] Tags
- [X] Ratings
- [X] Stats
- [X] Metadata
- [ ] Dependencies
- [ ] Teams
- [ ] General
- [ ] Reports
- [ ] Batch
- [ ] Me
