package leader

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	data := graphQLdata{}
	assert.NoError(t, json.Unmarshal([]byte(prResultFixture), &data))
	jsonStr, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.JSONEq(t, prResultFixture, string(jsonStr))
}

func TestChartDataFromGQL(t *testing.T) { //nolint:funlen
	gql := graphQLdata{}
	assert.NoError(t, json.Unmarshal([]byte(prResultFixture), &gql))
	prs := gql.Data.Repository.PullRequests.PRNodes
	config := ChartDataConfig{}
	got := ChartDataFromPRs(prs, config)
	expected := ChartData{
		ID: chartID(config),
		Authors: map[string]Author{
			"kasunfdo": {
				Login:     "kasunfdo",
				URL:       "https://github.com/kasunfdo",
				AvatarURL: "https://avatars3.githubusercontent.com/u/12541716?v=4",
			},
			"marcelocantos": {
				Login:     "marcelocantos",
				URL:       "https://github.com/marcelocantos",
				AvatarURL: "https://avatars3.githubusercontent.com/u/215143?v=4",
			},
			"dummy": {
				Login:     "dummy",
				URL:       "https://github.com/dummy",
				AvatarURL: "https://avatars3.githubusercontent.com/u/12541719?v=4",
			},
			"dummy2": {
				Login:     "dummy2",
				URL:       "https://github.com/dummy2",
				AvatarURL: "https://avatars3.githubusercontent.com/u/12541729?v=4",
			},
		},
		Charts: []Chart{
			{
				Title: "Merged Pull Requests",
				Points: []Point{
					{Author: "kasunfdo", Count: 1},
				},
				TotalCount: 1,
				MaxCount:   1,
			},
			{
				Title: "Merged or Open Pull Requests",
				Points: []Point{
					{Author: "dummy2", Count: 2},
					{Author: "dummy", Count: 1},
					{Author: "kasunfdo", Count: 1},
				},
				TotalCount: 4,
				MaxCount:   2,
			},
			{
				Title: "Code Reviews",
				Points: []Point{
					{Author: "marcelocantos", Count: 1},
				},
				TotalCount: 1,
				MaxCount:   1,
			},
			{
				Title: "Code Review Comments",
				Points: []Point{
					{Author: "marcelocantos", Count: 6},
				},
				TotalCount: 6,
				MaxCount:   6,
			},
		},
	}
	assert.Equal(t, expected, got)

	config = ChartDataConfig{
		LabelRegexp: regexp.MustCompile("^lab"),
	}
	got = ChartDataFromPRs(prs, config)
	expected.ID = chartID(config)
	expected.Config = config
	assert.Equal(t, expected, got)
}

func TestSkipWithoutLabLabel(t *testing.T) {
	gql := graphQLdata{}
	assert.NoError(t, json.Unmarshal([]byte(prResultFixtureForFilters), &gql))
	prs := gql.Data.Repository.PullRequests.PRNodes
	config := ChartDataConfig{
		LabelRegexp: regexp.MustCompile("^lab"),
	}
	got := ChartDataFromPRs(prs, config)
	expected := ChartData{
		ID:      chartID(config),
		Authors: map[string]Author{},
		Charts:  []Chart{},
		Config:  config,
	}
	assert.Equal(t, expected, got)
}

func TestSkipBeforeCreatedAfter(t *testing.T) {
	gql := graphQLdata{}
	assert.NoError(t, json.Unmarshal([]byte(prResultFixtureForFilters), &gql))
	prs := gql.Data.Repository.PullRequests.PRNodes
	createdAt, _ := time.Parse(time.RFC3339, "2019-05-24T11:24:48Z")
	config := ChartDataConfig{
		CreatedAfter: createdAt,
	}
	got := ChartDataFromPRs(prs, config)
	expected := ChartData{
		ID:      chartID(config),
		Authors: map[string]Author{},
		Charts:  []Chart{},
		Config:  config,
	}
	assert.Equal(t, expected, got)
}

