package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks that are not marked as done.`,
	Run: func(cmd *cobra.Command, args []string) {
		All, _ := cmd.Flags().GetBool("all")

		file, err := os.Open("tasks.csv")
		if err != nil {
			fmt.Println("No tasks to list")
			return
		}
		defer file.Close()

		csvReader := csv.NewReader(file)

		headers, err := csvReader.Read()
		if err != nil {
			panic(err)
		}

		// Print the headers
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, headers[0], "\t", headers[1], "\t", headers[2], "\t", headers[3], "\t", headers[4], "\t", headers[5])

		// Loop through the records and print the tasks that are not done or deleted, unless the all flag is set
		for {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			if !All && (record[2] == "done" || record[3] == "true") {
				continue
			}
			if record[5] == "" {
				record[5] = "N/A"
			}
			fmt.Fprintln(w, record[0], "\t", record[1], "\t", record[2], "\t", record[3], "\t", record[4], "\t", record[5])
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "List all tasks")
}
