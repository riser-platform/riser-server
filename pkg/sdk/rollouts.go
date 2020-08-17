package sdk

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/riser-platform/riser-server/api/v1/model"
)

const wildcardPercentCharacter = "*"
const wildcardPercentValue = math.MaxInt32

var trafficRuleExp = regexp.MustCompile(`r([0-9]+):(\*|[0-9]+)`)

type RolloutsClient interface {
	Save(deploymentName, namespace, envName string, trafficRule ...string) error
}

type rolloutsClient struct {
	client *Client
}

func (c *rolloutsClient) Save(deploymentName, namespace, envName string, trafficRules ...string) error {
	parsedRules, err := parseTrafficRules(trafficRules...)
	if err != nil {
		return err
	}

	rolloutRequest := model.RolloutRequest{
		Traffic: parsedRules,
	}

	request, err := c.client.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/rollout/%s/%s/%s", envName, namespace, deploymentName), rolloutRequest)
	if err != nil {
		return err
	}

	_, err = c.client.Do(request, nil)
	return err
}

// parseTrafficRules parses the traffic rules in human format "r(rev):(percentage)" into the API model.
// Minimal validation is performed in the client as the server does full validation.
func parseTrafficRules(trafficRules ...string) ([]model.TrafficRule, error) {
	parsedRules := []model.TrafficRule{}

	totalPercent := 0
	wildcardRule := false
	for _, rule := range trafficRules {
		if !trafficRuleExp.MatchString(rule) {
			return nil, errors.New("Rules must be in the format of \"r(rev):(percentage)\" e.g. \"r1:100\" routes 100% of traffic to rev 1")
		}
		ruleSplit := trafficRuleExp.FindStringSubmatch(rule)
		percent := ruleSplit[2]
		percentParsed := 0
		if percent == wildcardPercentCharacter {
			if wildcardRule {
				return nil, errors.New("You may only specify one wildcard rule")
			}
			wildcardRule = true
			percentParsed = wildcardPercentValue
		} else {
			percentParsed = int(mustParseInt(percent))
			totalPercent += percentParsed
		}
		parsedRules = append(parsedRules,
			model.TrafficRule{
				RiserRevision: mustParseInt(ruleSplit[1]),
				Percent:       percentParsed,
			})
	}

	if wildcardRule {
		for idx := range parsedRules {
			if parsedRules[idx].Percent == wildcardPercentValue {
				parsedRules[idx].Percent = 100 - totalPercent
			}
		}
	}

	return parsedRules, nil
}

// mustParseInt panics - validate input before using!
func mustParseInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}
