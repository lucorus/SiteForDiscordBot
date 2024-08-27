package models

import (
	"SiteForDsBot/conf"
	sql "database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

// открывает подлючение к бд
func CreateConnectToDiscordBotBD() (*sql.DB, error) {
  db, err := sql.Open("postgres", "user=" + conf.DB_user_ds_bot + " password=" + conf.DB_password_ds_bot + " host=" + 
	conf.DB_host_ds_bot + " dbname=" + conf.DB_name_ds_bot + " sslmode=disable")
  if err != nil {
      fmt.Println("Error: Unable to connect to database:", err)
      return nil, err
  }
  return db, nil
}

// закрывает подключение к бд
func CloseConnectToDiscordBotBD(db *sql.DB) {
  if db != nil {
    db.Close()
  }
}


type Guild struct {
	GuildId int64
	GuildName string
	GuildIcon string
}

// модель пользователя, которая используется для пользователей в боте (1 сервер - 1 модель)
type DsBotUser struct {
	UUID string
	UserId int64
	Points int64
	LastMessageTime string
	Payment int64
	UserGuild Guild // сервер, на котоом зарегистрирован данный пользователь
	Username string // ник в дискорде
	UserIcon string // иконка пользователя в дискорде
	Exp int64
}

type JoinDsBotUserWithGuild struct {
	DsBotUser
	Guild
}

// Возвращает информацию о пользователях, сортируя её по кол-ву баллов
func AllDsBotUsers(limit string) ([]DsBotUser, error) {
	db, err := CreateConnectToDiscordBotBD()
	if err != nil {
	    return nil, fmt.Errorf("Error: %v", err)
	}

	query := "SELECT * FROM users JOIN guilds ON users.guild = guild_id ORDER BY users.points DESC"
	if limit != "0" {
		query += " LIMIT " + limit
	}

  rows, err := db.Query(query)
  if err != nil {
    return nil, fmt.Errorf("error querying users: %v", err)
  }
  defer rows.Close()
	CloseConnectToDiscordBotBD(db)

  var users []DsBotUser
  for rows.Next() {
  	var user DsBotUser
  	if err := rows.Scan(&user.UUID, &user.UserId, &user.Points, &user.LastMessageTime, &user.Payment, &user.UserGuild.GuildId, 
			&user.Username, &user.UserIcon, &user.Exp, &user.UserGuild.GuildId, &user.UserGuild.GuildName, &user.UserGuild.GuildIcon); err != nil {
    	return nil, fmt.Errorf("error scanning user: %v", err)
  	}
  	users = append(users, user)
  }
  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("error iterating rows: %v", err)
  }

  return users, nil
}

// выводит все учётные записи пользователя с переданным discord_id
func FindDsBotUsers(discord_user_id string) ([]DsBotUser, error) {
	db, err := CreateConnectToDiscordBotBD()
	if err != nil {
	    return nil, fmt.Errorf("Error: %v", err)
	}

	user_id, err := strconv.Atoi(discord_user_id)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при преобразовании типа")
	}

	query := "SELECT * FROM users JOIN guilds ON users.guild = guild_id WHERE user_id = $1 ORDER BY users.points DESC"
  rows, err := db.Query(query, user_id)
  if err != nil {
    return nil, fmt.Errorf("error querying users: %v", err)
  }
  defer rows.Close()
	CloseConnectToDiscordBotBD(db)

  var users []DsBotUser
  for rows.Next() {
  	var user DsBotUser
  	if err := rows.Scan(&user.UUID, &user.UserId, &user.Points, &user.LastMessageTime, &user.Payment, &user.UserGuild.GuildId, 
			&user.Username, &user.UserIcon, &user.Exp, &user.UserGuild.GuildId, &user.UserGuild.GuildName, &user.UserGuild.GuildIcon); err != nil {
    	return nil, fmt.Errorf("error scanning user: %v", err)
  	}
  	users = append(users, user)
  }
  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("error iterating rows: %v", err)
  }

  return users, nil
}


// выводит все учётные записи пользователей на указанном дискорд сервере
func ListUsersInGuild(GuildId string) ([]DsBotUser, error) {
	db, err := CreateConnectToDiscordBotBD()
	if err != nil {
	    return nil, fmt.Errorf("Error: %v", err)
	}

	guild_id, err := strconv.Atoi(GuildId)
	if err != nil {
		return nil, fmt.Errorf("Error in change type")
	}

	query := "SELECT * FROM users JOIN guilds ON users.guild = guild_id WHERE guild_id = $1 ORDER BY users.points DESC"
  rows, err := db.Query(query, guild_id)
  if err != nil {
    return nil, fmt.Errorf("error querying users: %v", err)
  }
  defer rows.Close()
	CloseConnectToDiscordBotBD(db)

  var users []DsBotUser
  for rows.Next() {
  	var user DsBotUser
  	if err := rows.Scan(&user.UUID, &user.UserId, &user.Points, &user.LastMessageTime, &user.Payment, &user.UserGuild.GuildId, 
			&user.Username, &user.UserIcon, &user.Exp, &user.UserGuild.GuildId, &user.UserGuild.GuildName, &user.UserGuild.GuildIcon); err != nil {
    	return nil, fmt.Errorf("error scanning user: %v", err)
  	}
  	users = append(users, user)
  }
  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("error iterating rows: %v", err)
  }

  return users, nil
}

