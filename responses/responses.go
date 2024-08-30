package responses

import (
	"SiteForDsBot/models"
)

type ProfileResponse struct {
    User           *models.User          `json:"user"`
    UserDsAccounts []models.DsBotUser    `json:"user_ds_account"`
}

