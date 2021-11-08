cmake_minimum_required(VERSION 3.13)
project(socks5)

set(CMAKE_CXX_STANDARD 20)

add_executable(proxy main.cpp proxy.cpp)