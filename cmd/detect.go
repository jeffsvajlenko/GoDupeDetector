/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"GoDupeDetector/internal/detection"
	"GoDupeDetector/internal/parsing"
	"GoDupeDetector/internal/printer"
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// detectCmd represents the detect command
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Runs clone detection.",
	Long:  `Runs clone detection.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle input folder
		input, _ := cmd.Flags().GetString("input")
		finfo, _ := os.Stat(input)
		if finfo == nil {
			return errors.New("input directory does not exist or is invalid")
		}
		if !finfo.IsDir() {
			return errors.New("input must be a directory")
		}

		// Handle output file
		output, _ := cmd.Flags().GetString("output")
		finfo, _ = os.Stat(output)
		if finfo != nil && finfo.IsDir() {
			return errors.New("output exists and is a directory")
		}
		outf, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file with error: %s", err.Error())
		}
		defer outf.Close()

		// Handle threshold
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		if threshold < 0.0 || threshold > 1.0 {
			return errors.New("threshold needs to be a value between 0.0 and 1.0")
		}

		// Get go files
		fmt.Print("Collecting go files...")
		files, err := parsing.FileList(input)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(" %d files collected.\n", len(files))

		// Parse code files
		fmt.Print("Parsing...")
		pset, _ := parsing.Parse(files)
		fmt.Println(" Done!")

		// Detect clones
		fmt.Print("Detecting...")
		cset, _ := detection.DetectClones(pset, 0.70)
		fmt.Printf(" %d clones detected.\n", len(cset.Clones))

		// Output clones
		fmt.Println("Outputting...")
		bout := bufio.NewWriter(outf)
		err = printer.PrintCloneReport(pset, cset, bout)
		if err != nil {
			log.Fatal(err)
		}

		// Made it without error!
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.Flags().String("input", "", "Input source folder.")
	detectCmd.MarkFlagRequired("input")
	detectCmd.Flags().String("output", "", "Output report file.")
	detectCmd.MarkFlagRequired("output")
	detectCmd.Flags().Float64("threshold", 0.70, "Minimum similarity threshold for reported clones.")
}