const prResultFixture = `{
  "data": {
    "repository": {
      "url": "https://github.com/devopsacademyau/academy",
      "pullRequests": {
        "totalCount": 189,
        "pageInfo": {
          "endCursor": "Y3Vyc29yOnYyOpHOEPPZoQ==",
          "hasNextPage": false
        },
        "nodes": [
          {
            "number": 130,
            "url": "https://github.com/devopsacademyau/academy/pull/130",
            "state": "MERGED",
            "title": "Lab 1 - Fibonacci and Negafibonacci",
            "createdAt": "2019-02-24T11:24:48Z",
            "author": {
              "login": "kasunfdo",
              "url": "https://github.com/kasunfdo",
              "avatarUrl": "https://avatars3.githubusercontent.com/u/12541716?v=4"
            },
            "reviews": {
              "nodes": [
                {
                  "author": {
                    "login": "marcelocantos",
                    "url": "https://github.com/marcelocantos",
                    "avatarUrl": "https://avatars3.githubusercontent.com/u/215143?v=4"
                  },
                  "comments": { "totalCount": 5 }
                },
                {
                  "author": {
                    "login": "kasunfdo",
                    "url": "https://github.com/kasunfdo",
		    "avatarUrl": "https://avatars3.githubusercontent.com/u/12541716?v=4"
                  },
                  "comments": { "totalCount": 5 }
                }
              ]
            },
            "labels": { "nodes": [{"name": "lab1"} ] }
          },
          {
            "number": 132,
            "url": "https://github.com/devopsacademyau/academy/pull/132",
            "state": "OPEN",
            "title": "Dummy - with label",
            "createdAt": "2019-05-24T11:24:48Z",
            "author": {
              "login": "dummy",
              "url": "https://github.com/dummy",
              "avatarUrl": "https://avatars3.githubusercontent.com/u/12541719?v=4"
            },
            "reviews": {
              "nodes": [ ]
            },
            "labels": { "nodes": [{"name": "lab1"} ] }
          },
          {
            "number": 133,
            "url": "https://github.com/devopsacademyau/academy/pull/133",
            "state": "OPEN",
            "title": "Dummy2 - with label",
            "createdAt": "2019-05-24T11:24:48Z",
            "author": {
              "login": "dummy2",
              "url": "https://github.com/dummy2",
              "avatarUrl": "https://avatars3.githubusercontent.com/u/12541729?v=4"
            },
            "reviews": {
              "nodes": [ ]
            },
            "labels": { "nodes": [{"name": "lab1"} ] }
          },
          {
            "number": 134,
            "url": "https://github.com/devopsacademyau/academy/pull/134",
            "state": "OPEN",
            "title": "Dummy2 - lab 2",
            "createdAt": "2019-05-24T11:24:48Z",
            "author": {
              "login": "dummy2",
              "url": "https://github.com/dummy2",
              "avatarUrl": "https://avatars3.githubusercontent.com/u/12541729?v=4"
            },
            "reviews": {
              "nodes": [ ]
            },
            "labels": { "nodes": [{"name": "lab2"} ] }
          }
        ]
      }
    }
  }
}
`

const prResultFixtureForFilters = `{
  "data": {
    "repository": {
      "url": "https://github.com/devopsacademyau/academy",
      "pullRequests": {
        "totalCount": 189,
        "pageInfo": {
          "endCursor": "Y3Vyc29yOnYyOpHOEPPZoQ==",
          "hasNextPage": false
        },
        "nodes": [
          {
            "number": 131,
            "url": "https://github.com/devopsacademyau/academy/pull/131",
            "state": "MERGED",
            "title": "Dummy - no label",
            "createdAt": "2019-05-24T11:24:48Z",
            "author": {
              "login": "dummy",
              "url": "https://github.com/dummy",
              "avatarUrl": "https://avatars3.githubusercontent.com/u/12541719?v=4"
            }
      }
        ]
      }
    }
  }
}
`

const prResultFixtureWithNextPage = `{
  "data": {
    "repository": {
      "url": "https://github.com/devopsacademyau/academy",
      "pullRequests": {
        "totalCount": 189,
        "pageInfo": {
          "endCursor": "Y3Vyc29yOnYyOpHOEPPZoQ==",
          "hasNextPage": true
        },
        "nodes": [ ]
      }
    }
  }
}
`
