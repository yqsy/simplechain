cmake_minimum_required(VERSION 3.12)
project(cpphash)

set(CMAKE_CXX_STANDARD 11)

add_executable(cpphash main.cpp sha256.cpp sha256.h)