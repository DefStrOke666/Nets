#pragma once

#include <sys/stat.h>
#include <netdb.h>
#include <cstring>
#include <arpa/inet.h>
#include <fstream>
#include <chrono>
#include <filesystem>
#include <unistd.h>
#include <iomanip>
#include "exceptions.h"

using namespace std::chrono;

class Client {
    std::string serverAddress;
    std::string serverPort;
    std::string filePath;
    std::string fileName;

    int sock;
    int fileSize;
    int totalBytesSent = 0;
    int bytesLeft;

    time_point<high_resolution_clock> sendStart;

    struct addrinfo *server;

    std::ifstream file;

    int argc;
    char **argv;
private:

    void parseArgs() {
        bool gotFile, gotServAddr, gotServPort;

        int c;
        while ((c = getopt(argc, argv, "a:p:f:h")) != -1) {
            switch (c) {
                case 'a':
                    gotServAddr = true;
                    serverAddress.assign(optarg);
                    break;
                case 'p':
                    gotServPort = true;
                    serverPort.assign(optarg);
                    break;
                case 'f':
                    gotFile = true;
                    filePath.assign(optarg);
                    break;
                case '?':
                case 'h':
                    throw parseException(std::string("Usage: ") + *argv[0]
                                         + std::string(" -a [address] -p [port] -f [file path]"));
                default:
                    throw parseException("parseException");
            }
        }

        if (!gotServAddr) {
            throw parseException("Server address is not provided with -a [addr or domain name] option");
        } else if (!gotServPort) {
            throw parseException("Server port is not provided with -p [port] option");
        } else if (!gotFile) {
            throw parseException("File path is not provided with -f [file path] option");
        }

        fileName = filePath.substr(filePath.find_last_of("/\\") + 1);
        if (fileName.size() * sizeof(char) > 4096) {
            throw parseException("File name is larger than 4096 bytes");
        }
    }

    void getServer() {
        struct addrinfo hints{};
        memset(&hints, 0, sizeof(struct addrinfo));
        hints.ai_family = AF_UNSPEC;
        hints.ai_socktype = SOCK_STREAM;
        hints.ai_flags |= AI_CANONNAME;

        int err;
        if ((err = getaddrinfo(serverAddress.c_str(), serverPort.c_str(), &hints, &server)) != 0) {
            throw getServerException(std::string("getaddrinfo: ") + gai_strerror(err));
        }

        char addr[100];
        inet_ntop(server->ai_family, &((struct sockaddr_in *) server->ai_addr)->sin_addr, addr, 100);
        std::cout << "Server address: " << addr << " canon name: " << server->ai_canonname << std::endl;
    }

    void connectToServ() {
        if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
            throw connectionException(std::string("socket: ") + strerror(errno));
        }

        int optVal = 1;
        if (setsockopt(sock, SOL_SOCKET, SO_KEEPALIVE, &optVal, sizeof optVal) != 0) {
            close(sock);
            throw connectionException(std::string("setsockopt: SO_KEEPALIVE ") + strerror(errno));
        }

        if (connect(sock, (struct sockaddr *) server->ai_addr, server->ai_addrlen) != 0) {
            close(sock);
            throw connectionException(std::string("connect: ") + strerror(errno));
        }

        std::cout << "Connected to " << server->ai_canonname << std::endl;
    }

    void getFile() {
        file.open(filePath, std::ios::in | std::ios::binary);

        if (!file.is_open()) {
            throw fileException("Cannot open file " + fileName);
        }
        std::cout << "Opened " << fileName << std::endl;

        struct stat st{};
        stat(filePath.c_str(), &st);
        fileSize = st.st_size;
        if (fileSize > 1e12) { // 1TB
            throw fileException("File is larger than 1TB");
        }
    }

    void sendFileInfo() {
        std::string info = fileName + " " + std::to_string(fileSize) + "\n";
        int sent = send(sock, info.c_str(), info.size(), 0);
        if (sent == -1) {
            throw sendException(std::string("send: ") + strerror(errno));
        }
    }

    void showSpeed(microseconds elapsed, int bytesSent, int iter) {
        system("clear");
        printf("%s: %d / %d\n", fileName.c_str(), totalBytesSent, fileSize);
        double complete = (totalBytesSent / double(fileSize)) * 100;
        int columns = 40;
        double percentPerColumn = 100 / double(columns);
        for (int i = 0; i < columns; ++i) {
            if (i * percentPerColumn < complete) {
                printf("▓");
            } else {
                printf("░");
            }
        }
        printf("[%.0f%%]\n", complete);

        int megaByte = 1024 * 1024;
        double speed = (bytesSent / (elapsed.count() / 1e6 / iter)) / megaByte;
        auto now = high_resolution_clock::now();
        microseconds duration = duration_cast<microseconds>(now - sendStart);
        double bytesPerSecond = totalBytesSent / (duration.count() / 1e6);
        double avSpeed = bytesPerSecond / megaByte;
        printf("%15s %8.1f MB/s \n", "Speed:", speed);
        printf("%15s %8.1f MB/s \n", "Av. speed:", avSpeed);
        printf("%15s %8.3f sec\n", "Time elapsed:", duration.count() / 1e6);
        printf("%15s %8.3f sec\n", "Time left:", bytesLeft / bytesPerSecond);
    }

    void getServerResponse() {
        int bufSize = 100;
        char buf[bufSize];
        int received = recv(sock, buf, bufSize, 0);
        if (received == -1) {
            file.close();
            close(sock);
            throw connectionException(std::string("recv: ") + strerror(errno));
        }
        if (received == 0) {
            file.close();
            close(sock);
            throw connectionException(serverAddress + std::string(" disconnected"));
        }
        int responseBytes = std::stoi(buf);
        if (responseBytes == fileSize) {
            std::cout << "OK" << std::endl;
        } else {
            std::cout << "FAIL" << std::endl;
        }
    }

    void sendFile() {
        sendFileInfo();

        int bufSize = 1024 * 4;
        char buffer[bufSize];
        int iter = 0, sent;
        time_point<high_resolution_clock> start, end, updateTime = high_resolution_clock::now();
        microseconds timeElapsed{microseconds(0)};
        sendStart = high_resolution_clock::now();
        bytesLeft = fileSize;
        while (file.good() | !file.eof()) {
            iter++;

            start = high_resolution_clock::now();
            file.read(buffer, bufSize);
            sent = send(sock, buffer, std::min(bufSize, bytesLeft), 0);
            if (sent == -1) {
                throw sendException(std::string("send: ") + strerror(errno));
            }
            end = high_resolution_clock::now();

            totalBytesSent += sent;
            bytesLeft = fileSize - totalBytesSent;
            timeElapsed += duration_cast<microseconds>(end - start);

            time_point<high_resolution_clock> time = high_resolution_clock::now();
            microseconds dur = duration_cast<microseconds>(time - updateTime);
            if (dur.count() / 1e3 > 50) { // 50ms
                showSpeed(timeElapsed, sent, iter);
                updateTime = high_resolution_clock::now();
                timeElapsed = microseconds(0);
                iter = 0;
            }
        }
        showSpeed(timeElapsed, sent, iter);
        getServerResponse();
    }

public:

    Client(int argCount, char **args) {
        argc = argCount;
        argv = args;
    }

    void run() {
        parseArgs();

        getServer();

        connectToServ();

        getFile();

        freeaddrinfo(server);

        sendFile();
        file.close();
        close(sock);
    }
};
