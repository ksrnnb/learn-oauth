package resource

type AllowList struct {
	ClientUsers []*ClientUser
}

type ClientUser struct {
	ClientId string
	UserId   int
}

var allowList AllowList

// 許可リストに加える
func AddAllowListIfNeeded(clientId string, userId int) {
	if allowList.isDuplicate(clientId, userId) {
		return
	}

	allowList.Add(clientId, userId)
}

// 重複しているかどうか
func (a *AllowList) isDuplicate(clientId string, userId int) bool {
	for _, clientUser := range a.ClientUsers {
		if clientUser.ClientId != clientId {
			continue
		}

		if clientUser.UserId != userId {
			continue
		}

		return true
	}

	return false
}

// 許可リストに追加する
func (a *AllowList) Add(clientId string, userId int) {
	clientUser := &ClientUser{
		ClientId: clientId,
		UserId:   userId,
	}

	a.ClientUsers = append(a.ClientUsers, clientUser)
}
