# Документация к API:

_Небольшая справка: авторизованым считается пользователь, который имеет действительный JWT токен в headers по ключу `Authorization`_

## Endpoints:

### /delete_user/

- **Method**: DELETE
- **Description**: удаляет текущего пользователя (для использования нужно быть авторизованым)
- **Data**: Null

#### Answer on successful completion

```http
  success!
```

#### Answer on error completion

redirect `/user/`

---

### /update_user/

- **Method**: PUT
- **Description**: редактирует текущего пользователя (имя + пароль) (для использования нужно быть авторизованым)
- **Data**:

```json
{
	"username": "user1",
	"password": "pass"
}
```

#### Answer on successful completion

```http
  success!
```

#### Answer on error completion

```http
  empty fields
```

---

### /register/

- **Method**: POST
- **Description**: Регистрирует нового пользователя
- **Data**:

```json
{
	"username": "user1",
	"password": "pass"
}
```

#### Answer on successful completion

```http
  success!
```

#### Answer on error completion

```http
  empty fields
```

---

### /login/

- **Method**: POST
- **Description**: Выдаёт jwt токен пользователю с указанным именем и паролем
- **Data**:

```json
{
	"username": "user1",
	"password": "pass"
}
```

#### Answer on successful completion

```http
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
```

#### Answer on error completion

```http
  empty fields
```

---

### /profile/`uuid`

- **Method**: GET
- **Description**: Возвращает данные авторизованного пользователя (если не был предоставлен uuid) + список учётных записей, если привязан аккаунт в discord
- **Data**: Null

#### Answer on successful completion

```json
{
	"user": {
		"UUID": "cc87cdd0-8016-481f-901b-daf52f557d2b",
		"Username": "lucorus",
		"Password": "42c92dd0cf7e75f60924b83587fc56c5ea5d4d6d51976f1b345a1f7d39868395",
		"Discord_server_id": "854253015862607872",
		"Is_authorized": true,
		"Token": "*user_token*"
	},
	"user_ds_account": [
		{
			"UUID": "be544a26-6dc5-4540-97d5-c92dba807854",
			"UserId": 854253015862607900,
			"Points": 1353,
			"LastMessageTime": "2024/03/29/14/11",
			"Payment": 3,
			"UserGuild": {
				"GuildId": 1192624958124343300,
				"GuildName": "test_server",
				"GuildIcon": "https://cdn.discordapp.com/icons/1192624958124343306/0942b72185e15adab40acfd6626ef94f.png?size=1024"
			},
			"Username": "lic0ris",
			"UserIcon": "https://cdn.discordapp.com/avatars/854253015862607872/12caa248e30674935d29126b9ff7bef5.png?size=1024",
			"Exp": 30
		}
	]
}
```

#### Answer on error completion

```json
{
	"user": {
		"UUID": "",
		"Username": "",
		"Password": "",
		"Discord_server_id": "",
		"Is_authorized": false,
		"Token": "",
		"user_ds_account": null
	},
	"user_ds_account": null
}
```

---

### /main_page/`?page=1`

- **Method**: GET
- **Description**: Выводит топ учётных записей по баллам
- **Data**: Null

#### Answer on successful completion

```json
{
  [
    {
      "UUID": "6256e93f-70c2-42da-a8b7-bf62e4b4d1d3",
      "UserId": 854253015862607900,
      "Points":	3,
      "LastMessageTime": "2024/01/05/08/17",
      "Payment": 1,
      "UserGuild":
      {
        "GuildId": 1142194883684212900,
        "GuildName": "test server",
        "GuildIcon": "https://cdn.discordapp.com/icons/1142194883684212888/54b506b474b1a0a492fa22c95af6e369.png?size=1024",
      },
      "Username":	"lic0ris",
      "UserIcon": "https://cdn.discordapp.com/avatars/854253015862607872/12caa248e30674935d29126b9ff7bef5.png?size=1024",
      "Exp": 1
    },
  ]
}
```

#### Answer on error completion

```http
  *error*
```

---

### /guild/_id_

- **Method**: GET
- **Description**: Выводит топ учётных записях на данном discord сервере
- **Data**: Null

#### Answer on successful completion

```json
{
  [
    {
      "UUID": "6256e93f-70c2-42da-a8b7-bf62e4b4d1d3",
      "UserId": 854253015862607900,
      "Points":	3,
      "LastMessageTime": "2024/01/05/08/17",
      "Payment": 1,
      "UserGuild":
      {
        "GuildId": 1142194883684212900,
        "GuildName": "test server",
        "GuildIcon": "https://cdn.discordapp.com/icons/1142194883684212888/54b506b474b1a0a492fa22c95af6e369.png?size=1024",
      },
      "Username":	"lic0ris",
      "UserIcon": "https://cdn.discordapp.com/avatars/854253015862607872/12caa248e30674935d29126b9ff7bef5.png?size=1024",
      "Exp": 1
    },
  ]
}
```

#### Answer on error completion

```json
null
```

---

### /change_token/

- **Method**: PATCH
- **Description**: устанавливает новый токен авторизованному пользователю (для использования нужно быть авторизованым)
- **Data**: Null

#### Answer on successful completion

```json
Success
```

#### Answer on error completion

```json
Not success
```

---

### /authorize/

- **Method**: POST
- **Description**: добавляет асоциацию аккаунта в discord с аккаунтом на сайте
- **Data**:
  ```json
  {
  	"access": "*access_token*",
  	"user": "*user_id*",
  	"token": "*user_token*"
  }
  ```

#### Answer on successful completion

```json
success!
```

#### Answer on error completion

```json
null
```

---

### /anauthorize/

- **Method**: POST
- **Description**: убирает асоциацию аккаунта в discord с аккаунтом на сайте
- **Data**:
  ```json
  {
  	"access": "*access_token*",
  	"user": "*user_id*",
  	"token": "*user_token*"
  }
  ```

#### Answer on successful completion

```json
success!
```

#### Answer on error completion

```json
null
```

_Последние 2 метода могут быть использованы только ботом в дискорде, который имеет токен доступа_
