package coze

import (
	"context"
	"coze-agent-platform/config"
	"coze-agent-platform/utils"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

const (
	COZE_TOKEN_KEY       = "coze:access_token"
	TOKEN_EXPIRE_MINUTES = 14
)

func GetToken() (string, error) {
	ctx := context.Background()

	// 首先尝试从Redis获取缓存的token
	if utils.RDB != nil {
		cachedToken, err := utils.RDB.Get(ctx, COZE_TOKEN_KEY).Result()
		if err == nil && cachedToken != "" {
			fmt.Println("使用缓存的token")
			return cachedToken, nil
		}
	}

	// Redis中没有token或者Redis不可用，从API获取新token
	fmt.Println("从API获取新token")

	cozeConfig := config.GetCozeConfig()
	// 优先使用配置文件中的私钥字符串，如果为空则尝试读取文件
	var jwtOauthPrivateKey string
	if cozeConfig.PrivateKey != "" {
		jwtOauthPrivateKey = cozeConfig.PrivateKey
	} else if cozeConfig.PrivateKeyFilePath != "" {
		privateKeyBytes, err := os.ReadFile(cozeConfig.PrivateKeyFilePath)
		if err != nil {
			return "", fmt.Errorf("读取私钥文件失败: %v", err)
		}
		jwtOauthPrivateKey = string(privateKeyBytes)
	} else {
		return "", fmt.Errorf("未配置私钥")
	}

	oauth, err := coze.NewJWTOAuthClient(coze.NewJWTOAuthClientParam{
		PrivateKeyPEM: jwtOauthPrivateKey,
		ClientID:      cozeConfig.ClientID,
		PublicKey:     cozeConfig.PublicKeyID,
	}, coze.WithAuthBaseURL(cozeConfig.APIURL))

	if err != nil {
		return "", fmt.Errorf("创建JWT OAuth客户端失败: %v", err)
	}

	resp, err := oauth.GetAccessToken(ctx, nil)
	fmt.Println("resp", resp)
	if err != nil {
		return "", fmt.Errorf("获取AccessToken失败: %v", err)
	}

	// 将新token存储到Redis，设置14分钟过期时间
	if utils.RDB != nil {
		expiration := time.Duration(TOKEN_EXPIRE_MINUTES) * time.Minute
		err = utils.RDB.Set(ctx, COZE_TOKEN_KEY, resp.AccessToken, expiration).Err()
		if err != nil {
			fmt.Printf("警告: 无法将token存储到Redis: %v\n", err)
		} else {
			fmt.Printf("token已存储到Redis，有效期: %d分钟\n", TOKEN_EXPIRE_MINUTES)
		}
	}

	return resp.AccessToken, nil
}
