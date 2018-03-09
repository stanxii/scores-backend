from telegram import TelegramBot
import logging
logging.basicConfig(level=logging.DEBUG)

if __name__ == '__main__':
    t = TelegramBot()
    t.getMessages()
