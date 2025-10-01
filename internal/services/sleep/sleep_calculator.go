package sleep

import (
	"math"
	"time"

	"github.com/DieGopherLT/LatensBackend/internal/services/github"
	"github.com/samber/lo"
)

const (
	maxScore                = 100.0
	inactivityThresholdDays = 180.0
	fragmentationThreshold  = 3.0
	stalenessThresholdDays  = 90.0
	inactivityWeight        = 0.5
	fragmentationWeight     = 0.3
	stalenessWeight         = 0.2
)

func CalculateScore(repo *github.OwnedRepository) int {
	if repo.IsArchived || repo.IsDisabled || repo.IsFork {
		return int(maxScore)
	}

	inactivityScore := calculateInactivityScore(repo)
	fragmentationScore := calculateFragmentationScore(repo)
	stalenessScore := calculateStalenessScore(repo)

	// Adjust weights for repos with only default branch
	hasMultipleBranches := hasAdditionalBranches(repo)
	daysSince := time.Since(repo.DefaultBranchRef.Target.CommittedDate).Hours() / 24

	if !hasMultipleBranches || daysSince > inactivityThresholdDays {
		return int(math.Round(inactivityScore))
	}

	sleepScore := (inactivityWeight * inactivityScore) +
		(fragmentationWeight * fragmentationScore) +
		(stalenessWeight * stalenessScore)

	// Clamp between 0 and 100, then round to int
	sleepScore = math.Max(0, math.Min(maxScore, sleepScore))
	return int(math.Round(sleepScore))
}

func hasAdditionalBranches(repo *github.OwnedRepository) bool {
	return repo.Refs.TotalCount > 1
}

func calculateInactivityScore(repo *github.OwnedRepository) float64 {
	daysSinceLastCommit := time.Since(repo.DefaultBranchRef.Target.CommittedDate).Hours() / 24
	if daysSinceLastCommit >= inactivityThresholdDays {
		return maxScore
	}
	score := (daysSinceLastCommit / inactivityThresholdDays) * maxScore
	return math.Min(maxScore, score)
}

func calculateFragmentationScore(repo *github.OwnedRepository) float64 {
	var aheadBranches int
	defaultBranchCommitDate := repo.DefaultBranchRef.Target.CommittedDate

	for _, edge := range repo.Refs.Edges {
		if edge.Node.Name == repo.DefaultBranchRef.Name {
			continue
		}

		if edge.Node.Target.CommittedDate.After(defaultBranchCommitDate) {
			aheadBranches++
		}
	}

	score := (float64(aheadBranches) / fragmentationThreshold) * maxScore
	return math.Min(maxScore, score)
}

func calculateStalenessScore(repo *github.OwnedRepository) float64 {
	daysSinceLastCommitOnBranches := make([]float64, 0, len(repo.Refs.Edges))
	defaultBranchCommitDate := repo.DefaultBranchRef.Target.CommittedDate

	for _, edge := range repo.Refs.Edges {
		if edge.Node.Name == repo.DefaultBranchRef.Name {
			continue
		}

		branchCommitDate := edge.Node.Target.CommittedDate
		if branchCommitDate.Before(defaultBranchCommitDate) || branchCommitDate.Equal(defaultBranchCommitDate) {
			continue
		}

		daysSinceLastCommit := time.Since(branchCommitDate).Hours() / 24
		daysSinceLastCommitOnBranches = append(daysSinceLastCommitOnBranches, daysSinceLastCommit)
	}

	if len(daysSinceLastCommitOnBranches) == 0 {
		return 0
	}

	totalDays := lo.Reduce(daysSinceLastCommitOnBranches, func(acc float64, days float64, _ int) float64 {
		return acc + days
	}, 0.0)

	avgDays := totalDays / float64(len(daysSinceLastCommitOnBranches))
	score := (avgDays / stalenessThresholdDays) * maxScore
	return math.Min(maxScore, score)
}
