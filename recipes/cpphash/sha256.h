//
// Created by yq on 18-10-16.
//

#ifndef CPPHASH_SHA256_H
#define CPPHASH_SHA256_H

#include <algorithm>
#include <iterator>
#include <vector>
#include <string>
#include <sstream>
#include <fstream>

#include <stdint.h>
#include <endian.h>
#include <cassert>

#define PICOSHA2_BUFFER_SIZE_FOR_INPUT_ITERATOR 1048576


// hash到的目的容器的大小(字节) 32 * 8 = 256位
static const size_t k_digest_size = 32;

// 1. 初始散列值,取值为前8个质数的平方根的小数部分的前32位
// h0 ~ h7
const uint32_t initial_message_digest[8] = {0x6a09e667, 0xbb67ae85, 0x3c6ef372,
                                            0xa54ff53a, 0x510e527f, 0x9b05688c,
                                            0x1f83d9ab, 0x5be0cd19};

// 2. 64个32位的常数序列,取值为前64个质数的立方根的小数部分的前32位
// K[0..63]
const uint32_t add_constant[64] = {
        0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1,
        0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
        0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786,
        0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
        0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147,
        0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
        0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b,
        0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
        0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a,
        0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
        0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2};


inline uint8_t mask_8bit(uint8_t x) { return x & 0xff; }

inline uint32_t mask_32bit(uint32_t x) { return x & 0xffffffff; }

inline uint32_t ch(uint32_t x, uint32_t y, uint32_t z) { return (x & y) ^ ((~x) & z); }

inline uint32_t maj(uint32_t x, uint32_t y, uint32_t z) {
    return (x & y) ^ (x & z) ^ (y & z);
}

inline uint32_t rotr(uint32_t x, std::size_t n) {
    assert(n < 32);
    return mask_32bit((x >> n) | (x << (32 - n)));
}

inline uint32_t bsig0(uint32_t x) { return rotr(x, 2) ^ rotr(x, 13) ^ rotr(x, 22); }

inline uint32_t bsig1(uint32_t x) { return rotr(x, 6) ^ rotr(x, 11) ^ rotr(x, 25); }

inline uint32_t shr(uint32_t x, std::size_t n) {
    assert(n < 32);
    return x >> n;
}

inline uint32_t ssig0(uint32_t x) { return rotr(x, 7) ^ rotr(x, 18) ^ shr(x, 3); }

inline uint32_t ssig1(uint32_t x) { return rotr(x, 17) ^ rotr(x, 19) ^ shr(x, 10); }

template<typename RaIter1, typename RaIter2>
void hash256_block(RaIter1 message_digest, RaIter2 first, RaIter2 last) {
    assert(first + 64 == last);

    // 4. 512bit(64byte)的分组分成16个部分(每组4byte,大端法)生成w[0..16)
    // TODO: 这里我不太明白为什么要变成大端法? 因为本质原生的字节序也没有指明? 默认小?
    // Raiter2: std::vector<uint8_t>::iterator
    uint32_t w[64] = {};
    for (int i = 0; i < 16; ++i) {
        w[i] = (static_cast<uint32_t >(mask_8bit(*(first + i * 4))) << 24) |
               (static_cast<uint32_t >(mask_8bit(*(first + i * 4 + 1))) << 16) |
               (static_cast<uint32_t >(mask_8bit(*(first + i * 4 + 2))) << 8) |
               (static_cast<uint32_t >(mask_8bit(*(first + i * 4 + 3))));
    }

    // 5. 补充剩下w[16..64)
    // SSIG1(W(t-2)) + W(t-7) + SSIG0(w(t-15)) + W(t-16)
    for (int i = 16; i < 64; ++i) {
        w[i] = mask_32bit(ssig1(w[i - 2]) + w[i - 7] + ssig0(w[i - 15]) + w[i - 16]);
    }


    // 6. 初始化8个变量
    uint32_t a = *message_digest;
    uint32_t b = *(message_digest + 1);
    uint32_t c = *(message_digest + 2);
    uint32_t d = *(message_digest + 3);
    uint32_t e = *(message_digest + 4);
    uint32_t f = *(message_digest + 5);
    uint32_t g = *(message_digest + 6);
    uint32_t h = *(message_digest + 7);


    // 7. main loop, 64次计算
    for (std::size_t i = 0; i < 64; ++i) {
        uint32_t temp1 = h + bsig1(e) + ch(e, f, g) + add_constant[i] + w[i];
        uint32_t temp2 = bsig0(a) + maj(a, b, c);
        h = g;
        g = f;
        f = e;
        e = mask_32bit(d + temp1);
        d = c;
        c = b;
        b = a;
        a = mask_32bit(temp1 + temp2);
    }

    // 8. 更新h0 ~ h7
    *message_digest += a;
    *(message_digest + 1) += b;
    *(message_digest + 2) += c;
    *(message_digest + 3) += d;
    *(message_digest + 4) += e;
    *(message_digest + 5) += f;
    *(message_digest + 6) += g;
    *(message_digest + 7) += h;

    for (int i = 0; i < 8; ++i) {
        *(message_digest + i) = mask_32bit(*(message_digest + i));
    }

}


