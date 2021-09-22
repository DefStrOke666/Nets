#pragma once

#include <cstring>
#include <arpa/inet.h>
#include <fstream>
#include <filesystem>
#include <thread>
#include <unistd.h>
#include <iomanip>
#include "exceptions.h"
#include "uploadThread.h"

class Server {
    int listenPort;
    int maxClients = 10;
    int sock;

    struct sockaddr_in serverAddr{};

    int argc;
    char **argv;
private:

    void parseArgs() {
        bool gotServPort;

        int c;
        while ((c = getopt(argc, argv, "p:h")) != -1) {
            switch (c) {
                case 'p':
                    gotServPort = true;
                    listenPort = std::stoi(optarg);
                    break;
                case '?':
                case 'h':
                    throw parseException(std::string("Usage: ") + *argv[0] + std::string(" -p [port]"));
                default:
                    throw parseException("parseException");
            }
        }

        if (!gotServPort) {
            throw parseException("Server port is not provided with -p [port] option");
        }
    }

    void prepareFolder() {
        std::filesystem::create_directories("./uploads");
    }

    void listenSock() {
        if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
            throw socketException(std::string("socket: ") + strerror(errno));
        }

        bzero(&serverAddr, sizeof(serverAddr));
        serverAddr.sin_family = AF_INET;
        serverAddr.sin_addr.s_addr = INADDR_ANY;
        serverAddr.sin_port = htons(listenPort);
        if (bind(sock, (struct sockaddr *) &serverAddr, sizeof(serverAddr)) < 0) {
            close(sock);
            throw socketException(std::string("bind: ") + strerror(errno));
        }

        if (listen(sock, maxClients) != 0) {
            close(sock);
            throw socketException(std::string("listen: ") + strerror(errno));
        }
    }

    void acceptConnections() const {
        char addr[100];

        while (true) {
            struct sockaddr_in clientAddr{};
            socklen_t sockLen = sizeof(clientAddr);
            int newSock = accept(sock, (struct sockaddr *) &clientAddr, &sockLen);
            if (newSock == -1) {
                close(sock);
                throw connectionException(std::string("accept: ") + strerror(errno));
            }

            inet_ntop(clientAddr.sin_family, &clientAddr.sin_addr.s_addr, addr, 100);
            std::cout << addr << " connected" << std::endl;

            std::thread (upload, newSock, std::string(addr)).detach();
        }
    };

public:

    Server(int argCount, char **args) {
        argc = argCount;
        argv = args;
    }

    void run() {
        parseArgs();

        prepareFolder();

        listenSock();

        std::cout << "Listening on port " << listenPort << std::endl;

        acceptConnections();
        close(sock);
    }
};
