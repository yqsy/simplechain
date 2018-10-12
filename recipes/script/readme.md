<!-- TOC -->

- [1. 说明](#1-说明)
- [2. 一个认证](#2-一个认证)

<!-- /TOC -->


<a id="markdown-1-说明" name="1-说明"></a>
# 1. 说明

* http://book.8btc.com/books/6/masterbitcoin2cn/_book/ch06.html
* https://en.bitcoin.it/wiki/Script
* https://en.bitcoin.it/wiki/Transaction (P2PkH P2SH)



<a id="markdown-2-一个认证" name="2-一个认证"></a>
# 2. 一个认证

1. 压入签名
2. 压入公钥
3. 复制公钥(再次压入公钥)
4. 栈顶公钥变成公钥哈希
5. 压入(out的)公钥哈希
6. 判断该公钥是否是(上一个out)的公钥哈希对应的
7. 校验签名和公钥是配对的


提供:
* 加锁: pubkeyHash
* 解锁: 1. pubKey 2. sig

验证:  
1. 公钥和上一笔输出的公钥哈希对应
2. 公钥和签名是对应的

![](./pic/verify_script1.png)

![](./pic/verify_script2.png)

