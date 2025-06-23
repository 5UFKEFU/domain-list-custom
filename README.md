# 简介

运行 go run . -mode multi 
就可以在目录里生成相关 geosite 文件

以下python 文件可以获得内容:
jeffcheng@JEFFs-MacBook-Air-M2-15 domain-list-custom % python3 parse_geosite.py ./publish/geosite.dat      
=== GeoSite.dat 文件信息 ===
文件路径: ./publish/geosite.dat
文件大小: 720322 字节 (703.44 KB)
分类总数: 30

=== 所有分类列表 ===
  1. CATEGORY-PORN                  (12013 个域名)
  2. CN                             (5206 个域名)
  3. APPLE                          (4220 个域名)
  4. GOOGLE                         (2417 个域名)
  5. META                           (1636 个域名)
  6. MICROSOFT                      (1465 个域名)
  7. CATEGORY-ADS-ALL               ( 896 个域名)
  8. HUAWEI                         ( 876 个域名)
  9. BYTEDANCE                      ( 796 个域名)
 10. ALIBABA                        ( 518 个域名)
 11. CATEGORY-RU                    ( 491 个域名)
 12. JD                             ( 474 个域名)
 13. BAIDU                          ( 458 个域名)
 14. AMAZON                         ( 445 个域名)
 15. PRIVATE                        ( 252 个域名)
 16. TENCENT                        ( 201 个域名)
 17. SINA                           ( 193 个域名)
 18. CATEGORY-AI-!CN                ( 173 个域名)
 19. XIAOMI                         ( 110 个域名)
 20. BILIBILI                       (  83 个域名)
 21. NETEASE                        (  78 个域名)
 22. OPPO                           (  68 个域名)
 23. CATEGORY-AI-CN                 (  54 个域名)
 24. IQIYI                          (  53 个域名)
 25. SOHU                           (  52 个域名)
 26. CATEGORY-IR                    (  49 个域名)
 27. TIKTOK                         (  21 个域名)
 28. DIDI                           (  18 个域名)
 29. RUANMEI                        (  16 个域名)
 30. MEITUAN                        (  14 个域名)
 31. YOUKU                          (  12 个域名)
 32. VIVO                           (   9 个域名)
 

使用方法:
  python parse_geosite.py ./publish/geosite.dat <分类名>  # 查看指定分类的所有域名
  示例: python parse_geosite.py ./publish/geosite.dat tiktok

*结果示例:*
 
*jeffcheng@JEFFs-MacBook-Air-M2-15 domain-list-custom % python3 parse_geosite.py ./publish/geosite.dat 'CATEGORY-AI-!CN'*

=== CATEGORY-AI-!CN 域名列表 ===
文件路径: ./publish/geosite.dat
域名数量: 118

=== Full 类型域名 (13 个) ===
   1. ai.google.dev
   2. alkalicore-pa.clients6.google.com
   3. alkalimakersuite-pa.clients6.google.com
   4. webchannel-alkalimakersuite-pa.clients6.google.com
   5. openaiapi-site.azureedge.net
   6. openaicom-api-bdcpf8c6d2e9atf6.z01.azurefd.net
   7. openaicomproductionae4b.blob.core.windows.net
   8. production-openaicom-storage.azureedge.net
   9. pplx-res.cloudinary.com
  10. servd-anthropic-website.b-cdn.net
  11. o33249.ingest.sentry.io
  12. openaicom.imgix.net
  13. browser-intake-datadoghq.com

=== RootDomain 类型域名 (104 个) ===
   1. hf.co
   2. oaistatic.com
   3. clipdrop.co
   4. jasper.ai
   5. chutes.ai
   6. dify.ai
   7. claudeusercontent.com
   8. meta.ai
   9. mistral.ai
  10. openart.ai
  11. openrouter.ai
  12. gemini.google
  13. groq.com
  14. jules.google
  15. deepmind.com
  16. notebooklm.google
  17. generativeai.google
  18. claude.ai
  19. anthropic.com
  20. x.ai
  21. grok.com
  22. pplx.ai
  23. perplexity.ai
  24. sora.com
  25. openai.com
  26. oaiusercontent.com
  27. deepmind.google
  28. coderabbit.ai
  29. huggingface.co
  30. poe.com
  31. poecdn.net
  32. cursor-cdn.com
  33. cursor.com
  34. cursor.sh
  35. cursorapi.com
  36. chatgpt.com
  37. chat.com
  38. proactivebackend-pa.googleapis.com
  39. makersuite.google.com
  40. generativelanguage.googleapis.com
  41. notebooklm.google.com
  42. chatgpt.livekit.cloud
  43. geller-pa.googleapis.com
  44. turn.livekit.cloud
  45. jules.google.com
  46. host.livekit.cloud
  47. gemini.google.com
  48. aistudio.google.com
  49. bard.google.com
  50. coderabbit.gallery.vsassets.io
  51. gateway.ai.cloudflare.com
  52. openai.com.cdn.cloudflare.net
  53. hf.co
  54. oaistatic.com
  55. clipdrop.co
  56. jasper.ai
  57. chutes.ai
  58. dify.ai
  59. claudeusercontent.com
  60. meta.ai
  61. mistral.ai
  62. openart.ai
  63. openrouter.ai
  64. gemini.google
  65. groq.com
  66. jules.google
  67. deepmind.com
  68. notebooklm.google
  69. generativeai.google
  70. claude.ai
  71. anthropic.com
  72. x.ai
  73. grok.com
  74. pplx.ai
  75. perplexity.ai
  76. sora.com
  77. openai.com
  78. oaiusercontent.com
  79. deepmind.google
  80. coderabbit.ai
  81. huggingface.co
  82. poe.com
  83. poecdn.net
  84. cursor-cdn.com
  85. cursor.com
  86. cursor.sh
  87. cursorapi.com
  88. chatgpt.com
  89. chat.com
  90. proactivebackend-pa.googleapis.com
  91. makersuite.google.com
  92. generativelanguage.googleapis.com
  93. notebooklm.google.com
  94. chatgpt.livekit.cloud
  95. geller-pa.googleapis.com
  96. turn.livekit.cloud
  97. jules.google.com
  98. host.livekit.cloud
  99. gemini.google.com
 100. aistudio.google.com
 101. bard.google.com
 102. coderabbit.gallery.vsassets.io
 103. gateway.ai.cloudflare.com
 104. openai.com.cdn.cloudflare.net
