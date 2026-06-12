package main

import (
	"encoding/json"
	"github.com/elek/go-utils/github"
	"github.com/elek/go-utils/jira"
	jsonhelper "github.com/elek/go-utils/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

func CloseJira(jiraId string) error {
	jiraApi := jira.Jira{
		Url: "https://issues.apache.org/jira",
	}

	updated := map[string]interface{}{
		"fixVersions": []interface{}{
			map[string]interface{}{
				"add":
				map[string]interface{}{
					"name": "1.1.0",
				},
			},
		},
	}

	_, err := jiraApi.DoTransition(jiraId, "5", updated)
	return err
}

// TODO set patch available

func OpenJira(pullRequestId string, githubProject string) error {
	jiraProject := JiraNameFromGithubProject(githubProject)

	jiraApi := jira.Jira{
		Url: "https://issues.apache.org/jira",
	}

	pr, err := jsonhelper.AsJson(github.ReadGithubApiV3("https://api.github.com/repos/apache/" + githubProject + "/pulls/" + pullRequestId))
	if err != nil {
		return err
	}

	title := jsonhelper.MS(pr, "title")
	pullUrl := "https://github.com/apache/" + githubProject + "/pull/" + pullRequestId
	issuePattern, err := regexp.Compile(jiraProject + "-[0-9]+")
	if err != nil {
		return err
	}
	jiraId := issuePattern.FindString(title)

	if jiraId == "" {
		description := title + "\n\n" + pullUrl
		author := jsonhelper.MS(pr, "user", "login")
		if (strings.Contains(author, "dependabot")) {
			title = tweakDependabotTitle(title)
		}

		issue := map[string]interface{}{
			"project": map[string]string{
				"key": jiraProject,
			},
			"summary": title,
			"description": description,
			"issuetype": map[string]string{
				"name": "Task",
			},
		}
		log.Debug().Msg("Creating " + jiraProject + " issue: " + title + " for PR by " + author)
		resp, err := jiraApi.CreateJira(issue)
		respJson, err := jsonhelper.AsJson([]byte(resp), err)
		if err != nil {
			return err
		}
		//{"id":"13348103","key":"HDDS-4627","self":"https://issues.apache.org/jira/rest/api/2/issue/13348103"}
		jiraId = jsonhelper.MS(respJson, "key")
	}

	if jiraId == "" {
		return errors.New("Couldn't get or create jira Id")
	}

	return UpdatePullRequestTitle(pullRequestId, githubProject, jiraId, title)
}

func UpdatePullRequest(pullRequestId string, githubProject string, jiraId string) error {
	jiraApi := jira.Jira{
		Url: "https://issues.apache.org/jira",
	}

	resp, err := jiraApi.GetJira(jiraId)
	respJson, err := jsonhelper.AsJson([]byte(resp), err)
	if err != nil {
		return err
	}

	title := jsonhelper.MS(respJson, "fields", "summary")

	return UpdatePullRequestTitle(pullRequestId, githubProject, jiraId, title)
}

func UpdatePullRequestTitle(pullRequestId string, githubProject string, jiraId string, title string) error {
	patch := make(map[string]string)
	if !strings.Contains(title, jiraId) {
		patch["title"] = jiraId + ". " + title
	}
	log.Debug().Msg("Update PR " + pullRequestId + " title to " + patch["title"])
	if len(patch)>0 {
		patchJson, err := json.Marshal(patch)
		if err != nil {
			return err
		}
		resp, err := github.CallGithubApiV3WithBody("PATCH", "https://api.github.com/repos/apache/"+githubProject+"/pulls/"+pullRequestId, patchJson)
		if err != nil {
			return err
		}
		println(resp.StatusCode)
	}

	return nil
}

func tweakDependabotTitle(title string) string {
	// title := strings.Replace(title, "org.slf4j:slf4j-bom", "slf4j")
	title = strings.ReplaceAll(title, "software.amazon.awssdk:bom", "awssdk")

	bomRE := regexp.MustCompile(`[^: ]+:([^ ]+)-bom`)
	title = bomRE.ReplaceAllString(title, "$1")

	propRE := regexp.MustCompile(`([^. ]+)\.version`)
	title = propRE.ReplaceAllString(title, "$1")

	groupRE := regexp.MustCompile(`Bump ([^:]+):`)
	title = groupRE.ReplaceAllString(title, "Bump ")

	fromRE := regexp.MustCompile(` from [^ ]+ `)
	title = fromRE.ReplaceAllString(title, " ")

	deleteRE := regexp.MustCompile(`(shaded\.|\.Final$)`)
	title = deleteRE.ReplaceAllString(title, "")

	return title
}
