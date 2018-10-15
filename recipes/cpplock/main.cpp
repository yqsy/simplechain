#include <iostream>


#include "sync.h"


int main() {


    CCriticalSection cs;


    LOCK(cs);

    return 0;
}


