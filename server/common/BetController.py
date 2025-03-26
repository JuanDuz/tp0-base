import logging
from common.parser import parse_str_to_bets, parse_agency_id_from_get_winners
from common.BetService import BetService
from common.bet_client import BetClient

from common.parser import parse_bets_to_str


class BetController:
    def __init__(self, service: BetService):
        self.service = service

    def save_bets(self, raw_msg: str, client: BetClient):
        bets = parse_str_to_bets(raw_msg)

        if bets is None:
            logging.info("action: apuesta_recibida | result: fail | cantidad: 0")
            client.send_error("ERROR_INVALID_BATCH")
            return

        if len(bets) == 0:
            logging.info("action: apuesta_recibida | result: fail | cantidad: 0")
            client.send_error("ERROR_EMPTY_BATCH")
            return

        self.service.save_bets(bets)
        client.send_ack()

    def get_winners(self, raw_msg: str, client: BetClient):
        agency_id = parse_agency_id_from_get_winners(raw_msg)

        if agency_id is None:
            logging.info("action: consulta_ganadores | result: fail | reason: ERROR_INVALID_GET_WINNERS")
            client.send_error("ERROR_INVALID_GET_WINNERS")
            return

        winner_bets = self.service.get_winners(agency_id)
        if winner_bets is None:
            logging.info("action: consulta_ganadores | result: fail | reason: ERROR_LOTTERY_HASNT_ENDED")
            client.send_error("ERROR_LOTTERY_HASNT_ENDED")
        else:
            client.send_message(parse_bets_to_str(winner_bets))
