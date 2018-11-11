package main

import (
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"

	"github.com/raphi011/scores"
	"github.com/raphi011/scores/volleynet/client"
)

type volleynetHandler struct {
	volleynetRepository scores.VolleynetRepository
	userService      *scores.UserService
}

func (h *volleynetHandler) ladder(c *gin.Context) {
	gender := c.DefaultQuery("gender", "M")

	ladder, err := h.volleynetRepository.Ladder(gender)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	jsonn(c, http.StatusOK, ladder, "")
}

func (h *volleynetHandler) allTournaments(c *gin.Context) {
	gender := c.DefaultQuery("gender", "M")
	league := c.DefaultQuery("league", "AMATEUR TOUR")
	season := c.DefaultQuery("season", strconv.Itoa(time.Now().Year()))

	seasonNumber, err := strconv.Atoi(season)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	tournaments, err := h.volleynetRepository.GetTournaments(gender, league, seasonNumber)

	if err != nil {
		log.Print(err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	for _, t := range tournaments {
		t.Teams, _ = h.volleynetRepository.TournamentTeams(t.ID)
	}

	jsonn(c, http.StatusOK, tournaments, "")
}


func (h *volleynetHandler) tournament(c *gin.Context) {
	tournamentID, err := strconv.Atoi(c.Param("tournamentID"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tournament, err := h.volleynetRepository.Tournament(tournamentID)

	if err == scores.ErrorNotFound {
		c.AbortWithError(http.StatusNotFound, err)
	} else if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		tournament.Teams, err = h.volleynetRepository.TournamentTeams(tournament.ID)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		jsonn(c, http.StatusOK, tournament, "")
	}
}

type signupForm struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PartnerID    int    `json:"partnerId"`
	PartnerName  string `json:"partnerName"`
	TournamentID int    `json:"tournamentId"`
	RememberMe   bool   `json:"rememberMe"`
}

func (h *volleynetHandler) signup(c *gin.Context) {
	su := signupForm{}
	if err := c.ShouldBindWith(&su, binding.JSON); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if su.Username == "" ||
		su.Password == "" ||
		su.PartnerID <= 0 ||
		su.TournamentID <= 0 {

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vnClient := client.Default()
	loginData, err := vnClient.Login(su.Username, su.Password)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if su.RememberMe {
		session := sessions.Default(c)
		userID := session.Get("user-id")
		user, err := h.userService.ByEmail(userID.(string))

		if err != nil {
			log.Warnf("loading user by email: %s failed", userID.(string))
		}

		if user.VolleynetLogin != su.Username ||
			user.VolleynetUserID != loginData.ID {

			err = h.userService.SetVolleynetLogin(su.Username, loginData.ID)

			if err != nil {
				log.Warnf("updating volleynet user information failed for userID: %d", user.ID)
			}
		}
	}

	err = vnClient.TournamentEntry(su.PartnerName, su.PartnerID, su.TournamentID)

	if err != nil {
		log.Printf("entry to tournamentID %v with partnerID %v did not work: %s", su.TournamentID, su.PartnerID, err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *volleynetHandler) searchPlayers(c *gin.Context) {
	vnClient := client.Default()
	firstName := c.Query("fname")
	lastName := c.Query("lname")
	birthday := c.Query("bday")
	players, err := vnClient.SearchPlayers(firstName, lastName, birthday)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jsonn(c, http.StatusOK, players, "")
}
