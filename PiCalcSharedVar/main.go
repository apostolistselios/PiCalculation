package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// A struct that holds the calculated Pi in the result variable and a mutex
// so that it cant be incremented safely.
type sharedPi struct {
	result float64
	mux    sync.Mutex
}

// Global variable where the calculated Pi is stored.
var pi sharedPi

// Wait Group used to sunchronize the main Go routine with the others.
var wg = sync.WaitGroup{}

func main() {
	steps, noRoutines, err := parseArguments()
	if err != nil {
		fmt.Printf("Usage of %q:\n", os.Args[0])
		fmt.Printf("  -steps int\n    number of steps to used to calculate Pi\n")
		fmt.Printf("  -noRoutines int (optional)\n    number of Go routines to be created\n")
		fmt.Println(err)
		os.Exit(1)
	}

	// Initializes the number of Go routines to be created based on the command line args.
	if noRoutines == -1 {
		noRoutines = runtime.NumCPU()
	} else {
		runtime.GOMAXPROCS(noRoutines)
	}
	step := 1.0 / float64(steps)

	start := time.Now()
	wg.Add(noRoutines)
	for i := 0; i < noRoutines; i++ {
		start := (i * steps) / noRoutines
		end := ((i + 1) * steps) / noRoutines
		go calcPi(start, end, step)
	}
	wg.Wait()
	end := time.Now()

	fmt.Println("Pi calculated:", pi.result)
	fmt.Println("Math library Pi:", math.Pi)
	fmt.Println("Took:", end.Sub(start))

}

// Calculates pi and adds it to the field result.
func calcPi(start, end int, step float64) {
	sum := 0.0
	for i := start; i < end; i++ {
		x := step * (float64(i) + 0.5)
		sum += 4.0 / (1.0 + x*x)
	}
	pi.mux.Lock()
	pi.result += sum * step
	pi.mux.Unlock()
	wg.Done()
}

// Functions that helps parsing the command line arguments.

// Parses the command line arguments and checks for errors.
// If there are no errors, returns the proper values, else returns an error.
func parseArguments() (int, int, error) {
	switch noArgs := len(os.Args); noArgs {
	case 2:
		steps, err := parseStepsArgument()
		if err != nil {
			return 0, 0, err
		}
		return steps, -1, nil
	case 3:
		steps, err := parseStepsArgument()
		if err != nil {
			return 0, 0, err
		}

		noRoutines, err := parseNoRoutinesArgument()
		if err != nil {
			return 0, 0, err
		}
		return steps, noRoutines, nil
	default:
		return 0, 0, fmt.Errorf("error: too many or not enough arguments")
	}
}

// Parses the steps command line arguments and checks for errors.
// If there are no errors, returns the proper value, else returns an error.
func parseStepsArgument() (int, error) {
	steps, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return 0, fmt.Errorf("error: the number of steps has to be an integer")
	}
	if steps == 0 || steps < 0 {
		return 0, fmt.Errorf("error: the number of steps has to be a positive number")
	}
	return steps, nil
}

// Parses the noRoutines command line arguments and checks for errors.
// If there are no errors, returns the proper value, else returns an error.
func parseNoRoutinesArgument() (int, error) {
	noRoutines, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return 0, fmt.Errorf("error: the number of Go routines to be created has to be an integer")
	}

	if noRoutines == 0 || noRoutines < 0 {
		return 0, fmt.Errorf("error: the number of Go routines to be created has to be a positive number")
	}
	return noRoutines, nil
}
