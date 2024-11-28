package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var pomoCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Start a pomodoro timer",
	Long:  `Start a pomodoro timer for a flagged duration and task.`,
	Run: func(cmd *cobra.Command, args []string) {
		duration, _ := cmd.Flags().GetInt("duration")
		id, _ := cmd.Flags().GetInt("ID")

		pomo := time.NewTimer(time.Duration(duration) * time.Second)
		startTime := time.Now()
		countdown := time.NewTicker(5 * time.Second)

		file, err := os.OpenFile("tasks.csv", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("Error opening task csv file")
			return
		}

		defer file.Close()

		reader := csv.NewReader(file)
		rows, err := reader.ReadAll() // Read all records from the file
		if err != nil {
			fmt.Println("Error reading task csv file")
			return
		}
		var wg sync.WaitGroup
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-countdown.C:
					currentTime := time.Now()
					difference := duration - int(currentTime.Sub(startTime).Seconds())
					fmt.Println("Time remaining:", difference)
				case <-done:
					return
				}
			}
		}()

		fmt.Println("Pomodoro started!")
		recordCh := make(chan []string)
		var record []string

		if id != 0 {
			for i, row := range rows { // Iterate over all records
				if taskID, err := strconv.Atoi(row[0]); err == nil && taskID == id {
					wg.Add(1)
					fmt.Println("Task:", row[1])
					go func(i int, record []string) {
						defer wg.Done()
						rowValue, _ := strconv.Atoi(row[4])
						row[4] = strconv.Itoa(1 + rowValue)
						if row[5] != "" {
							pomNeeded, _ := strconv.Atoi(row[5])
							if pomNeeded == rowValue+1 {
								row[2] = "done"
							}
						}
						rows[i] = row
						recordCh <- row
					}(i, row)
					break
				}
			}
		}
		record = <-recordCh

		<-pomo.C // Wait for the timer to expire
		done <- true

		wg.Wait()

		// Close the file to open it in truncation mode
		file.Close()

		// Open the file in truncation mode to clear the contents
		file, err = os.OpenFile("tasks.csv", os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Open the file in read-write mode to write the updated records
		file, err = os.OpenFile("tasks.csv", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write the updated records to the file
		writer := csv.NewWriter(file)
		writer.WriteAll(rows)
		writer.Flush()

		if id != 0 {
			fmt.Println(record[4] + " Pomodoros completed for task " + string(id) + record[1])
		} else {
			fmt.Println("Pomodoro completed!")
		}
	},
}

func init() {
	rootCmd.AddCommand(pomoCmd)
	pomoCmd.Flags().IntP("ID", "i", 0, "ID of the task")
	pomoCmd.Flags().IntP("duration", "d", 25, "Duration of the pomodoro timer in minutes")
}
