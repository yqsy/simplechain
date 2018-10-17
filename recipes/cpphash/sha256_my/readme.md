<!-- TOC -->

- [1. cuckoo cycle](#1-cuckoo-cycle)
- [2. sha256](#2-sha256)

<!-- /TOC -->



<a id="markdown-1-cuckoo-cycle" name="1-cuckoo-cycle"></a>
# 1. cuckoo cycle
* https://github.com/bitcoin/bips/blob/master/bip-0154.mediawiki (bitcoin的提案)
* https://bc-2.jp/cuckoo-profile.pdf (评测)
* https://aeternity.com/aeternity-blockchain-whitepaper.pdf (欧洲以太坊)
* https://github.com/tromp/cuckoo (源码)

<a id="markdown-2-sha256" name="2-sha256"></a>
# 2. sha256 

* https://zh.wikipedia.org/wiki/SHA-2
* https://github.com/okdshin/PicoSHA2/blob/master/picosha2.h (源码)
* https://www.cnblogs.com/foxclever/p/8370712.html (虽然写的不是很好,但是入门学习下也好)

![](./pic/sha.png)

```bash
# sha256
消息摘要长度256位 --- 最终生成256位 -> 32字节 -> 64hex 串
消息长度小于2^64位 --- 被hash内容的位数 < 2^64
分组长度512位 --- 每个分组的大小为512位
计算字长度32位 --- 生成w0~w15,每个长度为32位
计算步骤数 --- 经过多少次的混淆?
```

