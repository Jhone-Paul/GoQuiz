package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	flag.Parse()
	file, err := os.Open(*csvFilename)
	if err != nil {

		exit(fmt.Sprintf("Failed to open a CSV file: %s\n", *csvFilename))
	}
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("could not read the given file")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(30 * time.Second)
	var score int

	for i := range problems {
		fmt.Printf("%s = ", problems[i].q)
		answerCh := make(chan string)

		go func() {
			var input string
			fmt.Scanln(&input)
			answerCh <- input
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime's up! Your score is %d\n", score)
			os.Exit(0)
		case input := <-answerCh:
			if input == problems[i].a {
				score++
			}
		}

	}
	fmt.Printf("\nYou finished before the time was up! Your score is %d\n", score)

}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}
