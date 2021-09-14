#include "ipv4.h"
#include "ipv6.h"
#include <cerrno>
#include <cstring>
#include <iostream>

int main(int argc, char *argv[]) {
    if (argc < 2) {
        std::cerr << "usage: " << argv[0] << " 224.0.0.0-224.255.255.255 or ff00:0:0:0:0:0:0:0-ffff:ffff..."
                  << std::endl;
        return -1;
    }

    char buf[sizeof(struct in6_addr)];
    const std::string myName = generateName();

    if (inet_pton(AF_INET, argv[1], buf) == 1) {
        IPV4 ipv4(argv[1], myName);
        try {
            std::cout << "Waiting for connections" << std::endl;
            ipv4.run();
        } catch (multicastException &e) {
            std::cerr << "Exception: " << e.what() << ": " << std::strerror(errno) << std::endl;
            return -1;
        }
    } else if (inet_pton(AF_INET6, argv[1], buf) == 1) {
        IPV6 ipv6(argv[1], myName);
        try {
            std::cout << "Waiting for connections" << std::endl;
            ipv6.run();
        } catch (multicastException &e) {
            std::cerr << "Exception: " << e.what() << ": " << std::strerror(errno) << std::endl;
            return -1;
        }
    } else {
        std::cout << "Address is not ipv4 or ipv6" << std::endl;
        std::cerr << "usage: " << argv[0] << " 224.0.0.0-224.255.255.255 or ff00:0:0:0:0:0:0:0-ffff:ffff..."
                  << std::endl;
    }

    return 0;
}

