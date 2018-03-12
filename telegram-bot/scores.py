import requests
import json
from auth import scores_user, scores_pw
import time
import logging

class ScoresAPI:
    scores_url = 'https://scores.raphi011.com/api'
    cookies = None

    filter = {
        'day': 'today',
        'month': 'month',
        'year': 'thisyear',
        None: 'all'
    }

    players = {
        505934690: 5,  # Dom
        26409079: 1,   # Raffi
        147596787: 4,  # Richie
        404242664: 9,  # Consti
        287950001: 3,  # Luki
        187903526: 7   # Gerli
    }

    def __init__(self):
        self.getScoresSession()

    def getRanks(self, filter=None):

        self.isSessionExpired()
        scores = self.loadStats(filter)

        ranks = ''
        if scores:
            counter = 1
            for player in scores['data']:
                ranks += '{}: {} ({} %)\n'.format(player['rank'], player['player']['name'], player['percentageWon'])
                if counter == 3: break
                counter += 1

        if ranks == '':
            ranks = 'Sorry my dear, I really tried, but cannot provide your requested data :('

        return ranks

    def getGoal(self, msg):
        if msg['message']['from']['id'] not in self.players:
            return 'Sorry, I don\'t know you yet'

        senderid = self.players[msg['message']['from']['id']]
        senderid = 1

        self.isSessionExpired()
        scores = self.loadStats()

        if not scores:
            return 'Sorry my dear, I really tried, but cannot provide your requested data :('

        senderindex = 0
        while scores['data'][senderindex]['playerId'] != senderid:
            senderindex += 1

        if senderindex == 0:
            return 'You are already the hero, {}! No goals for you...'.format(
                scores['data'][senderindex]['player']['name']
            )

        else:
            senderperc = [scores['data'][senderindex]['played'], scores['data'][senderindex]['gamesWon']]
            goalperc = scores['data'][senderindex - 1]['percentageWon']

            wins = 1

            while (senderperc[1] + wins) / senderperc[0] * 100 < goalperc:
                wins += 1

            return '{}, to reach rank {} you need to win {} match{}...'.format(
                scores['data'][senderindex]['player']['name'],
                senderindex,
                wins,
                'es' if wins > 1 else ''
            )

    def loadStats(self, filter=None):
        try:
            scores = requests.get('{}/statistics'.format(self.scores_url), cookies=self.cookies,
                                  params={'filter': self.filter[filter]}).json()
            logging.debug(json.dumps(scores, sort_keys=True, indent=4, separators=(',', ': ')))

            return scores

        except Exception as ex:
            logging.error(ex)

        return None

    def isSessionExpired(self):
        if self.cookies:
            for c in self.cookies:
                if c.expires < time.time():
                    self.getScoresSession()
        else:
            self.getScoresSession()

    def getScoresSession(self):
        try:
            logging.warning('requesting new cookie')
            self.cookies = requests.post('{}/pwAuth'.format(self.scores_url),
                                         json={'email': scores_user, 'password': scores_pw},
                                         allow_redirects=False).cookies
        except Exception as ex:
            logging.error(ex)
