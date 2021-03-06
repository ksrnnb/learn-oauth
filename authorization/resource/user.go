package resource

import "errors"

type User struct {
	Id         string
	Name       string
	Email      string
	Password   string
	PictureUrl string
}

func usersInMemory() []*User {
	return []*User{
		{
			Id:         "ke561cwh1o943",
			Name:       "テストユーザー",
			Email:      "test@example.com",
			Password:   "3jd8Ge30Qcw2h",
			PictureUrl: "https://hogehoge********.com/83hagbfahaaeg",
		},
		{
			Id:         "9dm40h2kd7vgzo",
			Name:       "攻撃者",
			Email:      "attacker@example.com",
			Password:   "09H$w63hdiHEDd",
			PictureUrl: "https://attacker********.com/83hagbfahaaeg",
		},
	}
}

// メールアドレスとパスワードが一致するユーザーを探す
func FindUser(email string, password string) (*User, error) {
	users := usersInMemory()

	for _, user := range users {
		if user.Email != email {
			continue
		}

		// 今回はローカルでの動作のみ想定しているため、ハッシュ化していない
		if user.Password != password {
			return nil, errors.New("user not found")
		}

		return user, nil
	}

	return nil, errors.New("user not found")
}

// ユーザーの存在有無チェック
func ExistsUser(userId string) bool {
	users := usersInMemory()

	for _, user := range users {
		if user.Id == userId {
			return true
		}
	}

	return false
}
