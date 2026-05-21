package jobman

import "strings"

type Filter struct {
	Keyword string
	Status  string
}

func JoinRequirements(requirements []string) string {
	cleaned := make([]string, 0, len(requirements))
	for _, requirement := range requirements {
		requirement = strings.TrimSpace(requirement)
		if requirement == "" {
			continue
		}
		cleaned = append(cleaned, requirement)
	}
	return strings.Join(cleaned, "\n")
}

func SplitRequirements(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}

	parts := strings.Split(raw, "\n")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		items = append(items, part)
	}
	return items
}
