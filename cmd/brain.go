package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	brainCmd.AddCommand(brainAllCmd)
}

var brainCmd = &cobra.Command{
	Use:   "brain",
	Short: "Checkout Brain",
	Long: `Brain Lookup Commands, give informations about the 
Local existing repositories.	
`,
}
