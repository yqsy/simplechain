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

int main() {
    SimpleTest();

//    CBlockHeader b1, b2;
//
//    b1.nVersion = 1;
//    b1.nTime = 2;
//    b1.nBits = 3;
//    b1.nNonce = 4;

//    CDataStream s;
//    s << b1;
//
//    s >> b2;



    return 0;
}
