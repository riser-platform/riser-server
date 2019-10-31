package model

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Ideally these rules would be in pkg/... for reuse with the service layer but this causes a circular dependency.
// Most validation happens in the API model so this works for now.

func RulesAppName() []validation.Rule {
	rules := []validation.Rule{
		// Max length takes into account possible deployment name suffixes
		validation.RuneLength(3, 50),
	}
	return append(rules, RulesNamingIdentifier()...)
}

// RulesNamingIdentifier returns rules for naming things (e.g. an app, stage) that are RFC 1035 subdomain compatible.
func RulesNamingIdentifier() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		validation.RuneLength(1, 63),
		// Change with care as we use naming identifiers for DNS names and this conforms to RFC 1035
		// Note that depending on the TLD the spec allows for more characters than allowed below. This restriction is
		// designed for maximum portability and simplicity.
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9-]+$")).Error("must be lowercase, alphanumeric, and start with a letter"),
	}
}