class hash256_one_by_one {
private:

    std::vector<uint8_t> buffer_;

    // 8个32位散列值, 也是最终结果
    uint32_t h_[8];

    // big-endian 数据长度
    uint64_t data_length_;
public:

    hash256_one_by_one() {
        buffer_.clear();
        std::copy(initial_message_digest, initial_message_digest + 8, h_);
    }

    template<typename RaIter>
    void process(RaIter first, RaIter last) {

        // 长度的大端字节序
        data_length_ = htobe64(std::distance(first, last));
        std::copy(first, last, std::back_inserter(buffer_));

        // 3. A. 填充1个'1'  B. 填充k bits '0'到达满 448 bit ([0,56) byte)C. 填充余下64bit为大端长度 ([56,64))

        buffer_.push_back(uint8_t('1'));

        // 56 byte = 448 bit, 64byte = 512 bit
        size_t remain = buffer_.size() % 64;

        if (remain < 56) {
            for (int i = 0; i < 56 - remain; ++i) {
                buffer_.push_back(uint8_t('0'));
            }
        } else if (remain > 56) {
            for (int i = 0; i < 64 - remain; ++i) {
                buffer_.push_back(uint8_t('0'));
            }
            for (int i = 0; i < 56 - remain; ++i) {
                buffer_.push_back(uint8_t('0'));
            }
        }

        assert(buffer_.size() % 56 == 0);

        uint8_t *pr = (uint8_t *) (&data_length_);
        for (int i = 0; i < 8; ++i) {
            buffer_.push_back(pr[i]);
        }

        assert(buffer_.size() % 64 == 0);

        // 根据分组长度512分割成一个个分组来进行处理
        for (int i = 0; i + 64 <= buffer_.size(); i += 64) {
            hash256_block(h_, buffer_.begin() + i, buffer_.begin() + i + 64);
        }

    }

    template<typename OutIter>
    void get_hash_bytes(OutIter first, OutIter last) const {
        for (const uint32_t *iter = h_; iter != h_ + 8; ++iter) {

            for (int i = 0; i < 4 && first != last; ++i) {
                *(first++) = mask_8bit(
                        static_cast<uint8_t >((*iter >> (24 - 8 * i))));
            }
        }
    }

};

template<typename InIter>
void output_hex(InIter first, InIter last, std::ostream &os) {

    // 16进制输出
    os.setf(std::ios::hex, std::ios::basefield);

    while (first != last) {
        // 右对齐
        os.width(2);

        // 填充'0' ?
        os.fill('0');

        // 4字节 -> 8hex
        os << static_cast<uint32_t>(*first);
        ++first;
    }
    // 十进制
    os.setf(std::ios::dec, std::ios::basefield);
}

