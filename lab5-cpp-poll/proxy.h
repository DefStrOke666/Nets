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

#define MAX_CONNECTIONS 2048
#define POLL_DELAY  3000
#define BUFFER_LENGTH 5000
#define HANDSHAKE_LENGTH 2
#define SOCKS_VERSION 5
#define SOCKS_SERVER_ERROR 1
#define INVALID_AUTHORISATION 0xff
#define SUPPORTED_AUTHORISATION 0
#define SUPPORTED_OPTION 1
#define IPv4_ADDRESS 1
#define IPv6_ADDRESS 4
#define DOMAIN_NAME 3
#define OPTION_NOT_SUPPORTED 7
#define PROTOCOL_ERROR 7
#define SOCKS_SUCCESS 0
#define ADDRESS_NOT_SUPPORTED 8
#define HOST_NOT_REACHABLE 4
#define MINIMUM_SOCKS_REQUEST_LENGTH 10
#define IPv4_ADDRESS_LENGTH 4
#define SOCKS5_OFFSET_BEFORE_ADDR 4

struct ResolverStructure {
    uint16_t port;
    gaicb *host;
    pollfd *waited;
};

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
    explicit Proxy(int port);

    void run();

    ~Proxy();

private:
    void pollManage();

    static void removeFromPoll(std::vector<pollfd>::iterator *it);

    void removeDeadDescriptors();

    void setupDNSSignal();

    void connectToIPv4Address(std::vector<pollfd>::iterator *clientIterator);

    void skipIPV6(std::vector<pollfd>::iterator *clientIterator);

    void resolveDomainName(std::vector<pollfd>::iterator *clientIterator);

    void getResolveResult(std::vector<pollfd>::iterator *clientIterator);

    void passIdentification(std::vector<pollfd>::iterator *clientIterator);

    void sendData(std::vector<pollfd>::iterator *clientIterator);

    void acceptConnection(pollfd *client);

    void handshake(std::vector<pollfd>::iterator *clientIterator);

    void readData(std::vector<pollfd>::iterator *clientIterator);

    bool checkSOCKSRequest(int client, ssize_t len);

    void printSocksBuffer(int size);

    static uint32_t constructIPv4Addr(char *addr);

    static uint16_t constructPort(char *port);
};

#endif
