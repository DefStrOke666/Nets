#include <iostream>
#include "serverClass.h"

int main(int argc, char *argv[]) {
    Server server(argc, argv);
    try {
        server.run();
    } catch (tcpException &e) {
        std::cerr << e.what() << std::endl;
    }

    return 0;
}