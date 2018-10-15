cmake_minimum_required(VERSION 3.12)
project(cpplock)

set(CMAKE_CXX_STANDARD 11)

add_executable(cpplock main.cpp sync.h)