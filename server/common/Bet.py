""" A lottery bet registry. """
import datetime


class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

    def __eq__(self, other):
        if not isinstance(other, Bet):
            return False
        return (
                self.agency == other.agency and
                self.first_name == other.first_name and
                self.last_name == other.last_name and
                self.document == other.document and
                self.birthdate == other.birthdate and
                self.number == other.number
        )

    def __hash__(self):
        return hash((self.agency, self.first_name, self.last_name, self.document, self.birthdate, self.number))

    def __repr__(self):
        return f"Bet({self.document}, {self.number}, agency={self.agency})"
