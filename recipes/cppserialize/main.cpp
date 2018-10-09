#include <iostream>
#include <stdint.h>


#include <serialize.h>
#include <streams.h>

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

int main() {
    SimpleTest();
    StructTest();
    return 0;
}
