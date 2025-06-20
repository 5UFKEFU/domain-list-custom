# GeoSite.dat 解析工具使用说明

这个工具可以直接解析 geosite.dat 二进制文件，查看其中包含的所有分类和域名。

## 文件说明

- `geosite.proto` - protobuf 定义文件
- `geosite_pb2.py` - 由 protobuf 生成的 Python 模块
- `parse_geosite.py` - 主要的解析脚本

## 安装依赖

```bash
# 安装 protobuf 编译器
brew install protobuf

# 安装 Python protobuf 库
pip install protobuf
```

## 使用方法

### 1. 查看所有分类

```bash
python3 parse_geosite.py ./publish/geosite.dat
```

这会显示：
- 文件基本信息（大小、分类总数）
- 所有分类列表（按域名数量排序）

### 2. 查看特定分类的所有域名

```bash
python3 parse_geosite.py ./publish/geosite.dat <分类名>
```

例如：
```bash
# 查看字节跳动域名
python3 parse_geosite.py ./publish/geosite.dat bytedance

# 查看伊朗广告域名
python3 parse_geosite.py ./publish/geosite.dat category-ads-ir

# 查看中国银行域名
python3 parse_geosite.py ./publish/geosite.dat category-bank-cn
```

### 3. 搜索包含关键词的域名

```bash
python3 parse_geosite.py ./publish/geosite.dat --search <关键词>
```

例如：
```bash
# 搜索包含 "tiktok" 的域名
python3 parse_geosite.py ./publish/geosite.dat --search tiktok

# 搜索包含 "ads" 的域名
python3 parse_geosite.py ./publish/geosite.dat --search ads
```

## 主要分类说明

根据你的 geosite.dat 文件，主要包含以下分类：

### 中国相关分类
- `CATEGORY-BANK-CN` (181 个域名) - 中国银行
- `CATEGORY-SECURITIES-CN` (161 个域名) - 中国证券
- `CATEGORY-LOGISTICS-CN` (151 个域名) - 中国物流
- `CATEGORY-DEV-CN` (116 个域名) - 中国开发者
- `CATEGORY-EDUCATION-CN` (110 个域名) - 中国教育
- `CATEGORY-SCHOLAR-CN` (89 个域名) - 中国学术
- `CATEGORY-NETDISK-CN` (60 个域名) - 中国网盘
- `CATEGORY-AI-CN` (29 个域名) - 中国AI

### 伊朗相关分类
- `CATEGORY-NEWS-IR` (55 个域名) - 伊朗新闻
- `CATEGORY-BOURSE-IR` (48 个域名) - 伊朗交易所
- `CATEGORY-TECH-IR` (44 个域名) - 伊朗科技
- `CATEGORY-BANK-IR` (32 个域名) - 伊朗银行
- `CATEGORY-GOV-IR` (29 个域名) - 伊朗政府
- `CATEGORY-MEDIA-IR` (26 个域名) - 伊朗媒体
- `CATEGORY-ADS-IR` (7 个域名) - 伊朗广告

### 俄罗斯相关分类
- `CATEGORY-MEDIA-RU` (131 个域名) - 俄罗斯媒体
- `CATEGORY-GOV-RU` (119 个域名) - 俄罗斯政府

### 成人内容分类
- `CATEGORY-PORN` (135 个域名) - 成人内容

### 主要公司分类
- `THESCOREGROUP` (97 个域名) - The Score Group
- `TENCENT-GAMES` (86 个域名) - 腾讯游戏
- `ICBC` (58 个域名) - 工商银行
- `LENOVO` (57 个域名) - 联想
- `UCLOUD` (54 个域名) - UCloud
- `BILIBILI` (50 个域名) - B站
- `CTRIP` (44 个域名) - 携程
- `DEWU` (44 个域名) - 得物
- `BYTEDANCE` (6 个域名) - 字节跳动

## 在 v2ray 中使用

脚本会自动生成 v2ray 配置示例。例如：

```json
{
  "routing": {
    "rules": [
      {
        "type": "field",
        "domain": ["geosite:BYTEDANCE"],
        "outboundTag": "proxy"
      }
    ]
  }
}
```

## 注意事项

1. **分类名称区分大小写**：在 geosite.dat 中，分类名称通常是大写的
2. **域名类型**：脚本会显示域名的类型（Plain、Regex、RootDomain、Full）
3. **文件路径**：确保 geosite.dat 文件路径正确

## 常见问题

### Q: 为什么找不到某些分类？
A: 某些分类可能不存在于你的 geosite.dat 文件中，或者分类名称大小写不匹配。使用第一个命令查看所有可用分类。

### Q: 如何获取所有广告域名？
A: 目前你的 geosite.dat 中只有 `CATEGORY-ADS-IR` (7个域名)。如果需要更多广告域名，可能需要重新生成 geosite.dat 文件。

### Q: 如何搜索特定域名？
A: 使用 `--search` 参数可以搜索包含特定关键词的域名，会在所有分类中查找。 