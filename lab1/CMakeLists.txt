cmake_minimum_required(VERSION 3.16)
project(lab1)

set(CMAKE_CXX_STANDARD 20)

add_executable(main multicast.cpp ipv4.h ipv6.h exceptions.h)