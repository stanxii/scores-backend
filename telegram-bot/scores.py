import requests
import json
from auth import scores_user, scores_pw
import time
import logging
import time
from datetime import datetime
import _thread
from dateutil.parser import parse
from jobs import Jobs

class Player:
    def __init__(self, name, telegramid, scoresid):
        self.name = name
        self.telegramid = telegramid
        self.scoresid = scoresid

class ScoresAPI:
    scores_url = 'https://scores.raphi011.com/api'
    cookies = None

    filter = {
        'day': 'today',
        'month': 'month',
        'year': 'thisyear',
        None: 'all'
    }

    def __init__(self, telegram):
        self.getScoresSession()
        self.telegram = telegram

        self.players = []
        self.players.append(Player('Raffi', 26409079, 1 ))
        self.players.append(Player('Dom', 505934690, 5 ))
        self.players.append(Player('Richie', 147596787, 4 ))
        self.players.append(Player('Consti', 404242664, 9 ))
        self.players.append(Player('Luki', 287950001, 3 ))
        self.players.append(Player('Gerli', 187903526, 7 ))

        self.jobs = Jobs(self, telegram)

    def handleTelegramMessage(self, msg):
        self.isSessionExpired()

        cmd = msg['message']['command']

        player = self.getPlayer(telegramid=msg['message']['from']['id'])
        chatid = msg['message']['chat']['id']

        if cmd[0] == 'help':
            self.telegram.sendMessage(self.getHelp(), chatid)

        elif cmd[0] == 'top3':
            filter = cmd[1] if len(cmd) > 1 and cmd[1] in ['day', 'month', 'year'] else None
            self.telegram.sendMessage(self.getRanks(filter), chatid)

        elif cmd[0] == 'setgoal' and player:
            self.telegram.sendMessage(self.setGoal(player, cmd, chatid), chatid)

        elif not player:
            self.telegram.sendMessage('Sorry, I don\'t know you yet', chatid)

    def getHelp(self):
        return '''Hey there! I accept the following commands:
/top3 - shows the current over all top3
/setgoal wins <count> - you can set a goal and I will write you when you reached it
'''


    def getRanks(self, filter=None):
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

    def watchGoal(self, job):
        matches = self.loadMatches()
        #stats = self.loadStats()

        if job['type'] == 'wins':
            wins = 0
            for m in matches['data']:
                if parse(m['createdAt']).replace(tzinfo=None) > datetime.now().replace(hour=0, minute=0):
                    if int(m['scoreTeam1']) > int(m['scoreTeam2']) and \
                            (self.getPlayer(scoresid=m['team1']['player1Id']) == job['player'] or \
                             self.getPlayer(scoresid=m['team1']['player2Id']) == job['player']):
                        wins += 1
                    elif int(m['scoreTeam1']) < int(m['scoreTeam2']) and \
                            (self.getPlayer(scoresid=m['team2']['player1Id']) == job['player'] or \
                             self.getPlayer(scoresid=m['team2']['player2Id']) == job['player']):
                        wins += 1

            if wins >= job['value']:
                self.telegram.sendMessage('Congratulations {}, you have reached your goal!'.format(job['player'].name), job['chatid'])
                return True

            return False

        return True


    def setGoal(self, player, cmd, chatid):

        if cmd[1] in ['wins', 'rank', 'percent'] and cmd[2].isdigit():

            goal = {
                'delegate': self.watchGoal,
                'type': cmd[1],
                'value': int(cmd[2]),
                'player': player,
                'until': datetime.now().replace(hour=22, minute=0),
                'chatid': chatid
            }

            self.jobs.appendJob(goal)

            return 'Thanks {}, your goal is saved! You\'ll hear from me ;)'.format(player.name)

        else:
            return 'Sorry {}, I can\'t set that goal'.format(player.name)

    def getPlayer(self, scoresid=None, telegramid=None):
        for p in self.players:
            if scoresid and p.scoresid == scoresid:
                return p
            if telegramid and p.telegramid == telegramid:
                return p

        return None

    def loadStats(self, filter=None):
        try:
            scores = requests.get('{}/statistics'.format(self.scores_url), cookies=self.cookies,
                                  params={'filter': self.filter[filter]}).json()
            logging.debug(json.dumps(scores, sort_keys=True, indent=4, separators=(',', ': ')))

            return scores

        except Exception as ex:
            logging.error(ex)

        return None

    def loadMatches(self):
        try:
            matches = requests.get('{}/matches'.format(self.scores_url), cookies=self.cookies).json()
            logging.debug(json.dumps(matches, sort_keys=True, indent=4, separators=(',', ': ')))

            return matches
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