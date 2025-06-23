package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
)

var (
	dataPath            = flag.String("datapath", filepath.Join("./", "data"), "Path to your custom 'data' directory")
	datName             = flag.String("datname", "geosite.dat", "Name of the generated dat file")
	outputPath          = flag.String("outputpath", "./publish", "Output path to the generated files")
	exportLists         = flag.String("exportlists", "google,apple,meta,facebook,facebook-dev,instagram,messenger,oculus,threads,whatsapp,microsoft,amazon,tiktok,baidu,baidu-ads,alibaba,alibaba-ads,alibabacloud,tencent,tencent-ads,bytedance,bytedance-ads,xiaomi,huawei,huaweicloud,oppo,vivo,meituan,didi,jd,netease,sina,sohu,iqiyi,youku,bilibili,category-ai-cn,category-ai-!cn,openai,anthropic,google-deepmind,groq,huggingface,perplexity,poe,xai,cursor,cn,tld-cn,category-ir,category-ru,mailru,ok,ozon,vk,yandex,category-gov-ru,category-porn,category-ads-all,taboola,category-ads,acfun-ads,adcolony-ads,adjust-ads,adobe-ads,amazon-ads,apple-ads,applovin-ads,atom-data-ads,category-ads-ir,clearbitjs-ads,dmm-ads,duolingo-ads,emogi-ads,flurry-ads,google-ads,growingio-ads,hiido-ads,hotjar-ads,hunantv-ads,inner-active-ads,iqiyi-ads,jd-ads,kuaishou-ads,kugou-ads,leanplum-ads,letv-ads,mixpanel-ads,mopub-ads,mxplayer-ads,netease-ads,newrelic-ads,ogury-ads,onesignal-ads,ookla-speedtest-ads,openx-ads,pocoiq-ads,pubmatic-ads,qihoo360-ads,ruanmei,segment-ads,sensorsdata-ads,sina-ads,sohu-ads,spotify-ads,supersonic-ads,tagtic-ads,tappx-ads,television-ads,tencent-ads,uberads-ads,umeng-ads,unity-ads,vivo,wteam-ads,xhamster-ads,xiaomi-ads,ximalaya-ads,yahoo-ads,zynga-ads,android,blogspot,dart,fastlane,firebase,flutter,golang,google-gemini,google-play,google-registry,google-scholar,google-trust-services,googlefcm,kaggle,opensourceinsights,polymer,v8,youtube,apple-dev,apple-pki,apple-tvplus,apple-update,beats,icloud,itunes,swift,azure,bing,github,microsoft-dev,microsoft-pki,msn,onedrive,xbox,amazontrust,aws,imdb,kindle,primevideo,wholefoodsmarket", "export lists")
	excludeAttrs        = flag.String("excludeattrs", "", "Exclude rules with certain attributes in certain lists, seperated by ',' comma, support multiple attributes in one list. Example: geolocation-!cn@cn@ads,geolocation-cn@!cn")
	toGFWList           = flag.String("togfwlist", "", "List to be exported in GFWList format")
	useRecursiveInclude = flag.Bool("recursive", false, "Use recursive include processing for Chinese domains")
	recursiveMode       = flag.String("mode", "cn", "Recursive mode: cn (Chinese only), multi (Multi-country including CN, IR, RU, Porn)")
	showInfo            = flag.Bool("info", false, "Show information about the generated geosite.dat file")
	getList             = flag.String("getlist", "", "Get all domains from a specific list (e.g., tiktok, category-ads-all)")
	getAllAds           = flag.Bool("getallads", false, "Get all advertising domains from all -ads files")
)

