package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Condition is the input of a Transition
type Condition struct {
	state  string
	symbol rune
}

func (c Condition) String() string {
	return "state=" + c.state + " symbol=" + strconv.QuoteRuneToASCII(c.symbol)
}

// Effect is the output of a Transition
type Effect struct {
	state     string
	symbol    rune
	direction string
}

func (c Effect) String() string {
	return "state=" + c.state + " symbol=" + strconv.QuoteRuneToASCII(c.symbol) + " direction=" + c.direction
}

// The set of transition functions
type Transitions map[Condition]Effect

func addTransition(transitions Transitions, line string, alphabet []rune, states []string) Transitions {
	parts := strings.Split(line, " ")

	cond := Condition{state: parts[0], symbol: []rune(parts[1])[0]}
	if !stringInSlice(cond.state, states) {
		log.Fatalf("Unknown state: [%s]", cond.state)
	}
	if !runeInSlice(cond.symbol, alphabet) {
		log.Fatalf("Unknown symbol! Got [%s]", cond.symbol)
	}
	if _, ko := transitions[cond]; ko {
		log.Fatalf("Duplicate transition! Got [%s] condition already", cond)
	}

	effect := Effect{state: parts[3], symbol: []rune(parts[4])[0], direction: parts[5]}
	if !stringInSlice(effect.state, states) {
		log.Fatalf("Unknown state: [%s]", effect.state)
	}
	if !runeInSlice(effect.symbol, alphabet) {
		log.Fatalf("Unknown symbol! Got [%s]", effect.symbol)
	}

	transitions[cond] = effect
	return transitions
}

// Move moves the readHead in the direction specified in the effect
func move(effect Effect, readHead int) int {
	switch effect.direction {
	case "<":
		return readHead - 1
	case ">":
		return readHead + 1
	default:
		log.Fatal("Bad direction")
	}
	return 0
}

// PrintTape prints to the standard output the current state of the Turing Machine. Zero is the number of padding whitespaces to prepend the symbol '|'.
func printTape(state string, readHead int, tape string, zero int) {
	rh := strings.Repeat(" ", zero) + "|"
	if readHead > 0 {
		rh = rh + strings.Repeat(" ", readHead-1) + "^"
	} else if readHead < 0 {
		rh = "^" + strings.Repeat(" ", -readHead-1) + rh
		tape = strings.Repeat(" ", -readHead) + tape
	}
	fmt.Printf("%s\n%s\n%s\n\n", state, tape, rh)
}

// RuneInSlice returns true if the rune a is in the slice list
func runeInSlice(a rune, list []rune) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// StringInSlice returns true if the string a is in the slice list
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to open file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var alphabet []rune
	var states []string
	startState := ""
	acceptState := ""
	var tape []rune
	transitions := Transitions{}
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		switch i {
		case 0:
			alphabet = []rune(line + "_")
		case 1:
			states = strings.Split(line, " ")
		case 2:
			startState = line
			if !stringInSlice(startState, states) {
				log.Fatalf("Unknown state: [%s]", startState)
			}
		case 3:
			acceptState = line
			if !stringInSlice(acceptState, states) {
				log.Fatalf("Unknown state: [%s]", acceptState)
			}
		case 4:
			tape = []rune(line)
		default:
			transitions = addTransition(transitions, line, alphabet, states)
		}
	}

	printTape(startState, 0, string(tape), 0)
	currentCondition, readHead, zero := Condition{state: startState, symbol: tape[0]}, 0, 0
	// Loops until the current state is the accepting state
	for currentCondition.state != acceptState {
		effect := transitions[currentCondition]
		// Apply the step
		currentCondition.state = effect.state
		if readHead < 0 {
			tape = append([]rune{effect.symbol}, tape...)
			zero++
		} else if readHead >= len(tape) {
			tape = append(tape, effect.symbol)
		} else {
			tape[readHead] = effect.symbol
		}
		readHead = move(effect, readHead)
		if readHead >= len(tape) {
			currentCondition.symbol = '_'
		} else if readHead < 0 {
			currentCondition.symbol = '_'
		} else {
			currentCondition.symbol = tape[readHead]
		}

		printTape(currentCondition.state, readHead, string(tape), zero)
	}
}
