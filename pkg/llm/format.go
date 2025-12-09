package llm

import "strings"

func formatAnswer(answer string) string {
	answer = strings.TrimSpace(answer)
	if after, ok := strings.CutPrefix(answer, "\n\n"); ok {
		answer = after
	}
	return answer
}
