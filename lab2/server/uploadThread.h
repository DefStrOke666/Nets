#pragma once

#include <chrono>
#include <sstream>
#include "exceptions.h"

using namespace std::chrono;

void showSpeed(microseconds elapsed, time_point<high_resolution_clock> recvStart, int bytesReceived,
               int totalBytesReceived, int iter, std::string &client) {
    int megaByte = 1024 * 1024;
    double speed = (bytesReceived / (elapsed.count() / 1e6 / iter)) / megaByte;
    auto now = high_resolution_clock::now();
    microseconds duration = duration_cast<microseconds>(now - recvStart);
    double bytesPerSecond = totalBytesReceived / (duration.count() / 1e6);
    double avSpeed = bytesPerSecond / megaByte;
    std::cout << std::setprecision(1) << std::fixed;
    std::cout << "\t" + client + " speed: " << speed << " MB/s" << " av. speed: " << avSpeed << " MB/s" << std::endl;
}

int prepareFile(int clientSock, int *fileSize, std::string &client, std::ofstream &file) {
    int bufSize = 4200;
    char buf[bufSize];
    int received = recv(clientSock, buf, bufSize, 0);
    if (received == -1) {
        close(clientSock);
        throw connectionException(std::string("recv: ") + strerror(errno));
    }
    if (received == 0) {
        close(clientSock);
        throw connectionException("Client disconnected");
    }

    std::string fileName, fileSizeStr;
    std::istringstream ss(buf);
    ss >> fileName;
    ss >> fileSizeStr;
    *fileSize = std::stoi(fileSizeStr);

    if (fileName.empty()) {
        throw fileException("No file name");
    }
    if (*fileSize == 0) {
        throw fileException("No file size");
    }

    std::cout << "\tFile name: " << fileName << std::endl;
    std::cout << "\tFile size: " << *fileSize << std::endl;

    std::string filePath("./uploads/" + fileName);
    int ret = unlink(filePath.c_str());
    if (!ret) {
        std::cout << "\tDeleted " << fileName << std::endl;
    } else {
        std::cout << "\tunlink " + fileName + ": " << strerror(errno) << std::endl;
    }

    file.open(filePath);
    if (!file.is_open()) {
        throw fileException("Cannot open " + fileName);
    }
    std::cout << "\tOpened " << fileName << std::endl;

    int pos;
    for (pos = 0; pos < sizeof(buf); ++pos) {
        if (buf[pos] == '\n') {
            pos++;
            break;
        }
    }
    char part[received - pos];
    for (int i = pos; i < received; ++i) {
        part[i - pos] = buf[i];
    }
    file.write(part, received - pos);

    return received - pos;
}

void sendResponse(int clientSock, const std::string& totalBytesReceived, std::ofstream &file, std::string &client) {
    int sent = send(clientSock, totalBytesReceived.c_str(), totalBytesReceived.size(), 0);
    if (sent == -1) {
        file.close();
        close(clientSock);
        throw connectionException(std::string("send: ") + strerror(errno));
    }
    if (sent == 0) {
        file.close();
        close(clientSock);
        throw connectionException(client + std::string(" disconnected"));
    }
}

void uploadFile(int clientSock, std::string &client) {
    int fileSize;
    std::ofstream file;
    int totalBytesReceived = prepareFile(clientSock, &fileSize, client, file);

    int bufSize = 1024 * 4;
    int iter = 0, received;
    char buf[bufSize];
    time_point<high_resolution_clock> start, end, recvStart, updateTime = high_resolution_clock::now();
    microseconds timeElapsed{microseconds(0)};
    recvStart = high_resolution_clock::now();
    for (int i = 0; totalBytesReceived < fileSize; ++i) {
        iter++;

        start = high_resolution_clock::now();
        int bytesLeft = fileSize - totalBytesReceived;
        received = recv(clientSock, buf, std::min(bufSize, bytesLeft), 0);
        if (received == -1) {
            file.close();
            close(clientSock);
            throw connectionException(std::string("recv: ") + strerror(errno));
        }
        if (received == 0) {
            file.close();
            close(clientSock);
            throw connectionException(client + std::string(" disconnected"));
        }
        file.write(buf, received);
        totalBytesReceived += received;

        end = high_resolution_clock::now();
        timeElapsed += duration_cast<microseconds>(end - start);

        time_point<high_resolution_clock> time = high_resolution_clock::now();
        microseconds dur = duration_cast<microseconds>(time - updateTime);
        if (dur.count() / 1e6 > 3) { // 3sec
            showSpeed(timeElapsed, recvStart, received, totalBytesReceived, iter, client);
            updateTime = high_resolution_clock::now();
            timeElapsed = microseconds(0);
            iter = 0;
        }
    }
    showSpeed(timeElapsed, recvStart, received, totalBytesReceived, iter, client);
    sendResponse(clientSock, std::to_string(totalBytesReceived), file, client);

    file.close();
    close(clientSock);

    std::cout << "Connection with " + client + " closed, total bytes received: " << totalBytesReceived << std::endl;
}

void upload(int clientSock, std::string client) {
    try {
        uploadFile(clientSock, client);
    } catch (tcpException &e) {
        std::cerr << e.what() << std::endl;
        pthread_exit(reinterpret_cast<void *>(EXIT_FAILURE));
    }
    pthread_exit(EXIT_SUCCESS);
}