package services

import (
	"douyin-simple/models"
	"douyin-simple/units"
	"douyin-simple/utils"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

const (
	TokenKey          = "66557773"
	TokenPreKey       = "DOUYIN"
	TokenPreKeyLength = 6
	TokenExpiration   = 24 * 60 * 60 * 1000
)

type TokenData struct {
	User     models.User `json:"user"`
	CreateAt time.Time   `json:"create_at"`
}

func CreateToken(user models.User) (string, error) {
	var tokenData = TokenData{
		User:     user,
		CreateAt: time.Now(),
	}
	// 生成序列化数据
	data, err := json.Marshal(tokenData)
	if err != nil {
		log.Println("MakeToken中在序列化takenData时出错：", err)
		utils.PrintLog(err, "[Error]")
		return "", err
	}
	// 拼接token字符串
	deToken := strings.Join([]string{TokenPreKey, string(data)}, "")
	// 加密字符串
	des, err := units.EncryptDes(deToken, []byte(TokenKey))
	if err != nil {
		log.Println("EncryptDes加密过程出错：", err)
		utils.PrintLog(err, "[Error]")
		return "", err
	}
	return des, nil

}

// ResolveToken 解析token。若error不为nil，则token无效
func ResolveToken(token string) (TokenData, error) {
	// 解码token
	des, err := units.DecryptDes(token, []byte(TokenKey))
	if err != nil {
		log.Println("token进行des解码时失败", err)
		utils.PrintLog(err, "[Info]")
		return TokenData{}, err
	}

	// 校验token前缀符
	if len(des) <= TokenPreKeyLength {
		err = errors.New("无效的token:" + des)
		log.Println("token前缀符错误:", des)
		utils.PrintLog(err, "[Info]")
		return TokenData{}, err
	} else {
		sli := des[0:TokenPreKeyLength]
		if strings.Compare(sli, TokenPreKey) != 0 {
			err = errors.New("无效的token:" + des)
			log.Println("token前缀符错误:", des)
			utils.PrintLog(err, "[Info]")
			return TokenData{}, err
		}
	}

	// 解析token
	var tokenData TokenData
	err = json.Unmarshal([]byte(des[TokenPreKeyLength:]), &tokenData)
	if err != nil {
		log.Println("token数据反序列化失败:", des)
		utils.PrintLog(err, "[Info]")
		return TokenData{}, errors.New("token解析失败")
	}
	//// 判断token是否过期
	//dif := time.Since(time.Now()) - time.Since(tokenData.CreateAt)
	//if dif > TokenExpiration {
	//	err = errors.New("无效的token:" + des)
	//	log.Println("已过期的token:", token)
	//	utils.PrintLog(err, "[Info]")
	//	return TokenData{}, err
	//}
	return tokenData, nil
}
