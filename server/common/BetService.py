import logging
from typing import Optional

from common.utils import Bet
from common.utils import has_won
from common.utils import store_bets, log_bets_stored, load_bets


class BetService:
    def __init__(
            self,
            total_agencies,
            agencies_ready,
            winners,
            lottery_ended,
            bets_file_monitor,
            lock
    ):
        self.agencies_ready = agencies_ready
        self.winners = winners
        self.lottery_ended = lottery_ended
        self.expected_agencies = total_agencies
        self.bets_file_monitor = bets_file_monitor
        self._lock = lock

    def save_bets(self, bets: list[Bet]):
        self.bets_file_monitor.safe_store_bets(bets)
        log_bets_stored(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")

    def get_winners(self, agency_id: int) -> Optional[list[Bet]]:
        with self._lock:
            if self.lottery_ended.value:
                return self._get_winners_by_agency_id(agency_id)

            if agency_id not in self.agencies_ready:
                self.agencies_ready.append(agency_id)

            if len(self.agencies_ready) == self.expected_agencies:
                self.__draw_lottery()

            if self.lottery_ended.value:
                return self._get_winners_by_agency_id(agency_id)
            return None

    def __draw_lottery(self):
        self.lottery_ended.value = True

        for bet in self.bets_file_monitor.safe_load_bets():
            if has_won(bet):
                self.winners.append(bet)

        logging.info("action: sorteo | result: success")

    def _get_winners_by_agency_id(self, agency_id):
        if not self.winners:
            return None
        return [winner for winner in self.winners if winner.agency == agency_id]
