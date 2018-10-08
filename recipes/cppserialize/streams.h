//
// Created by yq on 18-10-8.
//

#ifndef CPPSERIALIZE_STREAMS_H
#define CPPSERIALIZE_STREAMS_H

#include <serialize.h>

#include <vector>
#include <bits/ios_base.h>
#include <cstring>

class CDataStream {
public:

    CDataStream() : readPos(0) {

    }

    std::vector<char> serializedData;

    size_t readPos;

    template<typename T>
    CDataStream &operator<<(const T &obj) {
        ::Serialize(*this, obj);
        return (*this);
    }

    // TODO: 为什么bitcoin的operator>>用右值? 是bitcoin的语法错误吧!
    template<typename T>
    CDataStream &operator>>(T &obj) {
        ::Unserialize(*this, obj);
        return (*this);
    }

    void read(char *p, size_t len) {
        size_t readPosNext = readPos + len;

        if (readPosNext > serializedData.size()) {
            throw std::ios_base::failure("CDataStream::read(): end of data");
        }

        memcpy(p, &serializedData[readPos], len);

        if (readPosNext == serializedData.size()) {
            readPos = 0;
            serializedData.clear();
            return;
        }

        readPos = readPosNext;
    }

    void write(const char *p, size_t len) {
        serializedData.insert(serializedData.end(), p, p + len);
    }
};

#endif //CPPSERIALIZE_STREAMS_H
