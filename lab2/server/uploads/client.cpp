#include <unistd.h>
#include <iostream>

#include "clientClass.h"

int main(int argc, char *argv[]) {
    Client client(argc, argv);
    try {
        client.run();
    } catch (tcpException &e) {
        std::cerr << e.what() << std::endl;
    }

    return 0;
}