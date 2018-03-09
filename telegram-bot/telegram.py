from auth import telegram_token
import requests
import json
import time
from scores import ScoresAPI
import logging

class TelegramBot:
    telegram_url = 'https://api.telegram.org/bot{}'.format(telegram_token)
    messagequeue = []

    def __init__(self, debug=False):
        self.scores = ScoresAPI()
        self.debug = debug

    def handleMessage(self, msg):
        text = msg['message']['text'].split(' ')

        if text[0] == '/start' and msg['message']['chat']['type'] == 'private' \
                or text[0] == '/start@bvbscoresbot':

            self.sendMessage('Hello! I am your scores bot :)', msg['message']['chat']['id'])

        elif text[0] == '/rank' and msg['message']['chat']['type'] == 'private' \
                or text[0] == '/rank@bvbscoresbot':

            filter = text[1] if len(text) > 1 and text[1] in ['day', 'month', 'year'] else None

            self.sendMessage(self.scores.getRank(filter), msg['message']['chat']['id'])

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
                    logging.info('received update id {}'.format(update['update_id']))

                    if 'message' in update:
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
