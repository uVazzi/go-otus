package console

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

type appVersion struct {
	Release   string
	BuildDate string
	GitHash   string
}

func GetVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print app version",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := json.NewEncoder(os.Stdout).Encode(appVersion{
				Release:   release,
				BuildDate: buildDate,
				GitHash:   gitHash,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}
}
