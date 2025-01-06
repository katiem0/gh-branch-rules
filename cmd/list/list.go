package list

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/katiem0/gh-branch-rules/internal/data"
	"github.com/katiem0/gh-branch-rules/internal/log"
	"github.com/katiem0/gh-branch-rules/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
	token    string
	hostname string
	listFile string
	debug    bool
}

func NewCmdList() *cobra.Command {
	cmdFlags := cmdFlags{}
	var authToken string

	listCmd := &cobra.Command{
		Use:   "list [flags] <organization> [repo ...]",
		Short: "Generate a report of branch protection rules for repositories.",
		Long:  "Generate a report of branch protection rules for a list of repositories",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(listCmd *cobra.Command, args []string) error {
			var err error
			var gqlClient *api.GraphQLClient
			var restClient *api.RESTClient

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

			restClient, err = api.NewRESTClient(api.ClientOptions{
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

			gqlClient, err = api.NewGraphQLClient(api.ClientOptions{
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
			repos := args[1:]

			if _, err := os.Stat(cmdFlags.listFile); errors.Is(err, os.ErrExist) {
				return err
			}

			reportWriter, err := os.OpenFile(cmdFlags.listFile, os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				return err
			}

			return runCmdList(owner, repos, &cmdFlags, utils.NewAPIGetter(gqlClient, restClient), reportWriter)
		},
	}

	reportFileDefault := fmt.Sprintf("BranchRules-%s.csv", time.Now().Format("20060102150405"))

	// Configure flags for command

	listCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub Personal Access Token (default "gh auth token")`)
	listCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	listCmd.Flags().StringVarP(&cmdFlags.listFile, "output-file", "o", reportFileDefault, "Name of file to write CSV list to")
	listCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")

	return listCmd
}

func runCmdList(owner string, repos []string, cmdFlags *cmdFlags, g *utils.APIGetter, reportWriter io.Writer) error {
	var reposCursor *string
	var allRepos []data.RepoInfo
	zap.S().Infof("Gathering repositories in %s to list branch protection policies", owner)
	csvWriter := csv.NewWriter(reportWriter)

	err := csvWriter.Write([]string{
		"RepositoryName",
		"RepositoryID",
		"BranchProtectionRulePattern",
		"BranchProtectionRuleId",
		"AllowsDeletions",
		"AllowsForcePushes",
		"BlockCreations",
		"DismissesStaleReviews",
		"IsAdminEnforced",
		"LockAllowsFetchAndMerge",
		"LockBranch",
		"RequireLastPushApproval",
		"RequiredApprovingReviewCount",
		"RequiresApprovingReviews",
		"RequiresCodeOwnerReviews",
		"RequiresCommitSignatures",
		"RequiresConversationResolution",
		"RequiresDeployments",
		"RequiresLinearHistory",
		"RequiresStatusChecks",
		"RequiresStrictStatusChecks",
		"RestrictsPushes",
		"RestrictsReviewDismissals",
	})

	if err != nil {
		return err
	}
	zap.S().Infof("Gathering repositories and branch protection rules")

	if len(repos) > 0 {
		zap.S().Infof("Processing repos: %s", repos)

		for _, repo := range repos {

			zap.S().Debugf("Processing %s/%s", owner, repo)

			repoQuery, err := g.GetRepo(owner, repo)
			if err != nil {
				return err
			}
			allRepos = append(allRepos, repoQuery.Repository)
		}

	} else {
		// Prepare writer for outputting report
		for {
			zap.S().Debugf("Processing list of repositories for %s", owner)
			reposQuery, err := g.GetReposList(owner, reposCursor)

			if err != nil {
				return err
			}

			allRepos = append(allRepos, reposQuery.Organization.Repositories.Nodes...)

			reposCursor = &reposQuery.Organization.Repositories.PageInfo.EndCursor

			if !reposQuery.Organization.Repositories.PageInfo.HasNextPage {
				break
			}
		}
	}

	for _, singleRepo := range allRepos {
		zap.S().Debugf("Gathering Branch Protection Policies for repo %s", singleRepo.Name)
		var bpCursor *string
		var allBPPolicies []data.BranchProtectionRule
		for {
			branchProtectionList, err := g.GetBranchProtections(owner, singleRepo.Name, bpCursor)

			if err != nil {
				return err
			}

			allBPPolicies = append(allBPPolicies, branchProtectionList.Repository.BranchProtectionRules.Nodes...)
			bpCursor = &branchProtectionList.Repository.BranchProtectionRules.PageInfo.EndCursor
			if !branchProtectionList.Repository.BranchProtectionRules.PageInfo.HasNextPage {
				break
			}
		}
		for _, policy := range allBPPolicies {
			err = csvWriter.Write([]string{
				singleRepo.Name,
				strconv.Itoa(singleRepo.DatabaseId),
				policy.Pattern,
				policy.ID,
				strconv.FormatBool(policy.AllowsDeletions),
				strconv.FormatBool(policy.AllowsForcePushes),
				strconv.FormatBool(policy.BlocksCreations),
				strconv.FormatBool(policy.DismissesStaleReviews),
				strconv.FormatBool(policy.IsAdminEnforced),
				strconv.FormatBool(policy.LockAllowsFetchAndMerge),
				strconv.FormatBool(policy.LockBranch),
				strconv.FormatBool(policy.RequireLastPushApproval),
				strconv.Itoa(policy.RequiredApprovingReviewCount),
				strconv.FormatBool(policy.RequiresApprovingReviews),
				strconv.FormatBool(policy.RequiresCodeOwnerReviews),
				strconv.FormatBool(policy.RequiresCommitSignatures),
				strconv.FormatBool(policy.RequiresConversationResolution),
				strconv.FormatBool(policy.RequiresDeployments),
				strconv.FormatBool(policy.RequiresLinearHistory),
				strconv.FormatBool(policy.RequiresStatusChecks),
				strconv.FormatBool(policy.RequiresStrictStatusChecks),
				strconv.FormatBool(policy.RestrictsPushes),
				strconv.FormatBool(policy.RestrictsReviewDismissals),
			})

			if err != nil {
				zap.S().Error("Error raised in writing output", zap.Error(err))
			}
		}

	}
	fmt.Printf("Successfully listed repository level branch protection policies for %s", owner)
	csvWriter.Flush()

	return nil
}
