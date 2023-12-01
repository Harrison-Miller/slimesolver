# Slime Dungeon
Slime Dungeon is a puzzle g where you control slimes to solve sokoban like puzzles
The main difference is that when player input is received all slimes are moved simultaneously

## Problem
The simultaneous movement of slimes makes it difficult to determine the order in which the slimes should move
This is because the order in which the slimes move can affect the outcome of the puzzle
To solve this problem we create a graph of events and then solve steps starting from the leaves

## Rules
- Actors can only move in cardinal directions
- Actors can not move off the grid or into walls
- Slimes can push crates
- Slimes die in pits
- Crates block pits
- Crate can not push other crates
- Switches open all doors
- Boxes and slimes can move onto and activate switches
- Doors block slimes and crates when closed