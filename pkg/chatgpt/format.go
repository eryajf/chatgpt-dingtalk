package chatgpt

func formatAnswer(answer string) string {
	for len(answer) > 0 {
		if answer[:1] == "\n" || answer[0] == ' ' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return answer
}
