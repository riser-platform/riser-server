package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/riser-platform/riser-server/api/v1/model"
)

var trafficRuleExp = regexp.MustCompile("r([0-9]+):([0-9]+)")

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

func parseTrafficRules(trafficRules ...string) ([]model.TrafficRule, error) {
	parsedRules := []model.TrafficRule{}
	for _, rule := range trafficRules {
		if !trafficRuleExp.MatchString(rule) {
			return nil, errors.New("Rules must be in the format of \"r(rev):(percentage)\" e.g. \"r1:100\" routes 100% of traffic to rev 1")
		}
		ruleSplit := trafficRuleExp.FindStringSubmatch(rule)
		parsedRules = append(parsedRules,
			model.TrafficRule{
				RiserRevision: mustParseInt(ruleSplit[1]),
				Percent:       int(mustParseInt(ruleSplit[2])),
			})
	}

	return parsedRules, nil
}

// mustParseInt panics which should never happen - validate input before using!
func mustParseInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}
