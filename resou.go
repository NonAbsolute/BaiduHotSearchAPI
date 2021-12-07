package main

import (
	"fmt"
	"net/http"
	// "html/template"
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Value struct {
	Index     int    `json:"index"`
	WordQuery string `json:"wordQuery"`
	Desc      string `json:"desc"`
	Img       string `json:"img"`
	RawUrl    string `json:"rawUrl"`
	HotScore  int    `json:"hotScore"`
	HotTag    int    `json:"hotTag"`
	HotTagImg string `json:"hotTagImg"`
}
type Data struct {
	Text  string  `json:"text"`
	Value []Value `json:"value"`
}
type Datavalue struct {
	Data    []Data `json:"data"`
	Retcode int    `json:"retcode"`
	Status  string `json:"status"`
	Wording string `json:"wording"`
}

func grabReSouData(sum int) Datavalue {
	//抓取热榜数据
	resp, err := http.Get("https://top.baidu.com/board?tab=realtime")
	if err != nil {
		// handle error
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		panic(err)
	}

	//解析正则表达式，如果成功返回解释器
	reg := regexp.MustCompile(`<\!--s-data:({.*})-->{1}`)
	if reg == nil {
		fmt.Println("regexp err")
		panic(reg)
	}
	//根据规则提取关键信息
	result := reg.FindAllStringSubmatch(string(body), -1)
	text := gjson.Get(result[0][1], "data.cards.0.text").Str //热榜名
	var d Data
	for i := 0; ; i++ {
		index := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".index") //排名
		if !index.Exists() || i >= sum {
			break
		} else {
			wordQuery := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".word").Str      //描述
			desc := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".desc").Str           //摘要，为空则无摘要
			img := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".img").Str             //缩略图
			rawUrl := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".rawUrl").Str       //详情链接
			hotScore := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".hotScore")       //热度指数
			hotTag := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".hotTag")           //标签，为0则无标签 /*	3 热		2 商		1 新		0 无	*/
			hotTagImg := gjson.Get(result[0][1], "data.cards.0.content."+strconv.Itoa(i)+".hotTagImg").Str //标签图标，如果标签为0则此项为空
			d.Value = append(d.Value, Value{
				Index:     int(index.Int()),
				WordQuery: wordQuery,
				Desc:      desc,
				Img:       img,
				RawUrl:    rawUrl,
				HotScore:  int(hotScore.Int()),
				HotTag:    int(hotTag.Int()),
				HotTagImg: hotTagImg,
			})
		}
	}
	d.Text = text
	dv := Datavalue{
		Data: []Data{
			d,
		},
		Retcode: 200,
		Status:  "success",
		Wording: "ok",
	}
	return dv
}

func handleReSouGetIsPost(writer http.ResponseWriter, request *http.Request) {
	// 处理application/json类型的POST请求,根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(request.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	if params["sum"] != "" {
		fmt.Printf("POST json: sum=%s\n", params["sum"])
		sum, err := strconv.Atoi(params["sum"])
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(`{"data":null,"retcode":500,"status":"failed","wording":"参数错误"}`))
		} else {
			dv := grabReSouData(sum)
			json.NewEncoder(writer).Encode(dv)
		}
		// fmt.Fprintf(writer, `{"code":0}`)
	} else {
		//接收GET请求
		query := request.URL.Query()
		if query.Get("sum") != "" {
			fmt.Printf("GET: sum=%s\n", query.Get("sum"))
			sum, err := strconv.Atoi(query.Get("sum"))
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(`{"data":null,"retcode":500,"status":"failed","wording":"参数错误"}`))
			} else {
				dv := grabReSouData(sum)
				json.NewEncoder(writer).Encode(dv)
			}
			// fmt.Fprintf(writer, `{"code":0}`)
		} else {
			dv := grabReSouData(999)
			json.NewEncoder(writer).Encode(dv)
		}
	}
}

//默认首页
func handleIndex(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(`{"data":null,"retcode":404,"status":"failed","wording":"API不存在"}`))

	//用来返回html的
	// t, _ := template.ParseFiles("index.html")
	// t.Execute(writer, nil)
}

//CLI帮助
func inCliValue() (string, string) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p",
			Value: "8082",
			Usage: "server port",
		},
		cli.StringFlag{
			Name:  "httpsP",
			Value: "8083",
			Usage: "https server port",
		},
	}
	p := ""
	sp := ""
	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			fmt.Println("您想做什么？")
		}
		if c.String("p") != "" {
			p = c.String("p")
			sp = c.String("httpsP")
		} else {
			p = "8082"
			sp = "8083"
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if p == "" || sp == "" {
		os.Exit(0)
	}
	return p, sp
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/resou", handleReSouGetIsPost)

	port, Sport := inCliValue()
	fmt.Println("Running at port " + port + ",https port " + Sport + " ...")
	fmt.Println("Other options input -h")
	err := http.ListenAndServe(":"+port, nil)
	//未处理https请求
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
