package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
	"time"
)

func run(all bool, authorFilter string, reference Reference) error {
	var key string
	key = reference.Org + "-" + reference.Repo + "-"
	if all {
		key += "pr"
	} else {
		key += "review"
	}
	apiCall := func() ([]byte, error) {
		return readPrWithGraphql(reference)
	}
	body, err := cachedGet3min(apiCall, key)
	if err != nil {
		return err
	}

	result := make(map[string]interface{})
	json.Unmarshal(body, &result)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Author", "Summary", "Participants", "Check"})
	table.SetAutoWrapText(false)
	prs := m(result, "data", "repository", "pullRequests", "edges")

	for _, prNode := range l(prs) {

		pr := m(prNode, "node")

		if !all && !ready(pr) {
			continue
		}

		author := prAuthor(pr)
		participants := getParticipants(pr, author)
		statusMark := ""
		destMark := ""
		if ms(pr, "baseRefName") != "master" {
			destMark = "(->" + ms(pr, "baseRefName") + ")"
		}
		if ms(pr, "mergeable") == "CONFLICTING" {
			statusMark = "[C] "
		}
		if m(pr, "isDraft") == true {
			statusMark += "[D]"
		}
		if authorFilter == "" || authorFilter == author {
			table.Append([]string{
				fmt.Sprintf("%d", int(m(pr, "number").(float64))),
				">" + limit(author, 12),
				limit(statusMark+destMark+ms(pr, "title"), 50),
				limit(strings.Join(participants, ","), 35),
				buildStatus(pr),
			})
		}
	}
	table.Render() // Send output

	return nil
}



func ready(pr interface{}) bool {
	if ms(pr, "mergeable") == "CONFLICTING" {
		return false
	}
	for _, commitEdge := range l(m(pr, "commits", "edges")) {
		commit := m(commitEdge, "node", "commit")
		for _, suite := range l(m(commit, "checkSuites", "edges")) {
			for _, runs := range l(m(suite, "node", "checkRuns", "edges")) {
				conclusion := ms(runs, "node", "conclusion")
				if conclusion == "FAILURE" || conclusion == "CANCELLED" {
					return false
				}
			}
		}
		break
	}

	for _, review := range lastReviewsPerUser(pr) {
		state := ms(review, "state")
		if state == "CHANGES_REQUESTED" {
			return false
		}
	}

	return true
}

func getParticipants(pr interface{}, author string) []string {
	reviews := lastReviewsPerUser(pr)

	participants := make([]string, 0)

	participants = append(participants, filterReviews(reviews, "CHANGES_REQUESTED", "✕")...)
	participants = append(participants, filterReviews(reviews, "APPROVED", "✓")...)
	participants = append(participants, filterReviews(reviews, "COMMENTED", "")...)

	commenters := make(map[string]bool)
	for _, participant := range l(m(pr, "participants", "edges")) {
		login := ms(participant, "node", "login")
		if _, ok := reviews[login]; !ok && login != author {
			participants = append(participants, limit(login, 5))
			commenters[login] = true
		}
	}

	for _, login := range reviewRequests(pr) {
		_, reviewed := reviews[login]
		_, commented := commenters[login]
		if !reviewed && !commented && login != author {
			participants = append(participants, "?" + limit(strings.ToUpper(login), 5))
		}
	}

	return participants
}

func reviewRequests(pr interface{}) []string {
	requests := make([]string, 0)
	for _, request := range l(m(pr, "reviewRequests", "edges")) {
		requests = append(requests, ms(request, "node", "requestedReviewer", "login"))
	}
	return requests
}

func lastReviewsPerUser(pr interface{}) map[string]interface{} {
	prAuthor := prAuthor(pr)
	reviewers := make(map[string]interface{})
	for _, review := range l(m(pr, "reviews", "nodes")) {
		author := ms(review, "author", "login")
		if last_review, found := reviewers[author]; found {

			oldRecord, err := time.Parse(time.RFC3339, ms(last_review, "updatedAt"))
			if err != nil {
				panic(err)
			}

			newRecord, err := time.Parse(time.RFC3339, ms(review, "updatedAt"))
			if err != nil {
				panic(err)
			}

			if oldRecord.Before(newRecord) {
				reviewers[author] = review
			}

		} else if author != prAuthor {
			reviewers[author] = review
		}
	}
	return reviewers
}

func filterReviews(reviews map[string]interface{}, status string, symbol string) []string {
	result := make([]string, 0)
	for _, review := range reviews {
		state := ms(review, "state")
		if state == status {
			result = append(result, symbol+limit(strings.ToUpper(ms(review, "author", "login")), 5))
		}
	}
	return result
}

func prAuthor(pr interface{}) string {
	return ms(pr, "author", "login")
}

type statusTransform struct {
	position int
	abbrev   byte
}
