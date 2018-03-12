from auth import telegram_token
import requests
import json
import time
from scores import ScoresAPI
import logging

'''
telegram commands:

top3 - Show current rank with filter
goal - Show players goal
'''

class TelegramBot:
    telegram_url = 'https://api.telegram.org/bot{}'.format(telegram_token)
    messagequeue = []

    def __init__(self, debug=False):
        self.scores = ScoresAPI()
        self.debug = debug

    def handleMessage(self, msg):
        cmd = msg['message']['command']

        if cmd[0] == 'top3':
            filter = cmd[1] if len(cmd) > 1 and cmd[1] in ['day', 'month', 'year'] else None
            self.sendMessage(self.scores.getRanks(filter), msg['message']['chat']['id'])

        elif cmd[0] == 'goal':
            self.sendMessage(self.scores.getGoal(msg), msg['message']['chat']['id'])

    def getMessages(self):
        offset = 0

        while True:
            try:
                ret = requests.get('{}/getUpdates'.format(self.telegram_url),
                                   params={'offset': offset}).json()

                for update in ret['result']:
                    offset = update['update_id'] + 1    # offset is used to tell telegram that we got that message.
                                                        # early increasing the offset does not block the bot when
                                                        # the message leads to an exception

                    logging.debug(json.dumps(update, sort_keys=True, indent=4, separators=(',', ': ')))
                    logging.warning('received update id {}'.format(update['update_id']))

                    if 'message' in update:
                        if self.setArgs(update):
                            self.handleMessage(update)
                        # think about async handling of messages, but not necessary for now

            except Exception as ex:
                logging.error(ex)

            time.sleep(1)

    def sendMessage(self, msg, chatid):
        try:
            requests.post('{}/sendMessage'.format(self.telegram_url), data={'chat_id': chatid, 'text': msg})
        except Exception as ex:
            logging.error(ex)

    def setArgs(self, msg):
        t = msg['message']['text'].split(' ')

        if t[0][0] != '/':
            return False

        h = t[0].split('@')

        if not (len(h) == 1 and msg['message']['chat']['type'] == 'private' or
                len(h) >  1 and h[1] in ['bvbscoresbot', 'testbvbscoresbot']):
            return False

        t[0] = h[0][1:]

        msg['message']['command'] = t

        return True