#pragma once

#include <random>
#include <fstream>
#include <map>
#include <utility>

static int port = 9999;
static int bufsize = 256;

std::string generateName() {
    std::ifstream names("../names.txt");
    long lineCount = std::count(std::istreambuf_iterator<char>(names), std::istreambuf_iterator<char>(), '\n');
    names.seekg(0);

    std::random_device dev;
    std::mt19937 rng(dev());
    std::uniform_int_distribution<std::mt19937::result_type> distrib(1, lineCount);

    std::string name;
    unsigned long count = distrib(rng);
    for (int i = 0; i < count; ++i) {
        getline(names, name);
    }

    char hexStr[10];
    snprintf(hexStr, 10, "%lX", count);
    return name + "-" + hexStr;
}

class FriendList {
private:
    std::string myName;
    const int timeoutSeconds = 30;
    struct Friend {
        std::string name;
        time_t lastSeen{};
        std::string addr;
    };
    std::map<std::string, Friend> friends;

public:

    void setName(std::string name){
        myName = std::move(name);
    }

    void add(const std::string &addr, const std::string &name) {
        if (friends.find(name) != friends.end()) {
            friends.erase(friends.find(name));
        }

        Friend fr{
                name,
                time(nullptr),
                addr,
        };
        friends[name] = fr;
    }

    void removeExpired() {
        for (const auto &el: friends) {
            if (time(nullptr) - el.second.lastSeen > timeoutSeconds) {
                friends.erase(el.first);
            }
        }
    }

    void showFriendList() {
        system("clear");
        printf("My name is %s\n\n", myName.c_str());
        printf("%-15s |%-15s |%-15s\n", "Name", "Address", "Last seen");
        printf("------------------------------------------------------------\n");

        for (const auto &fr: friends) {
            char buf[100];
            const auto time = localtime(&fr.second.lastSeen);
            strftime(buf, sizeof(buf), "%H:%M:%S", time);
            std::string timeStr(buf);
            printf("%-15s |%-15s |%-15s\n", fr.second.name.c_str(), fr.second.addr.c_str(), buf);
        }
    }
};
