import requests
import json
from auth import scores_user, scores_pw
import time
import logging


class ScoresAPI:
    scores_url = 'https://scores.raphi011.com/api'

    ranks = {
        1: 'Master',
        2: 'Padawan',
        3: 'Puppet'
    }

    filter = {
        'day': 'today',
        'month': 'month',
        'year': 'thisyear',
        None: 'all'
    }

    def __init__(self):
        self.getScoresSession()

    def getRank(self, filter=None):
        rank = ''

        self.isSessionExpired()

        try:
            scores = requests.get('{}/statistics'.format(self.scores_url), cookies=self.cookies,
                                  params={'filter': self.filter[filter]}).json()
            logging.debug(json.dumps(scores, sort_keys=True, indent=4, separators=(',', ': ')))

            counter = 1
            for player in scores['data']:
                rank += '{}: {} ({} %)\n'.format(self.ranks[counter], player['player']['name'], player['percentageWon'])
                if counter == 3: break
                counter += 1

        except Exception as ex:
            logging.error(ex)

        if rank == '':
            rank = 'Sorry my dear, I really tried, but cannot provide your requested data'

        return rank

    def isSessionExpired(self):
        if self.cookies:
            for c in self.cookies:
                if c.expires < time.time():
                    self.getScoresSession()

    def getScoresSession(self):
        try:
            logging.info('requesting new cookie')
            self.cookies = requests.post('{}/pwAuth'.format(self.scores_url),
                                         json={'email': scores_user, 'password': scores_pw},
                                         allow_redirects=False).cookies
        except Exception as ex:
            logging.error(ex)
