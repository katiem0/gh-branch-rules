package create

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/katiem0/gh-branch-rules/internal/data"
	"github.com/katiem0/gh-branch-rules/internal/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
	sourceToken    string
	sourceOrg      string
	sourceHostname string
	token          string
	hostname       string
	fileName       string
	debug          bool
}

func NewCmdCreate() *cobra.Command {
	cmdFlags := cmdFlags{}
	var authToken string

	cmd := &cobra.Command{
		Use:   "create <target organization> [flags]",
		Short: "Create organization level webhooks",
		Long:  "Create organization level webhooks",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(createCmd *cobra.Command, args []string) error {
			if len(cmdFlags.fileName) == 0 && len(cmdFlags.sourceOrg) == 0 {
				return errors.New("A file or source organization must be specified where webhooks will be created from.")
			} else if len(cmdFlags.sourceOrg) > 0 && len(cmdFlags.sourceToken) == 0 {
				return errors.New("A Personal Access Token must be specified to access webhooks from the Source Organization.")
			} else if len(cmdFlags.fileName) > 0 && len(cmdFlags.sourceOrg) > 0 {
				return errors.New("Specify only one of `--source-organization` or `from-file`.")
			}
			return nil
		},
		RunE: func(createCmd *cobra.Command, args []string) error {
			var err error
			var restClient api.RESTClient

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

			owner := args[0]

			return runCmdCreate(owner, &cmdFlags, data.NewAPIGetter(restClient))
		},
	}
	// Configure flags for command
	cmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub personal access token for organization to write to (default "gh auth token")`)
	cmd.PersistentFlags().StringVarP(&cmdFlags.sourceToken, "source-token", "s", "", `GitHub personal access token for Source Organization (Required for --source-organization)`)
	cmd.PersistentFlags().StringVarP(&cmdFlags.sourceOrg, "source-organization", "o", "", `Name of the Source Organization to copy webhooks from (Requires --source-token)`)
	cmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	cmd.PersistentFlags().StringVarP(&cmdFlags.sourceHostname, "source-hostname", "", "github.com", "GitHub Enterprise Server hostname where webhooks are copied from")
	cmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to create webhooks from")
	cmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")

	return cmd
}

func runCmdCreate(owner string, cmdFlags *cmdFlags, g *data.APIGetter) error {
	var webhookData [][]string
	var webhooksList []data.CreatedWebhook
	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening webhooks csv file")
		}
		// remember to close the file at the end of the program
		defer f.Close()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		webhookData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading webhooks from csv file")
		}
		webhooksList = g.CreateWebhookList(webhookData)
		zap.S().Debugf("Identifying Webhook list to create under %s", owner)
	} else if len(cmdFlags.sourceOrg) > 0 {
		zap.S().Debugf("Reading in webhooks from %s", cmdFlags.sourceOrg)
		var authToken string
		var restSourceClient api.RESTClient

		if cmdFlags.sourceToken != "" {
			authToken = cmdFlags.sourceToken
		} else {
			t, _ := auth.TokenForHost(cmdFlags.sourceHostname)
			authToken = t
		}

		restSourceClient, err := gh.RESTClient(&api.ClientOptions{
			Headers: map[string]string{
				"Accept": "application/vnd.github+json",
			},
			Host:      cmdFlags.sourceHostname,
			AuthToken: authToken,
		})
		if err != nil {
			zap.S().Errorf("Error arose retrieving source rest client")
			return err
		}
		zap.S().Debugf("Gathering webhooks %s", cmdFlags.sourceOrg)

		webhookResponse, err := data.GetSourceOrganizationWebhooks(cmdFlags.sourceOrg, data.NewAPIGetter(restSourceClient))
		if err != nil {
			return err
		}
		err = json.Unmarshal(webhookResponse, &webhooksList)
		if err != nil {
			return err
		}
	} else {
		zap.S().Errorf("Error arose identifying webhooks")
	}
	zap.S().Debugf("Determining webhooks to create")
	for _, webhook := range webhooksList {
		if webhook.Config.Secret == "********" {
			zap.S().Debugf("Webhook with URL %s required a secret, and needs a new secret to be entered.", webhook.Config.Url)
			webhookString := fmt.Sprintf("Please enter the new secret to be created with webhook %s:", webhook.Config.Url)
			webhookSecret := data.SensitivePrompt(webhookString)
			webhook.Config.Secret = webhookSecret
		}
		createWebhook, err := json.Marshal(webhook)

		if err != nil {
			return err
		}

		reader := bytes.NewReader(createWebhook)
		zap.S().Debugf("Creating Webhooks under %s", owner)
		err = g.CreateOrganizationWebhook(owner, reader)
		if err != nil {
			zap.S().Errorf("Error arose creating webhook with %s", webhook.Config.Url)
		}
	}
	fmt.Printf("Successfully created webhooks for: %s.", owner)
	return nil
}
