#pragma once

#include <exception>
#include <utility>

class tcpException : public std::exception {
private:
    std::string errorString;
public:
    explicit tcpException(std::string errStr) {
        errorString = std::move(errStr);
    }

    [[nodiscard]] const char *what() const noexcept override {
        return errorString.c_str();
    }
};

class parseException : public tcpException {
public:
    explicit parseException(std::string errStr) : tcpException(std::move(errStr)) {}
};

class getServerException : public tcpException {
public:
    explicit getServerException(std::string errStr) : tcpException(std::move(errStr)) {}
};

class fileException : public tcpException {
public:
    explicit fileException(std::string errStr) : tcpException(std::move(errStr)) {}
};

class connectionException : public tcpException {
public:
    explicit connectionException(std::string errStr) : tcpException(std::move(errStr)) {}
};

class sendException : public tcpException {
public:
    explicit sendException(std::string errStr) : tcpException(std::move(errStr)) {}
};

class socketException : public tcpException {
public:
    explicit socketException(std::string errStr) : tcpException(std::move(errStr)) {}
};
