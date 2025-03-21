import socket
import logging


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.was_killed = False
        signal.signal(signal.SIGTERM, self.graceful_shutdown)
        signal.signal(signal.SIGINT, self.graceful_shutdown)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while not self.was_killed:
            client_sock = self.__accept_new_connection()
            if (client_sock != None):
                self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            if(self.was_killed):
                logging.info(f"action: client_socket_closed_by_acceptor_socket | result : success")
            else:
                logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            logging.info("action: close_client_socket | result: in_progress")
            client_sock.close()
            logging.info("action: close_client_socket | result: success")

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        except OSError as e:
            if(self.was_killed):
                logging.info(f"action: client_socket_closed_by_acceptor_socket | result : success")
            else:
                self._server_socket.close()
        return c
    
    def stop(self):
        logging.info("action: close_socket | result: in_progress")
        self.was_killed = True
        self._server_socket.close()
        logging.info("action: close_socket | result: success")

    def graceful_shutdown(signum, frame):
        logging.info("action: shutdown | result: in_progress | signal: %s", signum)
        server.stop()
        logging.info("action: shutdown | result: success")
        sys.exit(0)
