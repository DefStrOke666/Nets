#include "proxy.h"

Proxy::Proxy(int port) {
    signal(SIGPIPE, SIG_IGN);
    serverSocket = socket(AF_INET, SOCK_STREAM, 0);
    pollDescryptors = new std::vector<pollfd>;
    transferMap = new std::unordered_map<pollfd *, pollfd *>;
    dataPieces = new std::unordered_map<pollfd *, std::vector<char> >;

    int opt = 1;
    setsockopt(serverSocket, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));
    if (serverSocket == -1) {
        throw proxyException("Can't open server socket!");
    }

    serverAddr.sin_addr.s_addr = htonl(INADDR_ANY);
    serverAddr.sin_family = AF_INET;
    serverAddr.sin_port = htons(port);

    if (bind(serverSocket, (sockaddr *) &serverAddr, sizeof(serverAddr))) {
        throw proxyException("Can't bind server socket!");
    }

    if (listen(serverSocket, MAX_CONNECTIONS)) {
        throw proxyException("Can't listen this socket!");
    }

    if (fcntl(serverSocket, F_SETFL, fcntl(serverSocket, F_GETFL, 0) | O_NONBLOCK) == -1) {
        throw proxyException("Can't make server socket nonblock!");
    }

    pollfd server{};
    server.fd = serverSocket;
    server.events = POLLIN;
    pollDescryptors->reserve(MAX_CONNECTIONS);
    pollDescryptors->push_back(server);

    supportedAddressTypes.insert(IPv4_ADDRESS);
    supportedAddressTypes.insert(IPv6_ADDRESS);
    supportedAddressTypes.insert(DOMAIN_NAME);

    setupDNSSignal();
}

Proxy::~Proxy() {
    close(dnsSignal);
    close(serverSocket);
    delete this->transferMap;
    delete this->pollDescryptors;
}

void Proxy::run() {
    std::cerr << "Started accepting connections" << std::endl;
    while (true) {
        try {
            pollManage();
        } catch (proxyException &e) {
            std::cerr << "Exception: " << e.what() << std::endl;
        }
    }
}

void Proxy::removeFromPoll(std::vector<pollfd>::iterator *it) {
    if (close((*it)->fd)) {
        throw proxyException("close");
    }
    (*it)->fd = -(*it)->fd;
}

void Proxy::pollManage() {
    pollfd c{};
    c.fd = -1;
    c.events = POLLIN;
    c.revents = 0;

    poll(pollDescryptors->data(), pollDescryptors->size(), POLL_DELAY);

    for (auto it = pollDescryptors->begin(); it != pollDescryptors->end(); ++it) {
        if (it->fd > 0) {
            if (it->revents & POLLERR) {
                removeFromPoll(&it);
            } else if (it->revents & POLLOUT) {
                sendData(&it);
            } else if (it->revents & POLLIN) {
                if (it->fd == serverSocket) {
                    std::cerr << "Accepting connection" << std::endl;
                    acceptConnection(&c);
                } else if (it->fd == dnsSignal) {
                    std::cerr << "DNS resolved" << std::endl;
                    getResolveResult(&it);
                } else if (passedFullSOCKSprotocol.count(&*it)) {
                    readData(&it);
                } else if (passedHandshake.count(&*it)) {
                    std::cerr << "auth" << std::endl;
                    passIdentification(&it);
                    std::cerr << "pass" << std::endl;
                } else {
                    std::cerr << "handshake" << std::endl;
                    handshake(&it);
                }
            }
        }
    }

    if (c.fd != -1) {
        pollDescryptors->push_back(c);
    }

    if (!waitedCounter) {
        removeDeadDescriptors();
    }
}

void Proxy::passIdentification(std::vector<pollfd>::iterator *clientIterator) {
    pollfd *client = &**clientIterator;
    std::fill(buffer, buffer + BUFFER_LENGTH, 0);
    auto read = recv(client->fd, buffer, BUFFER_LENGTH, 0);

    if (!checkSOCKSRequest(client->fd, read)) {
        removeFromPoll(clientIterator);
        return;
    }

    char addrType = buffer[3];
    switch (addrType) {
        case IPv4_ADDRESS:
            connectToIPv4Address(clientIterator);
            return;
        case IPv6_ADDRESS:
            skipIPV6(clientIterator);
            return;
        case DOMAIN_NAME:
            resolveDomainName(clientIterator);
            return;
        default:
            throw proxyException("Invalid address type");
    }
}

