package github

import "github.com/DieGopherLT/mfc_backend/pkg/graphql"

const (
	// RepositoryMetadataQuery fetches essential repository metadata for database storage
	RepositoryMetadataQuery graphql.Query = `
		query GetRepositoryMetadata($owner: String!, $name: String!) {
			# Main repository metadata for database storage
			repository(owner: $owner, name: $name) {
				# Basic repository information
				name
				description
				primaryLanguage {
					name
				}

				# Default branch reference
				defaultBranchRef {
					name
				}

				# Temporal data for sleep analysis
				createdAt
				updatedAt
				pushedAt

				# Repository state flags
				isArchived
				isPrivate
				isFork

				# Owner information
				owner {
					login
					avatarUrl
					__typename  # Will be either "User" or "Organization"
				}

				# Repository metrics
				stargazerCount
				forkCount
				diskUsage

				# URLs and licensing
				url
				licenseInfo {
					name
					spdxId
				}
			}

			# Rate limiting information for API management
			rateLimit {
				limit
				cost
				remaining
				resetAt
			}
		}
	`

	// SleepAnalysisQuery fetches data for calculating repository sleep score
	SleepAnalysisQuery graphql.Query = `
		query GetSleepAnalysisData($owner: String!, $name: String!, $since: DateTime!) {
			repository(owner: $owner, name: $name) {
				# Default branch for commits analysis
				defaultBranchRef {
					name

					# Recent commits for temporal analysis (last 50 for performance)
					target {
						... on Commit {
							history(first: 50, since: $since) {
								edges {
									node {
										oid
										message
										committedDate
										author {
											name
											email
										}
										# Change statistics for activity assessment
										additions
										deletions
									}
								}
							}

							# Total commits since specified date for frequency calculation
							totalCount: history(since: $since) {
								totalCount
							}
						}
					}
				}

				# All branches analysis for active development detection
				refs(refPrefix: "refs/heads/", first: 100, orderBy: {field: TAG_COMMIT_DATE, direction: DESC}) {
					edges {
						node {
							name
							target {
								... on Commit {
									oid
									committedDate
									message
									author {
										name
									}
									# Compare with default branch to detect ahead/behind
									compareWithTip: compare(headRef: "HEAD") {
										aheadBy
										behindBy
										status
									}
								}
							}
						}
					}
				}

				# Context files for project understanding
				readme: object(expression: "HEAD:README.md") {
					... on Blob {
						text
					}
				}

				# Configuration files by technology (using aliases to avoid conflicts)
				packageJson: object(expression: "HEAD:package.json") {
					... on Blob {
						text
					}
				}

				goMod: object(expression: "HEAD:go.mod") {
					... on Blob {
						text
					}
				}

				requirements: object(expression: "HEAD:requirements.txt") {
					... on Blob {
						text
					}
				}

				gemfile: object(expression: "HEAD:Gemfile") {
					... on Blob {
						text
					}
				}

				cargoToml: object(expression: "HEAD:Cargo.toml") {
					... on Blob {
						text
					}
				}

				pomXml: object(expression: "HEAD:pom.xml") {
					... on Blob {
						text
					}
				}

				# Community activity indicators
				issues(states: OPEN) {
					totalCount
				}

				pullRequests(states: OPEN) {
					totalCount
				}

				# Latest release for versioning context
				latestRelease {
					name
					tagName
					publishedAt
					description
				}

				# Security alerts as activity indicator
				vulnerabilityAlerts(first: 1) {
					totalCount
				}
			}

			# Rate limiting information for API management
			rateLimit {
				limit
				cost
				remaining
				resetAt
			}
		}
	`

	// OwnedRepositoriesQuery fetches only repositories owned by the authenticated user (basic tier)
	OwnedRepositoriesQuery graphql.Query = `
		query GetOwnedRepositories($first: Int!, $after: String) {
			# Get only repositories owned by authenticated user (basic tier)
			viewer {
				login

				# Only owned repositories (no collaborations or organizations)
				repositories(
					first: $first
					after: $after
					orderBy: {field: PUSHED_AT, direction: DESC}
					ownerAffiliations: [OWNER]
					isFork: false  # Exclude forks for cleaner personal repos list
				) {
					# Pagination information
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}

					totalCount

					# Repository nodes with essential dashboard information
					nodes {
						# Basic identification
						id
						name
						nameWithOwner
						url

						# Repository metadata
						description
						primaryLanguage {
							name
							color
						}

						# Default branch for sleep analysis
						defaultBranchRef {
							name
						}

						# Repository state flags for filtering
						isPrivate
						isArchived
						isFork
						isDisabled
						isEmpty

						# Critical dates for sleep score calculation
						createdAt
						updatedAt
						pushedAt

						# Metrics for dashboard display
						stargazerCount
						forkCount
						diskUsage

						# Recent activity indicators (for quick sleep assessment)
						defaultBranchRef {
							name
							target {
								... on Commit {
									# Latest commit info for sleep calculation
									oid
									committedDate
									message
									author {
										name
									}
								}
							}
						}

						# Activity metrics for sleep score (simplified for individual repos)
						issues(states: OPEN) {
							totalCount
						}

						pullRequests(states: OPEN) {
							totalCount
						}

						# Repository topics (reduced for basic tier)
						repositoryTopics(first: 5) {
							nodes {
								topic {
									name
								}
							}
						}

						# License information
						licenseInfo {
							name
							spdxId
						}
					}
				}
			}

			# Rate limiting information for API management
			rateLimit {
				limit
				cost
				remaining
				resetAt
			}
		}
	`
)