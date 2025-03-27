from multiprocessing import Lock

from common import Bet

from common.utils import store_bets, load_bets

class BetsFileMonitor:
    def __init__(self, lock):
        self._lock = lock

    def safe_store_bets(self, bets: list[Bet]):
        with self._lock:
            store_bets(bets)

    def safe_load_bets(self) -> list[Bet]:
        with self._lock:
            return list(load_bets())
