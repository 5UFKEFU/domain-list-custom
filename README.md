# 简介

基于 [v2fly/domain-list-community#256](https://github.com/v2fly/domain-list-community/issues/256) 的提议，重构 [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) 的构建流程，并添加新功能。

## 与官方版 `dlc.dat` 不同之处

- 将 `dlc.dat` 重命名为 `geosite.dat`
- 去除 `cn` 列表里带有 `@ads`、`@!cn` 属性的规则
- 去除 `geolocation-cn` 列表里带有 `@ads`、`@!cn` 属性的规则
- 去除 `geolocation-!cn` 列表里带有 `@ads`、`@cn` 属性的规则，尽量避免在中国大陆有接入点的海外公司的域名走代理。例如，避免国区 Steam 游戏下载服务走代理。

## 下载地址

[https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geosite.dat](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geosite.dat)

## 使用本项目的项目

[@Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)

# Domain List Custom

基于 v2fly/domain-list-community 的极简 geosite.dat 生成器

## 功能特点

- **极简配置**: 只包含必要的域名分类，文件大小控制在合理范围内
- **递归include处理**: 自动递归展开所有分类文件的include关系，无需手动维护公司列表
- **多国家支持**: 支持中国、伊朗、俄罗斯、色情等分类的自动处理
- **完整覆盖**: 包含中国、伊朗、俄罗斯国家域名，TikTok，AI分类，所有广告域名

## 使用方法

### 1. 多国家递归include（推荐）

自动递归处理中国、伊朗、俄罗斯、色情域名：

```bash
go run main.go listinfo.go listinfomap.go trie.go common.go -recursive -mode=multi
```

**特点:**
- 自动发现399个被多国家分类文件include的实际域名文件
- 包含中国、伊朗、俄罗斯、色情等所有分类的域名
- 完全自动化，无需手动维护任何国家或分类列表
- 文件大小合理（86K）

### 2. 中国递归include

仅处理中国域名：

```bash
go run main.go listinfo.go listinfomap.go trie.go common.go -recursive -mode=cn
```

**特点:**
- 自动发现329个被中国分类文件include的实际域名文件
- 包含所有中国主流互联网公司、银行、媒体、教育、游戏等域名
- 完全自动化，无需手动维护公司列表
- 文件大小适中（63K）

### 3. 手动指定文件（传统方式）

```bash
go run main.go listinfo.go listinfomap.go trie.go common.go
```

**当前配置包含:**
- 主要互联网公司: google, apple, meta, microsoft, amazon
- 中国主要互联网公司: tiktok, baidu, alibaba, tencent, bytedance
- 中国手机厂商: xiaomi, huawei, oppo, vivo
- 中国其他重要公司: meituan, didi, jd, netease, sina, sohu, iqiyi, youku, bilibili
- 中国顶级域名: tld-cn
- 中国主要银行: boc, ccb, citic, cmb, icbc, unionpay
- 中国媒体: 36kr, cctv, chinanews, dgtle, geekpark, ifanr, jiemian, landian, phoenix, ruanmei
- AI服务: category-ai-cn, category-ai-!cn
- 广告域名: category-ads-all + 61个广告分类

## 输出文件

- `geosite.dat`: 主要的geosite文件
- `*.txt`: 各个分类的明文域名列表
- `gfwlist.txt`: GFWList格式的域名列表

## 参数说明

- `-recursive`: 启用递归include处理
- `-mode`: 递归模式选择
  - `cn`: 仅处理中国域名（默认）
  - `multi`: 处理多国家域名（中国、伊朗、俄罗斯、色情）
- `-exportlists`: 指定要导出的列表（用逗号分隔）
- `-excludeattrs`: 指定要排除的属性
- `-togfwlist`: 指定要转换为GFWList格式的列表
- `-datapath`: 数据目录路径（默认: ./data）
- `-outputpath`: 输出目录路径（默认: ./publish）
- `-datname`: 生成的dat文件名（默认: geosite.dat）

## 递归include处理原理

1. 从分类文件开始（中国、伊朗、俄罗斯、色情等）
2. 递归解析所有 `include:xxx` 语句
3. 收集所有被引用的实际域名文件
4. 自动处理所有发现的文件，生成完整的域名列表

## 多国家模式包含的分类

### 中国分类（30个）
- 主分类: cn, geolocation-cn, tld-cn
- 功能分类: category-ai-cn, category-bank-cn, category-media-cn, category-social-media-cn 等

### 伊朗分类（15个）
- 功能分类: category-ads-ir, category-bank-ir, category-media-ir, category-news-ir 等

### 俄罗斯分类（2个）
- 功能分类: category-gov-ru, category-media-ru

### 色情分类（1个）
- 功能分类: category-porn

## 优势

- **自动化**: 无需手动维护公司列表，自动发现所有被分类文件引用的域名
- **完整性**: 包含所有v2fly/domain-list-community项目中分类的域名
- **可维护性**: 当社区更新分类文件时，自动包含新的域名
- **效率**: 文件大小合理，加载速度快
- **灵活性**: 支持多种处理模式，满足不同需求

## 示例输出

### 多国家模式
```
启用多国家递归include处理，自动收集中国、伊朗、俄罗斯、色情域名...
递归发现 399 个被include的文件
最终将处理 447 个文件（多国家（中国、伊朗、俄罗斯、色情）模式）
geosite.dat has been generated successfully in './publish'.
```

### 中国模式
```
启用中国递归include处理，自动收集所有中国域名...
递归发现 329 个被include的文件
最终将处理 359 个文件（中国模式）
geosite.dat has been generated successfully in './publish'.
```

## 注意事项

- 递归include处理会包含所有被分类文件引用的域名，确保覆盖完整
- 多国家模式会显著增加文件大小，但提供更全面的域名覆盖
- 如果只需要特定分类，可以使用手动指定文件的方式
- 生成的文件大小会根据包含的域名数量自动调整

## 实际运行结果示例

### geosite.dat（部分内容，二进制格式，建议用解析工具查看）
```
ANTHROPIC
servd-anthropic-website.b-cdn.net
anthropic.com
claude.ai
claudeusercontent.com
...
IQIYI
71.am
iqiyi.com
iqiyipic.com
pps.tv
...
```

### category-ads-all.txt（前若干行）
```
full:pinggai0.caixin.com
full:pinggai1.caixin.com
full:pinggai2.caixin.com
full:pinggai3.caixin.com
full:pinggai4.caixin.com
full:pinggai5.caixin.com
full:pinggai6.caixin.com
full:pinggai7.caixin.com
full:pinggai8.caixin.com
full:pinggai9.caixin.com
domain:moat.com
domain:moatads.com
domain:aaxads.com
...
```

### category-ai-!cn.txt（前若干行）
```
full:openaiapi-site.azureedge.net
full:openaicom-api-bdcpf8c6d2e9atf6.z01.azurefd.net
full:openaicomproductionae4b.blob.core.windows.net
full:production-openaicom-storage.azureedge.net
full:pplx-res.cloudinary.com
full:servd-anthropic-website.b-cdn.net
full:ai.google.dev
full:alkalicore-pa.clients6.google.com
full:alkalimakersuite-pa.clients6.google.com
full:webchannel-alkalimakersuite-pa.clients6.google.com
domain:cursor.com
domain:oaiusercontent.com
domain:jasper.ai
...
```


jeffcheng@JEFFs-MacBook-Air-M2-15 domain-list-custom % python3 parse_geosite.py ./publish/geosite.dat
=== GeoSite.dat 文件信息 ===
文件路径: ./publish/geosite.dat
文件大小: 571370 字节 (557.98 KB)
分类总数: 110

=== 所有分类列表 ===
  1. CATEGORY-PORN                  (12013 个域名)
  2. GEOLOCATION-CN                 (4815 个域名)
  3. APPLE                          (1599 个域名)
  4. GOOGLE                         (1176 个域名)
  5. CATEGORY-ADS-ALL               ( 824 个域名)
  6. CATEGORY-ADS                   ( 801 个域名)
  7. BYTEDANCE                      ( 797 个域名)
  8. HUAWEI                         ( 594 个域名)
  9. MICROSOFT                      ( 581 个域名)
 10. HUAWEICLOUD                    ( 544 个域名)
 11. JD                             ( 466 个域名)
 12. ALIBABA                        ( 412 个域名)
 13. BAIDU                          ( 302 个域名)
 14. CATEGORY-RU                    ( 288 个域名)
 15. CATEGORY-GOV-RU                ( 238 个域名)
 16. AMAZON                         ( 176 个域名)
 17. TENCENT                        ( 165 个域名)
 18. BAIDU-ADS                      ( 156 个域名)
 19. SINA                           ( 122 个域名)
 20. CATEGORY-AI-!CN                ( 118 个域名)
 21. ALIBABA-ADS                    (  98 个域名)
 22. XIAOMI                         (  95 个域名)
 23. YANDEX                         (  86 个域名)
 24. BILIBILI                       (  83 个域名)
 25. SINA-ADS                       (  71 个域名)
 26. OPPO                           (  68 个域名)
 27. VK                             (  68 个域名)
 28. NETEASE                        (  66 个域名)
 29. CATEGORY-AI-CN                 (  54 个域名)
 30. CATEGORY-IR                    (  49 个域名)
 31. GOOGLE-ADS                     (  45 个域名)
 32. IQIYI                          (  36 个域名)
 33. TENCENT-ADS                    (  36 个域名)
 34. GOOGLE-DEEPMIND                (  34 个域名)
 35. SOHU                           (  31 个域名)
 36. OPENAI                         (  28 个域名)
 37. TIKTOK                         (  21 个域名)
 38. SOHU-ADS                       (  21 个域名)
 39. BYTEDANCE-ADS                  (  20 个域名)
 40. DIDI                           (  18 个域名)
 41. IQIYI-ADS                      (  17 个域名)
 42. RUANMEI                        (  16 个域名)
 43. XIAOMI-ADS                     (  15 个域名)
 44. MEITUAN                        (  14 个域名)
 45. CATEGORY-ADS-IR                (  14 个域名)
 46. LETV-ADS                       (  13 个域名)
 47. TELEVISION-ADS                 (  13 个域名)
 48. NETEASE-ADS                    (  12 个域名)
 49. YOUKU                          (  12 个域名)
 50. DMM-ADS                        (  11 个域名)
 51. ALIBABACLOUD                   (  10 个域名)
 52. VIVO                           (   9 个域名)
 53. JD-ADS                         (   8 个域名)
 54. CURSOR                         (   8 个域名)
 55. QIHOO360-ADS                   (   8 个域名)
 56. KUAISHOU-ADS                   (   7 个域名)
 57. ANTHROPIC                      (   7 个域名)
 58. MAILRU                         (   6 个域名)
 59. OZON                           (   6 个域名)
 60. UNITY-ADS                      (   6 个域名)
 61. HUNANTV-ADS                    (   6 个域名)
 62. OPENX-ADS                      (   5 个域名)
 63. PERPLEXITY                     (   5 个域名)
 64. POE                            (   4 个域名)
 65. HUGGINGFACE                    (   4 个域名)
 66. SUPERSONIC-ADS                 (   4 个域名)
 67. ADJUST-ADS                     (   4 个域名)
 68. META                           (   4 个域名)
 69. APPLE-ADS                      (   4 个域名)
 70. UMENG-ADS                      (   4 个域名)
 71. XAI                            (   4 个域名)
 72. XIMALAYA-ADS                   (   3 个域名)
 73. KUGOU-ADS                      (   3 个域名)
 74. SPOTIFY-ADS                    (   3 个域名)
 75. OGURY-ADS                      (   3 个域名)
 76. YAHOO-ADS                      (   3 个域名)
 77. ATOM-DATA-ADS                  (   3 个域名)
 78. DUOLINGO-ADS                   (   3 个域名)
 79. APPLOVIN-ADS                   (   2 个域名)
 80. HIIDO-ADS                      (   2 个域名)
 81. GROQ                           (   2 个域名)
 82. OK                             (   2 个域名)
 83. POCOIQ-ADS                     (   2 个域名)
 84. TAGTIC-ADS                     (   2 个域名)
 85. TABOOLA                        (   2 个域名)
 86. ADOBE-ADS                      (   2 个域名)
 87. AMAZON-ADS                     (   2 个域名)
 88. ZYNGA-ADS                      (   2 个域名)
 89. PUBMATIC-ADS                   (   1 个域名)
 90. EMOGI-ADS                      (   1 个域名)
 91. NEWRELIC-ADS                   (   1 个域名)
 92. UBERADS-ADS                    (   1 个域名)
 93. ADCOLONY-ADS                   (   1 个域名)
 94. MIXPANEL-ADS                   (   1 个域名)
 95. ONESIGNAL-ADS                  (   1 个域名)
 96. MXPLAYER-ADS                   (   1 个域名)
 97. OOKLA-SPEEDTEST-ADS            (   1 个域名)
 98. INNER-ACTIVE-ADS               (   1 个域名)
 99. SENSORSDATA-ADS                (   1 个域名)
100. XHAMSTER-ADS                   (   1 个域名)
101. ACFUN-ADS                      (   1 个域名)
102. CLEARBITJS-ADS                 (   1 个域名)
103. GROWINGIO-ADS                  (   1 个域名)
104. FLURRY-ADS                     (   1 个域名)
105. HOTJAR-ADS                     (   1 个域名)
106. LEANPLUM-ADS                   (   1 个域名)
107. MOPUB-ADS                      (   1 个域名)
108. TAPPX-ADS                      (   1 个域名)
109. SEGMENT-ADS                    (   1 个域名)
110. WTEAM-ADS                      (   1 个域名)

使用方法:
  python parse_geosite.py ./publish/geosite.dat <分类名>  # 查看指定分类的所有域名
  示例: python parse_geosite.py ./publish/geosite.dat tiktok


  go run . -mode multi 
  