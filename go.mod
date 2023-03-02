module github.com/eryajf/chatgpt-dingtalk

go 1.18

require (
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/solywsh/chatgpt v0.0.14
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/sashabaranov/go-gpt3 v1.3.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
)

replace github.com/solywsh/chatgpt => ./pkg/chatgpt
