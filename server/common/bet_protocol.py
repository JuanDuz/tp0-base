import logging


def send_string(sock, message):
    length = str(len(message.encode('utf-8'))) + "\n"
    encoded = length.encode('utf-8') + message.encode('utf-8')
    try:
        sock.sendall(encoded)
    except Exception as e:
        logging.error("error sending message | error: %s | message: %s", e, repr(message))


def receive_string(sock):
    buffer = b""
    while not buffer.endswith(b"\n"):
        data = sock.recv(1)
        if not data:
            raise ConnectionError("Connection closed while reading length")
        buffer += data
    length = int(buffer.strip())
    message = b""
    while len(message) < length:
        data = sock.recv(length - len(message))
        if not data:
            raise ConnectionError("Connection closed while reading message")
        message += data
    return message.decode('utf-8')
