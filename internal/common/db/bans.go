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
func (bt *BansTable) GetBan(identifier string) (Ban, error) {
	ban := Ban{}

	err := bt.gDB.First(
		&ban,
		"discord_id=? OR player_id=?",
		identifier, identifier,
	).Error

	return ban, err
}

// Unban a banned player.
func (bt *BansTable) Pardon(banned Ban) error {
	return bt.gDB.Where(
		"discord_id = ? OR player_id = ?",
		banned.DiscordID,
		banned.PlayerID,
	).Delete(&banned).Error
}
