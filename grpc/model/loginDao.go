package model

import (
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}

func generate_Incr(client *redis.Client, incrName string) (incr_Key string) {
	result, err := client.Incr("incrName").Result()
	if err != nil {
		panic(err)
	}

	incr_Key = incrName + ":" + strconv.FormatInt(result, 10)

	return incr_Key
}

func get_lately_Incr(client *redis.Client, incrName string) (incr_Key string) {
	incr_Key, err := client.Get(incrName).Result()
	if err != nil {
		panic(err)
	}

	return incr_Key
}

func password_md5(password, salt string) string {
	h := md5.New()
	h.Write([]byte(salt + password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type User struct {
	userName string
	email    string
	password string
	iphone   string
	other    string
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func register_user(client *redis.Client, user User) string {
	if user.userName == "" || user.password == "" {
		return "用户名或者密码不能为空"
	}
	if VerifyEmailFormat(user.email) == false {
		return "邮箱格式不对"
	}

	nameKey := "userName:" + strings.ToLower(strings.TrimSpace(user.userName))
	cmd, err := client.Exists(nameKey).Result()
	if err != nil {
		panic(err)
	}

	if 0 != cmd {
		return "该用户名已经被注册"
	}

	emailkey := "email:" + user.email
	cmd, err = client.Exists(emailkey).Result()
	if err != nil {
		panic(err)
	}

	if 0 != cmd {
		return "该邮箱已经被注册"
	}

	uid_Key := generate_Incr(client, "userID")
	sta, err := client.Set(nameKey, uid_Key, 0).Result()
	if err != nil {
		panic(err)
	}
	if sta != "OK" {
		panic(sta)
	}
	sta, err = client.Set(emailkey, uid_Key, 0).Result()
	if err != nil {
		panic(err)
	}
	if sta != "OK" {
		panic(sta)
	}

	mess := map[string]interface{}{
		"userName":     user.userName,
		"email":        user.email,
		"password":     user.password,
		"registertime": time.Now().Format("2006-01-02 15:04:05"),
	}

	if user.iphone != "" {
		mess["iphone"] = user.iphone
	}
	if user.other != "" {
		mess["other"] = user.other
	}

	_, err = client.HMSet(uid_Key, mess).Result()
	if err != nil {
		panic(err)
	}

	return "register successfully"
}

func UserLogin(client *redis.Client, userName, email, password, salt string) (isLogin, getName string) {
	if userName != "" {
		uidKey, err := client.Get("userName:" + strings.ToLower(strings.TrimSpace(userName))).Result()
		if err == redis.Nil {
			return "userName does not exist", "00000000"
		}

		getpwd, err := client.HGet(uidKey, "password").Result()
		if err != nil {
			panic(err)
		}

		if getpwd == password_md5(password, salt) {
			return "login sucessfully", userName
		}

	} else if email != "" {
		uidKey, err := client.Get("email:" + email).Result()
		if err == redis.Nil {
			return "email does not exist", "00000000"
		}

		getpwd, err := client.HGet(uidKey, "password").Result()
		if err != nil {
			panic(err)
		}

		if getpwd == password_md5(password, salt) {
			userName, err = client.HGet(uidKey, "userName").Result()
			if err != nil {
				panic(err)
			}
			return "login sucessfully", userName
		}

	}
	return "fail", "00000000"
}

func RedisRegister(client *redis.Client, userInfo map[string]string) string {
	var newResister User
	//newResister.userName = "Victor"
	//newResister.email = "985844987@qq.com"
	//newResister.password = password_md5("dys123", "yan")
	//newResister.iphone = "15200338626"
	newResister.userName = userInfo["username"]
	newResister.email = userInfo["email"]
	newResister.password = password_md5(userInfo["password"], "yan")
	newResister.iphone = userInfo["phone"]

	status := register_user(client, newResister)
	return status

}

func NewRedisClient() *redis.Client {
	return createClient()
}

//func main() {
//	client := createClient()
//	fmt.Println(client)
//
//	//register(client)
//	status, userName := userLogin(client, "Victor", "", "dys123", "yan")
//	fmt.Println(status, userName)
//	//if status == "login sucessfully" {
//	//	fmt.Println(userName)
//	//}
//
//}
