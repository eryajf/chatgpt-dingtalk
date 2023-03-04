module github.com/eryajf/chatgpt-dingtalk

go 1.18

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/solywsh/chatgpt v0.0.14
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/sashabaranov/go-gpt3 v1.3.0 // indirect
	golang.org/x/net v0.0.0-20211029224645-99673261e6eb // indirect
)

replace github.com/solywsh/chatgpt => ./pkg/chatgpt
