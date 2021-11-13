package resource

type ClientUser struct {
	ClientId string
	UserId   string
}

var authorizationList []*ClientUser

// 認可リストに加える
func AddAuthorizationListIfNeeded(clientId string, userId string) {
	clientUser := &ClientUser{
		ClientId: clientId,
		UserId:   userId,
	}

	if clientUser.isDuplicate() {
		return
	}

	clientUser.add()
}

// 重複しているかどうか
func (c *ClientUser) isDuplicate() bool {
	for _, clientUser := range authorizationList {
		if clientUser.ClientId != c.ClientId {
			continue
		}

		if clientUser.UserId != c.UserId {
			continue
		}

		return true
	}

	return false
}

// 認可リストに追加する
func (c *ClientUser) add() {
	clientUser := &ClientUser{
		ClientId: c.ClientId,
		UserId:   c.UserId,
	}

	authorizationList = append(authorizationList, clientUser)
}
