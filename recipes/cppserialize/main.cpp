#include <iostream>
#include <stdint.h>

#include <vector>

#include <serialize.h>
#include <streams.h>
#include <transaction.h>

class CBlockHeader {
public:
    int32_t nVersion;
    //uint256 hashPrevBlock;
    //uint256 hashMerkleRoot;
    uint32_t nTime;
    uint32_t nBits;
    uint32_t nNonce;

    ADD_SERIALIZE_METHODS;

    template<typename Stream, typename Operation>
    inline void SerializationOp(Stream &s, Operation ser_action) {
        READWRITE(nVersion);
        READWRITE(nTime);
        READWRITE(nBits);
        READWRITE(nNonce);
    }
};

void SimpleTest() {
    int32_t a, b;
    a = 1;
    CDataStream s;
    s << 1;
    s >> b;
    if (b != 1) {
        throw std::ios_base::failure("b do not equal 1");
    }
}

void StructTest() {
    CBlockHeader a, b;
    a.nVersion = 1;
    a.nTime = 2;
    a.nBits = 3;
    a.nNonce = 4;

    CDataStream s;
    s << a;
    s >> b;

    if (b.nVersion != 1 || b.nTime != 2 || b.nBits != 3 || b.nNonce != 4) {
        throw std::ios_base::failure("struct test error");
    }
}

void VectorTest() {
    std::vector<int> v{1, 2, 3};

    std::vector<int> v2;

    CDataStream s;

    s << v;
    s >> v2;

    for (auto &ele: v2) {
        std::cout << ele << std::endl;
    }
}

void TransactionTest() {
    std::vector<CTransactionRef> vtx;

    std::vector<CTransactionRef> vtx2;

    vtx.push_back(std::make_shared<CTransaction>(CMutableTransaction()));

    CDataStream s;

    s << vtx;

    s >> vtx2;

    std::cout << vtx2[0]->nVersion << std::endl;
}


int main() {
    SimpleTest();
    StructTest();
    VectorTest();
    TransactionTest();
    return 0;
}
