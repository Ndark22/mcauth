// An administrator can ban a Minecraft player or Discord user which will
// prevent them from using the bot and connection to the server.
package db

import (
	"github.com/jinzhu/gorm"
)

type BansTable struct {
	gDB *gorm.DB
}

type Ban struct {
	DiscordID string `gorm:"column:discord_id;type:text;unique; not null"`
	PlayerID  string `gorm:"column:player_id;type:text;unique; not null"`
}

func (Ban) TableName() string {
	return "bans"
}

func GetBansTable(gDB *gorm.DB) BansTable {
	gDB.AutoMigrate(&Ban{})

	return BansTable{gDB: gDB}
}

// Ban a linked account, this will prevent both the Discord user and player
// from interacting with the authenticator.
func (bt *BansTable) Ban(link LinkedAcc) error {
	ban := Ban{
		DiscordID: link.DiscordID,
		PlayerID:  link.PlayerID,
	}

	return bt.gDB.Create(&ban).Error
}

// The identifier can be their Discord user ID or Minecraft player UUID.
func (bt *BansTable) GetBan(identifier string) (*Ban, error) {
	ban := Ban{}

	err := bt.gDB.First(
		&ban,
		"discord_id=? OR player_id=?",
		identifier, identifier,
	).Error

	if len(ban.DiscordID) == 0 && len(ban.PlayerID) == 0 {
		return nil, err
	}

	return &ban, err
}

func (bt *BansTable) PardonPlayer(playerID string) error {
	return bt.gDB.Where(
		"player_id = ?",
		playerID,
	).Delete(&Ban{}).Error
}

func (bt *BansTable) PardonUser(userID string) error {
	return bt.gDB.Where(
		"discord_id = ?",
		userID,
	).Delete(&Ban{}).Error
}
