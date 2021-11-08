#ifndef SOCKS5_EXCEPTIONS_H
#define SOCKS5_EXCEPTIONS_H

#include <exception>
#include <utility>

class proxyException : public std::exception {
private:
    std::string errorString;
public:
    explicit proxyException(std::string errStr) {
        errorString = std::move(errStr);
    }

    [[nodiscard]] const char *what() const noexcept override {
        return errorString.c_str();
    }
};

#endif
