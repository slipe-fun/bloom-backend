package config

import "regexp"

var UsernameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
