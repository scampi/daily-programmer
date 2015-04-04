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

	// the condition of the transition
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

	// the effect of the transition
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

type TuringMachine struct {
	startState  string
	acceptState string
	tape        []rune
	transitions Transitions
}

func NewTuringMachine(fileName string) *TuringMachine {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	tm := &TuringMachine{transitions: Transitions{}}
	var alphabet []rune
	var states []string
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		switch i {
		case 0:
			alphabet = []rune(line + "_")
		case 1:
			states = strings.Split(line, " ")
		case 2:
			tm.startState = line
			if !stringInSlice(tm.startState, states) {
				log.Fatalf("Unknown state: [%s]", tm.startState)
			}
		case 3:
			tm.acceptState = line
			if !stringInSlice(tm.acceptState, states) {
				log.Fatalf("Unknown state: [%s]", tm.acceptState)
			}
		case 4:
			tm.tape = []rune(line)
		default:
			tm.transitions = addTransition(tm.transitions, line, alphabet, states)
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	return tm
}

func main() {
	turingMachine := NewTuringMachine(os.Args[1])

	printTape(turingMachine.startState, 0, string(turingMachine.tape), 0)
	currentCondition, readHead, zero := Condition{state: turingMachine.startState, symbol: turingMachine.tape[0]}, 0, 0
	// Loops until the current state is the accepting state
	for currentCondition.state != turingMachine.acceptState {
		effect := turingMachine.transitions[currentCondition]

		// Apply the step

		// 1. change current state
		currentCondition.state = effect.state
		// 2. move the readHead
		if readHead < 0 {
			turingMachine.tape = append([]rune{effect.symbol}, turingMachine.tape...)
			zero++
		} else if readHead >= len(turingMachine.tape) {
			turingMachine.tape = append(turingMachine.tape, effect.symbol)
		} else {
			turingMachine.tape[readHead] = effect.symbol
		}
		readHead = move(effect, readHead)
		// 3. update the current symbol under the readHead
		if readHead < 0 || readHead >= len(turingMachine.tape) {
			currentCondition.symbol = '_'
		} else {
			currentCondition.symbol = turingMachine.tape[readHead]
		}

		printTape(currentCondition.state, readHead, string(turingMachine.tape), zero)
	}
}