// mergeCategories 将子分类合并到主分类中，并删除子分类以减少文件大小
func mergeCategories(lm *ListInfoMap) {
	// 定义主分类和其对应的子分类
	categoryMappings := map[string][]string{
		"google": {
			"google-ads", "google-deepmind", "google-gemini", "google-play",
			"google-registry", "google-scholar", "google-trust-services",
			"googlefcm", "android", "blogspot", "dart", "fastlane", "firebase",
			"flutter", "golang", "kaggle", "opensourceinsights", "polymer",
			"v8", "youtube",
		},
		"apple": {
			"apple-ads", "apple-dev", "apple-pki", "apple-tvplus",
			"apple-update", "beats", "icloud", "itunes", "swift",
		},
		"microsoft": {
			"microsoft-dev", "microsoft-pki", "azure", "bing", "github",
			"msn", "onedrive", "xbox",
		},
		"amazon": {
			"amazon-ads", "amazontrust", "aws", "imdb", "kindle",
			"primevideo", "wholefoodsmarket",
		},
		"meta": {
			"facebook", "facebook-dev", "instagram", "messenger",
			"oculus", "threads", "whatsapp",
		},
		"baidu": {
			"baidu-ads", "zuoyebang",
		},
		"alibaba": {
			"alibaba-ads", "alibabacloud", "aliyun", "dingtalk", "eleme",
			"teambition", "amap", "cainiao", "uc", "umeng",
		},
		"tencent": {
			"tencent-ads", "tencent-games", "yuewen", "qcloud", "tencent-dev",
		},
		"bytedance": {
			"bytedance-ads", "bcy", "fqnovel", "juejin", "lark", "tiktok", "volcengine",
		},
		"xiaomi": {
			"xiaomi-ads",
		},
		"huawei": {
			"huaweicloud", "huawei-dev",
		},
		"jd": {
			"jd-ads",
		},
		"netease": {
			"netease-ads",
		},
		"sina": {
			"sina-ads",
		},
		"sohu": {
			"sohu-ads", "sogou",
		},
		"iqiyi": {
			"iqiyi-ads",
		},
		"youku": {
			"youku-ads",
		},
		"bilibili": {
			"bilibili-game",
		},
		"instagram": {
			"instagram-ads",
		},
		"whatsapp": {
			"whatsapp-ads",
		},
		"xbox": {
			"bethesda", "forza", "mojang", "asobo",
		},
		"aws": {
			"aws-cn",
		},
		"icloud": {
			"icloudprivaterelay",
		},
		"cn": {
			"tld-cn",
		},
		"category-ai-cn": {
			"deepseek",
		},
		"google-gemini": {
			"google-deepmind",
		},
		"category-porn": {
			"fans66", "lethalhardcore", "nudevista", "yunlaopo", "fansta", "hentaivn",
			"javcc", "jkf", "18comic", "boboporn", "zhimeishe", "camwhores", "chatwhores",
			"kubakuba", "54647", "anon-v", "bttzyw", "lisiku", "bongacams", "dlsite",
			"konachan", "moxing", "truyen-hentai", "xhamster", "youjizz", "haitang",
			"heyzo", "smtiaojiaoshi", "xingkongwuxianmedia", "hooligapps", "illusion-nonofficial",
			"jable", "shireyishunjian", "swag", "uu-chat", "cuinc", "dmm-porn", "javbus",
			"kemono", "pornpros", "tokyo-toshokan", "erolabs", "picacg", "sehuatang",
			"bdsmhub", "clips4sale", "johren", "justav", "metart", "thescoregroup",
			"xnxx", "cavporn", "javdb", "missav", "netflav", "avmoo", "bilibili2",
			"boylove", "ehentai", "illusion", "mindgeek-porn", "playboy", "spankbang",
			"japonx", "javwide", "theporndude", "xvideos",
		},
		"category-ads": {
			"emogi-ads", "hunantv-ads", "iqiyi-ads", "kuaishou-ads", "newrelic-ads",
			"tencent-ads", "tagtic-ads", "wteam-ads", "alibaba-ads", "flurry-ads",
			"hotjar-ads", "jd-ads", "kugou-ads", "onesignal-ads", "adjust-ads",
			"openx-ads", "supersonic-ads", "adobe-ads", "applovin-ads", "mopub-ads",
			"ookla-speedtest-ads", "spotify-ads", "youku-ads", "acfun-ads", "amazon-ads",
			"duolingo-ads", "pocoiq-ads", "segment-ads", "xhamster-ads", "adcolony-ads",
			"hiido-ads", "sensorsdata-ads", "xiaomi-ads", "tappx-ads", "yahoo-ads",
			"mixpanel-ads", "qihoo360-ads", "atom-data-ads", "leanplum-ads", "sohu-ads",
			"unity-ads", "bytedance-ads", "clearbitjs-ads", "zynga-ads", "apple-ads",
			"sina-ads", "uberads-ads", "baidu-ads", "category-ads-ir", "letv-ads",
			"ogury-ads", "television-ads", "umeng-ads", "google-ads", "netease-ads",
			"pubmatic-ads", "growingio-ads", "inner-active-ads", "mxplayer-ads",
			"ximalaya-ads", "dmm-ads",
		},
		"category-ads-all": {
			"category-ads", "taboola",
		},
		"category-ir": {
			"category-travel-ir", "category-bourse-ir", "category-education-ir",
			"category-insurance-ir", "category-forums-ir", "category-gov-ir",
			"category-news-ir", "category-social-media-ir", "category-tech-ir",
			"category-bank-ir", "category-media-ir", "category-payment-ir",
			"category-scholar-ir", "category-shopping-ir", "snapp",
		},
		"category-ru": {
			"ozon", "vk", "yandex", "category-gov-ru", "mailru", "ok",
		},
		"category-ai-!cn": {
			"huggingface", "openai", "anthropic", "groq", "perplexity", "poe", "xai", "cursor", "google-deepmind",
		},
		"oppo": {
			"oneplus",
		},
	}

	// 记录要删除的子分类
	subCategoriesToDelete := make(map[fileName]bool)

	// 遍历每个主分类，合并其子分类
	for mainCategory, subCategories := range categoryMappings {
		mainFileName := fileName(strings.ToUpper(mainCategory))
		mainListInfo, mainExists := (*lm)[mainFileName]

		if !mainExists {
			continue
		}

		// 合并每个子分类的内容
		for _, subCategory := range subCategories {
			subFileName := fileName(strings.ToUpper(subCategory))
			subListInfo, subExists := (*lm)[subFileName]

			if !subExists {
				continue
			}

			// 合并域名列表
			mainListInfo.DomainTypeList = append(mainListInfo.DomainTypeList, subListInfo.DomainTypeList...)
			mainListInfo.FullTypeList = append(mainListInfo.FullTypeList, subListInfo.FullTypeList...)
			mainListInfo.KeywordTypeList = append(mainListInfo.KeywordTypeList, subListInfo.KeywordTypeList...)
			mainListInfo.RegexpTypeList = append(mainListInfo.RegexpTypeList, subListInfo.RegexpTypeList...)
			mainListInfo.AttributeRuleUniqueList = append(mainListInfo.AttributeRuleUniqueList, subListInfo.AttributeRuleUniqueList...)

			// 合并属性规则映射
			for attr, rules := range subListInfo.AttributeRuleListMap {
				mainListInfo.AttributeRuleListMap[attr] = append(mainListInfo.AttributeRuleListMap[attr], rules...)
			}

			// 合并包含关系
			if subListInfo.HasInclusion {
				mainListInfo.HasInclusion = true
				for subFile, attrs := range subListInfo.InclusionAttributeMap {
					mainListInfo.InclusionAttributeMap[subFile] = append(mainListInfo.InclusionAttributeMap[subFile], attrs...)
				}
			}

			// 标记子分类为要删除
			subCategoriesToDelete[subFileName] = true
		}
	}

	// 删除所有子分类
	deletedCount := 0
	for subFileName := range subCategoriesToDelete {
		delete(*lm, subFileName)
		deletedCount++
	}

	fmt.Printf("已删除 %d 个子分类以减少文件大小\n", deletedCount)
}

