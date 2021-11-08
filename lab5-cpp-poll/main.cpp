#include <iostream>
#include "proxy.h"

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "usage: " << argv[0] << " port" << std::endl;
        return -1;
    }

    int port = std::stoi(argv[1]);
    auto *proxy = new Proxy(port);

    proxy->run();

    return 0;
}
