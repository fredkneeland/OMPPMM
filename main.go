package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

func main() {
	startTime := time.Now()
	args := os.Args[1:]

	p, err := ioutil.ReadFile(args[0])

	if err != nil {
		panic("Failed to read pattern file")
	}

	t, err := ioutil.ReadFile(args[1])

	if err != nil {
		panic("Failed to read text file")
	}

	fmt.Println("Parsing the files took: ", time.Since(startTime))

	indexes := Algorithm3Mutations(string(p), string(t))
	// indexes := Algorithm3(string(p), string(t))

	fmt.Printf("Indexes: %v", len(indexes))
}

//////////////////////////////////////////////////
//												//
// 				Variable Pattern				//
//												//
//////////////////////////////////////////////////

// Pattern is a struct for the variable pattern created as the first step in finding the palindromic structure of DNA
type Pattern struct {
	// indicates all of the
	Alpha           [][]int
	Beta            []int
	pattern         string
	currentAlphabet []string
	currentIndex    int

	// used for the mutation algorithm...
	mutations []bool
	letters   []string
}

// NewPattern creates a new variable pattern to be used in calculating palindromes
func NewPattern(p string) *Pattern {
	pattern := &Pattern{}

	pattern.pattern = p
	pattern.Alpha = make([][]int, 0)
	pattern.Beta = make([]int, len(p))
	pattern.currentAlphabet = make([]string, len(p))
	return pattern
}

// NewPatternMutation creates a new variable pattern to be used in calculating palindromes
func NewPatternMutation(p string) *Pattern {
	pattern := &Pattern{}

	pattern.pattern = p
	pattern.Alpha = make([][]int, 0)
	pattern.Beta = make([]int, len(p))
	pattern.currentAlphabet = make([]string, len(p))
	pattern.mutations = make([]bool, len(p))
	pattern.letters = make([]string, len(p))
	return pattern
}

// Print out the variable pattern in a human readable way
func (p *Pattern) Print() {
	fmt.Println("Variable Pattern: ")
	fmt.Printf("Beta: %d \n", len(p.Beta))
	for _, val := range p.Beta {
		fmt.Printf("variable: %d \n", val)
	}

	fmt.Printf("Alpha: %d \n", len(p.Alpha))
	for i, list := range p.Alpha {
		fmt.Printf("%d: ", i)
		for _, j := range list {
			fmt.Printf("%d, ", j)
		}
		fmt.Printf("\n")
	}
}

func (p *Pattern) addNewLetter(letter string) bool {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("\n Recovered in f", r)
	// 		fmt.Printf(" currentIndex: %d len Beta: %d", p.currentIndex, len(p.Beta))
	// 		fmt.Printf("\n current alphabet: %d beta value: %d\n", len(p.currentAlphabet), p.Beta[p.currentIndex-1])
	// 	}
	// }()

	// fmt.Printf("currentIndex: %d %d", p.currentIndex, p.Beta[p.currentIndex]-1)
	// fmt.Println("currentIndex addNew: ", p.currentIndex)
	if len(p.currentAlphabet[p.Beta[p.currentIndex]-1]) != 0 {
		if letter == p.currentAlphabet[p.Beta[p.currentIndex]-1] {
			p.letters[p.currentIndex] = letter
			p.currentIndex++
			return true
		}
		return false
	}

	for _, index := range p.Alpha[p.Beta[p.currentIndex]-1] {
		// fmt.Printf("\nindex: %d\n", index)
		if index-1 > p.currentIndex {
			break
		}

		if p.currentAlphabet[index-1] == letter {
			return false
		}
	}

	// fmt.Printf("\n set current Alphabet to letter: %d \n", p.Beta[p.currentIndex]-1)
	p.currentAlphabet[p.Beta[p.currentIndex]-1] = letter
	p.letters[p.currentIndex] = letter
	p.currentIndex++
	return true
}

// func (p *Pattern) otherContradiction() {
// 	for _, index := range p.Alpha {

// 	}
// }

// NoContradictions makes sure that the current letter doesn't violate any existing rules
func (p *Pattern) NoContradictions(letter string) bool {
	for _, index := range p.Alpha[p.Beta[p.currentIndex]-1] {
		// fmt.Printf("\nindex: %d\n", index)
		if index-1 > p.currentIndex {
			break
		}

		if p.currentAlphabet[index-1] == letter {
			return false
		}
	}

	return true
}

