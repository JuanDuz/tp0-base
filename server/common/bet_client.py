from common.bet_formatter import parse_bet_message
from common.bet_protocol import receive_string, send_string

from common.utils import Bet


class BetClient:
    def __init__(self, socket):
        self.socket = socket

    def receive_message(self) -> str:
        return receive_string(self.socket)

    def send_message(self, message: str):
        send_string(self.socket, message)

    def send_ack(self):
        self.send_message("ACK")

    def send_error(self, msg):
        self.send_message(msg)

