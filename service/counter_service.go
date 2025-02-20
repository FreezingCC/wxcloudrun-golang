package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

var SECRETKEY = os.Getenv("APP_SECRET_KEY")

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

const (
	success = 200
	failed  = 500
)

type GptResponse struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Choices []Message `json:"choices"`
}

type Message struct {
	Text         string         `json:"text"`
	Index        int            `json:"index"`
	FinishReason string         `json:"finish_reason"`
	Usage        map[string]int `json:"usage"`
}

// IndexHandler 计数器接口
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getIndex()
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	fmt.Fprint(w, data)
}

// CounterHandler 计数器接口
func CounterHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	if r.Method == http.MethodGet {
		counter, err := getCurrentCounter()
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = counter.Count
		}
	} else if r.Method == http.MethodPost {
		count, err := modifyCounter(r)
		if err != nil {
			res.Code = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = count
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

type scoreUploadReq struct {
	UserId string `json:"user_id"`
	Score  int32  `json:"score"`
}

func UploadScore(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{
		Code:     0,
		ErrorMsg: "",
		Data:     nil,
	}
	var req scoreUploadReq
	var bd []byte
	bd, err := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(bd, &req)

	err = dao.UserMaxScore.UpdateScoreByUserId(req.UserId, req.Score)
	if err != nil {
		msg, _ := json.Marshal(res)
		w.Header().Set("content-type", "application/json")
		w.Write(msg)
		return
	}
	res.Code = success
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

func GetScore(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{
		Code:     0,
		ErrorMsg: "",
		Data:     nil,
	}
	var req scoreUploadReq
	var bd []byte
	bd, err := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(bd, &req)

	userScore, err := dao.UserMaxScore.GetScoreByUserId(req.UserId)
	if err != nil {
		msg, _ := json.Marshal(res)
		w.Header().Set("content-type", "application/json")
		w.Write(msg)
		return
	}
	res.Code = success
	res.Data = map[string]int32{"score": userScore.Score}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

type request struct {
	Message string `json:"message"`
}

func Chat(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	var req request
	var bd []byte
	_, _ = r.Body.Read(bd)
	_ = json.Unmarshal(bd, &req)

	//if r.Method == http.MethodPost {
	//	client := resty.New()
	//	apiKey := "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"
	//	url := "https://api.openai.com/v1/completions"
	//
	//	response, err := client.R().
	//		SetHeader("Content-Type", "application/json").
	//		SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
	//		SetBody(map[string]interface{}{
	//			"model":       "text-davinci-003",
	//			"prompt":      req.Message,
	//			"max_tokens":  1024,
	//			"temperature": 0.5,
	//			"top_p":       1,
	//			"n":           1,
	//			"stream":      false,
	//		}).
	//		Post(url)
	//
	//	if err != nil {
	//		fmt.Printf("Error: %s\n", err.Error())
	//	}
	//	var gresp GptResponse
	//	err = json.Unmarshal(response.Body(), &gresp)
	//	if err != nil {
	//		return
	//	}
	//	if len(gresp.Choices) == 0 {
	//		return
	//	}
	//	if err != nil {
	//		res.Code = -1
	//		res.ErrorMsg = err.Error()
	//	} else {
	//		res.Data = gresp.Choices[0].Text
	//	}
	//	//res.Data = "我是你爹"
	//} else {
	//	res.Code = -1
	//	res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	//}

	res.Data = "我是你爹"
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

type JsCodeReq struct {
	Code string `json:"code"`
}

func GetUserId(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{
		Code:     0,
		ErrorMsg: "",
		Data:     nil,
	}
	var req JsCodeReq
	var bd []byte
	bd, _ = ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(bd, &req) // 从请求中获取登录凭证 code

	// 构建 GET 请求 URL
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=wxb3fa7c7530beebe5&secret="+SECRETKEY+"&js_code=%s&grant_type=authorization_code", req.Code)

	// 发送 GET 请求到微信服务器
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "登录失败", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// 读取响应体数据
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "登录失败", http.StatusInternalServerError)
		return
	}

	// 解析 JSON 响应
	var data struct {
		OpenID string `json:"openid"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "登录失败", http.StatusInternalServerError)
		return
	}
	res.Data = map[string]string{"openid": data.OpenID}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// modifyCounter 更新计数，自增或者清零
func modifyCounter(r *http.Request) (int32, error) {
	action, err := getAction(r)
	if err != nil {
		return 0, err
	}

	var count int32
	if action == "inc" {
		count, err = upsertCounter(r)
		if err != nil {
			return 0, err
		}
	} else if action == "clear" {
		err = clearCounter()
		if err != nil {
			return 0, err
		}
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// upsertCounter 更新或修改计数器
func upsertCounter(r *http.Request) (int32, error) {
	currentCounter, err := getCurrentCounter()
	var count int32
	createdAt := time.Now()
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		count = 1
		createdAt = time.Now()
	} else {
		count = currentCounter.Count + 1
		createdAt = currentCounter.CreatedAt
	}

	counter := &model.CounterModel{
		Id:        1,
		Count:     count,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = dao.Imp.UpsertCounter(counter)
	if err != nil {
		return 0, err
	}
	return counter.Count, nil
}

func clearCounter() error {
	return dao.Imp.ClearCounter(1)
}

// getCurrentCounter 查询当前计数器
func getCurrentCounter() (*model.CounterModel, error) {
	counter, err := dao.Imp.GetCounter(1)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

// getAction 获取action
func getAction(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", err
	}
	defer r.Body.Close()

	action, ok := body["action"]
	if !ok {
		return "", fmt.Errorf("缺少 action 参数")
	}

	return action.(string), nil
}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
