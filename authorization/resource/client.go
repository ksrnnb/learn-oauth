package resource

import "errors"

type Client struct {
	Name string
	ClientId string
	ClientSecret string
	RedirectUri string
}

func clientsInMemory() []*Client {
	return []*Client{
		{
			Name: "ぴよぴよ",
			ClientId: "abcde12345",
			ClientSecret: "abcde12345secert",
			RedirectUri: "http://localhost:3000/callback",
		},
		{
			Name: "dummy",
			ClientId: "dummy",
			ClientSecret: "dummy",
			RedirectUri: "dummy",
		},
	}
}

func FindClient(clientId string) (*Client, error) {
	clients := clientsInMemory()

	for _, client := range clients {
		if (client.ClientId != clientId) {
			continue
		}

		return client, nil
	}

	return nil, errors.New("client not found")
}