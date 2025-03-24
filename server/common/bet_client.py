from common.bet_formatter import parse_bet_message
from common.bet_protocol import receive_string, send_string

from common.utils import Bet


class BetClient:
    def __init__(self, socket):
        self.socket = socket

    def receive_bets(self) -> list[Bet]:
        raw_msg = receive_string(self.socket)
        lines = raw_msg.strip().split('\n')
        bets = []
        for line in lines:
            bet = parse_bet_message(line)
            bets.append(bet)
        self.send_ack()
        return bets

    def send_ack(self):
        send_string(self.socket, "ACK")
