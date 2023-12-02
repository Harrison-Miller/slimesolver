# Slime Dungeon
Slime Dungeon is a puzzle g where you control slimes to solve sokoban like puzzles
The main difference is that when player input is received all slimes are moved simultaneously

## Problem
The simultaneous movement of slimes makes it difficult to determine the order in which the slimes should move
This is because the order in which the slimes move can affect the outcome of the puzzle
To solve this problem we create a graph of events and then solve steps starting from the leaves

## Rules
### Actors
- Can only move in cardinal directions
- Can not move off the grid or into walls
- Block movement unless specified

### Slime
- Can push crates (only 1)
- Die in pits
- If damaged splits into two smaller slimes

### Small Slime
- Can not push crates
- Can not activate switches
- Combine into a Slime if they move into the same square

### Crate
- Block pits
- Can not push other crates

### Switches
- Can be activated by crates and slimes
- When activated open all doors
- Don't block movement

### Doors
- Doors block slimes and crates when closed

### Spikes
- Kill slimes when activated
- Activated every other turn
- Don't block movement

### Pusher
- Push slimes and crates in the direction they are facing when activated
- Activated every other turn
- Don't block movement