func (p *Pattern) addNewLetterMutation(letter string) bool {
	if len(p.currentAlphabet[p.Beta[p.currentIndex]-1]) != 0 {
		curr := p.currentAlphabet[p.Beta[p.currentIndex]-1]
		if letter == curr {
			p.letters[p.currentIndex] = letter
			p.currentIndex++
			return true
		} else if len(curr) > 1 {
			for _, l := range curr {
				if string(l) == letter && p.NoContradictions(letter) { // we also have to make a check to see if there are any incompatibilities...
					p.letters[p.currentIndex] = letter
					p.currentAlphabet[p.Beta[p.currentIndex]-1] = letter
					p.currentIndex++
					return true
				}
			}
		}

		return false
	}

	if !p.NoContradictions(letter) {
		return false
	}

	// fmt.Printf("\n set current Alphabet to letter: %d \n", p.Beta[p.currentIndex]-1)
	p.currentAlphabet[p.Beta[p.currentIndex]-1] = letter
	p.letters[p.currentIndex] = letter
	p.currentIndex++
	return true
}

func (p *Pattern) resetToPosition(index int) {
	for i := index; i < len(p.currentAlphabet); i++ {
		p.currentAlphabet[i] = ""
	}

	p.currentIndex = index
}

// Algorithm1 defines a variable pattern on p to be used in calculating the palindromic structure
func (p *Pattern) Algorithm1() {
	pals := Pals(p.pattern)

	c := float32(0.5)
	l := float32(0)
	s := 0

	for i := 1; i <= len(p.pattern); i++ {
		if l >= float32(i) {
			p.Beta[i-1] = p.Beta[int(2*c)-i-1]
		} else {
			s++
			p.Beta[i-1] = s
			p.Alpha = append(p.Alpha, make([]int, 0))

			for j := range pals {
				center := pals[j][0]
				radius := pals[j][1]
				if center > float32(i) {
					break
				}

				if (center+radius-0.5) == (float32(i)-1) && center-radius-0.5 >= 0 {
					lprime := int(center - radius - 0.5)

					if lprime == 0 || p.Beta[i-1] == p.Beta[lprime-1] {
						continue
					}

					p.Alpha[p.Beta[i-1]-1] = append(p.Alpha[p.Beta[i-1]-1], p.Beta[lprime-1])
					p.Alpha[p.Beta[lprime-1]-1] = append(p.Alpha[p.Beta[lprime-1]-1], p.Beta[i-1])
				}
			}
		}

		if l <= float32(i) {
			index := -1
			lowest := float32(9999999999.9)

			for j := range pals {
				cp := pals[j][0]
				rp := pals[j][1]

				if cp-rp+0.5 <= float32(i) && cp+rp-0.5 >= float32(i)+1 && cp+rp-0.5 < lowest {
					lowest = cp + rp - 0.5
					index = j
				}
			}

			if index != -1 && pals[index][0] <= float32(i)+0.5 {
				l = pals[index][0] + pals[index][1] - 0.5
				c = pals[index][0]
			} else if l == float32(i-1) {
				l++
			}
		}
	}

	for i := range p.Alpha {
		sort.Ints(p.Alpha[i])
	}
}

// Pals returns the palindromic structure of a string in the form {center, radius} from 1 to the end of the string in 0.5 increments to account for
// both odd and even palindromes
func Pals(p string) [][]float32 {
	pals := make([][]float32, len(p)*2-1)

	for i := range pals {
		pals[i] = make([]float32, 2)
	}

	for i := 0; i < len(pals); i++ {
		center := float32(i+2) / 2
		length := float32(0.0)

		// looking for an odd palindrome
		if i%2 == 0 {
			length = 1
			start := int(center) - 2
			stop := int(center)
			for start >= 0 && stop < len(p) && p[start] == p[stop] {
				length += 2
				start--
				stop++
			}
		} else { // even palindrome
			start := int(center) - 1
			stop := start + 1
			for start >= 0 && stop < len(p) && p[start] == p[stop] {
				length += 2
				start--
				stop++
			}
		}

		pals[i][0] = center
		pals[i][1] = length / 2 // find the radius by dividing the length in half
	}

	return pals
}

// EqualPals determines if two pals are equivalent or not
func EqualPals(p1, p2 [][]float32) bool {
	if len(p1) != len(p2) {
		return false
	}

	for i := range p1 {
		for j := range p1[i] {
			if p1[i][j] != p2[i][j] {
				return false
			}
		}
	}

	return true
}

//////////////////////////////////////////////////
//												//
// 				Finite Automaton				//
//												//
//////////////////////////////////////////////////

