# BaiduHotSearchAPI
百度热搜API，继[在你的网站上添加百度热搜榜一栏](https://gitee.com/deng_wenyi/Crawling_BaiDu_Hot_search_list)
By FJD团队-DengXianSheng
感谢百度风云榜提供数据支持！！
**本项目仅交流学习**

### 请求类型

GET/POST(JSON)

### 参数

| 参数值 | 类型 | 释义                     | 是否必填 |
| ------ | ---- | ------------------------ | -------- |
| sum    | int  | 返回的数量，默认返回全部 | 否       |

### 返回值

类型：JSON

| 参数      | 释义                                |
| --------- | ----------------------------------- |
| status    | 状态，failed或success               |
| rawUrl    | 详情链接                            |
| wordQuery | 描述                                |
| img       | 缩略图                              |
| hotTag    | 标签，3 热、2 商、1 新，为0则无标签 |
| hotScore  | 热度指数                            |
| hotTagImg | 标签图标，如果标签为0则此项为空     |
| desc      | 摘要，为空则无摘要                  |
| index     | 排名                                |
| wording   | 返回结果的释义                      |
| retcode   | HTTP状态码                          |



```json
{
  "status" : "success",
  "data" : [
    {
      "value" : [
        {
          "rawUrl" : "",
          "wordQuery" : "",
          "img" : "",
          "hotTag" : 3,
          "hotScore" : 4944936,
          "hotTagImg" : "",
          "desc" : "",
          "index" : 0
        },
        {
          "rawUrl" : "",
          "wordQuery" : "",
          "img" : "",
          "hotTag" : 3,
          "hotScore" : 4812168,
          "hotTagImg" : "",
          "desc" : "",
          "index" : 1
        }
      ],
      "text" : "热搜榜"
    }
  ],
  "wording" : "ok",
  "retcode" : 200
}
```
