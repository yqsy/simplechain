<!-- TOC -->

- [1. 相关知识](#1-相关知识)
- [2. 图片](#2-图片)
- [3. 实现列表](#3-实现列表)

<!-- /TOC -->

<a id="markdown-1-相关知识" name="1-相关知识"></a>
# 1. 相关知识

* http://www.youdzone.com/signature.html (签名+认证)

<a id="markdown-2-图片" name="2-图片"></a>
# 2. 图片

流程图1  
![](./pic/address-generation-scheme.png)

流程图1-补充  
![](./pic/address-generation-extra.png)

所有的钱包版本  
![](./pic/all_base58_version.png)

私钥,公钥,公钥哈希,钱包地址的关系  
![](./pic/relation.png)

签名交易流程  
![](./pic/sign_workflow.png)

签名验证  
![](./pic/sign_verify.png)

一笔转账  
![](./pic/transfer_confirm.png)

* privateKeyBytes: 64 hex (32 byte)
* publicKeyBytes: 128 hex (64 byte)
* publicKeyHash(SHA256+RIPEMD160): 40 hex (20 byte)
* walletAddress(Base58Encode): 34 hex (17 byte)
---
* signature: 128 hex (64 byte)
---
* txId(SHA256): 64 hex (32 byte)

<a id="markdown-3-实现列表" name="3-实现列表"></a>
# 3. 实现列表
* RANDOM 生成私钥
* SECP256K1 生成对应的公钥
* SHA256 + RIPEMD160 生成公钥哈希
* SHA256 + SHA256 生成公钥哈希的Checksum
* Base58Encode 版本号+公钥哈希+ChecksumCut 生成钱包地址
* 钱包地址提取公钥哈希
* 签名一笔交易
* 验证一笔交易
 