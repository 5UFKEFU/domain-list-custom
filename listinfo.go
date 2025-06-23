package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

// ListInfo is the information structure of a single file in data directory.
// It includes all types of rules of the file, as well as servel types of
// sturctures of same items for convenience in later process.
type ListInfo struct {
	Name                    fileName
	HasInclusion            bool
	InclusionAttributeMap   map[fileName][]attribute
	FullTypeList            []*router.Domain
	KeywordTypeList         []*router.Domain
	RegexpTypeList          []*router.Domain
	AttributeRuleUniqueList []*router.Domain
	DomainTypeList          []*router.Domain
	DomainTypeUniqueList    []*router.Domain
	AttributeRuleListMap    map[attribute][]*router.Domain
	GeoSite                 *router.GeoSite
}

// NewListInfo return a ListInfo
func NewListInfo() *ListInfo {
	return &ListInfo{
		InclusionAttributeMap:   make(map[fileName][]attribute),
		FullTypeList:            make([]*router.Domain, 0, 10),
		KeywordTypeList:         make([]*router.Domain, 0, 10),
		RegexpTypeList:          make([]*router.Domain, 0, 10),
		AttributeRuleUniqueList: make([]*router.Domain, 0, 10),
		DomainTypeList:          make([]*router.Domain, 0, 10),
		DomainTypeUniqueList:    make([]*router.Domain, 0, 10),
		AttributeRuleListMap:    make(map[attribute][]*router.Domain),
	}
}

