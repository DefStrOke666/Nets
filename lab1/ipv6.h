#pragma once

#include <sys/poll.h>
#include <sys/socket.h>
#include <net/if.h>
#include <unistd.h>

#include "utils.h"
#include "exceptions.h"

class IPV6 {
private:
    FriendList fl{};
    std::string myName;
    std::string addr;
    int inSock{}, outSock{};
    struct sockaddr_in6 bcAddr{};

    void initSockets() {
        int optVal = 1;
        if ((inSock = socket(AF_INET6, SOCK_DGRAM, 0)) < 0) {
            throw multicastException("input socket");
        }

        if (setsockopt(inSock, SOL_SOCKET, SO_REUSEADDR, &optVal, sizeof optVal) != 0) {
            throw multicastException("input: setsockopt: SO_REUSEADDR");
        }

        int hops = 255;
        if (setsockopt(inSock, IPPROTO_IPV6, IPV6_MULTICAST_HOPS, &hops, sizeof(hops)) != 0) {
            throw multicastException("input: setsockopt: IPV6_MULTICAST_HOPS");
        }

        struct sockaddr_in6 server_addr{};
        server_addr.sin6_family = AF_INET6;
        server_addr.sin6_port = htons(port);
        server_addr.sin6_addr = in6addr_any;
        if (bind(inSock, (const struct sockaddr *) &server_addr, sizeof(server_addr)) != 0) {
            throw multicastException("bind");
        }

        struct ipv6_mreq group{};
        group.ipv6mr_interface = if_nametoindex("eth0");
        inet_pton(AF_INET6, addr.c_str(), &group.ipv6mr_multiaddr);
        if (setsockopt(inSock, IPPROTO_IPV6, IPV6_JOIN_GROUP, &group, sizeof group) != 0) {
            throw multicastException("input: setsockopt: IPV6_JOIN_GROUP");
        }

        if ((outSock = socket(AF_INET6, SOCK_DGRAM, 0)) < 0) {
            throw multicastException("output socket");
        }
    }

public:
    IPV6(std::string ipAddr, std::string name) {
        addr = std::move(ipAddr);
        myName = std::move(name);
    }

    void run() {
        initSockets();

        bcAddr.sin6_family = AF_INET;
        inet_pton(AF_INET6, addr.c_str(), &bcAddr.sin6_addr);
        bcAddr.sin6_port = htons(port);

        int fdsCount = 1;
        struct pollfd fd[fdsCount];
        fd[0].fd = inSock;
        fd[0].events = POLLIN;

        char buf[bufsize];
        char hostname[INET6_ADDRSTRLEN];

        struct sockaddr_in6 someFriend{};
        int len = sizeof(someFriend);

        fl.setName(myName);

        while (true) {
            if (sendto(outSock, myName.c_str(), myName.length(), 0, (sockaddr *) &bcAddr, sizeof(bcAddr)) < 0) {
                throw multicastException("sendto");
            }

            fl.removeExpired();

            sleep(1);
            int ret = poll(fd, fdsCount, 2000);
            if (ret < 0) {
                throw multicastException("poll");
            }
            if (ret != 0) {
                long read = recvfrom(inSock, &buf, bufsize, MSG_WAITALL, (sockaddr *) &someFriend, (socklen_t *) &len);
                if (read < 0) {
                    throw multicastException("recvfrom");
                }

                buf[read] = '\0';
                std::string friendName = buf;
                if (friendName == myName) {
                    continue;
                }

                inet_ntop(AF_INET6, &(someFriend.sin6_addr), hostname, INET6_ADDRSTRLEN);

                fl.add(hostname, friendName);
            }
            fl.showFriendList();
        }
    }
};