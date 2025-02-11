package cmd

import (
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/kubeshop/kusk/spec"
	"github.com/kubeshop/kusk/wizard"
	"github.com/kubeshop/kusk/wizard/prompt"
)

func init() {
	var apiSpecPath string

	wizardCmd := &cobra.Command{
		Use:   "wizard",
		Short: "Connects to current Kubernetes cluster and lists available generators",
		Run: func(cmd *cobra.Command, args []string) {
			if isTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()); !isTTY {
				log.Fatal("the wizard is only supported in an interactive context i.e. TTY")
			}

			// parse OpenAPI spec
			apiSpec, err := spec.NewParser(openapi3.NewLoader()).Parse(apiSpecPath)
			if err != nil {
				log.Fatal(err)
			}

			wizard.Start(apiSpecPath, apiSpec, prompt.New())
		},
	}

	// add common required flags
	wizardCmd.Flags().StringVarP(
		&apiSpecPath,
		"in",
		"i",
		"",
		"file path to api spec file to generate mappings from. e.g. --in apispec.yaml",
	)
	wizardCmd.MarkFlagRequired("in")

	rootCmd.AddCommand(wizardCmd)
}
