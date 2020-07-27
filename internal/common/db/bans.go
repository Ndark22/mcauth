package db

import (
	"github.com/jinzhu/gorm"
)

type BansTable struct {
	gDB *gorm.DB
}

type Ban struct {
	DiscordID string `gorm:"column:discord_id;type:text;unique"`
	PlayerID  string `gorm:"column:player_id;type:text;unique"`
}

func (Ban) TableName() string {
	return "bans"
}

func GetBansTable(gDB *gorm.DB) BansTable {
	gDB.AutoMigrate(&Ban{})

	return BansTable{gDB: gDB}
}

func (bt *BansTable) BanLink(link LinkedAcc) error {
	ban := Ban{
		DiscordID: link.DiscordID,
		PlayerID:  link.PlayerID,
	}

	return bt.gDB.Create(&ban).Error
}

func (bt *BansTable) BanPlayer(playerID string) error {
	ban := Ban{
		PlayerID: playerID,
	}

	return bt.gDB.Create(&ban).Error
}

func (bt *BansTable) BanUser(userID string) error {
	ban := Ban{
		DiscordID: userID,
	}

	return bt.gDB.Create(&ban).Error
}

// The identifier can be their Discord user ID or Minecraft player UUID
func (bt *BansTable) GetBanned(identifier string) (Ban, error) {
	ban := Ban{}

	err := bt.gDB.First(
		&ban,
		"discord_id=? OR player_id=?",
		identifier, identifier,
	).Error

	return ban, err
}

func (bt *BansTable) Pardon(banned Ban) error {
	return bt.gDB.Where(
		"discord_id = ? OR player_id = ?",
		banned.DiscordID,
		banned.PlayerID,
	).Delete(&banned).Error
}
