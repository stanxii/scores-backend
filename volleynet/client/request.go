package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/raphi011/scores/volleynet"
	"github.com/raphi011/scores/volleynet/scrape"
)

// Client is the interface to the volleynet api, use DefaultClient()
// to get a new Client.
type Client interface {
	Login(username, password string) (*scrape.LoginData, error)

	Tournaments(gender, league string, year int) ([]*volleynet.TournamentInfo, error)
	Ladder(gender string) ([]*volleynet.Player, error)
	ComplementTournament(tournament *volleynet.TournamentInfo) (*volleynet.Tournament, error)

	WithdrawFromTournament(tournamentID int) error
	EnterTournament(playerName string, playerID, tournamentID int) error

	SearchPlayers(firstName, lastName, birthday string) ([]*scrape.PlayerInfo, error)
}

// Default implements the Client interface
type Default struct {
	PostURL string
	GetURL  string
	Cookie  string
	Log     logrus.FieldLogger
}

// DefaultClient returns a Client with the correct PostURL and GetURL fields set.
func DefaultClient() Client {
	return &Default{
		PostURL: "https://beach.volleynet.at",
		GetURL:  "http://www.volleynet.at",
		Log:     logrus.New(),
	}
}

// WithLogger returns a DefaultClient with a logger.
func WithLogger(log logrus.FieldLogger) Client {
	return &Default{
		PostURL: "https://beach.volleynet.at",
		GetURL:  "http://www.volleynet.at",
		Log:     log,
	}
}

// Login authenticates the user against the volleynet page, if
// successfull the Client cookie is set, else an error is returned.
func (c *Default) Login(username, password string) (*scrape.LoginData, error) {
	form := url.Values{}
	form.Add("login_name", username)
	form.Add("login_pass", password)
	form.Add("action", "Beach/Profile/ProfileLogin")
	form.Add("submit", "OK")
	form.Add("mode", "X")

	url := c.buildPostURL("/Admin/formular").String()
	resp, err := http.PostForm(url, form)

	if err != nil {
		return nil, errors.Wrap(err, "client login")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	loginData, err := scrape.Login(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "parse login")
	}

	c.Cookie = resp.Header.Get("Set-Cookie")

	semicolonIndex := strings.Index(c.Cookie, ";")

	if semicolonIndex > 0 {
		c.Cookie = c.Cookie[:semicolonIndex]
	}

	return loginData, nil
}

// Tournaments reads all tournaments of a certain gender, league and year.
// To get all details of a tournamnent use `Client.ComplementTournament`.
func (c *Default) Tournaments(gender, league string, year int) ([]*volleynet.TournamentInfo, error) {
	url := c.buildGetAPIURL(
		"/beach/bewerbe/%s/phase/%s/sex/%s/saison/%d/information/all",
		league,
		league,
		gender,
		year,
	)

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return scrape.TournamentList(resp.Body, c.GetURL)
}

// Ladder loads all ranked players of a certain gender.
func (c *Default) Ladder(gender string) ([]*volleynet.Player, error) {
	url := c.buildGetAPIURL(
		"/beach/bewerbe/Rangliste/phase/%s",
		genderLong(gender),
	).String()

	resp, err := http.Get(url)

	if err != nil {
		return nil, errors.Wrapf(err, "loading ladder %q failed", gender)
	}

	defer resp.Body.Close()

	return scrape.Ladder(resp.Body)
}

func genderLong(gender string) string {
	if gender == "M" {
		return "Herren"
	} else if gender == "W" {
		return "Damen"
	}

	return ""
}

// ComplementTournament adds the missing information from `Tournaments`.
func (c *Default) ComplementTournament(tournament *volleynet.TournamentInfo) (
	*volleynet.Tournament, error) {
	url := c.getAPITournamentLink(tournament)

	c.Log.Debugf("Downloading tournament: %s\n", url)
	resp, err := http.Get(url)

	if err != nil {
		return nil, errors.Wrapf(err, "loading tournament %d failed", tournament.ID)
	}

	t, err := scrape.Tournament(resp.Body, time.Now(), tournament)

	if err != nil {
		return nil, errors.Wrapf(err, "parsing tournament %d failed", tournament.ID)
	}

	t.Link = c.getTournamentLink(tournament)

	return t, nil
}

