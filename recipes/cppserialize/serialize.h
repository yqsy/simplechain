//
// Created by yq on 18-10-7.
//

#ifndef CPPSERIALIZE_SERIALIZE_H
#define CPPSERIALIZE_SERIALIZE_H


// http://www.cplusplus.com/reference/cstdint/
#include <stdint.h>


template<typename Stream>
inline void ser_writedate8(Stream &s, uint8_t obj) {
    s.Write();
}


template<typename Stream>
inline void Serialize(Stream &s, int8_t a) {};


#endif //CPPSERIALIZE_SERIALIZE_H
