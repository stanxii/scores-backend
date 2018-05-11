package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/raphi011/scores"
	"github.com/raphi011/scores/sqlite"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var store = sessions.NewCookieStore([]byte("ultrasecret"))

func initRouter(app app) *gin.Engine {
	var router *gin.Engine
	if app.production {
		router = gin.Default()
	} else {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		router.Use(gin.Recovery())
	}

	teamService := &sqlite.TeamService{DB: app.db}
	userService := &sqlite.UserService{DB: app.db, PW: &scores.PBKDF2PasswordService{
		SaltBytes:  16,
		Iterations: 10000,
	}}
	matchService := &sqlite.MatchService{DB: app.db}
	playerService := &sqlite.PlayerService{DB: app.db}
	statisticService := &sqlite.StatisticService{DB: app.db}
	groupService := &sqlite.GroupService{DB: app.db}

	authHandler := authHandler{playerService: playerService, userService: userService, conf: app.conf}
	playerHandler := playerHandler{playerService: playerService}
	matchHandler := matchHandler{
		matchService:  matchService,
		userService:   userService,
		playerService: playerService,
		teamService:   teamService,
		groupService:  groupService,
	}
	statisticHandler := statisticHandler{statisticService: statisticService}
	groupHandler := groupHandler{
		playerService:    playerService,
		groupService:     groupService,
		statisticService: statisticService,
		matchService:     matchService,
	}
	volleynetHandler := volleynetHandler{}

	router.Use(sessions.Sessions("goquestsession", store))

	router.GET("/volleynet/signup", volleynetHandler.signup)
	router.GET("/volleynet/tournaments", volleynetHandler.allTournaments)
	router.GET("/volleynet/tournaments/:tournamentID", volleynetHandler.tournament)
	router.GET("/volleynet/players/search", volleynetHandler.searchPlayers)

	router.GET("/userOrLoginRoute", authHandler.loginRouteOrUser)
	router.GET("/auth", authHandler.googleAuthenticate)
	router.POST("/pwAuth", authHandler.passwordAuthenticate)

	auth := router.Group("/")
	auth.Use(authRequired())
	{
		auth.POST("/logout", authHandler.logout)

		auth.GET("/groups", groupHandler.index)
		auth.GET("/groups/:groupID/matches", matchHandler.index)
		auth.POST("/groups/:groupID/matches", matchHandler.matchCreate)
		auth.GET("/groups/:groupID/players", playerHandler.playerIndex)
		auth.GET("/groups/:groupID/playerStatistics", statisticHandler.players)
		auth.GET("/groups/:groupID/teamStatistics", statisticHandler.players)
		auth.GET("/groups/:groupID", groupHandler.groupShow)

		auth.GET("/matches/:matchID", matchHandler.matchShow)
		auth.DELETE("/matches/:matchID", matchHandler.matchDelete)

		auth.POST("/players", playerHandler.playerCreate)
		auth.GET("/players/:playerID/playerStatistics", statisticHandler.player)
		auth.GET("/players/:playerID/matches", matchHandler.byPlayer)
		auth.GET("/players/:playerID/teamStatistics", statisticHandler.playerTeams)
	}

	return router
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user-id")

		if userID == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func jsonn(c *gin.Context, code int, data interface{}, message string) {

	out, _ := json.Marshal(gin.H{
		"status":  code,
		"message": message,
		"data":    data,
	})

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(out)))

	c.Writer.WriteHeader(code)
	c.Writer.Write(out)
}
