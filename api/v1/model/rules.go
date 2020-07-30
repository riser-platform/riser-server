package model

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

// Ideally these rules would be in pkg/... for reuse with the service layer but this causes a circular dependency.
// Most validation happens in the API model so this works for now.

func RulesAppName() []validation.Rule {
	rules := []validation.Rule{
		validation.Required,
		// Max length takes into account the RFC 1035 subdomain plus 8 characters reserved for prefix and suffix each.
		validation.RuneLength(3, 47),
	}
	return append(rules, RulesNamingIdentifier()...)
}

// RulesNamingIdentifier returns rules for naming things (e.g. an app, environment) that are RFC 1035 subdomain compatible.
func RulesNamingIdentifier() []validation.Rule {
	return []validation.Rule{
		validation.RuneLength(3, 63),
		// Change with care as we use naming identifiers for DNS names that must conform to RFC 1035
		// Note that depending on the TLD the spec allows for more characters than allowed below. This restriction is
		// designed for maximum portability.
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9-]*[a-z0-9]+$")).Error("must be lowercase, alphanumeric, and start with a letter"),
	}
}