void Proxy::handshake(std::vector<pollfd>::iterator *clientIterator) {
    pollfd *fd = &**clientIterator;
    auto read = recv(fd->fd, buffer, BUFFER_LENGTH, 0);
    buffer[read] = 0;

    char response[2] = {SOCKS_VERSION, char(INVALID_AUTHORISATION)};
    if (buffer[0] != SOCKS_VERSION) {
        send(fd->fd, response, HANDSHAKE_LENGTH, 0);
        removeFromPoll(clientIterator);
        return;
    }

    int availableAuthorisations = buffer[1];
    for (int j = 0, i = 2; j < availableAuthorisations; ++i, ++j) {
        if (buffer[i] == SUPPORTED_AUTHORISATION) {
            response[1] = SUPPORTED_AUTHORISATION;
        }
    }

    send(fd->fd, response, HANDSHAKE_LENGTH, 0);
    if (response[1] != SUPPORTED_AUTHORISATION) {
        removeFromPoll(clientIterator);
        return;
    }
    passedHandshake.insert(fd);
}

void Proxy::acceptConnection(pollfd *client) {
    size_t addSize = sizeof(addr);
    int newClient = accept(serverSocket, (sockaddr *) &addr, (socklen_t *) &addSize);
    if (newClient == -1) {
        throw proxyException("accept");
    }

    std::cerr << "New connection fd: " << newClient << ", addr:"
              << inet_ntoa(addr.sin_addr) << ":" << addr.sin_port << std::endl;
    client->fd = newClient;
}

void Proxy::sendData(std::vector<pollfd>::iterator *clientIterator) {
    auto client = &**clientIterator;
    ssize_t sent;
    while (not(*dataPieces)[client].empty()) {
        sent = send(client->fd, (*dataPieces)[client].data(), (*dataPieces)[client].size(), 0);
        if (sent == -1) {
            if (errno == EWOULDBLOCK)
                return;
            else {
                removeFromPoll(clientIterator);
                return;
            }
        }
        (*dataPieces)[client].erase((*dataPieces)[client].begin(), (*dataPieces)[client].begin() + sent);
    }
    dataPieces->erase(client);
    for (auto &pollDescryptor: *pollDescryptors) {
        if (&pollDescryptor == client) {
            pollDescryptor.events ^= POLLOUT;
            return;
        }
    }
}

void Proxy::removeDeadDescriptors() {
    auto npollDescriptor = new std::vector<pollfd>;
    npollDescriptor->reserve(MAX_CONNECTIONS);

    auto ntransferPipes = new std::unordered_map<pollfd *, pollfd *>;
    auto ndataPieces = new std::unordered_map<pollfd *, std::vector<char> >;
    std::unordered_map<pollfd *, pollfd *> oldNewMap;

    std::unordered_set<pollfd *> newPassedHandshake;
    std::unordered_set<pollfd *> newPassedSocks;

    for (auto &pollDescryptor: *pollDescryptors) {
        if (pollDescryptor.fd > 0) {
            npollDescriptor->push_back(pollDescryptor);
            oldNewMap[&pollDescryptor] = &npollDescriptor->back();
        }
    }

    for (auto &it: *transferMap) {
        if (it.second->fd > 0 and it.first->fd > 0) {
            (*ntransferPipes)[oldNewMap[it.first]] = oldNewMap[it.second];
            (*ntransferPipes)[oldNewMap[it.second]] = oldNewMap[it.first];
        }
    }

    for (const auto &handshake: passedHandshake) {
        if (handshake->fd > 0) {
            newPassedHandshake.insert(oldNewMap[handshake]);
        }
    }

    for (const auto &sock: passedFullSOCKSprotocol) {
        if (sock->fd > 0) {
            newPassedSocks.insert(oldNewMap[sock]);
        }
    }

    for (auto &dataPiece: *dataPieces) {
        if (dataPiece.first->fd > 0) {
            (*ndataPieces)[oldNewMap[dataPiece.first]] = dataPiece.second;
        }
    }

    dataPieces->swap(*ndataPieces);
    delete ndataPieces;

    transferMap->swap(*ntransferPipes);
    delete (ntransferPipes);

    pollDescryptors->swap(*npollDescriptor);
    delete npollDescriptor;

    passedFullSOCKSprotocol.swap(newPassedSocks);
    passedHandshake.swap(newPassedHandshake);
}

void Proxy::printSocksBuffer(int size) {
    for (int i = 0; i < size; ++i) {
        std::cerr << i << " " << (int) buffer[i] << std::endl;
    }
}

