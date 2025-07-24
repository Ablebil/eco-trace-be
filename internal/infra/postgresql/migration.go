package postgresql

import (
	"github.com/Ablebil/eco-sample/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
		&entity.Challenge{},
		&entity.UserChallenge{},
		&entity.Badge{},
		&entity.UserBadge{},
	)
}
