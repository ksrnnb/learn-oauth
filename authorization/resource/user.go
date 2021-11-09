package resource

import "errors"

type User struct {
	Id int
	Email string
	Password string
	PictureUrl string
}

func usersInMemory() []*User {
	return []*User{
		{
			Id: 1,
			Email: "test@example.com",
			Password: "3jd8Ge30Qcw2h",
			PictureUrl: "https://hogehoge********.com/83hagbfahaaeg",
		},
		{
			Id: 2,
			Email: "test2@example.com",
			Password: "09H$w63hdiHEDd",
			PictureUrl: "https://hogehoge2********.com/83hagbfahaaeg",
		},
	}
}

// メールアドレスとパスワードが一致するユーザーを探す
func FindUser(email string, password string) (*User, error) {
	users := usersInMemory()

	for _, user := range users {
		if (user.Email != email) {
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