bool Proxy::checkSOCKSRequest(int client, ssize_t len) {
    if (len < MINIMUM_SOCKS_REQUEST_LENGTH) {
        char notSupported[] = {SOCKS_VERSION, PROTOCOL_ERROR, 0};
        send(client, notSupported, sizeof(notSupported), 0);
        return false;
    }
    printSocksBuffer(len);
    char socksVersion = buffer[0];
    char socksCommand = buffer[1];
    char addrType = buffer[3];

    if (socksVersion != SOCKS_VERSION) {
        char notSupported[] = {SOCKS_VERSION, PROTOCOL_ERROR, 0};
        send(client, notSupported, sizeof(notSupported), 0);
        std::cerr << "Invalid socks request!" << std::endl;
        return false;
    }
    if (!supportedAddressTypes.count(addrType)) {
        char notSupported[] = {SOCKS_VERSION, ADDRESS_NOT_SUPPORTED, 0};
        send(client, notSupported, sizeof(notSupported), 0);
        return false;
    }
    if (socksCommand != SUPPORTED_OPTION) {
        char notSupported[] = {SOCKS_VERSION, OPTION_NOT_SUPPORTED, 0};
        send(client, notSupported, sizeof(notSupported), 0);
        return false;
    }

    return true;
}

uint16_t Proxy::constructPort(char *port) {
    return *reinterpret_cast<uint16_t *>(port);
}

uint32_t Proxy::constructIPv4Addr(char *port) {
    return *reinterpret_cast<uint32_t *>(port);
}

void Proxy::readData(std::vector<pollfd>::iterator *clientIterator) {

    auto client = &**clientIterator;

    if (!transferMap->count(client)) {
        removeFromPoll(clientIterator);
        return;
    }
    pollfd *to = (*transferMap)[client];

    while (true) {
        auto read = recv(client->fd, buffer, BUFFER_LENGTH - 1, 0);
        if (!read or (read == -1 and errno != EWOULDBLOCK)) {
            to->events |= POLLOUT;
            removeFromPoll(clientIterator);
            return;
        }
        if (errno == EWOULDBLOCK) {
            errno = EXIT_SUCCESS;
            to->events |= POLLOUT;
            return;
        }
        buffer[read] = 0;
        for (int i = 0; i < read; ++i) {
            (*dataPieces)[to].emplace_back(buffer[i]);
        }
    }
}

void Proxy::connectToIPv4Address(std::vector<pollfd>::iterator *clientIterator) {
    auto client = &**clientIterator;
    sockaddr_in targetAddr{};
    socklen_t addrSize = sizeof(targetAddr);
    targetAddr.sin_family = AF_INET;
    targetAddr.sin_port = constructPort(buffer + SOCKS5_OFFSET_BEFORE_ADDR + IPv4_ADDRESS_LENGTH);
    targetAddr.sin_addr.s_addr = constructIPv4Addr(buffer + SOCKS5_OFFSET_BEFORE_ADDR);
    char response[] = {SOCKS_VERSION, SOCKS_SUCCESS, 0, IPv4_ADDRESS, 0, 0, 0, 0, 0, 0};
    auto target = socket(AF_INET, SOCK_STREAM, 0);
    if (target == -1) {
        std::cerr << "Invalid addr or port" << std::endl;
        response[1] = SOCKS_SERVER_ERROR;
        send(client->fd, response, sizeof(response), 0);
        passedHandshake.erase(client);
        removeFromPoll(clientIterator);
        return;
    }
    fcntl(target, F_SETFL, fcntl(target, F_GETFL, 0) | O_NONBLOCK);
    if (connect(target, (sockaddr *) &targetAddr, addrSize) and errno != EINPROGRESS) {
        response[1] = HOST_NOT_REACHABLE;
        send(client->fd, response, sizeof(response), 0);
        passedHandshake.erase(client);
        removeFromPoll(clientIterator);
        return;
    }
    send(client->fd, response, sizeof(response), 0);
    pollfd fd{};
    fd.fd = target;
    fd.events = POLLIN;
    fd.revents = 0;
    *clientIterator = pollDescryptors->insert(*clientIterator + 1, fd);
    passedHandshake.erase(client);
    passedFullSOCKSprotocol.insert(client);
    passedFullSOCKSprotocol.insert(&**clientIterator);
    (*transferMap)[client] = &**clientIterator;
    (*transferMap)[&**clientIterator] = client;
}

void Proxy::skipIPV6(std::vector<pollfd>::iterator *clientIterator) {
    std::cerr << "IPv6 skip" << std::endl;
    auto client = &**clientIterator;
    char response[] = {SOCKS_VERSION, SOCKS_SUCCESS, 0, IPv6_ADDRESS};
    response[1] = SOCKS_SERVER_ERROR;
    send(client->fd, response, sizeof(response), 0);
    passedHandshake.erase(client);
    removeFromPoll(clientIterator);
}

