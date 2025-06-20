#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
GeoSite.dat 文件解析工具
用于查看 geosite.dat 文件中的所有分类和域名
"""

import sys
import os
import argparse
from typing import List, Dict, Any

# 导入生成的 protobuf 模块
import geosite_pb2

def load_geosite_dat(dat_path: str) -> geosite_pb2.GeoSiteList:
    """加载 geosite.dat 文件"""
    try:
        with open(dat_path, 'rb') as f:
            data = f.read()
        
        geo = geosite_pb2.GeoSiteList()
        geo.ParseFromString(data)
        return geo
    except FileNotFoundError:
        print(f"错误: 文件 {dat_path} 不存在")
        sys.exit(1)
    except Exception as e:
        print(f"错误: 解析文件失败 - {e}")
        sys.exit(1)

def list_all_categories(dat_path: str):
    """列出所有分类"""
    geo = load_geosite_dat(dat_path)
    
    print(f"=== GeoSite.dat 文件信息 ===")
    print(f"文件路径: {dat_path}")
    print(f"文件大小: {os.path.getsize(dat_path)} 字节 ({os.path.getsize(dat_path)/1024:.2f} KB)")
    print(f"分类总数: {len(geo.entry)}")
    print()
    
    print("=== 所有分类列表 ===")
    categories = []
    for entry in geo.entry:
        categories.append((entry.country_code, len(entry.domain)))
    
    # 按域名数量排序
    categories.sort(key=lambda x: x[1], reverse=True)
    
    for i, (category, domain_count) in enumerate(categories, 1):
        print(f"{i:3d}. {category:<30} ({domain_count:>4d} 个域名)")
    
    print()
    print("使用方法:")
    print(f"  python {sys.argv[0]} {dat_path} <分类名>  # 查看指定分类的所有域名")
    print(f"  示例: python {sys.argv[0]} {dat_path} tiktok")

def list_category_domains(dat_path: str, category_name: str):
    """列出指定分类的所有域名"""
    geo = load_geosite_dat(dat_path)
    
    # 查找指定分类
    target_entry = None
    for entry in geo.entry:
        if entry.country_code.lower() == category_name.lower():
            target_entry = entry
            break
    
    if not target_entry:
        print(f"错误: 未找到分类 '{category_name}'")
        print("可用的分类:")
        for entry in geo.entry:
            print(f"  - {entry.country_code}")
        sys.exit(1)
    
    print(f"=== {target_entry.country_code} 域名列表 ===")
    print(f"文件路径: {dat_path}")
    print(f"域名数量: {len(target_entry.domain)}")
    print()
    
    # 按域名类型分组
    domain_types = {
        0: "Plain",      # 普通域名
        1: "Regex",      # 正则表达式
        2: "RootDomain", # 根域名
        3: "Full"        # 完整域名
    }
    
    domains_by_type = {}
    for domain in target_entry.domain:
        domain_type = domain_types.get(domain.type, "Unknown")
        if domain_type not in domains_by_type:
            domains_by_type[domain_type] = []
        domains_by_type[domain_type].append(domain.value)
    
    # 显示域名列表
    for domain_type, domains in domains_by_type.items():
        print(f"=== {domain_type} 类型域名 ({len(domains)} 个) ===")
        for i, domain in enumerate(domains, 1):
            print(f"{i:4d}. {domain}")
        print()
    
    # 显示 v2ray 配置示例
    print("=== 在 v2ray 中使用 ===")
    print("在 v2ray 配置文件中使用:")
    print("{")
    print("  \"routing\": {")
    print("    \"rules\": [")
    print("      {")
    print("        \"type\": \"field\",")
    print(f"        \"domain\": [\"geosite:{target_entry.country_code}\"],")
    print("        \"outboundTag\": \"proxy\"")
    print("      }")
    print("    ]")
    print("  }")
    print("}")

def search_domains(dat_path: str, search_term: str):
    """搜索包含指定关键词的域名"""
    geo = load_geosite_dat(dat_path)
    
    print(f"=== 搜索域名: '{search_term}' ===")
    print(f"文件路径: {dat_path}")
    print()
    
    found_domains = []
    for entry in geo.entry:
        for domain in entry.domain:
            if search_term.lower() in domain.value.lower():
                found_domains.append((entry.country_code, domain.value, domain.type))
    
    if not found_domains:
        print(f"未找到包含 '{search_term}' 的域名")
        return
    
    print(f"找到 {len(found_domains)} 个匹配的域名:")
    print()
    
    # 按分类分组显示
    domains_by_category = {}
    for category, domain, domain_type in found_domains:
        if category not in domains_by_category:
            domains_by_category[category] = []
        domains_by_category[category].append((domain, domain_type))
    
    for category, domains in domains_by_category.items():
        print(f"分类: {category} ({len(domains)} 个)")
        for domain, domain_type in domains:
            print(f"  - {domain}")
        print()

def main():
    parser = argparse.ArgumentParser(description="GeoSite.dat 文件解析工具")
    parser.add_argument("dat_file", help="geosite.dat 文件路径")
    parser.add_argument("category", nargs="?", help="要查看的分类名称")
    parser.add_argument("--search", "-s", help="搜索包含指定关键词的域名")
    
    args = parser.parse_args()
    
    if args.search:
        search_domains(args.dat_file, args.search)
    elif args.category:
        list_category_domains(args.dat_file, args.category)
    else:
        list_all_categories(args.dat_file)

if __name__ == "__main__":
    main() 