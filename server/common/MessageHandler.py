
from common import BetController

from common.NetworkClient import NetworkClient


class MessageHandler:
    def __init__(self, bet_controller: BetController):
        self.bet_controller: BetController = bet_controller

    def handle(self, raw_msg: str, bet_client: NetworkClient):
        if raw_msg.startswith("GET_WINNERS"):
            return self.bet_controller.get_winners(raw_msg, bet_client)
        else:
            return self.bet_controller.save_bets(raw_msg, bet_client)
