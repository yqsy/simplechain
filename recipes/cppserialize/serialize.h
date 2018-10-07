//
// Created by yq on 18-10-7.
//

#ifndef CPPSERIALIZE_SERIALIZE_H
#define CPPSERIALIZE_SERIALIZE_H


// http://www.cplusplus.com/reference/cstdint/
#include <stdint.h>

#include <endian.h>


// 内存数字 -> 小端法输出二进制
template<typename Stream>
inline void ser_writedate8(Stream &s, uint8_t obj) {
    s.write((char *) &obj, 1);
}

template<typename Stream>
inline void ser_writedata16(Stream &s, uint16_t obj) {
    obj = htole16(obj);
    s.write((char *) &obj, 2);
}

template<typename Stream>
inline void ser_writedata32(Stream &s, uint32_t obj) {
    obj = htole32(obj);
    s.write((char *) &obj, 4);
}

template<typename Stream>
inline void ser_writedata64(Stream &s, uint64_t obj) {
    obj = htole64(obj);
    s.write((char *) &obj, 8);
}

// 小端法输入二进制 -> 内存数字
template<typename Stream>
inline uint8_t ser_readdata8(Stream &s) {
    uint8_t obj;
    s.read((char *) &obj, 1);
    return obj;
}

template<typename Stream>
inline uint16_t ser_readdata16(Stream &s) {
    uint16_t obj;
    s.read((char *) &obj, 2);
    return le16toh(obj);
}

template<typename Stream>
inline uint32_t ser_readdata32(Stream &s) {
    uint32_t obj;
    s.read((char *) &obj, 4);
    return le32toh(obj);
}


template<typename Stream>
inline uint64_t ser_readdata64(Stream &s) {
    uint64_t obj;
    s.read((char *) &obj, 8);
    return le64toh(obj);
}

// 序列化
template<typename Stream>
inline void Serialize(Stream &s, int8_t a) { ser_writedate8(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint8_t a) { ser_writedate8(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, int16_t a) { ser_writedate16(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint16_t a) { ser_writedate16(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, int32_t a) { ser_writedate32(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint32_t a) { ser_writedate32(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, int64_t a) { ser_writedate32(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint64_t a) { ser_writedate64(s, a); };


// 反序列化
template<typename Stream>
inline void Unserialize(Stream &s, int8_t &a) { a = ser_readdata8(s); };

template<typename Stream>
inline void Unserialize(Stream &s, uint8_t &a) { a = ser_readdata8(s); };



#endif //CPPSERIALIZE_SERIALIZE_H
