from threading import Timer
from datetime import datetime
import logging
from auth import jobtimeout

class Jobs:
    jobs = []

    def __init__(self, scores, telegram):
        self.scores = scores
        self.telegram = telegram
        self.startTimer()

    def appendJob(self, job):
        self.jobs.append(job)
        logging.warning('new job created')

    def handleJobs(self):
        for j in self.jobs:
            if 'until' in j and datetime.now() < j['until']:
                d = j['delegate']
                delete = d(j)
            else:
                delete = True

            if delete:
                self.jobs.remove(j)   # need testing
                logging.warning('job finished')

        self.startTimer()

    def startTimer(self):
        t = Timer(jobtimeout, self.handleJobs)
        t.start()