package data

import "github.com/shurcooL/graphql"

type ReposQuery struct {
	Organization struct {
		Repositories struct {
			Nodes    []RepoInfo
			PageInfo struct {
				EndCursor   string
				HasNextPage bool
			}
		} `graphql:"repositories(first: 100, after: $endCursor)"`
	} `graphql:"organization(login: $owner)"`
}

type RepoInfo struct {
	DatabaseId int    `json:"databaseId"`
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}

type BranchProtectionRulesQuery struct {
	Repository struct {
		BranchProtectionRules struct {
			Nodes    []BranchProtectionRule
			PageInfo struct {
				EndCursor   string
				HasNextPage bool
			}
		} `graphql:"branchProtectionRules(first: 100, after: $endCursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type BranchProtectionRule struct {
	AllowsDeletions                bool   `json:"allowsDeletions"`
	AllowsForcePushes              bool   `json:"allowsForcePushes"`
	BlocksCreations                bool   `json:"blocksCreations"`
	ID                             string `json:"id"`
	DismissesStaleReviews          bool   `json:"dismissesStaleReviews"`
	IsAdminEnforced                bool   `json:"isAdminEnforced"`
	LockAllowsFetchAndMerge        bool   `json:"lockAllowsFetchAndMerge"`
	LockBranch                     bool   `json:"lockBranch"`
	Pattern                        string `json:"pattern"`
	RequireLastPushApproval        bool   `json:"requireLastPushApproval"`
	RequiredApprovingReviewCount   int    `json:"requiredApprovingReviewCount"`
	RequiresApprovingReviews       bool   `json:"requiresApprovingReviews"`
	RequiresCodeOwnerReviews       bool   `json:"requiresCodeOwnerReviews"`
	RequiresCommitSignatures       bool   `json:"requiresCommitSignatures"`
	RequiresConversationResolution bool   `json:"requiresConversationResolution"`
	RequiresDeployments            bool   `json:"requiresDeployments"`
	RequiresLinearHistory          bool   `json:"requiresLinearHistory"`
	RequiresStatusChecks           bool   `json:"requiresStatusChecks"`
	RequiresStrictStatusChecks     bool   `json:"requiresStrictStatusChecks"`
	RestrictsPushes                bool   `json:"restrictsPushes"`
	RestrictsReviewDismissals      bool   `json:"restrictsReviewDismissals"`
}

type RepoSingleQuery struct {
	Repository RepoInfo `graphql:"repository(owner: $owner, name: $name)"`
}

type MutationBranchProtection struct {
	UpdateBranchProtectionRule struct {
		ClientMutationId graphql.String
	} `graphql:"updateBranchProtectionRule (input: $input)"`
}

type UpdateBranchProtectionRuleInput struct {
	AllowsDeletions                graphql.Boolean `json:"allowsDeletions"`
	AllowsForcePushes              graphql.Boolean `json:"allowsForcePushes"`
	BlocksCreations                graphql.Boolean `json:"blocksCreations"`
	BranchProtectionRuleId         graphql.String  `json:"branchProtectionRuleId"`
	DismissesStaleReviews          graphql.Boolean `json:"dismissesStaleReviews"`
	IsAdminEnforced                graphql.Boolean `json:"isAdminEnforced"`
	LockAllowsFetchAndMerge        graphql.Boolean `json:"lockAllowsFetchAndMerge"`
	LockBranch                     graphql.Boolean `json:"lockBranch"`
	Pattern                        graphql.String  `json:"pattern"`
	RequireLastPushApproval        graphql.Boolean `json:"requireLastPushApproval"`
	RequiredApprovingReviewCount   graphql.Int     `json:"requiredApprovingReviewCount"`
	RequiresApprovingReviews       graphql.Boolean `json:"requiresApprovingReviews"`
	RequiresCodeOwnerReviews       graphql.Boolean `json:"requiresCodeOwnerReviews"`
	RequiresCommitSignatures       graphql.Boolean `json:"requiresCommitSignatures"`
	RequiresConversationResolution graphql.Boolean `json:"requiresConversationResolution"`
	RequiresDeployments            graphql.Boolean `json:"requiresDeployments"`
	RequiresLinearHistory          graphql.Boolean `json:"requiresLinearHistory"`
	RequiresStatusChecks           graphql.Boolean `json:"requiresStatusChecks"`
	RequiresStrictStatusChecks     graphql.Boolean `json:"requiresStrictStatusChecks"`
	RestrictsPushes                graphql.Boolean `json:"restrictsPushes"`
	RestrictsReviewDismissals      graphql.Boolean `json:"restrictsReviewDismissals"`
}
