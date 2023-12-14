package update

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/katiem0/gh-branch-rules/internal/data"
	"github.com/katiem0/gh-branch-rules/internal/log"
	"github.com/katiem0/gh-branch-rules/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
	token    string
	hostname string
	fileName string
	debug    bool
}

func NewCmdUpdate() *cobra.Command {
	cmdFlags := cmdFlags{}
	var authToken string

	updateCmd := &cobra.Command{
		Use:   "update [flags] <organization>",
		Short: "update branch protection policies",
		Long:  "Update branch protection policies for repositories from a file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(createCmd *cobra.Command, args []string) error {
			var err error
			var restClient api.RESTClient
			var gqlClient api.GQLClient

			// Reinitialize logging if debugging was enabled
			if cmdFlags.debug {
				logger, _ := log.NewLogger(cmdFlags.debug)
				defer logger.Sync() // nolint:errcheck
				zap.ReplaceGlobals(logger)
			}

			if cmdFlags.token != "" {
				authToken = cmdFlags.token
			} else {
				t, _ := auth.TokenForHost(cmdFlags.hostname)
				authToken = t
			}

			restClient, err = gh.RESTClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving rest client")
				return err
			}

			gqlClient, err = gh.GQLClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github.hawkgirl-preview+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving graphql client")
				return err
			}
			owner := args[0]

			return runCmdUpdate(owner, &cmdFlags, utils.NewAPIGetter(gqlClient, restClient))
		},
	}
	// Configure flags for command
	updateCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub personal access token for organization to write to (default "gh auth token")`)
	updateCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	updateCmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to create webhooks from")
	updateCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	updateCmd.MarkFlagRequired("from-file")

	return updateCmd
}

func runCmdUpdate(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var policyData [][]string
	var importBranchPolicyList []data.BranchProtectionRule
	zap.S().Infof("Reading in file %s and updating branch protection policies", cmdFlags.fileName)
	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening branch protection policies csv file")
		}
		// remember to close the file at the end of the program
		defer f.Close()
		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		policyData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading assignments from csv file")
		}
		importBranchPolicyList = utils.CreateBranchProtectionPolicyData(policyData)
	} else {
		zap.S().Errorf("Error arose identifying users to add")
	}
	zap.S().Debugf("Determining permissions to create")

	for _, importBranchPolicy := range importBranchPolicyList {
		zap.S().Debugf("Updating branch policy %s with ID %s", importBranchPolicy.Pattern, importBranchPolicy.ID)

		err := g.UpdateBranchProtectionPolicies(importBranchPolicy)
		if err != nil {
			zap.S().Errorf("Error arose creating permission %s", importBranchPolicy.Pattern)
		}
	}

	fmt.Printf("Successfully updated branch protection policies from %s in org %s", cmdFlags.fileName, owner)
	return nil
}
