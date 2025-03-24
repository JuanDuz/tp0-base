from common.bet_formatter import parse_bet_message
from common.bet_protocol import receive_string, send_string

from common.utils import Bet


class BetClient:
    def __init__(self, socket):
        self.socket = socket

    def receive_bet(self) -> Bet:
        raw_msg = receive_string(self.socket)
        bet: Bet = parse_bet_message(raw_msg)
        self.send_ack()
        return bet

    def send_ack(self):
        send_string(self.socket, "ACK")