func main() {
	flag.Parse()

	// 如果只是显示信息，则只执行显示功能
	if *showInfo {
		showGeoSiteInfo()
		return
	}

	// 如果只是获取指定列表的域名，则只执行获取功能
	if *getList != "" {
		getListDomains(*getList)
		return
	}

	// 如果只是获取所有广告域名，则只执行获取功能
	if *getAllAds {
		getAllAdvertisingDomains()
		return
	}

	// showDefaultLists()

	dir := GetDataDir()
	listInfoMap := make(ListInfoMap)

	// Process and split *exportLists first
	var exportListsSlice []string
	if *exportLists != "" {
		tempSlice := strings.Split(*exportLists, ",")
		for _, exportList := range tempSlice {
			exportList = strings.TrimSpace(exportList)
			if len(exportList) > 0 {
				exportListsSlice = append(exportListsSlice, exportList)
			}
		}
	}

	// 如果启用递归include处理
	if *useRecursiveInclude {
		var categoryFiles []fileName
		var modeDescription string

		if *recursiveMode == "multi" {
			fmt.Println("启用多国家递归include处理，自动收集中国、伊朗、俄罗斯、色情域名...")
			categoryFiles = GetMultiCountryCategoryFiles()
			modeDescription = "多国家（中国、伊朗、俄罗斯、色情）"
		} else {
			fmt.Println("启用中国递归include处理，自动收集所有中国域名...")
			categoryFiles = GetChineseCategoryFiles()
			modeDescription = "中国"
		}

		// 首先处理所有分类文件
		for _, filename := range categoryFiles {
			filePath := filepath.Join(dir, string(filename))
			if err := listInfoMap.Marshal(filePath); err != nil {
				fmt.Printf("Warning: Failed to process %s: %v\n", filename, err)
				continue
			}
		}

		// 递归处理include关系
		processor := NewRecursiveIncludeProcessor()
		if err := processor.ProcessRecursiveIncludes(categoryFiles, &listInfoMap); err != nil {
			fmt.Printf("递归处理include时出错: %v\n", err)
			os.Exit(1)
		}

		// 获取所有被include的文件
		allIncludedFiles := processor.GetAllIncludedFiles()
		fmt.Printf("递归发现 %d 个被include的文件\n", len(allIncludedFiles))

		// 处理所有被include的文件
		for _, filename := range allIncludedFiles {
			filePath := filepath.Join(dir, string(filename))
			if err := listInfoMap.Marshal(filePath); err != nil {
				fmt.Printf("Warning: Failed to process included file %s: %v\n", filename, err)
				continue
			}
		}

		// 更新exportListsSlice为所有发现的文件
		exportListsSlice = nil
		for _, filename := range categoryFiles {
			exportListsSlice = append(exportListsSlice, strings.ToLower(string(filename)))
		}
		for _, filename := range allIncludedFiles {
			exportListsSlice = append(exportListsSlice, strings.ToLower(string(filename)))
		}

		fmt.Printf("最终将处理 %d 个文件（%s模式）\n", len(exportListsSlice), modeDescription)
	} else {
		// 原有的处理逻辑
		// Only process specified files instead of walking entire directory
		for _, filename := range exportListsSlice {
			filePath := filepath.Join(dir, filename)
			if err := listInfoMap.Marshal(filePath); err != nil {
				fmt.Printf("Warning: Failed to process %s: %v\n", filename, err)
				continue
			}
		}
	}

	if err := listInfoMap.FlattenAndGenUniqueDomainList(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

	// 合并子分类到主分类中
	mergeCategories(&listInfoMap)
	fmt.Println("已合并子分类到主分类中")

	// 过滤掉已经被合并的子分类，避免生成重复的独立文件
	var filteredExportLists []string
	mergedSubCategories := make(map[string]bool)

	// 收集所有被合并的子分类
	categoryMappings := map[string][]string{
		"google": {
			"android", "blogspot", "dart", "fastlane", "firebase", "flutter", "golang", "google-ads", "google-deepmind", "google-gemini", "google-play", "google-registry", "google-scholar", "google-trust-services", "googlefcm", "kaggle", "opensourceinsights", "polymer", "v8", "youtube",
		},
		"apple": {
			"apple-ads", "apple-dev", "apple-pki", "apple-tvplus", "apple-update", "beats", "icloud", "itunes", "swift",
		},
		"microsoft": {
			"azure", "bing", "github", "microsoft-dev", "microsoft-pki", "msn", "onedrive", "xbox",
		},
		"amazon": {
			"amazontrust", "aws", "imdb", "kindle", "primevideo", "wholefoodsmarket",
		},
		"meta": {
			"facebook", "facebook-dev", "instagram", "messenger", "oculus", "threads", "whatsapp",
		},
		"category-porn": {
			"bilibili2", "lethalhardcore", "nudevista", "yunlaopo", "fansta", "hentaivn", "javcc", "jkf", "18comic", "boboporn", "zhimeishe", "camwhores", "chatwhores", "kubakuba", "54647", "anon-v", "bttzyw", "lisiku", "bongacams", "dlsite", "konachan", "moxing", "truyen-hentai", "xhamster", "youjizz", "haitang", "heyzo", "smtiaojiaoshi", "xingkongwuxianmedia", "hooligapps", "illusion-nonofficial", "jable", "shireyishunjian", "swag", "uu-chat", "cuinc", "dmm-porn", "javbus", "kemono", "pornpros", "tokyo-toshokan", "erolabs", "picacg", "sehuatang", "bdsmhub", "clips4sale", "johren", "justav", "metart", "thescoregroup", "xnxx", "cavporn", "javdb", "missav", "netflav", "avmoo", "boylove", "ehentai", "illusion", "mindgeek-porn", "playboy", "spankbang", "japonx", "javwide", "theporndude", "xvideos",
		},
		"category-ads": {
			"emogi-ads", "hunantv-ads", "iqiyi-ads", "kuaishou-ads", "newrelic-ads", "tencent-ads", "tagtic-ads", "wteam-ads", "alibaba-ads", "flurry-ads", "hotjar-ads", "jd-ads", "kugou-ads", "onesignal-ads", "adjust-ads", "openx-ads", "supersonic-ads", "adobe-ads", "applovin-ads", "mopub-ads", "ookla-speedtest-ads", "spotify-ads", "youku-ads", "acfun-ads", "amazon-ads", "duolingo-ads", "pocoiq-ads", "segment-ads", "xhamster-ads", "adcolony-ads", "hiido-ads", "sensorsdata-ads", "xiaomi-ads", "tappx-ads", "yahoo-ads", "mixpanel-ads", "qihoo360-ads", "atom-data-ads", "leanplum-ads", "sohu-ads", "unity-ads", "bytedance-ads", "clearbitjs-ads", "zynga-ads", "apple-ads", "sina-ads", "uberads-ads", "baidu-ads", "category-ads-ir", "letv-ads", "ogury-ads", "television-ads", "umeng-ads", "google-ads", "netease-ads", "pubmatic-ads", "growingio-ads", "inner-active-ads", "mxplayer-ads", "ximalaya-ads", "dmm-ads",
		},
		"category-ads-all": {
			"category-ads", "taboola",
		},
		"category-ir": {
			"category-travel-ir", "category-bourse-ir", "category-education-ir", "category-insurance-ir", "category-forums-ir", "category-gov-ir", "category-news-ir", "category-social-media-ir", "category-tech-ir", "category-bank-ir", "category-media-ir", "category-payment-ir", "category-scholar-ir", "category-shopping-ir", "snapp",
		},
		"category-ru": {
			"ozon", "vk", "yandex", "category-gov-ru", "mailru", "ok",
		},
		"category-ai-!cn": {
			"huggingface", "openai", "anthropic", "groq", "perplexity", "poe", "xai", "cursor", "google-deepmind",
		},
		"oppo": {
			"oneplus",
		},
	}

	for _, subCategories := range categoryMappings {
		for _, subCategory := range subCategories {
			mergedSubCategories[strings.ToLower(subCategory)] = true
		}
	}

	// 过滤掉已经被合并的子分类
	for _, exportList := range exportListsSlice {
		if !mergedSubCategories[exportList] {
			filteredExportLists = append(filteredExportLists, exportList)
		}
	}

	fmt.Printf("过滤后剩余 %d 个分类文件（已移除 %d 个被合并的子分类）\n", len(filteredExportLists), len(exportListsSlice)-len(filteredExportLists))

	// Process and split *excludeRules
	excludeAttrsInFile := make(map[fileName]map[attribute]bool)
	if *excludeAttrs != "" {
		exFilenameAttrSlice := strings.Split(*excludeAttrs, ",")
		for _, exFilenameAttr := range exFilenameAttrSlice {
			exFilenameAttr = strings.TrimSpace(exFilenameAttr)
			exFilenameAttrMap := strings.Split(exFilenameAttr, "@")
			filename := fileName(strings.ToUpper(strings.TrimSpace(exFilenameAttrMap[0])))
			excludeAttrsInFile[filename] = make(map[attribute]bool)
			for _, attr := range exFilenameAttrMap[1:] {
				attr = strings.TrimSpace(attr)
				if len(attr) > 0 {
					excludeAttrsInFile[filename][attribute(attr)] = true
				}
			}
		}
	}

	// Generate dlc.dat
	if geositeList := listInfoMap.ToProto(excludeAttrsInFile); geositeList != nil {
		protoBytes, err := proto.Marshal(geositeList)
		if err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		}
		if err := os.MkdirAll(*outputPath, 0755); err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(filepath.Join(*outputPath, *datName), protoBytes, 0644); err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		} else {
			fmt.Printf("%s has been generated successfully in '%s'.\n", *datName, *outputPath)
		}
	}

	// Generate plaintext list files
	if filePlainTextBytesMap, err := listInfoMap.ToPlainText(filteredExportLists); err == nil {
		for filename, plaintextBytes := range filePlainTextBytesMap {
			filename += ".txt"
			if err := os.WriteFile(filepath.Join(*outputPath, filename), plaintextBytes, 0644); err != nil {
				fmt.Println("Failed:", err)
				os.Exit(1)
			} else {
				fmt.Printf("%s has been generated successfully in '%s'.\n", filename, *outputPath)
			}
		}
	} else {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

	// Generate gfwlist.txt
	if gfwlistBytes, err := listInfoMap.ToGFWList(*toGFWList); err == nil {
		if f, err := os.OpenFile(filepath.Join(*outputPath, "gfwlist.txt"), os.O_RDWR|os.O_CREATE, 0644); err != nil {
			fmt.Println("Failed:", err)
			os.Exit(1)
		} else {
			encoder := base64.NewEncoder(base64.StdEncoding, f)
			defer encoder.Close()
			if _, err := encoder.Write(gfwlistBytes); err != nil {
				fmt.Println("Failed:", err)
				os.Exit(1)
			}
			fmt.Printf("gfwlist.txt has been generated successfully in '%s'.\n", *outputPath)
		}
	} else {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
}

// showGeoSiteInfo 显示geosite.dat文件的信息
func showGeoSiteInfo() {
	fmt.Println("=== GeoSite.dat 文件信息 ===")

	// 检查geosite.dat文件
	geositePath := filepath.Join(*outputPath, *datName)
	if _, err := os.Stat(geositePath); os.IsNotExist(err) {
		fmt.Printf("错误: 文件 %s 不存在\n", geositePath)
		fmt.Println("请先运行程序生成 geosite.dat 文件")
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(geositePath)
	if err != nil {
		fmt.Printf("错误: 无法获取文件信息: %v\n", err)
		return
	}

	fmt.Printf("文件路径: %s\n", geositePath)
	fmt.Printf("文件大小: %d 字节 (%.2f KB)\n", fileInfo.Size(), float64(fileInfo.Size())/1024)
	fmt.Printf("修改时间: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))

	// 显示生成该文件的源文件列表
	fmt.Println("\n=== 生成该文件的源文件列表 ===")

	// 检查publish目录中的txt文件
	publishDir := *outputPath
	files, err := os.ReadDir(publishDir)
	if err != nil {
		fmt.Printf("错误: 无法读取目录 %s: %v\n", publishDir, err)
		return
	}

	var txtFiles []string
	var totalSize int64
	var totalDomains int

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			// 过滤掉以"._"开头的临时文件
			if strings.HasPrefix(file.Name(), "._") {
				continue
			}

			txtFiles = append(txtFiles, file.Name())

			// 获取文件大小
			if fileInfo, err := file.Info(); err == nil {
				totalSize += fileInfo.Size()
			}

			// 统计域名数量
			domainCount := countDomainsInFile(filepath.Join(publishDir, file.Name()))
			totalDomains += domainCount
		}
	}

	fmt.Printf("源文件总数: %d (已过滤临时文件)\n", len(txtFiles))
	fmt.Printf("源文件总大小: %d 字节 (%.2f KB)\n", totalSize, float64(totalSize)/1024)
	fmt.Printf("总域名数量: %d\n", totalDomains)

	// 按分类显示文件
	fmt.Println("\n=== 按分类显示源文件 ===")

	categories := categorizeFiles(txtFiles)
	for category, files := range categories {
		fmt.Printf("\n%s (%d 个文件):\n", category, len(files))
		for _, file := range files {
			domainCount := countDomainsInFile(filepath.Join(publishDir, file))
			fmt.Printf("  - %s (%d 个域名)\n", file, domainCount)
		}
	}

	// 显示所有文件列表
	fmt.Println("\n=== 所有源文件列表 ===")
	for i, file := range txtFiles {
		domainCount := countDomainsInFile(filepath.Join(publishDir, file))
		fmt.Printf("%3d. %s (%d 个域名)\n", i+1, file, domainCount)
	}
}

func countDomainsInFile(filepath string) int {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return 0
	}

	lines := strings.Split(string(content), "\n")
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "include:") {
			count++
		}
	}
	return count
}

