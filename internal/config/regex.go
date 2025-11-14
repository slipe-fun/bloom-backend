package config

import "regexp"

var UsernameRegex = regexp.MustCompile(`^[a-z][a-z0-9._]*[a-z0-9]$`)
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*\.[a-zA-Z]{2,}$`)
