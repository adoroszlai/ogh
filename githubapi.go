package main

func GetWorkflowRunJobs(org string, repo string, runId string) (map[string]interface{}, error) {
	apiPath := org + "/" + repo + "/actions/runs/" + runId + "/jobs"
	apiGetter := func() ([]byte, error) {
		return readGithubApiV3("https://api.github.com/repos/" + apiPath)
	}
	return asJson(cachedGet(apiGetter, toCacheKey(apiPath), buildResultCache))
}

func GetArtifacts(org string, repo string, runId string) (map[string]interface{}, error) {
	apiPath := org + "/" + repo + "/actions/runs/" + runId + "/artifacts"
	apiGetter := func() ([]byte, error) {
		return readGithubApiV3("https://api.github.com/repos/" + apiPath)
	}
	return asJson(cachedGet3min(apiGetter, toCacheKey(apiPath)))
}

func GetWorkflowRunsOfBranch(org string, repo string, workflowId string, branch string) (map[string]interface{}, error) {
	apiPath := org + "/" + repo + "/actions/workflows/" + workflowId + "/runs"
	cacheKey := toCacheKey(apiPath)
	url := "https://api.github.com/repos/" + apiPath + "?per_page=100"
	if branch != "" {
		cacheKey += "-" + branch
		url += "&branch=" + branch
	}
	apiGetter := func() ([]byte, error) {
		return readGithubApiV3(url)
	}

	return asJson(cachedGet3min(apiGetter, cacheKey))
}

func GetWorkflowRuns(org string, repo string, workflowId string) (map[string]interface{}, error) {
	return GetWorkflowRunsOfBranch(org, repo, workflowId, "")
}

func GetPr(org string, repo string, pullId string) (map[string]interface{}, error) {
	apiPath := org + "/" + repo + "/pulls/" + pullId
	apiGetter := func() ([]byte, error) {
		return readGithubApiV3("https://api.github.com/repos/" + apiPath)
	}
	return asJson(cachedGet3min(apiGetter, toCacheKey(apiPath)))
}

func GetPrCommits(org string, repo string, pullId string) ([]interface{}, error) {
	apiPath := org + "/" + repo + "/pulls/" + pullId + "/commits"
	apiGetter := func() ([]byte, error) {
		return readGithubApiV3("https://api.github.com/repos/" + apiPath)
	}
	return asJsonList(cachedGet3min(apiGetter, toCacheKey(apiPath)))
}

func GetChecksForCommits(org string, repo string, commitId string) (map[string]interface{}, error) {
	apiPath := org + "/" + repo + "/commits/" + commitId + "/check-runs"
	apiGetter := func() ([]byte, error) {
		return readGithubApiV3("https://api.github.com/repos/" + apiPath)
	}
	return asJson(cachedGet3min(apiGetter, toCacheKey(apiPath)))
}

func GetAllWorkflowRuns(org string, repo string) (map[string]interface{}, error) {
	apiPath := org + "/" + repo + "/actions/runs"
	apiGetter := func() ([]byte, error) {
		return readGithubApiV3("https://api.github.com/repos/" + apiPath + "?per_page=100")
	}
	return asJson(cachedGet(apiGetter, toCacheKey(apiPath), buildResultCache))
}