func categorizeFiles(files []string) map[string][]string {
	categories := make(map[string][]string)

	for _, file := range files {
		category := getCategoryFromFilename(file)
		categories[category] = append(categories[category], file)
	}

	return categories
}

func getCategoryFromFilename(filename string) string {
	// 移除.txt后缀
	name := strings.TrimSuffix(filename, ".txt")

	// 根据文件名判断分类
	if strings.HasPrefix(name, "category-") {
		return "分类文件"
	} else if strings.Contains(name, "-ads") {
		return "广告相关"
	} else if strings.Contains(name, "-cn") {
		return "中国公司"
	} else if strings.Contains(name, "-ir") {
		return "伊朗相关"
	} else if strings.Contains(name, "-ru") {
		return "俄罗斯相关"
	} else if strings.Contains(name, "porn") || strings.Contains(name, "xvideos") ||
		strings.Contains(name, "youjizz") || strings.Contains(name, "boboporn") ||
		strings.Contains(name, "sehuatang") || strings.Contains(name, "lethalhardcore") {
		return "成人内容"
	} else {
		return "其他"
	}
}

// getListDomains 获取指定列表的所有域名
func getListDomains(listName string) {
	// 首先尝试从publish目录读取
	publishFilePath := filepath.Join(*outputPath, listName+".txt")
	dataFilePath := filepath.Join(*dataPath, listName)

	var filePath string
	var isFromData bool

	// 检查文件是否存在
	if _, err := os.Stat(publishFilePath); err == nil {
		filePath = publishFilePath
		isFromData = false
	} else if _, err := os.Stat(dataFilePath); err == nil {
		filePath = dataFilePath
		isFromData = true
	} else {
		fmt.Printf("错误: 文件 %s 或 %s 都不存在\n", publishFilePath, dataFilePath)
		fmt.Println("可用的列表文件:")
		showAvailableLists()
		return
	}

	// 递归获取所有域名（包括include引用的）
	allDomains, allComments, includedFiles := getDomainsRecursively(filePath)

	// 去重
	uniqueDomains := make(map[string]bool)
	var finalDomains []string
	for _, domain := range allDomains {
		if !uniqueDomains[domain] {
			uniqueDomains[domain] = true
			finalDomains = append(finalDomains, domain)
		}
	}

	// 显示结果
	fmt.Printf("=== %s 域名列表 ===\n", listName)
	fmt.Printf("文件路径: %s\n", filePath)
	if isFromData {
		fmt.Printf("数据源: 原始数据文件\n")
	} else {
		fmt.Printf("数据源: 生成的txt文件\n")
	}
	fmt.Printf("直接域名数量: %d\n", len(allDomains))
	fmt.Printf("去重后域名数量: %d\n", len(finalDomains))

	if len(allComments) > 0 {
		fmt.Printf("注释/引用数量: %d\n", len(allComments))
	}

	if len(includedFiles) > 0 {
		fmt.Printf("包含的文件数量: %d\n", len(includedFiles))
		fmt.Println("包含的文件:")
		for _, file := range includedFiles {
			fmt.Printf("  - %s\n", file)
		}
	}

	fmt.Println()

	// 显示域名列表
	if len(finalDomains) > 0 {
		fmt.Println("域名列表:")
		for i, domain := range finalDomains {
			fmt.Printf("%3d. %s\n", i+1, domain)
		}
	}

	// 显示注释和引用
	if len(allComments) > 0 {
		fmt.Println("\n注释和引用:")
		for _, comment := range allComments {
			fmt.Printf("  %s\n", comment)
		}
	}

	// 显示使用示例
	fmt.Println("\n=== 在 v2ray 中使用 ===")
	fmt.Printf("在 v2ray 配置文件中使用:\n")
	fmt.Printf("{\n")
	fmt.Printf("  \"routing\": {\n")
	fmt.Printf("    \"rules\": [\n")
	fmt.Printf("      {\n")
	fmt.Printf("        \"type\": \"field\",\n")
	fmt.Printf("        \"domain\": [\"geosite:%s\"],\n", listName)
	fmt.Printf("        \"outboundTag\": \"proxy\"\n")
	fmt.Printf("      }\n")
	fmt.Printf("    ]\n")
	fmt.Printf("  }\n")
	fmt.Printf("}\n")
}

