cmake_minimum_required(VERSION 3.12)
project(cppserialize)

set(CMAKE_CXX_STANDARD 11)

include_directories(/mnt/disk1/go/src/github.com/yqsy/simplechain/recipes/cppserialize)

add_executable(cppserialize main.cpp serialize.h)

