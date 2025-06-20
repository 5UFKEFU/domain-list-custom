package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

// ListInfoMap is the map of files in data directory and ListInfo
type ListInfoMap map[fileName]*ListInfo

// Marshal processes a file in data directory and generates ListInfo for it.
func (lm *ListInfoMap) Marshal(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	list := NewListInfo()
	listName := fileName(strings.ToUpper(filepath.Base(path)))
	list.Name = listName
	if err := list.ProcessList(file); err != nil {
		return err
	}

	(*lm)[listName] = list
	return nil
}

// FlattenAndGenUniqueDomainList flattens the included lists and
// generates a domain trie for each file in data directory to
// make the items of domain type list unique.
func (lm *ListInfoMap) FlattenAndGenUniqueDomainList() error {
	// 添加调试信息
	fmt.Println("=== 调试信息：所有加载的文件 ===")
	for filename, listinfo := range *lm {
		fmt.Printf("文件: %s, 有include: %v, include数量: %d\n",
			filename, listinfo.HasInclusion, len(listinfo.InclusionAttributeMap))
		if listinfo.HasInclusion {
			for includedFile := range listinfo.InclusionAttributeMap {
				fmt.Printf("  - include: %s\n", includedFile)
			}
		}
	}
	fmt.Println()

	inclusionLevel := make([]map[fileName]bool, 0, 20)
	okayList := make(map[fileName]bool)
	inclusionLevelAllLength, loopTimes := 0, 0
	maxLoops := 100 // 防止死循环的最大循环次数

	for inclusionLevelAllLength < len(*lm) && loopTimes < maxLoops {
		inclusionMap := make(map[fileName]bool)

		if loopTimes == 0 {
			for _, listinfo := range *lm {
				if listinfo.HasInclusion {
					continue
				}
				inclusionMap[listinfo.Name] = true
			}
		} else {
			for _, listinfo := range *lm {
				if !listinfo.HasInclusion || okayList[listinfo.Name] {
					continue
				}

				var passTimes int
				for filename := range listinfo.InclusionAttributeMap {
					if !okayList[filename] {
						break
					}
					passTimes++
				}
				if passTimes == len(listinfo.InclusionAttributeMap) {
					inclusionMap[listinfo.Name] = true
				}
			}
		}

		for filename := range inclusionMap {
			okayList[filename] = true
		}

		inclusionLevel = append(inclusionLevel, inclusionMap)
		inclusionLevelAllLength += len(inclusionMap)
		loopTimes++

		// 如果本轮没有处理任何文件，说明有循环依赖或者还有未处理的文件
		if len(inclusionMap) == 0 {
			// 检查是否还有未处理的文件
			remainingFiles := 0
			for _, listinfo := range *lm {
				if !okayList[listinfo.Name] {
					remainingFiles++
					fmt.Printf("[警告] 未处理的文件: %s (有include: %v)\n",
						listinfo.Name, listinfo.HasInclusion)
				}
			}
			if remainingFiles > 0 {
				fmt.Printf("[警告] 还有 %d 个文件未处理，可能存在循环依赖\n", remainingFiles)
			}
			break
		}
	}

	if loopTimes >= maxLoops {
		return fmt.Errorf("达到最大循环次数 %d，可能存在死循环", maxLoops)
	}

	for idx, inclusionMap := range inclusionLevel {
		fmt.Printf("Level %d:\n", idx+1)
		fmt.Println(inclusionMap)
		fmt.Println()

		for inclusionFilename := range inclusionMap {
			if err := (*lm)[inclusionFilename].Flatten(lm); err != nil {
				return err
			}
		}
	}

	// 强制处理所有未处理的文件，确保所有 include 关系都被展开
	fmt.Println("=== 强制处理未处理的文件 ===")
	forceProcessed := false
	for _, listinfo := range *lm {
		if !okayList[listinfo.Name] {
			fmt.Printf("强制处理文件: %s\n", listinfo.Name)
			if err := listinfo.Flatten(lm); err != nil {
				return err
			}
			forceProcessed = true
		}
	}
	if forceProcessed {
		fmt.Println("强制处理完成")
		fmt.Println()
	}

	return nil
}

// ToProto generates a router.GeoSite for each file in data directory
// and returns a router.GeoSiteList
func (lm *ListInfoMap) ToProto(excludeAttrs map[fileName]map[attribute]bool) *router.GeoSiteList {
	protoList := new(router.GeoSiteList)
	for _, listinfo := range *lm {
		listinfo.ToGeoSite(excludeAttrs)
		protoList.Entry = append(protoList.Entry, listinfo.GeoSite)
	}
	return protoList
}

// ToPlainText returns a map of exported lists that user wants
// and the contents of them in byte format.
func (lm *ListInfoMap) ToPlainText(exportListsMap []string) (map[string][]byte, error) {
	filePlainTextBytesMap := make(map[string][]byte)
	for _, filename := range exportListsMap {
		if listinfo := (*lm)[fileName(strings.ToUpper(filename))]; listinfo != nil {
			// 新增：递归收集所有域名
			collected := make(map[string]bool)
			allRules := lm.collectAllRulesRecursive(listinfo, collected)
			plaintextBytes := make([]byte, 0, 1024*512)
			for _, rule := range allRules {
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
				plaintextBytes = append(plaintextBytes, []byte(ruleString+"\n")...)
			}
			filePlainTextBytesMap[filename] = plaintextBytes
		} else {
			fmt.Println("Notice: " + filename + ": no such exported list in the directory, skipped.")
		}
	}
	return filePlainTextBytesMap, nil
}

// collectAllRulesRecursive 递归收集所有域名规则，防止循环引用
func (lm *ListInfoMap) collectAllRulesRecursive(listinfo *ListInfo, visited map[string]bool) []*router.Domain {
	if listinfo == nil {
		return nil
	}
	if visited[strings.ToUpper(string(listinfo.Name))] {
		return nil // 防止循环
	}
	visited[strings.ToUpper(string(listinfo.Name))] = true
	allRules := make([]*router.Domain, 0)

	// 收集本地规则 - 包括所有类型的域名
	allRules = append(allRules, listinfo.FullTypeList...)
	allRules = append(allRules, listinfo.DomainTypeList...)
	allRules = append(allRules, listinfo.KeywordTypeList...)
	allRules = append(allRules, listinfo.RegexpTypeList...)
	allRules = append(allRules, listinfo.AttributeRuleUniqueList...)

	// 递归 include
	if listinfo.HasInclusion {
		for filename := range listinfo.InclusionAttributeMap {
			included := (*lm)[filename]
			if included != nil {
				allRules = append(allRules, lm.collectAllRulesRecursive(included, visited)...)
			}
		}
	}
	return allRules
}

// ToGFWList returns the content of the list to be generated into GFWList format
// that user wants in bytes format.
func (lm *ListInfoMap) ToGFWList(togfwlist string) ([]byte, error) {
	if togfwlist != "" {
		if listinfo := (*lm)[fileName(strings.ToUpper(togfwlist))]; listinfo != nil {
			return listinfo.ToGFWList(), nil
		}
		return nil, errors.New("no such list: " + togfwlist)
	}
	return nil, nil
}
