package postgresdao

type User struct {
	UserId   int
	UserName string
	Email    string
}

func (dao *DAO) GetUserById(userId int) (*User, error) {
	user := &User{}
	err := dao.db.QueryRow("SELECT user_id, user_name, email FROM users WHERE user_id = $1", userId).Scan(&user.UserId, &user.UserName, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *DAO) AddUser(user *User) error {
	_, err := dao.db.Exec("INSERT INTO users (user_name, email) VALUES ($1, $2)", user.UserName, user.Email)
	return err
}
