package collector

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "collector",
	Short: "This is the todo collector",
	Long:  "This binary will poll kafka todo topic, and collect them",
	RunE:  action,
}

func action(cmd *cobra.Command, _ []string) (err error) {
	return nil

}