// ProcessList processes each line of every single file in the data directory
// and generates a ListInfo of each file.
func (l *ListInfo) ProcessList(file *os.File) error {
	scanner := bufio.NewScanner(file)
	// Parse a file line by line to generate ListInfo
	for scanner.Scan() {
		line := scanner.Text()
		if isEmpty(line) {
			continue
		}
		line = removeComment(line)
		if isEmpty(line) {
			continue
		}
		parsedRule, err := l.parseRule(line)
		if err != nil {
			return err
		}
		if parsedRule == nil {
			continue
		}
		l.classifyRule(parsedRule)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// parseRule parses a single rule
func (l *ListInfo) parseRule(line string) (*router.Domain, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, errors.New("empty line")
	}

	// Parse `include` rule first, eg: `include:google`, `include:google @cn @gfw`
	if strings.HasPrefix(line, "include:") {
		l.parseInclusion(line)
		return nil, nil
	}

	parts := strings.Split(line, " ")
	ruleWithType := strings.TrimSpace(parts[0])
	if ruleWithType == "" {
		return nil, errors.New("empty rule")
	}

	var rule router.Domain
	if err := l.parseTypeRule(ruleWithType, &rule); err != nil {
		return nil, err
	}

	for _, attrString := range parts[1:] {
		if attrString = strings.TrimSpace(attrString); attrString != "" {
			attr, err := l.parseAttribute(attrString)
			if err != nil {
				return nil, err
			}
			rule.Attribute = append(rule.Attribute, attr)
		}
	}

	return &rule, nil
}

func (l *ListInfo) parseInclusion(inclusion string) {
	inclusionVal := strings.TrimPrefix(strings.TrimSpace(inclusion), "include:")
	l.HasInclusion = true
	inclusionValSlice := strings.Split(inclusionVal, "@")
	filename := fileName(strings.ToUpper(strings.TrimSpace(inclusionValSlice[0])))
	switch len(inclusionValSlice) {
	case 1: // Inclusion without attribute
		// Use '@' as the placeholder attribute for 'include:filename'
		l.InclusionAttributeMap[filename] = append(l.InclusionAttributeMap[filename], attribute("@"))
	default: // Inclusion with attribute(s)
		// support new inclusion syntax, eg: `include:google @cn @gfw`
		for _, attr := range inclusionValSlice[1:] {
			attr = strings.ToLower(strings.TrimSpace(attr))
			if attr != "" {
				// Added in this format: '@cn'
				l.InclusionAttributeMap[filename] = append(l.InclusionAttributeMap[filename], attribute("@"+attr))
			}
		}
	}
}

func (l *ListInfo) parseTypeRule(domain string, rule *router.Domain) error {
	kv := strings.Split(domain, ":")
	switch len(kv) {
	case 1: // line without type prefix
		rule.Type = router.Domain_RootDomain
		rule.Value = strings.ToLower(strings.TrimSpace(kv[0]))
	case 2: // line with type prefix
		ruleType := strings.TrimSpace(kv[0])
		ruleVal := strings.TrimSpace(kv[1])
		rule.Value = strings.ToLower(ruleVal)
		switch strings.ToLower(ruleType) {
		case "full":
			rule.Type = router.Domain_Full
		case "domain":
			rule.Type = router.Domain_RootDomain
		case "keyword":
			rule.Type = router.Domain_Plain
		case "regexp":
			rule.Type = router.Domain_Regex
			rule.Value = ruleVal
		default:
			return errors.New("unknown domain type: " + ruleType)
		}
	}
	return nil
}

func (l *ListInfo) parseAttribute(attr string) (*router.Domain_Attribute, error) {
	if attr[0] != '@' {
		fmt.Printf("[警告] 无效属性: %s，已跳过\n", attr)
		return nil, nil // 跳过无效属性
	}
	attr = attr[1:] // Trim out attribute prefix `@` character

	// 允许小写字母、数字、-、_、!，否则警告
	for _, c := range attr {
		if !(c >= 'a' && c <= 'z') && !(c >= '0' && c <= '9') && c != '-' && c != '_' && c != '!' {
			fmt.Printf("[警告] 未知属性: @%s，已跳过\n", attr)
			return nil, nil // 跳过未知属性
		}
	}

	var attribute router.Domain_Attribute
	attribute.Key = strings.ToLower(attr)
	attribute.TypedValue = &router.Domain_Attribute_BoolValue{BoolValue: true}
	return &attribute, nil
}

// classifyRule classifies a single rule and write into *ListInfo
func (l *ListInfo) classifyRule(rule *router.Domain) {
	if len(rule.Attribute) > 0 {
		l.AttributeRuleUniqueList = append(l.AttributeRuleUniqueList, rule)
		var attrsString attribute
		for _, attr := range rule.Attribute {
			attrsString += attribute("@" + attr.GetKey()) // attrsString will be "@cn@ads" if there are more than one attributes
		}
		l.AttributeRuleListMap[attrsString] = append(l.AttributeRuleListMap[attrsString], rule)
	} else {
		switch rule.Type {
		case router.Domain_Full:
			l.FullTypeList = append(l.FullTypeList, rule)
		case router.Domain_RootDomain:
			l.DomainTypeList = append(l.DomainTypeList, rule)
		case router.Domain_Plain:
			l.KeywordTypeList = append(l.KeywordTypeList, rule)
		case router.Domain_Regex:
			l.RegexpTypeList = append(l.RegexpTypeList, rule)
		}
	}
}

// Flatten flattens the rules in a file that have "include" syntax
// in data directory, and adds those need-to-included rules into it.
// This feature supports the "include:filename@attribute" syntax.
// It also generates a domain trie of domain-typed rules for each file
// to remove duplications of them.
func (l *ListInfo) Flatten(lm *ListInfoMap) error {
	if l.HasInclusion {
		for filename, attrs := range l.InclusionAttributeMap {
			for _, attrWanted := range attrs {
				includedList := (*lm)[filename]
				if includedList == nil {
					fmt.Printf("[警告] 文件 %s 不存在于 ListInfoMap 中，跳过\n", filename)
					continue
				}
				switch string(attrWanted) {
				case "@":
					l.FullTypeList = append(l.FullTypeList, includedList.FullTypeList...)
					l.DomainTypeList = append(l.DomainTypeList, includedList.DomainTypeList...)
					l.KeywordTypeList = append(l.KeywordTypeList, includedList.KeywordTypeList...)
					l.RegexpTypeList = append(l.RegexpTypeList, includedList.RegexpTypeList...)
					l.AttributeRuleUniqueList = append(l.AttributeRuleUniqueList, includedList.AttributeRuleUniqueList...)
					for attr, domainList := range includedList.AttributeRuleListMap {
						l.AttributeRuleListMap[attr] = append(l.AttributeRuleListMap[attr], domainList...)
					}

				default:
					for attr, domainList := range includedList.AttributeRuleListMap {
						// If there are more than one attribute attached to the rule,
						// the attribute key of AttributeRuleListMap in ListInfo
						// will be like: "@cn@ads".
						// So if to extract rules with a specific attribute, it is necessary
						// also to test the multi-attribute keys of AttributeRuleListMap.
						// Notice: if "include:google @cn" and "include:google @ads" appear
						// at the same time in the parent list. There are chances that the same
						// rule with that two attributes(`@cn` and `@ads`) will be included twice in the parent list.
						if strings.Contains(string(attr)+"@", string(attrWanted)+"@") {
							l.AttributeRuleListMap[attr] = append(l.AttributeRuleListMap[attr], domainList...)
							l.AttributeRuleUniqueList = append(l.AttributeRuleUniqueList, domainList...)
						}
					}
				}
			}
		}
	}

	sort.Slice(l.DomainTypeList, func(i, j int) bool {
		return len(strings.Split(l.DomainTypeList[i].GetValue(), ".")) < len(strings.Split(l.DomainTypeList[j].GetValue(), "."))
	})

	trie := NewDomainTrie()
	for _, domain := range l.DomainTypeList {
		success, err := trie.Insert(domain.GetValue())
		if err != nil {
			return err
		}
		if success {
			l.DomainTypeUniqueList = append(l.DomainTypeUniqueList, domain)
		}
	}

	return nil
}

// ToGeoSite converts every ListInfo into a router.GeoSite structure.
// It also excludes rules with certain attributes in certain files that
// user specified in command line when runing the program.
func (l *ListInfo) ToGeoSite(excludeAttrs map[fileName]map[attribute]bool) {
	geosite := new(router.GeoSite)
	geosite.CountryCode = string(l.Name)
	geosite.Domain = append(geosite.Domain, l.FullTypeList...)
	geosite.Domain = append(geosite.Domain, l.DomainTypeList...)
	geosite.Domain = append(geosite.Domain, l.DomainTypeUniqueList...)
	geosite.Domain = append(geosite.Domain, l.RegexpTypeList...)

	for _, keywordRule := range l.KeywordTypeList {
		if len(strings.TrimSpace(keywordRule.GetValue())) > 0 {
			geosite.Domain = append(geosite.Domain, keywordRule)
		}
	}

	if excludeAttrs != nil && excludeAttrs[l.Name] != nil {
		excludeAttrsMap := excludeAttrs[l.Name]
		for _, domain := range l.AttributeRuleUniqueList {
			ifKeep := true
			for _, attr := range domain.GetAttribute() {
				if excludeAttrsMap[attribute(attr.GetKey())] {
					ifKeep = false
					break
				}
			}
			if ifKeep {
				geosite.Domain = append(geosite.Domain, domain)
			}
		}
	} else {
		geosite.Domain = append(geosite.Domain, l.AttributeRuleUniqueList...)
	}
	l.GeoSite = geosite
}

// ToPlainText convert router.GeoSite structure to plaintext format.
func (l *ListInfo) ToPlainText() []byte {
	plaintextBytes := make([]byte, 0, 1024*512)

	for _, rule := range l.GeoSite.Domain {
		ruleVal := strings.TrimSpace(rule.GetValue())
		if len(ruleVal) == 0 {
			continue
		}

		var ruleString string
		switch rule.Type {
		case router.Domain_Full:
			ruleString = "full:" + ruleVal
		case router.Domain_RootDomain:
			ruleString = "domain:" + ruleVal
		case router.Domain_Plain:
			ruleString = "keyword:" + ruleVal
		case router.Domain_Regex:
			ruleString = "regexp:" + ruleVal
		}

		if len(rule.Attribute) > 0 {
			ruleString += ":"
			for _, attr := range rule.Attribute {
				ruleString += "@" + attr.GetKey() + ","
			}
			ruleString = strings.TrimRight(ruleString, ",")
		}
		// Output format is: type:domain.tld:@attr1,@attr2
		plaintextBytes = append(plaintextBytes, []byte(ruleString+"\n")...)
	}

	return plaintextBytes
}

// ToGFWList converts router.GeoSite to GFWList format.
func (l *ListInfo) ToGFWList() []byte {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeString := fmt.Sprintf("! Last Modified: %s\n", time.Now().In(loc).Format(time.RFC1123))

	gfwlistBytes := make([]byte, 0, 1024*512)
	gfwlistBytes = append(gfwlistBytes, []byte("[AutoProxy 0.2.9]\n")...)
	gfwlistBytes = append(gfwlistBytes, []byte(timeString)...)
	gfwlistBytes = append(gfwlistBytes, []byte("! Expires: 24h\n")...)
	gfwlistBytes = append(gfwlistBytes, []byte("! HomePage: https://github.com/Loyalsoldier/domain-list-custom\n")...)
	gfwlistBytes = append(gfwlistBytes, []byte("! GitHub URL: https://raw.githubusercontent.com/Loyalsoldier/domain-list-custom/release/gfwlist.txt\n")...)
	gfwlistBytes = append(gfwlistBytes, []byte("! jsdelivr URL: https://cdn.jsdelivr.net/gh/Loyalsoldier/domain-list-custom@release/gfwlist.txt\n")...)
	gfwlistBytes = append(gfwlistBytes, []byte("\n")...)

	for _, rule := range l.GeoSite.Domain {
		ruleVal := strings.TrimSpace(rule.GetValue())
		if len(ruleVal) == 0 {
			continue
		}

		switch rule.Type {
		case router.Domain_Full:
			gfwlistBytes = append(gfwlistBytes, []byte("|http://"+ruleVal+"\n")...)
			gfwlistBytes = append(gfwlistBytes, []byte("|https://"+ruleVal+"\n")...)
		case router.Domain_RootDomain:
			gfwlistBytes = append(gfwlistBytes, []byte("||"+ruleVal+"\n")...)
		case router.Domain_Plain:
			gfwlistBytes = append(gfwlistBytes, []byte(ruleVal+"\n")...)
		case router.Domain_Regex:
			gfwlistBytes = append(gfwlistBytes, []byte("/"+ruleVal+"/\n")...)
		}
	}

	return gfwlistBytes
}

// RecursiveIncludeProcessor 递归处理include关系，收集所有被引用的文件
type RecursiveIncludeProcessor struct {
	processedFiles   map[fileName]bool
	allIncludedFiles map[fileName]bool
}

// NewRecursiveIncludeProcessor 创建新的递归处理器
func NewRecursiveIncludeProcessor() *RecursiveIncludeProcessor {
	return &RecursiveIncludeProcessor{
		processedFiles:   make(map[fileName]bool),
		allIncludedFiles: make(map[fileName]bool),
	}
}

// ProcessRecursiveIncludes 递归处理指定文件的include关系
func (r *RecursiveIncludeProcessor) ProcessRecursiveIncludes(startFiles []fileName, lm *ListInfoMap) error {
	for _, startFile := range startFiles {
		if err := r.processFileRecursively(startFile, lm); err != nil {
			return err
		}
	}
	return nil
}

// processFileRecursively 递归处理单个文件
func (r *RecursiveIncludeProcessor) processFileRecursively(fileName fileName, lm *ListInfoMap) error {
	// 如果已经处理过，跳过
	if r.processedFiles[fileName] {
		return nil
	}

	r.processedFiles[fileName] = true

	// 获取文件信息
	listInfo, exists := (*lm)[fileName]
	if !exists {
		fmt.Printf("[警告] 文件 %s 不存在，跳过\n", fileName)
		return nil
	}

	// 如果文件有include，递归处理
	if listInfo.HasInclusion {
		for includedFile := range listInfo.InclusionAttributeMap {
			r.allIncludedFiles[includedFile] = true
			if err := r.processFileRecursively(includedFile, lm); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetAllIncludedFiles 获取所有被递归include的文件
func (r *RecursiveIncludeProcessor) GetAllIncludedFiles() []fileName {
	var files []fileName
	for file := range r.allIncludedFiles {
		files = append(files, file)
	}
	return files
}

// GetMultiCountryCategoryFiles 获取多国家相关的分类文件
func GetMultiCountryCategoryFiles() []fileName {
	return []fileName{
		// 私有网络分类
		"PRIVATE",

		// 中国分类
		"CN",
		"TLD-CN",
		"CATEGORY-AI-CN",
		"CATEGORY-BANK-CN",
		"CATEGORY-BLOG-CN",
		"CATEGORY-COLLABORATE-CN",
		"CATEGORY-DEV-CN",
		"CATEGORY-DOCUMENTS-CN",
		"CATEGORY-EDUCATION-CN",
		"CATEGORY-ELECTRONIC-CN",
		"CATEGORY-ENTERPRISE-QUERY-PLATFORM-CN",
		"CATEGORY-ENTERTAINMENT-CN",
		"CATEGORY-FOOD-CN",
		"CATEGORY-GAME-ACCELERATOR-CN",
		"CATEGORY-GAMES-CN",
		"CATEGORY-HOSPITAL-CN",
		"CATEGORY-HTTPDNS-CN",
		"CATEGORY-LOGISTICS-CN",
		"CATEGORY-MEDIA-CN",
		"CATEGORY-MOOC-CN",
		"CATEGORY-NETDISK-CN",
		"CATEGORY-NETWORK-SECURITY-CN",
		"CATEGORY-NTP-CN",
		"CATEGORY-NUMBER-VERIFICATION-CN",
		"CATEGORY-OUTSOURCE-CN",
		"CATEGORY-SCHOLAR-CN",
		"CATEGORY-SECURITIES-CN",
		"CATEGORY-SOCIAL-MEDIA-CN",
		"CATEGORY-WIKI-CN",

		// 伊朗分类
		"CATEGORY-ADS-IR",
		"CATEGORY-BANK-IR",
		"CATEGORY-BOURSE-IR",
		"CATEGORY-EDUCATION-IR",
		"CATEGORY-FORUMS-IR",
		"CATEGORY-GOV-IR",
		"CATEGORY-INSURANCE-IR",
		"CATEGORY-MEDIA-IR",
		"CATEGORY-NEWS-IR",
		"CATEGORY-PAYMENT-IR",
		"CATEGORY-SCHOLAR-IR",
		"CATEGORY-SHOPPING-IR",
		"CATEGORY-SOCIAL-MEDIA-IR",
		"CATEGORY-TECH-IR",
		"CATEGORY-TRAVEL-IR",

		// 俄罗斯分类
		"CATEGORY-GOV-RU",
		"CATEGORY-MEDIA-RU",

		// 色情分类
		"CATEGORY-PORN",
	}
}

// GetChineseCategoryFiles 获取中国相关的分类文件（保持向后兼容）
func GetChineseCategoryFiles() []fileName {
	return []fileName{
		// 私有网络分类
		"PRIVATE",

		// 中国分类
		"CN",
		"TLD-CN",
		"CATEGORY-AI-CN",
		"CATEGORY-BANK-CN",
		"CATEGORY-BLOG-CN",
		"CATEGORY-COLLABORATE-CN",
		"CATEGORY-DEV-CN",
		"CATEGORY-DOCUMENTS-CN",
		"CATEGORY-EDUCATION-CN",
		"CATEGORY-ELECTRONIC-CN",
		"CATEGORY-ENTERPRISE-QUERY-PLATFORM-CN",
		"CATEGORY-ENTERTAINMENT-CN",
		"CATEGORY-FOOD-CN",
		"CATEGORY-GAME-ACCELERATOR-CN",
		"CATEGORY-GAMES-CN",
		"CATEGORY-HOSPITAL-CN",
		"CATEGORY-HTTPDNS-CN",
		"CATEGORY-LOGISTICS-CN",
		"CATEGORY-MEDIA-CN",
		"CATEGORY-MOOC-CN",
		"CATEGORY-NETDISK-CN",
		"CATEGORY-NETWORK-SECURITY-CN",
		"CATEGORY-NTP-CN",
		"CATEGORY-NUMBER-VERIFICATION-CN",
		"CATEGORY-OUTSOURCE-CN",
		"CATEGORY-SCHOLAR-CN",
		"CATEGORY-SECURITIES-CN",
		"CATEGORY-SOCIAL-MEDIA-CN",
		"CATEGORY-WIKI-CN",
	}
}