// getDomainsRecursively 递归获取所有域名，包括include引用的
func getDomainsRecursively(filePath string) ([]string, []string, []string) {
	var allDomains []string
	var allComments []string
	var includedFiles []string

	// 防止循环引用
	processedFiles := make(map[string]bool)

	var processFile func(string) error
	processFile = func(currentPath string) error {
		if processedFiles[currentPath] {
			return nil // 避免循环引用
		}
		processedFiles[currentPath] = true

		// 读取文件内容
		content, err := os.ReadFile(currentPath)
		if err != nil {
			return fmt.Errorf("无法读取文件 %s: %v", currentPath, err)
		}

		// 解析域名
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "#") {
				allComments = append(allComments, line)
				continue
			}

			if strings.HasPrefix(line, "include:") {
				allComments = append(allComments, line)
				// 处理include引用
				includedFile := strings.TrimPrefix(line, "include:")
				// 移除属性标记
				if idx := strings.Index(includedFile, ":"); idx != -1 {
					includedFile = includedFile[:idx]
				}

				// 首先尝试在publish目录中查找
				includedPath := filepath.Join(*outputPath, includedFile+".txt")
				if _, err := os.Stat(includedPath); err == nil {
					includedFiles = append(includedFiles, includedFile)
					// 递归处理include的文件
					if err := processFile(includedPath); err != nil {
						fmt.Printf("警告: 处理include文件 %s 时出错: %v\n", includedPath, err)
					}
					continue
				}

				// 如果publish目录中没有，尝试在data目录中查找
				dataIncludedPath := filepath.Join(*dataPath, includedFile)
				if _, err := os.Stat(dataIncludedPath); err == nil {
					includedFiles = append(includedFiles, includedFile)
					// 递归处理include的文件
					if err := processFile(dataIncludedPath); err != nil {
						fmt.Printf("警告: 处理include文件 %s 时出错: %v\n", dataIncludedPath, err)
					}
					continue
				}

				fmt.Printf("警告: 找不到include文件 %s\n", includedFile)
				continue
			}

			// 解析域名
			if strings.HasPrefix(line, "domain:") {
				domain := strings.TrimPrefix(line, "domain:")
				// 移除属性标记
				if idx := strings.Index(domain, ":"); idx != -1 {
					domain = domain[:idx]
				}
				allDomains = append(allDomains, domain)
			} else if strings.HasPrefix(line, "full:") {
				domain := strings.TrimPrefix(line, "full:")
				// 移除属性标记
				if idx := strings.Index(domain, ":"); idx != -1 {
					domain = domain[:idx]
				}
				allDomains = append(allDomains, domain)
			} else if strings.HasPrefix(line, "regexp:") {
				allComments = append(allComments, line)
			} else {
				// 假设是普通域名（检查是否有属性标记）
				if idx := strings.Index(line, " @"); idx != -1 {
					domain := line[:idx]
					allDomains = append(allDomains, domain)
				} else {
					allDomains = append(allDomains, line)
				}
			}
		}

		return nil
	}

	// 开始处理
	if err := processFile(filePath); err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	return allDomains, allComments, includedFiles
}

