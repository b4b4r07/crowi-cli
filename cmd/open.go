package cmd

import (
	"github.com/b4b4r07/crowi/cli"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a crowi page",
	Long:  "Open a crowi page",
	RunE:  open,
}

func open(cmd *cobra.Command, args []string) error {
	screen, err := cli.NewScreen()
	if err != nil {
		return err
	}

	selectedLines, err := screen.Filter()
	if err != nil {
		return err
	}

	line, err := screen.ParseLine(selectedLines[0])
	if err != nil {
		return err
	}

	return cli.OpenURL(line.URL)
}

func init() {
	RootCmd.AddCommand(openCmd)
}
