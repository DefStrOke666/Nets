#pragma once

#include <exception>
#include <utility>

class multicastException : public std::exception {
private:
    std::string errorString;
public:
    explicit multicastException(std::string errStr) {
        errorString = std::move(errStr);
    }

    [[nodiscard]] const char *what() const noexcept override {
        return errorString.c_str();
    }
};