// showAvailableLists 显示可用的列表文件
func showAvailableLists() {
	publishDir := *outputPath
	files, err := os.ReadDir(publishDir)
	if err != nil {
		fmt.Printf("错误: 无法读取目录 %s: %v\n", publishDir, err)
		return
	}

	var availableLists []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			// 过滤掉以"._"开头的临时文件
			if strings.HasPrefix(file.Name(), "._") {
				continue
			}

			listName := strings.TrimSuffix(file.Name(), ".txt")
			availableLists = append(availableLists, listName)
		}
	}

	// 按分类显示
	categories := categorizeFiles(availableLists)
	for category, files := range categories {
		fmt.Printf("\n%s (%d 个文件):\n", category, len(files))
		for _, file := range files {
			fmt.Printf("  - %s\n", file)
		}
	}

	fmt.Printf("\n使用方法: go run . -getlist <列表名称>\n")
	fmt.Printf("示例: go run . -getlist tiktok\n")
	fmt.Printf("示例: go run . -getlist category-ads-all\n")
}

// getAllAdvertisingDomains 获取所有广告域名
func getAllAdvertisingDomains() {
	fmt.Println("=== 所有广告域名汇总 ===")

	publishDir := *outputPath
	files, err := os.ReadDir(publishDir)
	if err != nil {
		fmt.Printf("错误: 无法读取目录 %s: %v\n", publishDir, err)
		return
	}

	var adFiles []string
	var allDomains []string
	var totalFiles int
	var totalDomains int

	// 收集所有 -ads.txt 文件
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			// 过滤掉以"._"开头的临时文件
			if strings.HasPrefix(file.Name(), "._") {
				continue
			}

			if strings.Contains(file.Name(), "-ads") {
				adFiles = append(adFiles, file.Name())
			}
		}
	}

	fmt.Printf("找到 %d 个广告相关文件\n", len(adFiles))
	fmt.Println()

	// 处理每个广告文件
	for _, fileName := range adFiles {
		filePath := filepath.Join(publishDir, fileName)
		domains, _, _ := getDomainsRecursively(filePath)

		if len(domains) > 0 {
			fmt.Printf("%s: %d 个域名\n", fileName, len(domains))
			allDomains = append(allDomains, domains...)
			totalFiles++
			totalDomains += len(domains)
		}
	}

	// 去重
	uniqueDomains := make(map[string]bool)
	var finalDomains []string
	for _, domain := range allDomains {
		if !uniqueDomains[domain] {
			uniqueDomains[domain] = true
			finalDomains = append(finalDomains, domain)
		}
	}

	fmt.Println()
	fmt.Printf("=== 汇总统计 ===\n")
	fmt.Printf("处理的文件数量: %d\n", totalFiles)
	fmt.Printf("原始域名总数: %d\n", totalDomains)
	fmt.Printf("去重后域名总数: %d\n", len(finalDomains))
	fmt.Printf("重复域名数量: %d\n", totalDomains-len(finalDomains))

	// 显示所有域名
	if len(finalDomains) > 0 {
		fmt.Println("\n=== 所有广告域名列表 ===")
		for i, domain := range finalDomains {
			fmt.Printf("%4d. %s\n", i+1, domain)
		}
	}

	// 显示使用示例
	fmt.Println("\n=== 在 v2ray 中使用 ===")
	fmt.Printf("在 v2ray 配置文件中使用:\n")
	fmt.Printf("{\n")
	fmt.Printf("  \"routing\": {\n")
	fmt.Printf("    \"rules\": [\n")
	fmt.Printf("      {\n")
	fmt.Printf("        \"type\": \"field\",\n")
	fmt.Printf("        \"domain\": [\n")

	// 显示所有广告分类
	for i, fileName := range adFiles {
		listName := strings.TrimSuffix(fileName, ".txt")
		if i == 0 {
			fmt.Printf("          \"geosite:%s\"", listName)
		} else {
			fmt.Printf(",\n          \"geosite:%s\"", listName)
		}
	}

	fmt.Printf("\n        ],\n")
	fmt.Printf("        \"outboundTag\": \"block\"\n")
	fmt.Printf("      }\n")
	fmt.Printf("    ]\n")
	fmt.Printf("  }\n")
	fmt.Printf("}\n")
}
