import socket
import logging
import signal
import sys

from common.bet_client import BetClient
from common.utils import store_bets, log_bets_stored

from common.utils import Bet

from common.BetController import BetController
from common.BetService import BetService
from common.MessageHandler import MessageHandler


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.was_killed = False
        self.agencies_ready = set()
        self.message_handler = MessageHandler(BetController(BetService()))
        signal.signal(signal.SIGTERM, self.graceful_shutdown)
        signal.signal(signal.SIGINT, self.graceful_shutdown)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # the server
        while not self.was_killed:
            client_sock = self.__accept_new_connection()
            if client_sock is not None:
                self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        try:
            bet_client = BetClient(client_sock)
            raw_msg = bet_client.receive_message()
            self.message_handler.handle(raw_msg, bet_client)

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: %s", e)

        finally:
            logging.info("action: close_client_socket | result: in_progress")
            client_sock.close()
            logging.info("action: close_client_socket | result: success")


    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        except OSError as e:
            if self.was_killed:
                logging.info(f"action: client_socket_closed_by_acceptor_socket | result : success")
            else:
                self._server_socket.close()
        return c
    
    def stop(self):
        logging.info("action: close_socket | result: in_progress")
        self.was_killed = True
        self._server_socket.close()
        logging.info("action: close_socket | result: success")

    def graceful_shutdown(self, signum, frame):
        logging.info("action: shutdown | result: in_progress | signal: %s", signum)
        self.stop()
        logging.info("action: shutdown | result: success")
        sys.exit(0)
