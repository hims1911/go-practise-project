package dbmodels

import (
	"context"
	"database/sql"
	"golang-practise-project/utils/database"
	"log"
)

const (
	UPDATE = "update"
)

type UserDB struct {
	UserId           int64          `gorm:"primaryKey, column:user_id"`
	AreaCode         sql.NullInt64  `gorm:"column:area_code"`
	CurrentLoginTime sql.NullInt64  `gorm:"column:current_login_time"`
	GameMode         sql.NullString `gorm:"column:game_mode"`
}

func (UserDB) TableName() string {
	return "user"
}

// SaveUserDetails will save the user record for the first timed user
func (t *UserDB) SaveUserDetails(ctx context.Context) error {
	err := database.Get().WithContext(ctx).Create(&t).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserDetails will update the UserDetails
func (t *UserDB) UpdateUserDetails(ctx context.Context) error {
	log.Println(t)
	err := database.Get().WithContext(ctx).Debug().Model(&UserDB{}).Where("user_id = ?", t.UserId).Updates(&t).Error
	if err != nil {
		return err
	}

	return nil
}

// FindUserDetailsByID will give the user details based on the UserId
func FindUserDetailsByID(ctx context.Context, userId string) (UserDB, error) {
	var user UserDB
	err := database.Get().WithContext(ctx).Model(UserDB{}).Where("user_id = ?", userId).Take(&user).Error
	if err != nil {
		return UserDB{}, err
	}
	return user, nil
}

// FindHighestPlayingModeByAreaCode it will return the highest playing mode
func FindHighestPlayingModeByAreaCode(ctx context.Context, areaCode int64) (gameMode string, err error) {
	err = database.Get().Raw("with a as (\n    select user.game_mode, count(user.game_mode) as count\n    from user\n    where game_mode != \"\" and area_code=?\n    group by user.game_mode\n    order by count desc\n    limit 1\n)\nselect game_mode from a;", areaCode).Scan(&gameMode).Error
	if err != nil {
		log.Println("Error Occurred While Fetching The Highest Playing Game Mode")
		return
	}
	return
}
