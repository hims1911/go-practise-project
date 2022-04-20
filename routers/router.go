package routers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"golang-practise-project/buisness"
	"golang-practise-project/models"
	"log"
	"strconv"
)

// GetRouter Returns the router
func GetRouter(ctx context.Context) *gin.Engine {
	router := gin.New()

	server := socketio.NewServer(nil)

	// connecting with socket we can add authentication layer over the socket
	server.OnConnect("/", func(s socketio.Conn) error {

		// fetching the query values for the particular user
		url := s.URL()
		userId := url.Query().Get("user_id")

		// check In Redis If UserId is present or Not
		// for first time user just add it to the DB and Redis

		err := buisness.UserExistOrNot(ctx, userId)
		if err != nil {
			log.Println("UserId is not available ", userId)
			return err
		}

		// setting up the context
		s.SetContext(userId)

		log.Println("connected to the server:", s.ID())
		log.Println("user_id received ", userId)

		return nil
	})

	// Event fetchPopularMode will fetch the popular mode
	server.OnEvent("/", "fetchPopularMode", func(s socketio.Conn, msg string) {
		log.Println("fetchPopularMode:", msg)

		// fetching the query values for the particular user
		url := s.URL()
		userId, _ := strconv.Atoi(url.Query().Get("user_id"))
		log.Println("Session ID and Log ID ", s.ID())
		log.Println("fetching the currently popular mode for user ID ", userId)

		// converting the string into struct
		var currentUserData models.CurrentPlayerSocketMessage
		json.Unmarshal([]byte(msg), &currentUserData)

		mode, err := buisness.GetHighestGameMode(ctx, currentUserData)
		if err != nil {
			log.Println("Error While Fetching the highest gaming modes")
			s.Emit("fetchPopularMode", err)
		}

		log.Println("got the mode as ", mode)

		if mode == "" {
			s.Emit("fetchPopularMode", "Select Any")
		} else {
			s.Emit("fetchPopularMode", mode)
		}
	})

	// startPlayingMode will start the gamingMode and it will update the same in the DB
	server.OnEvent("/", "startPlayingMode", func(s socketio.Conn, msg string) {
		log.Println("startPlayingMode:", msg)

		// fetching the query values for the particular user
		url := s.URL()
		userId, _ := strconv.Atoi(url.Query().Get("user_id"))

		// converting the string into struct
		var userPlayingData models.UserGameMessage
		json.Unmarshal([]byte(msg), &userPlayingData)

		userPlayingModeStruct := models.User{
			UserGame: userPlayingData,
			UserId:   int64(userId),
		}

		// saving the game mode to the DB
		err := buisness.SetGameMode(ctx, userPlayingModeStruct)
		if err != nil {
			log.Println("Error occurred while saving the current game mode")
			s.Emit("startPlayingMode", "Error Occurred While Setting Up The Mode ", err)
		}
		s.Emit("startPlayingMode", "Mode Has Been Setup ", userPlayingModeStruct.UserGame.GameMode)
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// fetching the query values for the particular user
		url := s.URL()
		userId := url.Query().Get("user_id")

		err := buisness.CloseCurrentlyPlayingMode(ctx, userId)
		if err != nil {
			log.Println("error while closing the mode and updating in the DB")
		}
		log.Println("closed "+s.ID(), reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	defer server.Close()

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	if err := router.Run(":8000"); err != nil {
		log.Fatal("failed run app: ", err)
	}

	return router
}
