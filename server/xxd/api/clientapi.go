/**
 * The aes file of api current module of xxd.
 *
 * @copyright   Copyright 2009-2017 青岛易软天创网络科技有限公司(QingDao Nature Easy Soft Network Technology Co,LTD, www.cnezsoft.com)
 * @license     ZPL (http://zpl.pub/page/zplv12.html)
 * @author      Archer Peng <pengjiangxiu@cnezsoft.com>
 * @package     util
 * @link        http://www.zentao.net
 */
package api

import (
	"xxd/hyperttp"
	"xxd/util"
)

var newline = []byte{'\n'}

// 从客户端发来的登录请求，通过该函数转发到后台服务器进行登录验证
func ChatLogin(clientData ParseData) ([]byte, int64, bool) {
	ranzhiServer, ok := util.Config.RanzhiServer[clientData.ServerName()]
	if !ok {
		util.LogError().Println("no ranzhi server name")
		return nil, -1, false
	}

	// 到http服务器请求，返回加密的结果
	retMessage, err := hyperttp.RequestInfo(ranzhiServer.RanzhiAddr, ApiUnparse(clientData, ranzhiServer.RanzhiToken))
	if err != nil || retMessage == nil {
		util.LogError().Println("hyperttp request info error:", err)
		return nil, -1, false
	}

	// 解析http服务器的数据,返回 ParseData 类型的数据
	retData, err := ApiParse(retMessage, ranzhiServer.RanzhiToken)
	if err != nil {
		util.LogError().Println("api parse error:", err)
		return nil, -1, false
	}

	retMessage, err = SwapToken(retMessage, ranzhiServer.RanzhiToken, util.Token)
	if err != nil {
		return nil, -1, false
	}

	// 返回值：
	// 1、返回给客户端加密后的数据
	// 2、返回用户的ID
	// 3、返回登录的结果
	return retMessage, retData.UserID(), retData.Result() == "success"
}

func ChatLogout(serverName string, userID int64) ([]byte, []int64, error) {
	ranzhiServer, ok := util.Config.RanzhiServer[serverName]
	if !ok {
		util.LogError().Println("no ranzhi server name")
		return nil, nil, util.Errorf("%s\n", "no ranzhi server name")
	}

	request := []byte(`{"module":"chat","method":"logout",userID:` + util.Int642String(userID) + `}`)
	message, err := aesEncrypt(request, ranzhiServer.RanzhiToken)
	if err != nil {
		util.LogError().Println("aes encrypt error:", err)
		return nil, nil, err
	}

	// 到http服务器请求user get list数据
	r2xMessage, err := hyperttp.RequestInfo(ranzhiServer.RanzhiAddr, message)
	if err != nil {
		util.LogError().Println("hyperttp request info error:", err)
		return nil, nil, err
	}

	// 解析http服务器的数据,返回 ParseData 类型的数据
	parseData, err := ApiParse(r2xMessage, ranzhiServer.RanzhiToken)
	if err != nil {
		util.LogError().Println("api parse error", err)
		return nil, nil, err
	}

	sendUsers := parseData.SendUsers()

	x2cMessage := ApiUnparse(parseData, util.Token)
	if x2cMessage == nil {
		return nil, nil, err
	}

	return x2cMessage, sendUsers, nil

}

func RepeatLogin() []byte {
	repeatLogin := []byte(`{module:  'null',method:  'null',message: 'This account logined in another place.'}`)
	repeatLogin = append(repeatLogin, newline...)

	message, err := aesEncrypt(repeatLogin, util.Token)
	if err != nil {
		util.LogError().Println("aes encrypt error:", err)
		return nil
	}

	return message
}

func TestLogin() []byte {
	loginData := []byte(`{"result":"success","data":{"id":12,"account":"demo8","realname":"\u6210\u7a0b\u7a0b","avatar":"","role":"hr","dept":0,"status":"online","admin":"no","gender":"f","email":"ccc@demo.com","mobile":"","site":"","phone":""},"sid":"18025976a786ec78194e491e7b790731","module":"chat","method":"login"}`)

	loginData = append(loginData, newline...)
	message, err := aesEncrypt(loginData, util.Token)
	if err != nil {
		util.LogError().Println("aes encrypt error:", err)
		return nil
	}

	return message
}

