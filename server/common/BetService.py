import logging
from typing import Optional

from common.bet_client import BetClient

from common.utils import Bet
from common.utils import has_won
from common.utils import store_bets, log_bets_stored, load_bets


class BetService:
    def __init__(self):
        self.agencies_ready: set[str] = set()
        self.winners_by_agency: dict[str, set[Bet]] = {}
        self.lottery_ended: bool = False

    def save_bets(self, bets: list[Bet]):
        store_bets(bets)
        log_bets_stored(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")

    def get_winners(self, agency_id: str) -> Optional[set[Bet]]:
        if self.lottery_ended:
            return self.winners_by_agency.get(agency_id, set())

        self.agencies_ready.add(agency_id)
        if len(self.agencies_ready) == 5:
            self.__end_lottery()

        if self.lottery_ended:
            return self.winners_by_agency.get(agency_id, set())

        return None


    def __end_lottery(self):
        self.lottery_ended = True

        for bet in load_bets():
            if has_won(bet):
                agency = bet.agency
                if agency not in self.winners_by_agency:
                    self.winners_by_agency[agency] = set()
                self.winners_by_agency[agency].add(bet)

        logging.info("action: sorteo | result: success")