// FiniteAutomaton is an automaton that excepts a single letter and progresses through a stream of dna data
type FiniteAutomaton struct {
	p  string
	vP *Pattern

	currentState       int
	currentAssignments []string

	states []*FAState

	// values for handling mutations
	mutationStates []*FAState
	state          *FAState
}

// NewFiniteAutomaton creates a FiniteAutomaton struct
func NewFiniteAutomaton(p string) *FiniteAutomaton {
	vP := NewPattern(p)
	vP.Algorithm1()
	// vP.Print()
	return &FiniteAutomaton{
		p:            p,
		vP:           vP,
		currentState: 0,
	}
}

// NewFiniteAutomatonMutation constructs a finite automaton that handles mutations
func NewFiniteAutomatonMutation(p string) *FiniteAutomaton {
	vP := NewPatternMutation(p)
	vP.Algorithm1()
	vP.Print()
	return &FiniteAutomaton{
		p:            p,
		vP:           vP,
		currentState: 0,
	}
}

// Print prints out the automaton
func (fA *FiniteAutomaton) Print() {
	fmt.Println("Finite Automaton: ")
	fmt.Printf(" pattern Len: %d, statesLen: %d \n", len(fA.p), len(fA.states))

	for _, state := range fA.states {
		fmt.Printf(" pos: %d, final: %v \n", state.position, state.FinalState)
		if state.PreviousState != nil {
			fmt.Printf("   prev: %d %v\n", state.PreviousState.position, state.PreviousState.isMutant)
		}
	}

	fmt.Println("\n Mutant states: ")
	for _, state := range fA.mutationStates {
		fmt.Println("pos: ", (state.position), " prev: ", (state.PreviousState.position), " prev non-mutant: ", (state.PreviousNonMutationSate.position))
	}
}

// NextLetter moves the automaton one letter
func (fA *FiniteAutomaton) NextLetter(letter string) bool {
	// fmt.Println("Next Letter: " + letter)
	if fA.vP.addNewLetter(letter) {
		fA.currentState++

		if fA.states[fA.currentState].FinalState {
			fA.states[fA.currentState].transFunc()
			fA.currentState = fA.states[fA.currentState].PreviousState.position
			return true
		}
		return false
	}

	if fA.currentState == 0 {
		fA.vP.resetToPosition(0)
		return fA.NextLetter(letter)
	}

	// we failed to transition so back up and try again
	fA.states[fA.currentState].transFunc()
	fA.currentState = fA.states[fA.currentState].PreviousState.position
	return fA.NextLetter(letter)
}

// NoMutants checks to see if there are no mutations in the array
func (fA *FiniteAutomaton) NoMutants() bool {
	for _, mutant := range fA.vP.mutations {
		if mutant {
			return false
		}
	}

	return true
}

// NextLetterMutation moves the automaton one letter while allowing mutations
func (fA *FiniteAutomaton) NextLetterMutation(letter string) bool {
	// fmt.Println("currentState: ", fA.currentState)
	if fA.vP.addNewLetterMutation(letter) {
		fA.currentState++

		if fA.state.isMutant {
			fA.state = fA.mutationStates[fA.currentState-1]
		} else {
			fA.state = fA.states[fA.currentState]
		}

		if fA.state.FinalState {
			fmt.Println("currentAlphabet: ", fA.vP.currentAlphabet, " letters: ", fA.vP.letters)
			// fmt.Println("FINAL STATE!!!!")
			fA.state.transFunc()
			if !fA.state.isMutant {
				fA.state = fA.mutationStates[len(fA.mutationStates)-1].PreviousNonMutationSate
			} else if fA.NoMutants() {
				fA.state = fA.state.PreviousNonMutationSate
			} else {
				fA.state = fA.state.PreviousState
			}
			fA.currentState = fA.state.position

			return true
		}
		return false
	} else if !fA.state.isMutant && fA.state.position != 0 {
		// fmt.Println("Pos: ", fA.state.position)
		if fA.state.PreviousState.FinalState {
			fA.vP.mutations[len(fA.vP.mutations)-1] = true

			fA.state = fA.mutationStates[len(fA.mutationStates)-1]
			fA.state.transFunc()
			fA.state = fA.state.PreviousState
			fA.currentState = fA.state.position
			// fmt.Println("currentIndex: ", fA.vP.currentIndex)
			// fA.vP.currentIndex = fA.state.position-1
			return true
		}

		// we want to add the possible mutation... and then verify the alphabet later...
		fA.vP.currentAlphabet[fA.vP.Beta[fA.vP.currentIndex]-1] += letter

		fA.vP.currentIndex++
		// fmt.Println("currentIndex3: ", fA.vP.currentIndex)
		fA.vP.mutations[fA.state.position-1] = true
		// fmt.Println("nextPosition: ", fA.state, " pos: ", fA.state.position, " len: ", len(fA.vP.mutations))
		// fmt.Println("CurrentState: ", fA.state.position, " mutant: ", fA.state.isMutant)
		fA.state = fA.state.PreviousState // this will move us to the next mutation state
		// fmt.Println("IsMutant: ", fA.state.isMutant, " idx: ", fA.state.position)
		fA.currentState = fA.state.position
		// fmt.Println("nextPosition: ", fA.state, " pos: ", fA.state.position, " len: ", len(fA.vP.mutations))
		return false
	}

	if fA.currentState == 0 {
		fA.vP.resetToPosition(0)
		return fA.NextLetterMutation(letter)
	}

	// we failed to transition so back up and try again
	fA.state.transFunc()

	if fA.NoMutants() {
		fA.state = fA.state.PreviousNonMutationSate
	} else {
		fA.state = fA.state.PreviousState
	}
	fA.currentState = fA.state.position
	return fA.NextLetterMutation(letter)
}

