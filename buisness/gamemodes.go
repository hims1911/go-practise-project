package buisness

import (
	"context"
	"database/sql"
	"errors"
	dbmodels "golang-practise-project/entity"
	"golang-practise-project/models"
	"gorm.io/gorm"
	"log"
	"strconv"
)

// SetGameMode will set the gameMode user is playing currently
func SetGameMode(ctx context.Context, userModel models.User) error {
	log.Print("received call to SetGameMode", userModel)

	var userDB dbmodels.UserDB

	userDB.UserId = userModel.UserId
	userDB.CurrentLoginTime = sql.NullInt64{
		Int64: userModel.UserGame.CurrentLoginTime,
		Valid: true,
	}
	userDB.AreaCode = sql.NullInt64{
		Int64: userModel.UserGame.AreaCode,
		Valid: true,
	}

	userDB.GameMode = sql.NullString{
		String: userModel.UserGame.GameMode,
		Valid:  true,
	}

	// saving in db or redis
	err := userDB.UpdateUserDetails(ctx)
	if err != nil {
		log.Println("Error occurred while Saving the UserDetails")
		return err
	}

	return nil
}

// GetHighestGameMode will fetch the Currently Popular Mode
func GetHighestGameMode(ctx context.Context, currentUser models.CurrentPlayerSocketMessage) (data string, err error) {
	// fetching it from the db
	data, err = dbmodels.FindHighestPlayingModeByAreaCode(ctx, currentUser.AreaCode)
	if err != nil {
		log.Println("error occurred while fetching ", err)
		return
	}

	return data, nil
}

// CloseCurrentlyPlayingMode will close the Currently playing User Session On Disconnect
func CloseCurrentlyPlayingMode(ctx context.Context, userId string) error {
	log.Print("Closing Currently playing mode for user ", userId)

	var userDB dbmodels.UserDB
	userDB, err := dbmodels.FindUserDetailsByID(ctx, userId)
	if err != nil {
		log.Println("Error Occurred While Fetching The UserData ", userId)
		return err
	}

	userDB.CurrentLoginTime = sql.NullInt64{
		Int64: 0,
		Valid: true,
	}

	userDB.AreaCode = sql.NullInt64{
		Int64: 0,
		Valid: true,
	}

	userDB.GameMode = sql.NullString{
		String: "",
		Valid:  true,
	}

	err = userDB.UpdateUserDetails(ctx)
	if err != nil {
		log.Println("Error Occured While Updating The DB for UserData ", userId)
		return err
	}

	return nil
}

// UserExistOrNot checks if User Exist or Not - If not then then it will create and add it to redis as well
func UserExistOrNot(ctx context.Context, userId string) error {
	log.Println("Check If user exist or not")

	// checking in db
	_, err := dbmodels.FindUserDetailsByID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Record Not Exist For User ID ", userId)
			var userDB dbmodels.UserDB
			userID, _ := strconv.Atoi(userId)
			userDB.UserId = int64(userID)

			userDB.GameMode = sql.NullString{
				String: "",
				Valid:  true,
			}

			// saving in db or redis
			err1 := userDB.SaveUserDetails(ctx)
			if err1 != nil {
				log.Println("Error occurred while Saving the UserDetails")
				return err1
			}
		} else {
			log.Println("Error Occurred While Fetching the User Data ", userId)
			return err
		}
	}

	return nil
}
