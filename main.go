package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	csvfilename = "todolist.csv"
)

var data [][]string

func main() {
	fmt.Println("Welcome to My to do list.")
	for {
		showlist()
		command, index := input()
		switch command {
		case "*":
			markdone(index)
		case "!":
			markundone(index)
		case "-":
			delete(index)
		case "":
			return
		default:
			addtodo(command)
		}
	}
}

func showlist() {
	lines, err := readcsv()
	if err != nil {
		log.Fatal(err)
	}
	data = lines
	if len(data) == 0 {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("Nothing to do")
		fmt.Println(strings.Repeat("=", 50))
		return
	}
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  To do list")
	fmt.Printf("---|%s\n", strings.Repeat("-", 44))
	var lineno int64
	for _, line := range lines {
		lineno++
		if line[1] == "1" {
			fmt.Print("[x]")
		} else {
			fmt.Print("[ ]")
		}
		fmt.Print("| ")
		fmt.Printf(" %d. %s\n", lineno, line[0])
	}
	fmt.Println(strings.Repeat("=", 50))
}

func input() (command string, index int) {
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Command [done=*#,undone=!#,delete=-#,enter=exit,other=add]:")
		input, err := consoleReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			return
		}
		command = input[:1]
		if command == "*" || command == "!" || command == "-" {
			if !is_numeric(input[1:]) {
				fmt.Println("Invalid Command!")
				continue
			}
			index, err = strconv.Atoi(input[1:])
			return command, index
		}
		return input, 0
	}
}

func addtodo(todo string) {
	data = append(data, []string{todo, "0"})
	updatedata()
}

func markdone(no int) {
	for i, v := range data {
		_ = v
		if i+1 == no {
			data[i][1] = "1"
			break
		}
	}
	updatedata()
}

func markundone(no int) {
	for i, v := range data {
		_ = v
		if i+1 == no {
			data[i][1] = "0"
			break
		}
	}
	updatedata()
}

func delete(no int) {
	s := make([][]string, 0, 4)
	s = data[:no-1]
	if len(data[no:]) > 0 {
		t := data[no:]
		s = append(s, t...)
	}
	data = s
	updatedata()
}

func updatedata() {
	err := updatecsv(data)
	if err != nil {
		log.Fatal(err)
	}
}

func is_numeric(word string) bool {
	return regexp.MustCompile(`\d`).MatchString(word)
}

func updatecsv(lines [][]string) error {
	log.Printf("%#v\n", lines)
	f, err := os.OpenFile(csvfilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	err = w.WriteAll(lines)
	if err != nil {
		return err
	}
	return nil
}

func readcsv() ([][]string, error) {
	// read the file
	f, err := os.Open(csvfilename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