func (fA *FiniteAutomaton) createTransFunc(index, i int) func() {
	return func() {
		used := make([]bool, len(fA.vP.currentAlphabet))
		for j := index; j < i; j++ {
			fA.vP.currentAlphabet[fA.vP.Beta[j-index]-1] = fA.vP.currentAlphabet[fA.vP.Beta[j]-1]
			used[fA.vP.Beta[j-index]-1] = true
		}

		// clear out the old alphabet naively
		for j := index; j < i; j++ {
			if used[fA.vP.Beta[j]-1] {
				continue
			}

			fA.vP.currentAlphabet[fA.vP.Beta[j]-1] = ""
		}

		fA.vP.currentIndex = index
		// fmt.Printf("Reset to index: %d", index)
	}
}

func (fA *FiniteAutomaton) createMutationFunc(i int) func() {
	return func() {
		fA.vP.mutations[i] = true
	}
}

func (fA *FiniteAutomaton) createTransFuncMutation(index, i int) func() {
	return func() {
		used := make([]bool, len(fA.vP.currentAlphabet))
		for j := index; j < i; j++ {
			idx := fA.vP.Beta[j-index] - 1
			idx2 := fA.vP.Beta[j] - 1
			fA.vP.currentAlphabet[idx] = fA.vP.currentAlphabet[idx2]
			used[idx] = true
			fA.vP.mutations[idx] = fA.vP.mutations[idx2]
		}

		// clear out the old alphabet naively
		for j := index; j < i; j++ {
			fA.vP.mutations[j] = false
			if used[fA.vP.Beta[j]-1] {
				continue
			}

			fA.vP.currentAlphabet[fA.vP.Beta[j]-1] = ""
		}

		// fmt.Println("currentIndex mutation: ", index)
		fA.vP.currentIndex = index
		// fmt.Printf("Reset to index: %d", index)
	}
}

// Algorithm2 creates a basic finite Automaton for the given pattern which does NOT account for mutations
func (fA *FiniteAutomaton) Algorithm2() {
	fA.states = make([]*FAState, len(fA.vP.Beta)+1)
	initialState := &FAState{
		position:   0,
		FinalState: false,
	}

	fA.states[0] = initialState

	for i := 0; i < len(fA.p); i++ {
		currentState := &FAState{
			position:   i + 1,
			FinalState: i == len(fA.p)-1,
		}

		if i < 2 {
			currentState.PreviousState = fA.states[i]
			currentState.transFunc = fA.createTransFunc(i, i)
		} else {
			index := 0
			for j := 1; j < i; j++ {
				pal1 := Pals(fA.p[:i-j+1])
				pal2 := Pals(fA.p[j : i+1])

				if EqualPals(pal1, pal2) {
					index = i - j
					break
				}
			}

			currentState.PreviousState = fA.states[index+1]
			currentState.transFunc = fA.createTransFunc(index, i)
		}

		fA.states[i+1] = currentState
	}
}

