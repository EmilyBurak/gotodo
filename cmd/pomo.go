/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var pomoCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Start a pomodoro timer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		duration, _ := cmd.Flags().GetInt("duration")
		id, _ := cmd.Flags().GetInt("ID")

		pomo := time.NewTimer(time.Duration(duration) * time.Second)
		startTime := time.Now()
		countdown := time.NewTicker(5 * time.Second)

		go func() {
			for {
				select {
				case <-countdown.C:
					currentTime := time.Now()
					difference := duration - int(currentTime.Sub(startTime).Seconds())
					fmt.Println("Time remaining:", difference)
				}
			}
		}()

		fmt.Println("Pomodoro started!")

		if id != 0 {
			file, err := os.Open("tasks.csv")
			if err != nil {
				fmt.Println("No tasks to complete")
				return
			}

			reader := csv.NewReader(file)

			for i := 0; ; i++ {
				record, err := reader.Read() // Read the next record
				if err == io.EOF {           // If we have reached the end of the file
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				if taskID, err := strconv.Atoi(record[0]); err == nil && taskID == id {
					fmt.Println("Task:", record[1])
				}
			}

			file.Close()
		}

		<-pomo.C // Wait for the timer to expire
		fmt.Println("Pomodoro complete!")
	},
}

func init() {
	rootCmd.AddCommand(pomoCmd)
	pomoCmd.Flags().IntP("ID", "i", 0, "ID of the task")
	pomoCmd.Flags().IntP("duration", "d", 25, "Duration of the pomodoro timer in minutes")
}
