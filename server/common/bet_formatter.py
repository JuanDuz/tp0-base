from common.utils import Bet


def format_bet_message(bet):
    return f"{bet.first_name}|{bet.last_name}|{bet.document_number}|{bet.dob}|{bet.number}"


def parse_bet_message(message):
    fields = message.split("|")
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