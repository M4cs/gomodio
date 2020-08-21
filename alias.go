package gomodio

import (
	"github.com/M4cs/gomodio/pkg/files"
	"github.com/M4cs/gomodio/pkg/games"

	"github.com/M4cs/gomodio/pkg/mods"
	"github.com/M4cs/gomodio/pkg/user"
)

// Files Export Response
type Files = files.Files

// File Export Response
type File = files.File

// Games Export Response
type Games = *games.Games

// Game Export Response
type Game = *games.Game

// Mods Export Response
type Mods = mods.Mods

// Mod Export Response
type Mod = mods.Mod

// User Export Type
type User = user.User

var (
	// NewUser Alias
	NewUser = user.NewUser
)

var (
	// GetGames makes a request to search games and returns a Games struct.
	// Parameters:
	// query - query map (map[string]string) |
	// user - your user object (User) |
	GetGames = games.GetGames
	// GetGame makes a request to grab a game and returns a Game struct.
	// Parameters:
	// gameID - game's ID (int) |
	// query - query map (map[string]string) |
	// user - your user object (User) |
	GetGame = games.GetGame
)

var (
	// GetMods makes a request to search mods and returns a Mods struct.
	// Parameters:
	GetMods = mods.GetMods
	// GetMod makes a request to search a mod and returns a Mod struct.
	// Parameters:
	GetMod = mods.GetMod
	// AddMod makes a request to add a mod and returns a Mod struct.
	// Parameters:
	AddMod = mods.AddMod
	// EditMod makes a request to edit a mod and returns a Mod struct.
	EditMod = mods.EditMod
	// DeleteMod makes a request to delete a mod and reutrns an error if any.
	DeleteMod = mods.DeleteMod
)

var (
	// GetModfile makes a request to grab a file and returns a File struct.
	GetModfile = files.GetModfile
	// GetModfiles makes a request to grab a group of modfiles and returns a Files struct.
	GetModfiles = files.GetModfiles
	// AddModfile sends a POST request to upload a mod file.
	AddModfile = files.AddModfile
	// EditModfile sends a PUT request to edit a mod file.
	EditModfile = files.EditModfile
)