func TransitData(clientData []byte, serverName string) ([]byte, []int64, error) {
	ranzhiServer, ok := util.Config.RanzhiServer[serverName]
	if !ok {
		util.LogError().Println("no ranzhi server name")
		return nil, nil, util.Errorf("%s\n", "no ranzhi server name")
	}

	message, err := SwapToken(clientData, util.Token, ranzhiServer.RanzhiToken)
	if err != nil {
		return nil, nil, err
	}

	// ranzhi to xxd message
	r2xMessage, err := hyperttp.RequestInfo(ranzhiServer.RanzhiAddr, message)
	if err != nil {
		util.LogError().Println("hyperttp request info error:", err)
		return nil, nil, err
	}

	parseData, err := ApiParse(r2xMessage, ranzhiServer.RanzhiToken)
	if err != nil {
		util.LogError().Println("api parse error:", err)
		return nil, nil, err
	}

	sendUsers := parseData.SendUsers()

	// xxd to client message
	x2cMessage := ApiUnparse(parseData, util.Token)
	if x2cMessage == nil {
		return nil, nil, err
	}

	return x2cMessage, sendUsers, nil
}

func UserGetlist(serverName string, userID int64) ([]byte, error) {
	ranzhiServer, ok := util.Config.RanzhiServer[serverName]
	if !ok {
		util.LogError().Println("no ranzhi server name")
		return nil, util.Errorf("%s\n", "no ranzhi server name")
	}

	// 固定的json格式
	request := []byte(`{"module":"chat","method":"userGetlist",userID:` + util.Int642String(userID) + `}`)
	message, err := aesEncrypt(request, ranzhiServer.RanzhiToken)
	if err != nil {
		util.LogError().Println("aes encrypt error:", err)
		return nil, err
	}

	// 到http服务器请求user get list数据
	retMessage, err := hyperttp.RequestInfo(ranzhiServer.RanzhiAddr, message)
	if err != nil {
		util.LogError().Println("hyperttp request info error:", err)
		return nil, err
	}

	//由于http服务器和客户端的token不一致，所以需要进行交换
	retData, err := SwapToken(retMessage, ranzhiServer.RanzhiToken, util.Token)
	if err != nil {
		return nil, err
	}

	return retData, nil
}

func Getlist(serverName string, userID int64) ([]byte, error) {
	ranzhiServer, ok := util.Config.RanzhiServer[serverName]
	if !ok {
		util.LogError().Println("no ranzhi server name")
		return nil, util.Errorf("%s\n", "no ranzhi server name")
	}

	// 固定的json格式
	request := []byte(`{"module":"chat","method":"getlist",userID:` + util.Int642String(userID) + `}`)
	message, err := aesEncrypt(request, ranzhiServer.RanzhiToken)
	if err != nil {
		util.LogError().Println("aes encrypt error:", err)
		return nil, err
	}

	// 到http服务器请求get list数据
	retMessage, err := hyperttp.RequestInfo(ranzhiServer.RanzhiAddr, message)
	if err != nil {
		util.LogError().Println("hyperttp request info error:", err)
		return nil, err
	}

	//由于http服务器和客户端的token不一致，所以需要进行交换
	retData, err := SwapToken(retMessage, ranzhiServer.RanzhiToken, util.Token)
	if err != nil {
		return nil, err
	}

	return retData, nil

}

func (pd ParseData) ServerName() string {
	params, ok := pd["params"]
	if !ok {
		return ""
	}

	// api中server name在数组固定位置为0
	ret := params.([]interface{})
	return ret[0].(string)
}

func (pd ParseData) Account() string {
	params, ok := pd["params"]
	if !ok {
		return ""
	}

	// api中account在数组固定位置为1
	ret := params.([]interface{})
	return ret[1].(string)
}

func (pd ParseData) Password() string {
	params, ok := pd["params"]
	if !ok {
		return ""
	}

	// api中password在数组固定位置为2
	ret := params.([]interface{})
	return ret[2].(string)
}

func (pd ParseData) Status() string {
	params, ok := pd["params"]
	if !ok {
		return ""
	}

	// api中status在数组固定位置为3
	ret := params.([]interface{})
	return ret[3].(string)
}
