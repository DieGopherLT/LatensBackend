package github

import "time"

// RepositoryMetadataResponse represents the response from the repository metadata query
type RepositoryMetadataResponse struct {
	Data struct {
		Repository struct {
			Name           string `json:"name"`
			Description    string `json:"description"`
			PrimaryLanguage struct {
				Name string `json:"name"`
			} `json:"primaryLanguage"`
			DefaultBranchRef struct {
				Name string `json:"name"`
			} `json:"defaultBranchRef"`
			CreatedAt       time.Time `json:"createdAt"`
			UpdatedAt       time.Time `json:"updatedAt"`
			PushedAt        time.Time `json:"pushedAt"`
			IsArchived      bool      `json:"isArchived"`
			IsPrivate       bool      `json:"isPrivate"`
			IsFork          bool      `json:"isFork"`
			Owner struct {
				Login     string `json:"login"`
				AvatarUrl string `json:"avatarUrl"`
				Typename  string `json:"__typename"`
			} `json:"owner"`
			StargazerCount int `json:"stargazerCount"`
			ForkCount      int `json:"forkCount"`
			DiskUsage      int `json:"diskUsage"`
			URL            string `json:"url"`
			LicenseInfo    struct {
				Name   string `json:"name"`
				SpdxId string `json:"spdxId"`
			} `json:"licenseInfo"`
		} `json:"repository"`
		RateLimit struct {
			Limit     int       `json:"limit"`
			Cost      int       `json:"cost"`
			Remaining int       `json:"remaining"`
			ResetAt   time.Time `json:"resetAt"`
		} `json:"rateLimit"`
	} `json:"data"`
}

// SleepAnalysisResponse represents the response from the sleep analysis query
type SleepAnalysisResponse struct {
	Data struct {
		Repository struct {
			DefaultBranchRef struct {
				Name   string `json:"name"`
				Target struct {
					History struct {
						Edges []struct {
							Node struct {
								Oid           string    `json:"oid"`
								Message       string    `json:"message"`
								CommittedDate time.Time `json:"committedDate"`
								Author struct {
									Name  string `json:"name"`
									Email string `json:"email"`
								} `json:"author"`
								Additions int `json:"additions"`
								Deletions int `json:"deletions"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"history"`
					TotalCount struct {
						TotalCount int `json:"totalCount"`
					} `json:"totalCount"`
				} `json:"target"`
			} `json:"defaultBranchRef"`
			Refs struct {
				Edges []struct {
					Node struct {
						Name   string `json:"name"`
						Target struct {
							Oid           string    `json:"oid"`
							CommittedDate time.Time `json:"committedDate"`
							Message       string    `json:"message"`
							Author struct {
								Name string `json:"name"`
							} `json:"author"`
							CompareWithTip struct {
								AheadBy  int    `json:"aheadBy"`
								BehindBy int    `json:"behindBy"`
								Status   string `json:"status"`
							} `json:"compareWithTip"`
						} `json:"target"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"refs"`
			Readme struct {
				Text string `json:"text"`
			} `json:"readme"`
			PackageJson struct {
				Text string `json:"text"`
			} `json:"packageJson"`
			GoMod struct {
				Text string `json:"text"`
			} `json:"goMod"`
			Requirements struct {
				Text string `json:"text"`
			} `json:"requirements"`
			Gemfile struct {
				Text string `json:"text"`
			} `json:"gemfile"`
			CargoToml struct {
				Text string `json:"text"`
			} `json:"cargoToml"`
			PomXml struct {
				Text string `json:"text"`
			} `json:"pomXml"`
			Issues struct {
				TotalCount int `json:"totalCount"`
			} `json:"issues"`
			PullRequests struct {
				TotalCount int `json:"totalCount"`
			} `json:"pullRequests"`
			LatestRelease struct {
				Name        string    `json:"name"`
				TagName     string    `json:"tagName"`
				PublishedAt time.Time `json:"publishedAt"`
				Description string    `json:"description"`
			} `json:"latestRelease"`
			VulnerabilityAlerts struct {
				TotalCount int `json:"totalCount"`
			} `json:"vulnerabilityAlerts"`
		} `json:"repository"`
		RateLimit struct {
			Limit     int       `json:"limit"`
			Cost      int       `json:"cost"`
			Remaining int       `json:"remaining"`
			ResetAt   time.Time `json:"resetAt"`
		} `json:"rateLimit"`
	} `json:"data"`
}


// OwnedRepositoriesResponse represents the response from the owned repositories query (basic tier)
type OwnedRepositoriesResponse struct {
	Viewer struct {
		Login        string `json:"login"`
		Repositories struct {
			PageInfo struct {
				HasNextPage     bool   `json:"hasNextPage"`
				HasPreviousPage bool   `json:"hasPreviousPage"`
				StartCursor     string `json:"startCursor"`
				EndCursor       string `json:"endCursor"`
			} `json:"pageInfo"`
			TotalCount int               `json:"totalCount"`
			Nodes      []OwnedRepository `json:"nodes"`
		} `json:"repositories"`
	} `json:"viewer"`
	RateLimit struct {
		Limit     int       `json:"limit"`
		Cost      int       `json:"cost"`
		Remaining int       `json:"remaining"`
		ResetAt   time.Time `json:"resetAt"`
	} `json:"rateLimit"`
}

// OwnedRepository represents essential information for user-owned repositories (basic tier)
type OwnedRepository struct {
	// Basic identification
	ID            string `json:"id"`
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	URL           string `json:"url"`

	// Repository metadata
	Description     string `json:"description"`
	PrimaryLanguage struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"primaryLanguage"`

	// Default branch for sleep analysis
	DefaultBranchRef struct {
		Name   string `json:"name"`
		Target struct {
			Oid           string    `json:"oid"`
			CommittedDate time.Time `json:"committedDate"`
			Message       string    `json:"message"`
			Author struct {
				Name string `json:"name"`
			} `json:"author"`
			History struct {
				TotalCount int `json:"totalCount"`
			} `json:"history"`
		} `json:"target"`
	} `json:"defaultBranchRef"`

	// Repository state flags
	IsPrivate  bool `json:"isPrivate"`
	IsArchived bool `json:"isArchived"`
	IsFork     bool `json:"isFork"`
	IsDisabled bool `json:"isDisabled"`
	IsEmpty    bool `json:"isEmpty"`

	// Critical dates for sleep score calculation
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	PushedAt  time.Time `json:"pushedAt"`

	// Metrics for dashboard
	StargazerCount int `json:"stargazerCount"`
	ForkCount      int `json:"forkCount"`
	DiskUsage      int `json:"diskUsage"`

	// Quick activity metrics
	Issues struct {
		TotalCount int `json:"totalCount"`
	} `json:"issues"`

	PullRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"pullRequests"`

	// Repository topics (reduced for basic tier)
	RepositoryTopics struct {
		Nodes []struct {
			Topic struct {
				Name string `json:"name"`
			} `json:"topic"`
		} `json:"nodes"`
	} `json:"repositoryTopics"`

	// License information
	LicenseInfo struct {
		Name   string `json:"name"`
		SpdxId string `json:"spdxId"`
	} `json:"licenseInfo"`

	// All branches for activity analysis
	Refs struct {
		TotalCount int `json:"totalCount"`
		Edges      []struct {
			Node struct {
				Name   string `json:"name"`
				Target struct {
					Oid           string    `json:"oid"`
					CommittedDate time.Time `json:"committedDate"`
					Message       string    `json:"message"`
					Author        struct {
						Name string `json:"name"`
					} `json:"author"`
				} `json:"target"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"refs"`
}