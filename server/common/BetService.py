import logging
import os
from typing import Optional

from common.bet_client import BetClient

from common.utils import Bet
from common.utils import has_won
from common.utils import store_bets, log_bets_stored, load_bets


class BetService:
    def __init__(self):
        self.agencies_ready: set[int] = set()
        self.winners_by_agency: dict[int, set[Bet]] = {}
        self.lottery_ended: bool = False
        self.expected_agencies = int(os.environ.get("TOTAL_AGENCIES", "5"))

    def save_bets(self, bets: list[Bet]):
        store_bets(bets)
        log_bets_stored(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")

    def get_winners(self, agency_id: int) -> Optional[set[Bet]]:
        if self.lottery_ended:
            return self.winners_by_agency.get(agency_id, set())

        self.agencies_ready.add(agency_id)
        if len(self.agencies_ready) == self.expected_agencies:
            self.__end_lottery()

        if self.lottery_ended:
            winners_per_agency = self.winners_by_agency.get(agency_id, set())
            return winners_per_agency

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
