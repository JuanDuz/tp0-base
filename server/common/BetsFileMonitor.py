import logging

from common import Bet

from common.utils import store_bets, load_bets

class BetsFileMonitor:
    def __init__(self, lock):
        self._lock = lock

    def safe_store_bets(self, bets: list[Bet]):
        with self._lock:
            try:
                store_bets(bets)
            except Exception as e:
                logging.error(f"action: safe_store_bets | result: fail | error: {e} | bets: {bets}")

    def safe_load_bets(self) -> list[Bet]:
        with self._lock:
            try:
                return list(load_bets())
            except Exception as e:
                logging.error(f"action: safe_load_bets | result: fail | error: {e}")
                return []
