package querying_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/underdog-tech/vulnbot/config"
	"github.com/underdog-tech/vulnbot/querying"
)

func TestCollectFindingsSingleProjectSingleFinding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyJson map[string]string
		var data []byte
		_ = json.NewDecoder(r.Body).Decode(&bodyJson)
		vulnQuery := strings.Contains(bodyJson["query"], "vulnerabilityAlerts")
		if vulnQuery {
			data, _ = os.ReadFile("testdata/single_project_single_finding_vulns.json")
		} else {
			data, _ = os.ReadFile("testdata/single_project_no_owners.json")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(data))
	}))
	defer server.Close()

	conf := config.Config{}
	env := config.Env{}
	env.GithubOrg = "heart-of-gold"
	env.GithubToken = "pangalactic-gargleblaster"

	ds := querying.NewGithubDataSource(conf, env)
	ds.GhClient = githubv4.NewEnterpriseClient(server.URL, &http.Client{})

	projects := querying.NewProjectCollection()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	err := ds.CollectFindings(projects, wg)
	if err != nil {
		t.Error(err)
	}
	expected := querying.ProjectCollection{
		Projects: []*querying.Project{
			{
				Name: "zaphod",
				Links: map[string]string{
					"GitHub": "https://heart-of-gold/zaphod",
				},
				Findings: []*querying.Finding{
					{
						Ecosystem:   config.FindingEcosystemGo,
						Severity:    config.FindingSeverityCritical,
						Description: "The Improbability Drive is far too improbable.",
						PackageName: "improbability-drive",
						Identifiers: querying.FindingIdentifierMap{
							querying.FindingIdentifierCVE: "CVE-42",
						},
					},
				},
				Owners: mapset.NewSet[config.TeamConfig](),
			},
		},
	}
	assert.Equal(t, &expected, projects)
}

// TestCollectFindingsOwnerNotConfigured is nearly identical to TestCollectFindingsSingleProjectSingleFinding
// The only difference is that this test simulates receiving an owning team from GitHub
// which is not present in config. This is to ensure that we don't end up with empty
// TeamConfig instances in our project owners set.
func TestCollectFindingsOwnerNotConfigured(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyJson map[string]string
		var data []byte
		_ = json.NewDecoder(r.Body).Decode(&bodyJson)
		vulnQuery := strings.Contains(bodyJson["query"], "vulnerabilityAlerts")
		if vulnQuery {
			data, _ = os.ReadFile("testdata/single_project_single_finding_vulns.json")
		} else {
			data, _ = os.ReadFile("testdata/single_project_single_owner.json")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(data))
	}))
	defer server.Close()

	conf := config.Config{}
	env := config.Env{}
	env.GithubOrg = "heart-of-gold"
	env.GithubToken = "pangalactic-gargleblaster"

	ds := querying.NewGithubDataSource(conf, env)
	ds.GhClient = githubv4.NewEnterpriseClient(server.URL, &http.Client{})

	projects := querying.NewProjectCollection()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	err := ds.CollectFindings(projects, wg)
	if err != nil {
		t.Error(err)
	}
	expected := querying.ProjectCollection{
		Projects: []*querying.Project{
			{
				Name: "zaphod",
				Links: map[string]string{
					"GitHub": "https://heart-of-gold/zaphod",
				},
				Findings: []*querying.Finding{
					{
						Ecosystem:   config.FindingEcosystemGo,
						Severity:    config.FindingSeverityCritical,
						Description: "The Improbability Drive is far too improbable.",
						PackageName: "improbability-drive",
						Identifiers: querying.FindingIdentifierMap{
							querying.FindingIdentifierCVE: "CVE-42",
						},
					},
				},
				Owners: mapset.NewSet[config.TeamConfig](),
			},
		},
	}
	assert.Equal(t, &expected, projects)
}

func TestCollectFindingsOwnerIsConfigured(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyJson map[string]string
		var data []byte
		_ = json.NewDecoder(r.Body).Decode(&bodyJson)
		vulnQuery := strings.Contains(bodyJson["query"], "vulnerabilityAlerts")
		if vulnQuery {
			data, _ = os.ReadFile("testdata/single_project_single_finding_vulns.json")
		} else {
			data, _ = os.ReadFile("testdata/single_project_single_owner.json")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(data))
	}))
	defer server.Close()

	crewTeam := config.TeamConfig{
		Name:        "Heart of Gold Crew",
		Github_slug: "crew",
	}
	conf := config.Config{
		Team: []config.TeamConfig{crewTeam},
	}
	env := config.Env{}
	env.GithubOrg = "heart-of-gold"
	env.GithubToken = "pangalactic-gargleblaster"

	ds := querying.NewGithubDataSource(conf, env)
	ds.GhClient = githubv4.NewEnterpriseClient(server.URL, &http.Client{})

	projects := querying.NewProjectCollection()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	err := ds.CollectFindings(projects, wg)
	if err != nil {
		t.Error(err)
	}
	owners := mapset.NewSet[config.TeamConfig]()
	owners.Add(crewTeam)
	expected := querying.ProjectCollection{
		Projects: []*querying.Project{
			{
				Name: "zaphod",
				Links: map[string]string{
					"GitHub": "https://heart-of-gold/zaphod",
				},
				Findings: []*querying.Finding{
					{
						Ecosystem:   config.FindingEcosystemGo,
						Severity:    config.FindingSeverityCritical,
						Description: "The Improbability Drive is far too improbable.",
						PackageName: "improbability-drive",
						Identifiers: querying.FindingIdentifierMap{
							querying.FindingIdentifierCVE: "CVE-42",
						},
					},
				},
				Owners: owners,
			},
		},
	}
	assert.Equal(t, &expected, projects)
}
