package authdb

//existing user

var UserExists = "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"

var RegisterUser = "INSERT INTO users (username, password) VALUES (?, ?)"

var LoginUser = "SELECT password, role, uuid FROM users Where username = ?"

//admin

var AdminGetUsers = "SELECT username, role FROM users"

var AdminGetUuid = "SELECT uuid FROM users WHERE username = ? AND role = ?"

var AdminUpdateRole = "UPDATE users SET role = ? WHERE uuid = ?"
