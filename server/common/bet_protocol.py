def send_string(sock, message):
    length = str(len(message)) + "\n"
    sock.sendall(length.encode('utf-8') + message.encode('utf-8'))


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
