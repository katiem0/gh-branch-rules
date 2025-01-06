package utils

import (
	"strconv"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/katiem0/gh-branch-rules/internal/data"
	"github.com/shurcooL/graphql"
)

type Getter interface {
	GetRepo(owner string, name string) ([]data.RepoSingleQuery, error)
	GetReposList(owner string, endCursor *string) ([]data.ReposQuery, error)
	GetBranchProtections(owner string, name string, endCursor *string) (*data.BranchProtectionRulesQuery, error)
	UpdateBranchProtectionPolicies(branchPolicy data.BranchProtectionRule) error
}

type APIGetter struct {
	gqlClient  api.GraphQLClient
	restClient api.RESTClient
}

func NewAPIGetter(gqlClient *api.GraphQLClient, restClient *api.RESTClient) *APIGetter {
	return &APIGetter{
		gqlClient:  *gqlClient,
		restClient: *restClient,
	}
}

func (g *APIGetter) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	query := new(data.ReposQuery)
	variables := map[string]interface{}{
		"endCursor": (*graphql.String)(endCursor),
		"owner":     graphql.String(owner),
	}

	err := g.gqlClient.Query("getRepos", &query, variables)

	return query, err
}

func (g *APIGetter) GetRepo(owner string, name string) (*data.RepoSingleQuery, error) {
	query := new(data.RepoSingleQuery)
	variables := map[string]interface{}{
		"owner": graphql.String(owner),
		"name":  graphql.String(name),
	}

	err := g.gqlClient.Query("getRepo", &query, variables)
	return query, err
}

func (g *APIGetter) GetBranchProtections(owner string, name string, endCursor *string) (*data.BranchProtectionRulesQuery, error) {
	query := new(data.BranchProtectionRulesQuery)
	variables := map[string]interface{}{
		"endCursor": (*graphql.String)(endCursor),
		"owner":     graphql.String(owner),
		"name":      graphql.String(name),
	}

	err := g.gqlClient.Query("getBranchProtectionPolicies", &query, variables)

	return query, err
}

func CreateBranchProtectionPolicyData(fileData [][]string) []data.BranchProtectionRule {
	var importBranchRules []data.BranchProtectionRule
	var branchPolicy data.BranchProtectionRule
	for _, each := range fileData[1:] {
		branchPolicy.Pattern = each[2]
		branchPolicy.ID = each[3]
		branchPolicy.AllowsDeletions, _ = strconv.ParseBool(each[4])
		branchPolicy.AllowsForcePushes, _ = strconv.ParseBool(each[5])
		branchPolicy.BlocksCreations, _ = strconv.ParseBool(each[6])
		branchPolicy.DismissesStaleReviews, _ = strconv.ParseBool(each[7])
		branchPolicy.IsAdminEnforced, _ = strconv.ParseBool(each[8])
		branchPolicy.LockAllowsFetchAndMerge, _ = strconv.ParseBool(each[9])
		branchPolicy.LockBranch, _ = strconv.ParseBool(each[10])
		branchPolicy.RequireLastPushApproval, _ = strconv.ParseBool(each[11])
		branchPolicy.RequiredApprovingReviewCount, _ = strconv.Atoi(each[12])
		branchPolicy.RequiresApprovingReviews, _ = strconv.ParseBool(each[13])
		branchPolicy.RequiresCodeOwnerReviews, _ = strconv.ParseBool(each[14])
		branchPolicy.RequiresCommitSignatures, _ = strconv.ParseBool(each[15])
		branchPolicy.RequiresConversationResolution, _ = strconv.ParseBool(each[16])
		branchPolicy.RequiresDeployments, _ = strconv.ParseBool(each[17])
		branchPolicy.RequiresLinearHistory, _ = strconv.ParseBool(each[18])
		branchPolicy.RequiresStatusChecks, _ = strconv.ParseBool(each[19])
		branchPolicy.RequiresStrictStatusChecks, _ = strconv.ParseBool(each[20])
		branchPolicy.RestrictsPushes, _ = strconv.ParseBool(each[21])
		branchPolicy.RestrictsReviewDismissals, _ = strconv.ParseBool(each[22])
		importBranchRules = append(importBranchRules, branchPolicy)
	}
	return importBranchRules
}

func (g *APIGetter) UpdateBranchProtectionPolicies(branchPolicy data.BranchProtectionRule) error {
	mutation := new(data.MutationBranchProtection)
	input := data.UpdateBranchProtectionRuleInput{
		AllowsDeletions:                graphql.Boolean(branchPolicy.AllowsDeletions),
		AllowsForcePushes:              graphql.Boolean(branchPolicy.AllowsForcePushes),
		BlocksCreations:                graphql.Boolean(branchPolicy.BlocksCreations),
		BranchProtectionRuleId:         graphql.String(branchPolicy.ID),
		DismissesStaleReviews:          graphql.Boolean(branchPolicy.DismissesStaleReviews),
		IsAdminEnforced:                graphql.Boolean(branchPolicy.IsAdminEnforced),
		LockAllowsFetchAndMerge:        graphql.Boolean(branchPolicy.LockAllowsFetchAndMerge),
		LockBranch:                     graphql.Boolean(branchPolicy.LockBranch),
		Pattern:                        graphql.String(branchPolicy.Pattern),
		RequireLastPushApproval:        graphql.Boolean(branchPolicy.RequireLastPushApproval),
		RequiredApprovingReviewCount:   graphql.Int(branchPolicy.RequiredApprovingReviewCount),
		RequiresApprovingReviews:       graphql.Boolean(branchPolicy.RequiresApprovingReviews),
		RequiresCodeOwnerReviews:       graphql.Boolean(branchPolicy.RequiresCodeOwnerReviews),
		RequiresCommitSignatures:       graphql.Boolean(branchPolicy.RequiresCommitSignatures),
		RequiresConversationResolution: graphql.Boolean(branchPolicy.RequiresConversationResolution),
		RequiresDeployments:            graphql.Boolean(branchPolicy.RequiresDeployments),
		RequiresLinearHistory:          graphql.Boolean(branchPolicy.RequiresLinearHistory),
		RequiresStatusChecks:           graphql.Boolean(branchPolicy.RequiresStatusChecks),
		RequiresStrictStatusChecks:     graphql.Boolean(branchPolicy.RequiresStrictStatusChecks),
		RestrictsPushes:                graphql.Boolean(branchPolicy.RestrictsPushes),
		RestrictsReviewDismissals:      graphql.Boolean(branchPolicy.RestrictsReviewDismissals),
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := g.gqlClient.Mutate("getBranchProtectionPolicies", &mutation, variables)
	return err

}
