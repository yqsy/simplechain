//
// Created by yq on 18-10-7.
//

#ifndef CPPSERIALIZE_SERIALIZE_H
#define CPPSERIALIZE_SERIALIZE_H


// http://www.cplusplus.com/reference/cstdint/
#include <stdint.h>

#include <endian.h>

#include <vector>
#include <limits>
#include <ios>
#include <memory>


struct deserialize_type {
};
constexpr deserialize_type deserialize{};

// 读的时候最多读 32MB 笔交易??
static const unsigned int MAX_SIZE = 0x02000000;

template<typename T>
inline T *NCONST_PTR(const T *val) {
    return const_cast<T *>(val);
}

#define ADD_SERIALIZE_METHODS                                        \
    template<typename Stream>                                        \
    void Serialize(Stream &s) const {                                \
        NCONST_PTR(this)->SerializationOp(s, CSerActionSerialize()); \
    }                                                                \
    template<typename Stream>                                        \
    void Unserialize(Stream &s) {                                      \
        SerializationOp(s, CSerActionUnserialize());                    \
    }

#define READWRITE(...) (::SerReadWriteMany(s, ser_action, __VA_ARGS__))


// 内存整数 -> 小端法输出二进制
template<typename Stream>
inline void ser_writedata8(Stream &s, uint8_t obj) {
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

// 小端法输入二进制 -> 内存整数
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

// 内存浮点数 -> 整数输出
inline uint32_t ser_float_to_uint32(float x) {
    union {
        float x;
        uint32_t y;
    } tmp{};
    tmp.x = x;
    return tmp.y;
}

inline uint64_t ser_double_to_uint64(double x) {
    union {
        double x;
        uint64_t y;
    } tmp{};

    tmp.x = x;
    return tmp.y;
}


// 整数输入 -> 内存浮点数
inline float ser_uint32_to_float(uint32_t y) {
    union {
        float x;
        uint32_t y;
    } tmp{};
    tmp.y = y;
    return tmp.x;
}

inline double ser_uint64_to_double(uint64_t y) {
    union {
        double x;
        uint64_t y;
    } tmp{};
    tmp.y = y;
    return tmp.x;
}


// 基础类型序列化
template<typename Stream>
inline void Serialize(Stream &s, int8_t a) { ser_writedata8(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint8_t a) { ser_writedata8(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, int16_t a) { ser_writedata16(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint16_t a) { ser_writedata16(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, int32_t a) { ser_writedata32(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint32_t a) { ser_writedata32(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, int64_t a) { ser_writedata32(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, uint64_t a) { ser_writedata64(s, a); };

template<typename Stream>
inline void Serialize(Stream &s, float a) { ser_writedata32(s, ser_float_to_uint32(a)); };

template<typename Stream>
inline void Serialize(Stream &s, double a) { ser_writedata64(s, ser_double_to_uint64(a)); };


// 基础类型反序列化
template<typename Stream>
inline void Unserialize(Stream &s, int8_t &a) { a = ser_readdata8(s); };

template<typename Stream>
inline void Unserialize(Stream &s, uint8_t &a) { a = ser_readdata8(s); };

template<typename Stream>
inline void Unserialize(Stream &s, int16_t &a) { a = ser_readdata16(s); };

template<typename Stream>
inline void Unserialize(Stream &s, uint16_t &a) { a = ser_readdata16(s); };

template<typename Stream>
inline void Unserialize(Stream &s, int32_t &a) { a = ser_readdata32(s); };

template<typename Stream>
inline void Unserialize(Stream &s, uint32_t &a) { a = ser_readdata32(s); };

template<typename Stream>
inline void Unserialize(Stream &s, int64_t &a) { a = ser_readdata64(s); };

template<typename Stream>
inline void Unserialize(Stream &s, uint64_t &a) { a = ser_readdata64(s); };

template<typename Stream>
inline void Unserialize(Stream &s, float &a) { a = ser_uint32_to_float(ser_readdata32(s)); };

template<typename Stream>
inline void Unserialize(Stream &s, double &a) { a = ser_uint64_to_double(ser_readdata64(s)); };

// shared_ptr
template<typename Stream, typename T>
void Serialize(Stream &os, const std::shared_ptr<const T> &p) {
    Serialize(os, *p);
}

template<typename Stream, typename T>
void Unserialize(Stream &is, std::shared_ptr<const T> &p) {
    p = std::make_shared<const T>(deserialize, is);
}


// vector

/**
 * Compact Size
 * size <  253        -- 1 byte
 * size <= USHRT_MAX  -- 3 bytes  (253 + 2 bytes)
 * size <= UINT_MAX   -- 5 bytes  (254 + 4 bytes)
 * size >  UINT_MAX   -- 9 bytes  (255 + 8 bytes)
 */

template<typename Stream>
void WriteCompactSize(Stream &os, uint64_t nSize) {
    if (nSize < 253) {
        ser_writedata8(os, nSize);
    } else if (nSize <= std::numeric_limits<unsigned short>::max()) {
        ser_writedata8(os, 253);
        ser_writedata16(os, nSize);
    } else if (nSize <= std::numeric_limits<unsigned int>::max()) {
        ser_writedata8(os, 254);
        ser_writedata32(os, nSize);
    } else {
        ser_writedata8(os, 255);
        ser_writedata64(os, nSize);
    }
}

template<typename Stream>
uint64_t ReadCompactSize(Stream &is) {
    uint8_t chSize = ser_readdata8(is);
    uint64_t nSizeRet = 0;

    if (chSize < 253) {
        nSizeRet = chSize;
    } else if (chSize == 253) {
        nSizeRet = ser_readdata16(is);
        if (nSizeRet < 253)
            throw std::ios_base::failure("non-canonical ReadCompactSize()");
    } else if (chSize == 254) {
        nSizeRet = ser_readdata32(is);
        if (nSizeRet < 0x10000u)
            throw std::ios_base::failure("non-canonical ReadCompactSize()");
    } else {
        nSizeRet = ser_readdata64(is);
        if (nSizeRet < 0x100000000ULL)
            throw std::ios_base::failure("non-canonical ReadCompactSize()");
    }
    if (nSizeRet > (uint64_t) MAX_SIZE)
        throw std::ios_base::failure("ReadCompactSize(): size too large");
    return nSizeRet;
}

template<typename Stream, typename T, typename A>
void SerializeImpl(Stream &os, const std::vector<T, A> &v, const unsigned char &) {
    WriteCompactSize(os, v.size());
    if (!v.empty()) {
        os.write((char *) v.data(), v.size() * sizeof(T));
    }
}

template<typename Stream, typename T, typename A, typename V>
void SerializeImpl(Stream &os, const std::vector<T, A> &v, const V &) {
    WriteCompactSize(os, v.size());
    for (typename std::vector<T, A>::const_iterator vi = v.begin(); vi != v.end(); ++vi) {
        ::Serialize(os, (*vi));
    }
}

template<typename Stream, typename T, typename A>
inline void Serialize(Stream &os, const std::vector<T, A> &v) {
    SerializeImpl(os, v, T());
}

template<typename Stream, typename T, typename A>
void UnserializeImpl(Stream &is, std::vector<T, A> &v, const unsigned char &) {
    v.clear();

    unsigned int nSize = ReadCompactSize(is);
    unsigned int i = 0;

    // i: 已读的下标
    // blk: 要读的数量 (每次最大不超过5MB)
    // nSize: 总计要读的数量
    while (i < nSize) {
        unsigned int blk = std::min(nSize - i, (unsigned int) (1 + 4999999 / sizeof(T)));
        v.resize(i + blk);
        is.read((char *) &v[i], blk * sizeof(T));
        i += blk;
    }
}


template<typename Stream, typename T, typename A, typename V>
void UnserializeImpl(Stream &is, std::vector<T, A> &v, const V &) {
    v.clear();

    unsigned int nSize = ReadCompactSize(is);
    unsigned int i = 0;
    unsigned int nMid = 0;


    // nMid: 要读的到的下标 (每次最大不超过5MB)
    while (nMid < nSize) {

        nMid += 5000000 / sizeof(T);
        if (nMid > nSize)
            nMid = nSize;

        v.resize(nMid);
        for (; i < nMid; i++)
            Unserialize(is, v[i]);
    }
}

template<typename Stream, typename T, typename A>
inline void Unserialize(Stream &is, std::vector<T, A> &v) {
    UnserializeImpl(is, v, T());
}


// 模板匹配
template<typename Stream, typename T>
inline void Serialize(Stream &os, const T &a) {
    a.Serialize(os);
}

template<typename Stream, typename T>
inline void Unserialize(Stream &is, T &&a) {
    a.Unserialize(is);
}

// 编译期动作模板
struct CSerActionSerialize {
    constexpr bool ForRead() const { return false; }
};

struct CSerActionUnserialize {
    constexpr bool ForRead() const { return true; }
};

template<typename Stream>
void SerializeMany(Stream &s) {

}

template<typename Stream, typename Arg, typename... Args>
void SerializeMany(Stream &s, const Arg &arg, const Args &... args) {
    ::Serialize(s, arg);
    ::SerializeMany(s, args...);
}


template<typename Stream>
void UnserializeMany(Stream &s) {

}


template<typename Stream, typename Arg, typename ... Args>
void UnserializeMany(Stream &s, Arg &&arg, Args &&... args) {
    ::Unserialize(s, arg);
    ::UnserializeMany(s, args...);
}

template<typename Stream, typename... Args>
inline void SerReadWriteMany(Stream &s, CSerActionSerialize ser_action, const Args &... args) {
    ::SerializeMany(s, args...);
}


template<typename Stream, typename... Args>
inline void SerReadWriteMany(Stream &s, CSerActionUnserialize ser_action, Args &&... args) {
    ::UnserializeMany(s, args...);
}

#endif //CPPSERIALIZE_SERIALIZE_H
