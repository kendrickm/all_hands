package game

func (room *Room) bfsFloor(start Pos) rune {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	for len(frontier) > 0 {
		current := frontier[0]
		currentTile := room.Map[current.Y][current.X]
		switch currentTile.Rune {
		case ShipFloor:
			return ShipFloor
		default:

		}

		frontier = frontier[1:]

		for _, next := range getNeighbors(room, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}

	return ShipFloor
}
