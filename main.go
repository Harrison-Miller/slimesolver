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
		case "w", "up":
			dir = game.Up
		case "a", "left":
			dir = game.Left
		case "s", "down":
			dir = game.Down
		case "d", "right":
			dir = game.Right
		case "q", "exit", "quit":
			os.Exit(0)
		case "r", "restart", "reset":
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
