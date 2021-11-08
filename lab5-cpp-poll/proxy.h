#ifndef SOCKS5_PROXY_H
#define SOCKS5_PROXY_H

#include <arpa/inet.h>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <iostream>
#include <unistd.h>
#include <netdb.h>
#include <poll.h>
#include <fcntl.h>
#include <csignal>
#include <sys/signalfd.h>
#include "exceptions.h"

struct ResolverStructure {
    uint16_t port;
    gaicb *host;
    pollfd *waited;
};

const static int MAX_CLIENTS = 2048;
const static int POLL_DELAY = 3000;
const static int BUFFER_LENGTH = 5000;
const static int HANDSHAKE_LENGTH = 2;
const static char SOCKS_VERSION = 5;
const static char SOCKS_SERVER_ERROR = 1;
const static char INVALID_AUTHORISATION = 0xff;
const static char SUPPORTED_AUTHORISATION = 0;
const static char SUPPORTED_OPTION = 1;
constexpr static char IPv4_ADDRESS = 1;
constexpr static char IPv6_ADDRESS = 4;
constexpr static char DOMAIN_NAME = 3;
const static char OPTION_NOT_SUPPORTED = 7;
const static char PROTOCOL_ERROR = 7;
const static char SOCKS_SUCCESS = 0;
const static char ADDRESS_NOT_SUPPORTED = 8;
const static char HOST_NOT_REACHABLE = 4;
const static int MINIMUM_SOCKS_REQUEST_LENGTH = 10;
const static int IPv4_ADDRESS_LENGTH = 4;
const static int SOCKS5_OFFSET_BEFORE_ADDR = 4;

class Proxy {
private:
    int serverSocket;
    int dnsSignal;

    std::unordered_set<char> supportedAddressTypes;
    char buffer[BUFFER_LENGTH];
    sockaddr_in serverAddr;
    sockaddr_in addr;
    std::vector<pollfd> *pollDescryptors;
    std::unordered_set<pollfd *> passedHandshake;
    std::unordered_set<pollfd *> passedFullSOCKSprotocol;
    std::unordered_map<pollfd *, pollfd *> *transferMap;
    std::unordered_map<pollfd *, std::vector<char> > *dataPieces;
    int waitedCounter = 0;

public:
    Proxy();

    ~Proxy();

    explicit Proxy(int port);

    void run();

private:
    void pollManage();

    void removeDeadDescriptors();

    void init(int port);

    void setupDnsSignal();

    void connectToIPv4Address(std::vector<pollfd>::iterator *clientIterator);

    void skipIPV6(std::vector<pollfd>::iterator *clientIterator);

    void resolveDomainName(std::vector<pollfd>::iterator *clientIterator);

    void getResolveResult(std::vector<pollfd>::iterator *clientIterator);

    void registerForWrite(pollfd *fd) { fd->events |= POLLOUT; }

    void passIdentification(std::vector<pollfd>::iterator *clientIterator);

    void sendData(std::vector<pollfd>::iterator *clientIterator);

    void acceptConnection(pollfd *client);

    void makeHandshake(std::vector<pollfd>::iterator *clientIterator);

    void readData(std::vector<pollfd>::iterator *clientIterator);

    bool checkSOCKSRequest(int client, ssize_t len);

    void printSocksBuffer(int size);

    static uint32_t constructIPv4Addr(char *addr);

    static in6_addr constructIPv6Addr(char *addr);

    static uint16_t constructPort(char *port);
};

#endif
