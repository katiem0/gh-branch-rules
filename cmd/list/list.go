package list

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/katiem0/gh-branch-rules/internal/data"
	"github.com/katiem0/gh-branch-rules/internal/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type listCmdFlags struct {
	token    string
	hostname string
	listFile string
	debug    bool
}

func NewCmdList() *cobra.Command {
	listCmdFlags := listCmdFlags{}
	var authToken string

	listCmd := &cobra.Command{
		Use:   "list <source organization> [flags]",
		Short: "List organization level webhooks",
		Long:  "List organization level webhooks",
		Args:  cobra.ExactArgs(1),
		RunE: func(listCmd *cobra.Command, args []string) error {

			var err error
			var restClient api.RESTClient

			// Reinitialize logging if debugging was enabled
			if listCmdFlags.debug {
				logger, _ := log.NewLogger(listCmdFlags.debug)
				defer logger.Sync() // nolint:errcheck
				zap.ReplaceGlobals(logger)
			}

			if listCmdFlags.token != "" {
				authToken = listCmdFlags.token
			} else {
				t, _ := auth.TokenForHost(listCmdFlags.hostname)
				authToken = t
			}

			restClient, err = gh.RESTClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github+json",
				},
				Host:      listCmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving rest client")
				return err
			}

			owner := args[0]

			if _, err := os.Stat(listCmdFlags.listFile); errors.Is(err, os.ErrExist) {
				return err
			}

			reportWriter, err := os.OpenFile(listCmdFlags.listFile, os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				return err
			}

			return runCmdList(owner, &listCmdFlags, data.NewAPIGetter(restClient), reportWriter)
		},
	}

	reportFileDefault := fmt.Sprintf("WebhookReport-%s.csv", time.Now().Format("20060102150405"))

	// Configure flags for command

	listCmd.PersistentFlags().StringVarP(&listCmdFlags.token, "token", "t", "", `GitHub personal access token for reading source organization (default "gh auth token")`)
	listCmd.PersistentFlags().StringVarP(&listCmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	listCmd.Flags().StringVarP(&listCmdFlags.listFile, "output-file", "o", reportFileDefault, "Name of file to write CSV list to")
	listCmd.PersistentFlags().BoolVarP(&listCmdFlags.debug, "debug", "d", false, "To debug logging")

	return listCmd
}

func runCmdList(owner string, listCmdFlags *listCmdFlags, g *data.APIGetter, reportWriter io.Writer) error {
	csvWriter := csv.NewWriter(reportWriter)

	err := csvWriter.Write([]string{
		"Type",
		"ID",
		"Name",
		"Active",
		"Events",
		"Config_ContentType",
		"Config_InsecureSSL",
		"Config_Secret",
		"Config_URL",
		"Updated_At",
		"Created_At",
	})

	if err != nil {
		return err
	}

	zap.S().Debugf("Gathering Webooks for %s", owner)
	orgWebhooks, err := g.GetOrganizationWebhooks(owner)
	if err != nil {
		return err
	}

	var responseWebhooks []data.Webhook
	err = json.Unmarshal(orgWebhooks, &responseWebhooks)
	if err != nil {
		return err
	}
	zap.S().Debugf("Writing data for %d webhook(s) to output for organization %s", len(responseWebhooks), owner)
	for _, webhook := range responseWebhooks {
		err = csvWriter.Write([]string{
			webhook.HookType,
			strconv.Itoa(webhook.ID),
			webhook.Name,
			strconv.FormatBool(webhook.Active),
			fmt.Sprint(strings.Join(webhook.Events, ";")),
			webhook.Config.ContentType,
			webhook.Config.InsecureSSL,
			webhook.Config.Secret,
			webhook.Config.Url,
			webhook.UpdatedAt.Format(time.RFC3339),
			webhook.CreatedAt.Format(time.RFC3339),
		})

		if err != nil {
			zap.S().Error("Error raised in writing output", zap.Error(err))
		}
	}
	fmt.Printf("Successfully listed organizational webhooks for %s", owner)
	csvWriter.Flush()

	return nil
}
