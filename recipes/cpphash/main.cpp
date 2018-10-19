#include <iostream>

#include <string>
#include <algorithm>
#include <stdint.h>

#include "sha256.h"

int main() {

//    std::string tmp("1");
//    auto hashed = hash256_hex_string(tmp);
//    std::cout << hashed << std::endl;

    std::string tmp = "11111111111111111111111111111111111111111111111111111111111111";
    auto  hashed = hash256_hex_string(tmp);
    std::cout << hashed << std::endl;


//    tmp = "1111111111111111111111111111111111111111111111111111111111111111";
//    hashed = hash256_hex_string(tmp);
//    std::cout << hashed << std::endl;



    return 0;
}
