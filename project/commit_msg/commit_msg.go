package commit_msg

import (
	"fmt"
	"regexp"
	"strings"
)

var RxCommit = regexp.MustCompile("^(FIX|TASK)(\\s+)?([A-Z]+-[0-9]+)?")

func CheckCommit(commit string) error {
	lines := strings.Split(commit, "\n")
	matched := RxCommit.FindStringSubmatch(lines[0])
	if len(matched) > 0 {
		summary := strings.TrimSpace(RxCommit.ReplaceAllString(lines[0], ""))
		if len(summary) < 6 {
			return fmt.Errorf("commit summary is too short, at least 6 characters, now %d", len(summary))
		}
		return nil
	}
	return fmt.Errorf("commit format incorrect: %s\nmust be (FIX|TASK) ([A-Z]+-[0-9]+)?", commit)
}
