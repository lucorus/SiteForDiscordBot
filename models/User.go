package models

import (
	"SiteForDsBot/conf"
	"crypto/sha256"
	sql "database/sql"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// создаёт бд, если её ещё нет
func CreateDB() error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        uuid VARCHAR(100) PRIMARY KEY,
        username VARCHAR(100) NOT NULL,
        password VARCHAR(100) NOT NULL,
        discord_server_id VARCHAR(50),
        is_authorized BOOL DEFAULT False,
        token VARCHAR(100) UNIQUE
    );
    `
	db, err := CreateConnect()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	_, err = db.Exec(createTableSQL)
    if err != nil {
      return fmt.Errorf("Error creating table: %v", err)
    }
	CloseConnect(db)
	return nil
}

// открывает подлючение к бд
func CreateConnect() (*sql.DB, error) {
    db, err := sql.Open("postgres", "user=" + conf.DB_user + " password=" + conf.DB_password + " host=" + conf.DB_host +
    " dbname=" + conf.DB_name + " sslmode=disable")
    if err != nil {
        fmt.Println("Error: Unable to connect to database:", err)
        return nil, err
    }
    return db, nil
}

// закрывает подключение к бд
func CloseConnect(db *sql.DB) {
    if db != nil {
        db.Close()
    }
}


type User struct {
	UUID    string
	Username string
	Password string   // пароль будет храниться в зашифрованом виде
    Discord_server_id string // id аккаунта в discord
    Is_authorized bool // имеет ли пользователь связь с аккаунтом в дискорде
    Token string // уникальный токен, с помощью которого пользователь будет привязывать аккаунт на сайте с акком в ds
}


// кодирует переданный пароль в sha256
func EncodePassword(password string) string {
	hasher := sha256.New()
	data := password + conf.Salt
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}


// Создаёт нового пользователя
func NewUser(username, password string) error {
    db, err := CreateConnect()
    if err != nil {
        return fmt.Errorf("error creating connection: %v", err)
    }

    query := "INSERT INTO users (uuid, username, password, discord_server_id, is_authorized, token) VALUES ($1, $2, $3, $4, $5, $6)"

    _, err = db.Exec(query, uuid.New().String(), username, EncodePassword(password), "", false, uuid.New().String())

    if err != nil {
        return fmt.Errorf("error inserting user: %v", err)
    }
	CloseConnect(db) 
    return nil
}


func LoginUser(username, password string) (string, error) {
    db, err := CreateConnect()
    if err != nil {
        return "", fmt.Errorf("error creating connection: %v", err)
    }

    query := "SELECT uuid FROM users WHERE username = $1 AND password = $2"

    rows, err := db.Query(query, username, EncodePassword(password))
    if err != nil {
        return "", fmt.Errorf("error querying users: %v", err)
    }
    defer rows.Close()
	var user string
    for rows.Next() {
        if err := rows.Scan(&user); err != nil {
            return "", fmt.Errorf("error scanning user: %v", err)
        }
    }
    if err := rows.Err(); err != nil {
        return "", fmt.Errorf("error iterating rows: %v", err)
    }
    return user, nil
}


// УДаляет пользователя с переданным uuid
func DeleteUser(uuid string) error {
    db, err := CreateConnect()
    if err != nil {
        return fmt.Errorf("error creating connection: %v", err)
    }

    _, err = db.Query("DELETE FROM users WHERE uuid = $1", uuid)
    if err != nil {
        return fmt.Errorf("error inserting user: %v", err)
    }

	CloseConnect(db) 
    return nil
}


// устанавливает значения пароля и ника для указанного юзера
func UpdateUser(username, password, UserUUID string) bool {
    db, err := CreateConnect()
	if err != nil {
	    return false
	}

    query := "UPDATE users SET username = $1, password = $2 WHERE uuid = $3"
    _, err = db.Exec(query, username, EncodePassword(password), UserUUID)
    if err != nil {
        return false
    }
    CloseConnect(db)

    return true
}


func All() ([]User, error) {
	db, err := CreateConnect()
	if err != nil {
	    return nil, fmt.Errorf("Error: %v", err)
	}
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        return nil, fmt.Errorf("error querying users: %v", err)
    }
    defer rows.Close()
	CloseConnect(db)

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.UUID, &user.Username, &user.Password, &user.Discord_server_id, &user.Is_authorized, 
            &user.Token); err != nil {
            return nil, fmt.Errorf("error scanning user: %v", err)
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %v", err)
    }
    return users, nil
}


// Возвращаает пользователя с указанным uuid
func Find(UUID string) (*User, error) {
	db, err := CreateConnect()
	if err != nil {
	    return nil, fmt.Errorf("Error: %v", err)
	}
    rows, err := db.Query("SELECT * FROM users WHERE uuid = $1;", UUID)
    if err != nil {
        return nil, fmt.Errorf("error querying users: %v", err)
    }
    defer rows.Close()
    CloseConnect(db)

	var user User
    for rows.Next() {
        if err := rows.Scan(&user.UUID, &user.Username, &user.Password, &user.Discord_server_id, &user.Is_authorized, 
        &user.Token); err != nil {
            return nil, fmt.Errorf("error scanning user: %v", err)
        }
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %v", err)
    }
    return &user, nil
}


// добавляет асоциацию пользователя на сайте с аккаунтом в ds (true - успешное выполнение операции)
func Authorize(Token string, UserIdInDiscord int) bool {
	db, err := CreateConnect()
	if err != nil {
	    return false
	}

    rows, err := db.Query("SELECT * FROM users WHERE discord_server_id = $1;", UserIdInDiscord)
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.UUID, &user.Username, &user.Password, &user.Discord_server_id, &user.Is_authorized, 
            &user.Token); err != nil {
            return false
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        return false
    }

    if len(users) > 0 {
        // если уже есть пользователи, которые привязаны к данному аккаунту, то отменяем привязку
        return false
    }

    query := "UPDATE users SET discord_server_id = $1, is_authorized = $2 WHERE token = $3"
    _, err = db.Exec(query, UserIdInDiscord, true, Token)
    if err != nil {
        return false
    }
    CloseConnect(db)

    return true
}


// убирает асоциацию всех пользователей на сайте с аккаунтом в ds (true - успешное выполнение операции)
func AnAuthorize(UserIdInDiscord int) bool {
	db, err := CreateConnect()
	if err != nil {
	    return false
	}

    query := "UPDATE users SET discord_server_id = $1, is_authorized = $2 WHERE discord_server_id = $3"
    _, err = db.Exec(query, "", false, UserIdInDiscord)
    if err != nil {
        return false
    }
    CloseConnect(db)

    return true
}


// изменяет токен пользователя с переданным uuid
func ChangeToken(UserUUID string) bool {
    db, err := CreateConnect()
	if err != nil {
	    return false
	}

    query := "UPDATE users SET token = $1 WHERE uuid = $2"
    _, err = db.Exec(query, uuid.New().String(), UserUUID)
    if err != nil {
        return false
    }
    CloseConnect(db)

    return true
}
