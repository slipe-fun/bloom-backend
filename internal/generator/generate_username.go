package generator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/mozillazg/go-unidecode"
)

var latinLetters = regexp.MustCompile("[^a-zA-Z]+")

func GenerateUsername(fullName string) string {
	parts := strings.Fields(fullName)

	var firstName, lastName string

	if len(parts) > 0 {
		firstName = unidecode.Unidecode(parts[0])
		firstName = strings.ToLower(latinLetters.ReplaceAllString(firstName, ""))
	}
	if len(parts) > 1 {
		lastName = unidecode.Unidecode(parts[1])
		lastName = strings.ToLower(latinLetters.ReplaceAllString(lastName, ""))
	}

	if firstName == "" {
		firstName = "user"
	}
	if lastName == "" {
		lastName = "user"
	}

	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(10000)
	randomStr := fmt.Sprintf("%04d", randomNum)

	return fmt.Sprintf("%s_%s_%s", firstName, lastName, randomStr)
}
