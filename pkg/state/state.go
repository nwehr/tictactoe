package state

import "fmt"

func InitialGameState() GameState {
	return GameState{
		"A1": " ",
		"A2": " ",
		"A3": " ",

		"B1": " ",
		"B2": " ",
		"B3": " ",

		"C1": " ",
		"C2": " ",
		"C3": " ",
	}
}

type GameState map[string]string

func (s GameState) Print() {
	state := "   A   B   C \n" +
		fmt.Sprintf("1  %s | %s | %s \n", s["A1"], s["B1"], s["C1"]) +
		"  -----------\n" +
		fmt.Sprintf("2  %s | %s | %s \n", s["A2"], s["B2"], s["C2"]) +
		"  -----------\n" +
		fmt.Sprintf("3  %s | %s | %s \n", s["A3"], s["B3"], s["C3"])

	fmt.Println(state)
}

func (s GameState) Capture(position string) GameState {
	next := GameState{}

	for key, value := range s {
		next[key] = value
	}

	if next[position] == " " {
		next[position] = s.CurrentPlayer()
	}

	return next
}

func (s GameState) CurrentPlayer() string {
	xNumCaptures := 0
	oNumCaptures := 0

	for _, square := range s {
		switch square {
		case "x":
			xNumCaptures++
		case "o":
			oNumCaptures++
		}
	}

	if xNumCaptures > oNumCaptures {
		return "o"
	}

	return "x"
}

func (s GameState) AvailableCaptures() []string {
	availableCaptures := []string{}
	for position, square := range s {
		if square == " " {
			availableCaptures = append(availableCaptures, position)
		}
	}
	return availableCaptures
}
func (s GameState) Winner() (string, bool) {
	for _, row := range s.Rows() {
		first := row[0]
		count := 1

		for _, cell := range row[1:] {
			if cell == first {
				count++
			}
		}

		if first != " " && count == 3 {
			return first, true
		}
	}
	return "", false
}

func (s GameState) Rows() [][]string {
	return [][]string{
		[]string{s["A1"], s["B1"], s["C1"]},
		[]string{s["A2"], s["B2"], s["C2"]},
		[]string{s["A3"], s["B3"], s["C3"]},

		// cols
		[]string{s["A1"], s["A2"], s["A3"]},
		[]string{s["B1"], s["B2"], s["B3"]},
		[]string{s["C1"], s["C2"], s["C3"]},

		// diag
		[]string{s["A1"], s["B2"], s["C3"]},
		[]string{s["A3"], s["B2"], s["C1"]},
	}
}
