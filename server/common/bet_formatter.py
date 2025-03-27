from common.utils import Bet


def format_bet_message(bet):
    return f"{bet.first_name}|{bet.last_name}|{bet.document_number}|{bet.dob}|{bet.number}|{bet.agency}"


def parse_bet_message(message):
    fields = message.split("|")
    if len(fields) != 6:
        print(message)
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