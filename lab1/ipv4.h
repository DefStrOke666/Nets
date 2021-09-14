#pragma once

#include <sys/poll.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netdb.h>

#include "utils.h"
#include "exceptions.h"

class IPV4 {
private:
    FriendList fl{};
    std::string myName;
    std::string addr;
    int inSock{}, outSock{};
    struct sockaddr_in bcAddr{};

    void initSockets() {
        int optVal = 1;
        if ((inSock = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) < 0) {
            throw multicastException("input socket");
        }

        if (setsockopt(inSock, SOL_SOCKET, SO_REUSEADDR, &optVal, sizeof optVal) != 0) {
            throw multicastException("input: setsockopt: SO_REUSEADDR");
        }

        struct sockaddr_in sockAddr{};
        sockAddr.sin_family = AF_INET;
        sockAddr.sin_port = htons(port);
        sockAddr.sin_addr.s_addr = htonl(INADDR_ANY);

        if (bind(inSock, (const struct sockaddr *) &sockAddr, sizeof(sockAddr)) != 0) {
            throw multicastException("bind");
        }

        struct ip_mreq mreq{};
        inet_aton(addr.c_str(), &(mreq.imr_multiaddr));
        mreq.imr_interface.s_addr = htonl(INADDR_ANY);
        if (setsockopt(inSock, IPPROTO_IP, IP_ADD_MEMBERSHIP, &mreq, sizeof(mreq)) != 0) {
            throw multicastException("input: setsockopt: IP_ADD_MEMBERSHIP");
        }

        if ((outSock = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) < 0) {
            throw multicastException("output socket");
        }
    }

public:
    IPV4(std::string ipAddr, std::string name) {
        addr = std::move(ipAddr);
        myName = std::move(name);
    }

    void run() {
        initSockets();

        bcAddr.sin_family = AF_INET;
        bcAddr.sin_addr.s_addr = inet_addr(addr.c_str());
        bcAddr.sin_port = htons(port);

        int fdsCount = 1;
        struct pollfd fd[fdsCount];
        fd[0].fd = inSock;
        fd[0].events = POLLIN;

        char buf[bufsize];
        char hostname[NI_MAXHOST];

        struct sockaddr_in someFriend{};
        int len = sizeof(someFriend);
        
        fl.setName(myName);

        while (true) {
            if (sendto(outSock, myName.c_str(), myName.length(), 0, (sockaddr *) &bcAddr, sizeof(bcAddr)) < 0) {
                throw multicastException("sendto");
            }

            fl.removeExpired();

            int ret = poll(fd, fdsCount, 500);
            if (ret < 0) {
                throw multicastException("poll");
            }
            if (ret != 0) {
                long read = recvfrom(inSock, &buf, bufsize, MSG_WAITALL, (sockaddr *) &someFriend, (socklen_t *) &len);
                if (read < 0) {
                    throw multicastException("recvfrom");
                }

                buf[read] = '\0';
                std::string friendName(buf);
                if (friendName == myName || read < 1) {
                    continue;
                }

                if (getnameinfo((sockaddr *) &someFriend, len, hostname, NI_MAXHOST, nullptr, 0, NI_NUMERICHOST) != 0) {
                    throw multicastException("getnameinfo");
                }

                fl.addFriend(hostname, friendName);
            }
            fl.showFriendList();
        }
    }
};