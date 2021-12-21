#include <iostream>
#include <cstring>
#include "proxy.h"

int main(int argc, char *argv[]) {
    if (argc < 2) {
        std::cerr << "usage: " << argv[0] << " port" << std::endl;
        return -1;
    }

    int port = std::stoi(argv[1]);
    try {
        auto *proxy = new Proxy(port);
        proxy->run();
    } catch (proxyException &e) {
        std::cerr << "Exception: " << e.what() << ": " << std::strerror(errno) << std::endl;
        return -1;
    }

    return 0;
}
