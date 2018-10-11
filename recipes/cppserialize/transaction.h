//
// Created by yq on 18-10-11.
//

#ifndef CPPSERIALIZE_TRANSACTION_H
#define CPPSERIALIZE_TRANSACTION_H

#include <memory>
#include <stdint.h>

#include <serialize.h>

/**
 * Basic transaction serialization format:
 * - int32_t nVersion
 * - std::vector<CTxIn> vin
 * - std::vector<CTxOut> vout
 * - uint32_t nLockTime
 *
 * Extended transaction serialization format:
 * - int32_t nVersion
 * - unsigned char dummy = 0x00
 * - unsigned char flags (!= 0)
 * - std::vector<CTxIn> vin
 * - std::vector<CTxOut> vout
 * - if (flags & 1):
 *   - CTxWitness wit;
 * - uint32_t nLockTime
 */


template<typename Stream, typename TxType>
inline void UnserializeTransaction(TxType &tx, Stream &s) {
    s >> tx.nVersion;
    s >> tx.nLockTime;
}

template<typename Stream, typename TxType>
inline void SerializeTransaction(const TxType &tx, Stream &s) {
    // 1. 基础信息输出方式
    // 2. 隔离见证输出方式

    s << tx.nVersion;
    s << tx.nLockTime;
}


struct CMutableTransaction;

class CTransaction {
public:

    static const int32_t CURRENT_VERSION = 2;

    static const int32_t MAX_STANDARD_VERSION = 2;


    // 基础序列化:
    //    const std::vector<CTxIn> vin;
    //    const std::vector<CTxOut> vout;
    const int32_t nVersion;
    const uint32_t nLockTime;

    // 只在内存中有:
    //const uint256 hash;
    //const uint256 m_witness_hash;

    // 其实这个没有用
    CTransaction() : nVersion(CTransaction::CURRENT_VERSION), nLockTime(0) {
        // 日常初始化
    }


    CTransaction(const CMutableTransaction &tx);

    CTransaction(CMutableTransaction &&tx);

    // 由于成员变量都是const,所以只能serialize

    template<typename Stream>
    inline void Serialize(Stream& s) const {
        SerializeTransaction(*this, s);
    }

    template<typename Stream>
    CTransaction(deserialize_type, Stream &s) : CTransaction(CMutableTransaction(deserialize, s)) {

    }

};

struct CMutableTransaction {
public:
    //    std::vector<CTxIn> vin;
    //    std::vector<CTxOut> vout;
    int32_t nVersion;
    uint32_t nLockTime;

    CMutableTransaction() : nVersion(CTransaction::CURRENT_VERSION), nLockTime(1) {

    }

    // 从CTransaction也可转过来
    explicit CMutableTransaction(const CTransaction &tx) : nVersion(tx.nVersion), nLockTime(tx.nLockTime) {

    }

    template<typename Stream>
    inline void Serialize(Stream &s) const {
        SerializeTransaction(*this, s);
    }

    template<typename Stream>
    inline void Unserialize(Stream &s) {
        UnserializeTransaction(*this, s);
    }

    template<typename Stream>
    CMutableTransaction(deserialize_type, Stream &s) {
        Unserialize(s);
    }
};


typedef std::shared_ptr<const CTransaction> CTransactionRef;

#endif //CPPSERIALIZE_TRANSACTION_H
