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
package 

import (
    "github.com/M4cs/gomodio"
)

func main() {
    // Create a User
    user := gomodio.NewUser("YOUR_API_KEY", "YOUR_EMAIL")

    // Search for games
    games, _ := gomodio.GetGames(map[string]string{"q": "Skater XL"}, user)
    for _, g := range games.Data {
        fmt.Println("Name:", g.Name)
        fmt.Println("Summary:", g.Summary)
    }
    // Grab Game Object
    game, _ := gomodio.GetGame(games.Data[0].ID, nil, user)
    fmt.Println(game.ID)
}
```

## Completion

[X] Authentication (No Third-Party)
[X] Games
[X] Mods
[X] Files
[ ] Subscribe
[ ] Comments
[ ] Media
[ ] Events
[ ] Tags
[ ] Ratings
[ ] Stats
[ ] Metadata
[ ] Dependencies
[ ] Teams
[ ] General
[ ] Reports
[ ] Batch
[ ] Me