// Algorithm2Mutation creates a state machine that handles mutations
func (fA *FiniteAutomaton) Algorithm2Mutation() {
	fA.states = make([]*FAState, len(fA.vP.Beta)+1)
	fA.mutationStates = make([]*FAState, len(fA.vP.Beta))
	initialState := &FAState{
		position:   0,
		FinalState: false,
		isMutant:   false,
	}

	fA.state = initialState

	fA.states[0] = initialState

	for i := 0; i < len(fA.p); i++ {
		currentState := &FAState{
			position:   i + 1,
			FinalState: i == len(fA.p)-1,
			isMutant:   false,
		}

		currentMutationState := &FAState{
			position:   i + 1,
			FinalState: i == len(fA.p)-1,
			isMutant:   true,
		}

		if i == 0 {
			currentState.PreviousState = currentMutationState
			currentMutationState.PreviousState = currentState
			currentMutationState.PreviousNonMutationSate = currentState
			currentMutationState.transFunc = fA.createTransFuncMutation(0, 0)

		} else if i == 1 {
			currentMutationState.PreviousState = currentState
			currentMutationState.PreviousNonMutationSate = currentState
			currentMutationState.transFunc = fA.createTransFuncMutation(1, 1)
			fA.states[0].PreviousState = currentMutationState
		} else {
			index := 0
			for j := 1; j < i; j++ {
				pal1 := Pals(fA.p[:i-j+1])
				pal2 := Pals(fA.p[j : i+1])

				if EqualPals(pal1, pal2) {
					index = i - j
					break
				}
			}

			currentMutationState.PreviousState = fA.mutationStates[index]
			currentMutationState.PreviousNonMutationSate = fA.states[index+1]
			currentMutationState.transFunc = fA.createTransFuncMutation(index, i)
			fA.states[i].PreviousState = currentMutationState

			if i == len(fA.p)-1 {
				currentState.PreviousState = fA.states[index+1]
				currentState.transFunc = fA.createTransFunc(index, i)
			}

			// currentState.transFunc = fA.createMutationFunc(i) -> update mutate state in the handle function
		}

		fA.states[i+1] = currentState
		fA.mutationStates[i] = currentMutationState
	}
}

// AGGAGCGTCTTCCCAAACCCG

// Algorithm3 returns the indexes with a matching
func Algorithm3(p, t string) []int {
	alg1 := time.Now()
	fa := NewFiniteAutomaton(p)
	fmt.Println("completed Algorithm1 in: ", time.Since(alg1))
	alg2 := time.Now()
	fa.Algorithm2()
	fmt.Println("completed Algorithm2 in: ", time.Since(alg2))

	queryTime := time.Now()

	// fa.Print()

	indexes := make([]int, 0)
	for i, l := range t {
		if i%100000000 == 0 {
			fmt.Println("\nAt Index: ", i, " query time is: ", time.Since(queryTime))
		}
		// fmt.Printf("i: %d, l: %s currentAlphabet: %s %s\n", i, string(l), fa.vP.currentAlphabet[0], fa.vP.currentAlphabet[1])
		if fa.NextLetter(string(l)) {
			fmt.Printf("\n Found an index: %d string: %s ", (i - len(p)), t[i+1-len(p):i+1])
			indexes = append(indexes, i+1-len(p))
		}

		// fmt.Printf(" curr: %d \n", fa.currentState)

		// fmt.Println("---------------------")
	}

	fmt.Println("Query time: ", time.Since(queryTime))

	return indexes
}

// TTTAAACGGCAAATTT

// Algorithm3Mutations executes the code necessary to
func Algorithm3Mutations(p, t string) []int {
	alg1 := time.Now()
	fa := NewFiniteAutomatonMutation(p)
	fmt.Println("completed Algorithm1 in: ", time.Since(alg1))
	alg2 := time.Now()
	fa.Algorithm2Mutation()
	fmt.Println("completed Algorithm2 in: ", time.Since(alg2))

	queryTime := time.Now()

	// fa.Print()

	indexes := make([]int, 0)
	for i, l := range t {
		if i%100000000 == 0 {
			fmt.Println("\nAt Index: ", i, " query time is: ", time.Since(queryTime))
		}
		// fmt.Printf("i: %d, l: %s currentAlphabet: %s %s %s mutant: %v\n", i, string(l), fa.vP.currentAlphabet[0], fa.vP.currentAlphabet[1], fa.vP.currentAlphabet[2], fa.state.isMutant)
		if fa.NextLetterMutation(string(l)) {
			fmt.Printf("\n Found an index: %d string: %s ", (i - len(p)), t[i+1-len(p):i])
			indexes = append(indexes, i+1-len(p))
		}

		// fmt.Printf(" curr: %d \n", fa.currentState)

		// fmt.Println("---------------------")
	}

	fmt.Println("Query time: ", time.Since(queryTime))

	return indexes
}

// FAState defines a struct which determines which state to go to next
type FAState struct {
	FinalState              bool
	NextState               *FAState
	PreviousState           *FAState
	PreviousNonMutationSate *FAState
	position                int
	transFunc               func()

	// mutation specific fields
	isMutant bool
}
