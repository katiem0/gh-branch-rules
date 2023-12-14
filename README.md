# gh-branch-rules

A GitHub `gh` [CLI](https://cli.github.com/) extension to create a report containing branch protections for a single repository or list of repositories, as well as create branch protections from a file.

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation) instructions.

2. Install the extension:
   ```sh
   gh extension install katiem0/gh-branch-rules
   ```

For more information: [`gh extension install`](https://cli.github.com/manual/gh_extension_install).

## Usage

The `gh-branch-rules` extension supports `GitHub.com` and GitHub Enterprise Server, through the use of `--hostname` and the following commands:

```sh
$ gh branch-rules -h
List and update branch protection rules for repositories in an organization.

Usage:
  branch-rules [command]

Available Commands:
  list        Generate a report of branch protection rules for repositories.
  update      Create and/or update branch protection policies

Flags:
  -h, --help   help for branch-rules

Use "branch-rules [command] --help" for more information about a command.
```

### List Branch Protection Policies

This extension will create a csv report of branch protection policies for specified repositories or all repositories in an organization.

```sh
$ gh branch-rules list -h
Generate a report of branch protection rules for a list of repositories

Usage:
  branch-rules list [flags] <organization> [repo ...]

Flags:
  -d, --debug                To debug logging
  -h, --help                 help for list
      --hostname string      GitHub Enterprise Server hostname (default "github.com")
  -o, --output-file string   Name of file to write CSV list to (default "BranchRules-20231214102016.csv")
  -t, --token string         GitHub Personal Access Token (default "gh auth token")
```

The output `csv` file contains the following information:

| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository where the data is extracted from |
|`RepositoryID`| The `ID` associated with the Repository, for API usage |
|`BranchProtectionRulePattern`| Identifies the protection rule pattern |
|`BranchProtectionRuleId`| The branch protection policy ID that is needed for updating policies |
|`AllowsDeletions`| If the branch associated to the policy can be deleted |
|`AllowsForcePushes`| If force pushes are allowed on the branch |
|`BlockCreations`| If branch creation matching the rule pattern is a protected operation |
|`DismissesStaleReviews`| If new commits pushed to matching branches dismiss pull request review approvals |
|`IsAdminEnforced`| If admins override branch protection |
|`LockAllowsFetchAndMerge`| If users can pull changes from upstream when the branch is locked. Set to `true` allows fork syncing. Set to false prevents fork syncing|
|`LockBranch`| If the branch is set as `read-only`. If this is `true`, users will not be able to push to the branch|
|`RequireLastPushApproval`| If the most recent push must be approved by someone other than the person who pushed it |
|`RequiredApprovingReviewCount`| Number of approving reviews required to update matching branches|
|`RequiresApprovingReviews`| If approving reviews are required to update matching branches |
|`RequiresCodeOwnerReviews`| If reviews from code owners are required to update matching branches |
|`RequiresCommitSignatures`| If commits are required to be signed|
|`RequiresConversationResolution`| If conversations are required to be resolved before merging |
|`RequiresDeployments`| If this branch requires deployment to specific environments before merging |
|`RequiresLinearHistory`| If merge commits are prohibited from being pushed to this branch|
|`RequiresStatusChecks`| If status checks are required to update matching branches|
|`RequiresStrictStatusChecks`| If branches are required to be up to date before merging|
|`RestrictsPushes`| If pushing to matching branches is restricted|
|`RestrictsReviewDismissals`| If dismissal of pull request reviews is restricted|

### Update Branch Protection Policies

Branch protection policies for specified repositories defined in a **required** csv file for an organization.

```sh
$ gh branch-rules update -h
Update branch protection policies for repositories from a file.

Usage:
  branch-rules update [flags] <organization>

Flags:
  -d, --debug              To debug logging
  -f, --from-file string   Path and Name of CSV file to create webhooks from
  -h, --help               help for update
      --hostname string    GitHub Enterprise Server hostname (default "github.com")
  -t, --token string       GitHub personal access token for organization to write to (default "gh auth token")
```

The required csv file should contain the following information:


| Field Name | Description |
|:-----------|:------------|
|`RepositoryName` | The name of the repository where the data is extracted from |
|`RepositoryID`| The `ID` associated with the Repository, for API usage |
|`BranchProtectionRulePattern`| Identifies the protection rule pattern |
|`BranchProtectionRuleId`| The branch protection policy ID that is needed for updating policies |
|`AllowsDeletions`| If the branch associated to the policy can be deleted |
|`AllowsForcePushes`| If force pushes are allowed on the branch |
|`BlockCreations`| If branch creation matching the rule pattern is a protected operation |
|`DismissesStaleReviews`| If new commits pushed to matching branches dismiss pull request review approvals |
|`IsAdminEnforced`| If admins override branch protection |
|`LockAllowsFetchAndMerge`| If users can pull changes from upstream when the branch is locked. Set to `true` allows fork syncing. Set to false prevents fork syncing|
|`LockBranch`| If the branch is set as `read-only`. If this is `true`, users will not be able to push to the branch|
|`RequireLastPushApproval`| If the most recent push must be approved by someone other than the person who pushed it |
|`RequiredApprovingReviewCount`| Number of approving reviews required to update matching branches|
|`RequiresApprovingReviews`| If approving reviews are required to update matching branches |
|`RequiresCodeOwnerReviews`| If reviews from code owners are required to update matching branches |
|`RequiresCommitSignatures`| If commits are required to be signed|
|`RequiresConversationResolution`| If conversations are required to be resolved before merging |
|`RequiresDeployments`| If this branch requires deployment to specific environments before merging |
|`RequiresLinearHistory`| If merge commits are prohibited from being pushed to this branch|
|`RequiresStatusChecks`| If status checks are required to update matching branches|
|`RequiresStrictStatusChecks`| If branches are required to be up to date before merging|
|`RestrictsPushes`| If pushing to matching branches is restricted|
|`RestrictsReviewDismissals`| If dismissal of pull request reviews is restricted|