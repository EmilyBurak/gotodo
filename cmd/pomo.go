package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func pomodoroCountdown(duration int, done chan bool) error {
	bar := progressbar.NewOptions(duration*60,
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Progress"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	var wg sync.WaitGroup
	wg.Add(1)
	go func(bar *progressbar.ProgressBar, duration int) error {
		defer wg.Done()
		progressBySecond(bar, duration*60)
		return nil
	}(bar, duration)
	for {
		select {
		case <-done:
			return nil
		}
	}
}

func recordLoop(rows [][]string, id int, recordCh chan []string, wg *sync.WaitGroup) {
	for i, row := range rows { // Iterate over all records
		if taskID, err := strconv.Atoi(row[0]); err == nil && taskID == id {
			wg.Add(1)
			fmt.Println("\nWorking on Pomodoros for Task:", row[1])
			go func(i int, record []string) error {
				defer wg.Done()
				rowValue, _ := strconv.Atoi(row[4])
				row[4] = strconv.Itoa(1 + rowValue)
				if row[5] != "" {
					pomNeeded, _ := strconv.Atoi(row[5])
					if pomNeeded == rowValue+1 {
						row[2] = "done"
					}
				}
				recordCh <- row
				return err
			}(i, row)
			if err != nil {
				fmt.Println(err)
				return
			}
			break
		}
	}
}

func progressBySecond(bar *progressbar.ProgressBar, duration int) {
	for i := 0; i < duration; i++ {
		bar.Add(1)
		time.Sleep(1 * time.Second)
	}
}

func takeBreak(breakDuration int, sessionPomos int) {
	breakTimer := time.NewTimer(time.Duration(breakDuration) * time.Minute)
	done := make(chan bool) // Create a channel to signal when the timer is done
	fmt.Println("\nSession Break #" + strconv.Itoa(sessionPomos) + " started for " + strconv.Itoa(breakDuration) + " minutes")
	go func() {
		err := pomodoroCountdown(breakDuration, done)
		if err != nil {
			fmt.Println(err)
		}
	}()
	if breakTimer == nil {
		breakTimer = time.NewTimer(time.Duration(breakDuration) * time.Minute)
	}
	<-breakTimer.C // Wait for the timer to expire
	done <- true   // Send a signal to the done channel
	close(done)    // Close the done channel
	fmt.Println("Break over!")
}

func updateTaskFile(rows [][]string) error {
	// Open the file in truncation mode to clear the contents
	fmt.Println("Updating task file")
	file, err := os.OpenFile("tasks.csv", os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the updated records to the file
	writer := csv.NewWriter(file)
	writer.WriteAll(rows)
	writer.Flush()
	fmt.Println("Task file updated")
	return writer.Error()
}

var pomoCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Start a pomodoro timer",
	Long:  `Start a pomodoro timer for a flagged duration and task.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Break out notification logic into a separate function
		// TODO: More desktop notifs?
		// TODO: Is it pomodoros or sessions completed to track?
		// TODO: Clean up print statements and code
		// TODO: Check against Pomodoros Needed for task completion
		duration, _ := cmd.Flags().GetInt("duration")
		id, _ := cmd.Flags().GetInt("ID")
		breakDuration, _ := cmd.Flags().GetInt("break")
		longBreak, _ := cmd.Flags().GetInt("longbreak")
		pomoAmount, _ := cmd.Flags().GetInt("pomos")

		sessionPomos := 0

		friendlyID := strconv.Itoa(id)
		if id == 0 {
			friendlyID = "No task specified"
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "Pomodoro Timer for Task ID:", friendlyID)
		w.Flush()

		for i := 0; i < pomoAmount; i++ {
			pomo := time.NewTimer(time.Duration(duration) * time.Minute)
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

			go pomodoroCountdown(duration, done)
			if err != nil {
				fmt.Println(err)
			}

			go fmt.Println("Work session started!")

			recordCh := make(chan []string, 10)
			var record []string

			if id != 0 {
				recordLoop(rows, id, recordCh, &wg)
			}

			<-pomo.C // Wait for the timer to expire
			done <- true
			close(done) // Close the done channel

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

			err = updateTaskFile(rows)

			if err != nil {
				fmt.Println(err)
				return
			}

			// Fyne notification for desktop and completion message for CLI
			notifApp := app.New()
			if id != 0 {
				// var breakWg sync.WaitGroup
				record = <-recordCh
				fmt.Println("Work session completed!")
				fmt.Println("\n" + record[4] + " Pomodoro work sessions completed for task " + strconv.Itoa(id) + " " + record[1])
				notif := fyne.NewNotification("Pomodoro completed!", record[4]+" total Pomodoros completed for task "+strconv.Itoa(id)+record[1]+"\nTime to take a break!")
				close(recordCh)
				sessionPomos = 1 + sessionPomos
				if sessionPomos%4 == 0 {
					notif = fyne.NewNotification("Long break!", "Time to take a long break!")
					notifApp.SendNotification(notif)
					takeBreak(longBreak, sessionPomos)
				} else if sessionPomos%4 != 0 {
					notifApp.SendNotification(notif)
					takeBreak(breakDuration, sessionPomos)
				}
			} else {
				close(recordCh)
				fmt.Println("Pomodoro completed!")
				notif := fyne.NewNotification("Pomodoro completed!", "Time to take a break!")
				sessionPomos = 1 + sessionPomos
				if sessionPomos%4 == 0 {
					notif = fyne.NewNotification("Long break!", "Time to take a long break!")
					notifApp.SendNotification(notif)
					takeBreak(longBreak, sessionPomos)
					return // Exit the loop
				} else if sessionPomos%4 != 0 {
					notif = fyne.NewNotification("Break time!", "Time to take a break!")
					notifApp.SendNotification(notif)
					takeBreak(breakDuration, sessionPomos)
				}
				notifApp.SendNotification(notif)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pomoCmd)
	pomoCmd.Flags().IntP("ID", "i", 0, "ID of the task")
	pomoCmd.Flags().IntP("duration", "d", 25, "Duration of the pomodoro timer in minutes")
	pomoCmd.Flags().IntP("break", "b", 2, "Duration of the break in minutes")
	pomoCmd.Flags().IntP("longbreak", "l", 3, "Duration of the long break in minutes")
	pomoCmd.Flags().IntP("pomos", "p", 4, "Number of pomodoros before a long break")
}
