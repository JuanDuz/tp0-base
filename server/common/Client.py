import logging

from common.BetController import BetController
from common.BetService import BetService
from common.MessageHandler import MessageHandler
from common.NetworkClient import NetworkClient


class Client:
    def __init__(
            self,
            socket,
            bet_service,
    ):
        self.was_killed = False
        self.network_client = NetworkClient(socket)
        self.message_handler = MessageHandler(BetController(bet_service))

    def run(self):
        try:
            while not self.was_killed:
                raw_msg = self.network_client.receive_message()
                self.message_handler.handle(raw_msg, self.network_client)

        except OSError as e:
            logging.error("")
            # logging.error("action: receive_message | result: fail | error: %s", e)

        finally:
            try:
                self.stop()
            except Exception as e:
                logging.error("action: close_client_socket | result: fail | error: %s", e)

    def stop(self):
        self.was_killed = True
        self.network_client.close()
        logging.error("action: closing_client_socket | result: success")
