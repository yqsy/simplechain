//
// Created by yq on 18-10-11.
//

#include <transaction.h>


CTransaction::CTransaction(const CMutableTransaction &tx) : nVersion(tx.nVersion), nLockTime(tx.nLockTime) {

}



CTransaction::CTransaction(CMutableTransaction &&tx): nVersion(tx.nVersion), nLockTime(tx.nLockTime) {

}

