package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/slipe-fun/skid-backend/internal/config"
)

func GenerateNickname() string {
	rand.Seed(time.Now().UnixNano())
	adj := config.Adjectives[rand.Intn(len(config.Adjectives))]
	noun := config.Nouns[rand.Intn(len(config.Nouns))]
	return fmt.Sprintf("%s %s", capitalize(adj), capitalize(noun))
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
