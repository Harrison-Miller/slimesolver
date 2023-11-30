package main

import (
	"fmt"
	"log"
	"os"
	"slimesolver/game"
	"strings"
)

func main() {
	// load the level data
	levelData, err := os.ReadFile("level.txt")
	if err != nil {
		log.Fatal(err)
	}

	runGame(string(levelData))
}

const helpText = `w|up, s|down, a|left, d|right, q|quit, r|restart`

func runGame(levelData string) {
	g := game.Game{}
	err := g.Parse(levelData)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println(g.String())
		fmt.Printf("> ")

		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		input = strings.ToLower(input)

		dir := game.Zero
		restart := false
		switch input {
		case "w":
			fallthrough
		case "up":
			dir = game.Up
		case "a":
			fallthrough
		case "left":
			dir = game.Left
		case "s":
			fallthrough
		case "down":
			dir = game.Down
		case "d":
			fallthrough
		case "right":
			dir = game.Right
		case "q":
			fallthrough
		case "quit":
			os.Exit(0)
		case "r":
			fallthrough
		case "restart":
			fallthrough
		case "reset":
			restart = true
		}

		if restart {
			err := g.Parse(levelData)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		if dir != game.Zero {
			g.Move(dir)
		} else {
			fmt.Println(helpText)
		}
	}
}
