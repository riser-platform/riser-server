package validation

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

func RulesAppName() []validation.Rule {
	rules := []validation.Rule{
		// Max length takes into account possible deployment name suffixes
		validation.RuneLength(3, 50),
	}
	return append(rules, RulesNamingIdentifier()...)
}

// RulesNamingIdentifier returns a base set of rules for naming things (e.g. an app name, stage name, etc.). These identifiers must be DNS RFC1035 compatible for use as a subdomain
func RulesNamingIdentifier() []validation.Rule {
	return []validation.Rule{
		validation.Required,
		validation.RuneLength(1, 63),
		// Change with care as we use naming identifiers for DNS names and this conforms to RFC 1035
		// Note that depending on the TLD the spec allows for more characters than allowed below. This restriction is
		// designed for maximum portability and simplicity.
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9]+")).Error("must be lowercase, alphanumeric, and start with a letter"),
	}
}
