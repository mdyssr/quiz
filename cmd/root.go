/*
Copyright © 2022 Muhammad Yasser <mhmmdyssraamr@gmail.com>

*/
package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/mdyssr/quiz/helpers"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A quiz game in your terminal",
	Long:  `A quiz game in your terminal`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if limit <= 0 {
			err := errors.New("limit must be greater than zero")
			return err
		}

		startQuiz()

		return nil
	},
}

type Quiz struct {
	Total int
	Score int
}

var csvPath string
var limit int
var shuffle bool

func startQuiz() {

	f, err := helpers.OpenFile(csvPath)
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	timerCh := make(chan bool)
	questionsCh := make(chan bool)
	quiz := Quiz{}

	fmt.Print("Press enter key to start the quiz...")
	fmt.Scanln()
	fmt.Println()

	go timer(timerCh)
	go askQuestions(f, &quiz, questionsCh)

	select {
	case <-questionsCh:
		endQuiz(quiz, "\nQuiz completed!")
		return
	case <-timerCh:
		endQuiz(quiz, "\n\nTime over!")
	}
	//
	// <-timerCh
	// fmt.Println("\n\ntime over")
	// fmt.Printf("You scored %d out of %d\n\n", quiz.Score, quiz.Total)
}

func endQuiz(quiz Quiz, message string) {
	fmt.Println(message)
	fmt.Printf("You scored %d out of %d\n\n", quiz.Score, quiz.Total)
}

// timer function
func timer(ch chan bool) {
	time.Sleep(time.Duration(limit) * time.Second)
	ch <- true
}

func askQuestions(file *os.File, quiz *Quiz, ch chan bool) {
	var answer string

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		log.Fatalf("error reading %s\n", file.Name())
	}

	quiz.Total = len(lines)

	if shuffle {
		helpers.Shuffle(lines)
	}

	for i, line := range lines {
		q := line[0]
		a := line[1]

		fmt.Printf("Problem #%d: %s = ", i+1, q)
		fmt.Scan(&answer)

		if helpers.CleanText(answer) == helpers.CleanText(a) {
			quiz.Score++
		}
	}

	ch <- true
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.quiz.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// csv flag
	rootCmd.Flags().StringVarP(&csvPath, "csv", "c", "problems.csv", "a csv file int the format of question,answer")

	// limit flag
	rootCmd.Flags().IntVarP(&limit, "limit", "l", 30, "the time limit for the quiz in seconds")
	rootCmd.Flags().BoolVarP(&shuffle, "shuffle", "s", false, "shuffle questions")
}
