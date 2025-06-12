package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	showCol1        bool
	showCol2        bool
	showCol3        bool
	caseInsensitive bool
	delimiter       string
)

var rootCmd = &cobra.Command{
	Use:   "gocomm [OPTION]... FILE1 FILE2",
	Short: "Compare sorted files FILE1 and FILE2 line by line.",
	Long: `Compare sorted files FILE1 and FILE2 line by line.

With no options, produce three-column output:
  column 1: lines only in FILE1
  column 2: lines only in FILE2
  column 3: lines in both FILE1 and FILE2

Options:
  -1                      suppress column 1 (lines unique to FILE1)
  -2                      suppress column 2 (lines unique to FILE2)
  -3                      suppress column 3 (lines that appear in both files)
  -i, --insensitive       case-insensitive comparison
  -d, --output-delimiter  separate columns with STR

Examples:
  gocomm file1.txt file2.txt
  gocomm -12 file1.txt file2.txt
  gocomm -i -d "|" file1.txt file2.txt`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showCol1, "suppress1", "1", false, "suppress column 1 (lines unique to FILE1)")
	rootCmd.Flags().BoolVarP(&showCol2, "suppress2", "2", false, "suppress column 2 (lines unique to FILE2)")
	rootCmd.Flags().BoolVarP(&showCol3, "suppress3", "3", false, "suppress column 3 (lines that appear in both files)")
	rootCmd.Flags().BoolVarP(&caseInsensitive, "insensitive", "i", false, "case-insensitive comparison")
	rootCmd.Flags().StringVarP(&delimiter, "output-delimiter", "d", "\t", "separate columns with STR")
}
