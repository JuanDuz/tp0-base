import os
import socket
import logging
import signal
import sys
from multiprocessing import Process, Manager

from common.utils import store_bets, log_bets_stored
from common.utils import Bet
from common.BetController import BetController
from common.BetService import BetService
from common.MessageHandler import MessageHandler
from common.Client import Client
from common.BetsFileMonitor import BetsFileMonitor


class Server:
    def __init__(self, port, listen_backlog):
        signal.signal(signal.SIGTERM, self.graceful_shutdown)
        signal.signal(signal.SIGINT, self.graceful_shutdown)

        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.was_killed = False
        self._clients = []
        self._processes = []
        self.manager = Manager()
        self.agencies_ready = self.manager.list()
        self.winners = self.manager.list()
        self.lottery_ended =  self.manager.Value('b', False)


    def run(self):

        bets_file_monitor = BetsFileMonitor(self.manager.Lock())
        shared_bet_service = BetService(
            total_agencies=int(os.environ.get("TOTAL_AGENCIES", "1")),
            agencies_ready=self.agencies_ready,
            winners=self.winners,
            lottery_ended=self.lottery_ended,
            bets_file_monitor=bets_file_monitor,
            lock=self.manager.Lock()
        )

        while not self.was_killed:
            client_sock = self.__accept_new_connection()
            if client_sock is not None:
                client = Client(client_sock, shared_bet_service)
                self._clients.append(client)
                p = Process(target=client.run)
                p.daemon = True
                self._processes.append(p)
                p.start()

            self._clean_up_dead_processes()

            if not self._processes:
                logging.info("action: no_more_clients | result: success")
                self.stop()


    def __accept_new_connection(self):
        logging.info('action: accept_connections | result: in_progress')
        c = None
        try:
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        except OSError as e:
            logging.error("action: accept_connections | result: fail | error: %s", e)
        return c


    def stop(self):
        if self.was_killed:
            return

        logging.info("action: close_socket | result: in_progress")
        self.was_killed = True

        try:
            self._server_socket.close()

            for client in self._clients:
                try:
                    client.stop()
                except Exception as e:
                    logging.warning("action: client_stop | result: already_closed | error: %s", e)

            for p in self._processes:
                try:
                    if p.is_alive():
                        p.terminate()
                    p.join()
                    logging.info("action: process_joined | result: success")
                except Exception as e:
                    logging.warning("action: process_join | result: fail | error: %s", e)


        except OSError as e:
            logging.error("action: close_socket | result: fail | error: %s", e)

        logging.info("action: close_socket | result: success")


    def graceful_shutdown(self, signum, frame):
        logging.info("action: shutdown | result: in_progress | signal: %s", signum)
        self.stop()
        logging.info("action: shutdown | result: success")
        sys.exit(0)

    def _clean_up_dead_processes(self):
        alive = []
        for p in self._processes:
            if not p.is_alive():
                p.join()
                # logging.info("action: process_cleaned | result: success")
            else:
                alive.append(p)
        self._processes = alive