// in: iter out: string&
template<typename InIter>
void bytes_to_hex_string(InIter first, InIter last, std::string &hex_str) {
    std::ostringstream oss;
    output_hex(first, last, oss);
    hex_str.assign(oss.str());
}

// in: container out: string&
template<typename InContainer>
void bytes_to_hex_string(const InContainer &bytes, std::string &hex_str) {
    bytes_to_hex_string(bytes.begin(), bytes.end(), hex_str);
}

// in: iter out: return string
template<typename InIter>
std::string bytes_to_hex_string(InIter first, InIter last) {
    std::string hex_str;
    bytes_to_hex_string(first, last, hex_str);
    return hex_str;
}

// in: container out: return string
template<typename InContainer>
std::string bytes_to_hex_string(const InContainer &bytes) {
    std::string hex_str;
    bytes_to_hex_string(bytes, hex_str);
    return hex_str;
}

template<typename RaIter, typename OutIter>
void hash256_impl(RaIter first, RaIter last, OutIter first2, OutIter last2, int,
                  std::random_access_iterator_tag) {
    hash256_one_by_one hasher;

    hasher.process(first, last);
    hasher.get_hash_bytes(first2, last2);
}


template<typename InputIter, typename OutIter>
void hash256_impl(InputIter first, InputIter last, OutIter first2, OutIter last2, int buffer_size,
                  std::input_iterator_tag) {

}

// in: iter, out: iter, buffer_size: 1MB
template<typename InIter, typename OutIter>
void hash256(InIter first, InIter last, OutIter first2, OutIter last2,
             int buffer_size = PICOSHA2_BUFFER_SIZE_FOR_INPUT_ITERATOR) {

    hash256_impl(
            first, last, first2, last2, buffer_size,
            typename std::iterator_traits<InIter>::iterator_category()
    );
}

// in: iter, out: cotainer
template<typename InIter, typename OutContainer>
void hash256(InIter first, InIter last, OutContainer &dst) {
    hash256(first, last, dst.begin(), dst.end());
}

// in: cotainer, out: iter
template<typename InContainer, typename OutIter>
void hash256(const InContainer &src, OutIter first, OutIter last) {
    hash256(src.begin(), src.end(), first, last);
}

// in: cotainer, out: cotainer
template<typename InContainer, typename OutContainer>
void hash256(const InContainer &src, OutContainer &dst) {
    hash256(src.begin(), src.end(), dst.begin(), dst.end());
}


// in: InIter out: string&
template<typename InIter>
void hash256_hex_string(InIter first, InIter last, std::string &hex_str) {
    // 1. hash到hashed数组
    uint8_t hashed[k_digest_size];
    hash256(first, last, hashed, hashed + k_digest_size);

    // 2. hashed数组转换成hex (1byte = 2hex)
    std::ostringstream oss;
    output_hex(hashed, hashed + k_digest_size, oss);
    hex_str.assign(oss.str());
}

// in: InIter out: return string
template<typename InIter>
std::string hash256_hex_string(InIter first, InIter last) {
    std::string hex_str;
    hash256_hex_string(first, last, hex_str);
    return hex_str;
}

// in: string out: string&
inline void hash256_hex_string(const std::string &src, std::string &hex_str) {
    hash256_hex_string(src.begin(), src.end(), hex_str);
}

// in: InContainer out: string&
template<typename InContainer>
void hash256_hex_string(const InContainer &src, std::string &hex_str) {
    hash256_hex_string(src.begin(), src.end(), hex_str);
}

// in: InContainer out: return string
template<typename InContainer>
std::string hash256_hex_string(const InContainer &src) {
    return hash256_hex_string(src.begin(), src.end());
}

template<typename OutIter>
void hash256(std::ifstream &f, OutIter first, OutIter last) {
    hash256(std::istreambuf_iterator<char>(f), std::istreambuf_iterator<char>(), first, last);
}

#endif //CPPHASH_SHA256_H