func (c *Default) loadUniqueWriteCode(tournamentID int) (string, error) {
	url := c.buildPostURL(
		"/Admin/index.php?screen=Beach/Profile/TurnierAnmeldung&parent=0&prev=0&next=0&cur=%d",
		tournamentID,
	).String()

	req, err := http.NewRequest(
		"GET",
		url,
		nil)

	if err != nil {
		return "", errors.Wrap(err, "creating request failed")
	}

	req.Header.Add("Cookie", c.Cookie)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return "", errors.Wrap(err, "loading unique writecode failed")
	}

	code, err := scrape.UniqueWriteCode(resp.Body)

	return code, errors.Wrap(err, "parsing unique writecode failed")
}

// WithdrawFromTournament withdraws a player from a tournament.
// A valid session Cookie must be set.
func (c *Default) WithdrawFromTournament(tournamentID int) error {
	url := c.buildPostURL("/Abmelden/0-%d-00-0", tournamentID).String()

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", c.Cookie)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return errors.Wrapf(err, "tournamentwithdrawal request for tournamentID: %d failed", tournamentID)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tournamentwithdrawal request for tournamentID: %d failed with code %d",
			tournamentID,
			resp.StatusCode)
	}

	return nil
}

// EnterTournament enters a player at a tournament.
// A valid session Cookie must be set.
func (c *Default) EnterTournament(playerName string, playerID, tournamentID int) error {
	if c.Cookie == "" {
		return errors.New("cookie must be set")
	}

	form := url.Values{}

	code, err := c.loadUniqueWriteCode(tournamentID)

	if err != nil {
		return errors.Wrapf(err, "could not load writecode for tournamentID: %d", tournamentID)
	}

	form.Add("action", "Beach/Profile/TurnierAnmeldung")
	form.Add("XX_unique_write_XXBeach/Profile/TurnierAnmeldung", code)
	form.Add("parent", "0")
	form.Add("prev", "0")
	form.Add("next", "0")
	form.Add("cur", strconv.Itoa(tournamentID))
	form.Add("name_b", playerName)
	form.Add("bte_per_id_b", strconv.Itoa(playerID))
	form.Add("submit", "Anmelden")

	url := c.buildPostURL("/Admin/formular").String()

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(form.Encode()))

	if err != nil {
		return errors.Wrap(err, "creating tournamententry request failed")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", c.Cookie)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return errors.Wrapf(err, "tournamententry request for tournamentID: %d failed", tournamentID)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tournamententry request for tournamentID: %d failed with code %d",
			tournamentID,
			resp.StatusCode)
	}

	entryData, err := scrape.Entry(resp.Body)

	if err != nil || !entryData.Successfull {
		return errors.Wrapf(err, "tournamententry request for tournamentID: %d failed", tournamentID)
	}

	return nil
}

func traceResponse(resp *http.Response) {
	outFile, err := os.Create("/home/raphi/login-response.html")
	if err != nil {
		return
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return
	}
}

// SearchPlayers searches for players via firstName, lastName and their birthdate in dd.mm.yyyy format.
func (c *Default) SearchPlayers(firstName, lastName, birthday string) ([]*scrape.PlayerInfo, error) {
	form := url.Values{}

	form.Add("XX_unique_write_XXAdmin/Search", "0.50981600 1525795371")
	form.Add("popup", "1")
	form.Add("add", "")
	form.Add("target", "bte_per_id_b")
	form.Add("txm_language", "de")
	form.Add("sai_id", "")
	form.Add("action", "Admin/Search")
	form.Add("submit", "Suchen")
	form.Add("search", "Person")
	form.Add("per_name", lastName)
	form.Add("per_vorname", firstName)
	form.Add("per_geburtsdatum", birthday)
	form.Add("doit", "1")
	form.Add("text", "0")

	url := c.buildPostURL("/Admin/formular")

	response, err := http.PostForm(url.String(), form)

	if err != nil {
		return nil, err
	}

	return scrape.Players(response.Body)
}