void Proxy::resolveDomainName(std::vector<pollfd>::iterator *clientIterator) {
    struct gaicb *host;
    struct addrinfo *hints;
    struct sigevent sig{};
    host = (gaicb *) calloc(sizeof(gaicb), 1);
    hints = (addrinfo *) calloc(sizeof(addrinfo), 1);

    int domainLength = buffer[4];
    char *domainName = new char[domainLength + 1];
    std::copy(buffer + 5, buffer + domainLength + 6, domainName);
    domainName[domainLength] = '\0';
    std::cerr << "Domain name: " << domainName << std::endl;

    hints->ai_family = AF_INET;
    hints->ai_socktype = SOCK_STREAM;
    hints->ai_flags = AI_PASSIVE;
    host->ar_name = domainName;
    host->ar_request = hints;

    auto *resolver = new ResolverStructure;
    resolver->host = host;
    resolver->port = constructPort(buffer + SOCKS5_OFFSET_BEFORE_ADDR + 1 + domainLength);
    std::cerr << "Port: " << htons(resolver->port) << std::endl;
    resolver->waited = &**clientIterator;
    sig.sigev_notify = SIGEV_SIGNAL;
    sig.sigev_value.sival_ptr = resolver;
    sig.sigev_signo = SIGRTMIN;
    getaddrinfo_a(GAI_NOWAIT, &host, 1, &sig);
    (&**clientIterator)->events = 0;
    waitedCounter++;
}

void Proxy::getResolveResult(std::vector<pollfd>::iterator *clientIterator) {
    auto resolverDescriptor = (&**clientIterator);
    ssize_t s;
    struct signalfd_siginfo fdsi{};
    ResolverStructure *resolver;

    s = read(resolverDescriptor->fd, &fdsi, sizeof(struct signalfd_siginfo));
    if (s < sizeof(struct signalfd_siginfo)) {
        throw proxyException("read");
    }

    char response[] = {SOCKS_VERSION, SOCKS_SUCCESS, 0, IPv4_ADDRESS, 0, 0, 0, 0, 0, 0};

    resolver = (ResolverStructure *) fdsi.ssi_ptr;
    auto client = resolver->waited;

    if (!resolver->host->ar_result) {
        waitedCounter--;
        response[1] = HOST_NOT_REACHABLE;
        send(client->fd, response, sizeof(response), 0);
        close(resolver->waited->fd);
        resolver->waited->fd = -resolver->waited->fd;
        return;
    }

    sockaddr_in targetAddr = *(sockaddr_in *) resolver->host->ar_result->ai_addr;
    targetAddr.sin_family = AF_INET;
    targetAddr.sin_port = resolver->port;
    socklen_t addrSize = sizeof(targetAddr);

    std::cerr << "Resolved: " << inet_ntoa(targetAddr.sin_addr);
    auto target = socket(AF_INET, SOCK_STREAM, 0);
    if (target == -1) {
        std::cerr << "Invalid addr or port" << std::endl;
        response[1] = SOCKS_SERVER_ERROR;
        send(client->fd, response, sizeof(response), 0);
        passedHandshake.erase(client);
        close(resolver->waited->fd);
        resolver->waited->fd = -resolver->waited->fd;
        waitedCounter--;
        return;
    }

    fcntl(target, F_SETFL, fcntl(target, F_GETFL, 0) | O_NONBLOCK);
    if (connect(target, (sockaddr *) &targetAddr, addrSize) and errno != EINPROGRESS) {
        response[1] = HOST_NOT_REACHABLE;
        send(client->fd, response, sizeof(response), 0);
        passedHandshake.erase(client);
        close(resolver->waited->fd);
        resolver->waited->fd = -resolver->waited->fd;
        waitedCounter--;
        return;
    }

    client->events = POLLIN;
    send(client->fd, response, sizeof(response), 0);

    pollfd fd{};
    fd.fd = target;
    fd.events = POLLIN;
    fd.revents = 0;
    pollDescryptors->emplace_back(fd);
    *clientIterator = pollDescryptors->end() - 1;

    passedHandshake.erase(client);
    passedFullSOCKSprotocol.insert(client);
    passedFullSOCKSprotocol.insert(&**clientIterator);

    (*transferMap)[client] = &**clientIterator;
    (*transferMap)[&**clientIterator] = client;

    delete resolver->host->ar_name;
    freeaddrinfo(resolver->host->ar_result);
    free((void *) resolver->host->ar_request);
    free(resolver->host);
    delete resolver;
    waitedCounter--;
}

void Proxy::setupDNSSignal() {
    sigset_t mask;
    sigemptyset(&mask);
    sigaddset(&mask, SIGRTMIN);
    sigprocmask(SIG_BLOCK, &mask, nullptr);
    dnsSignal = signalfd(-1, &mask, 0);
    pollfd fd{};
    fd.fd = dnsSignal;
    fd.events = POLLIN;
    fd.revents = 0;
    pollDescryptors->emplace_back(fd);
}
