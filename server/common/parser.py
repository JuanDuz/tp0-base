from typing import Optional

from common.utils import Bet


def parse_bet_to_str(bet):
    return f"{bet.first_name}|{bet.last_name}|{bet.document}|{bet.birthdate}|{bet.number}|{bet.agency}"

def parse_bets_to_str(bets: set[Bet]) -> str:
    return "\n".join(parse_bet_to_str(bet) for bet in bets)

def _parse_str_to_bet(raw_bet):
    fields = raw_bet.split("|")
    if len(fields) != 6:
        raise ValueError("Invalid bet message format")
    return Bet(
        agency=fields[5],
        first_name=fields[0],
        last_name=fields[1],
        document=fields[2],
        birthdate=fields[3],
        number=fields[4],
    )

def parse_str_to_bets(raw_msg: str) -> list[Bet]:
    lines = raw_msg.strip().split('\n')
    bets = []
    for line in lines:
        bet = _parse_str_to_bet(line)
        bets.append(bet)
    return bets

def parse_agency_id_from_get_winners(raw_msg: str) -> Optional[int]:
    try:
        parts = raw_msg.strip().split("|")
        if len(parts) != 2 or parts[0] != "GET_WINNERS":
            return None
        return int(parts[1])
    except Exception:
        return None