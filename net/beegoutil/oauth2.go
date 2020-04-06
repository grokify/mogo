package beegoutil

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/grokify/gotilla/type/stringsutil"
	ms "github.com/grokify/oauth2more/multiservice"
)

const (
	BeegoOauth2ProvidersCfgVar    string = "oauth2providers"
	BeegoOauth2ConfigCfgVarPrefix string = "oauth2config"
)

func InitOAuth2Config(o2ConfigSet *ms.ConfigMoreSet) error {
	oauth2providersraw := beego.AppConfig.String(BeegoOauth2ProvidersCfgVar)
	oauth2providers := stringsutil.SplitTrimSpace(oauth2providersraw, ",")
	for _, providerKey := range oauth2providers {
		providerKey = strings.TrimSpace(providerKey)
		if len(providerKey) == 0 {
			continue
		}
		oauth2ConfigParam := BeegoOauth2ConfigCfgVarPrefix + providerKey
		configJson := strings.TrimSpace(beego.AppConfig.String(oauth2ConfigParam))
		if len(configJson) == 0 {
			return fmt.Errorf("E_NO_CONFIG_FOR_OAUTH_PROVIDER_KEY [%v]", providerKey)
		}
		err := o2ConfigSet.AddConfigMoreJson(providerKey, []byte(configJson))
		if err != nil {
			return err
		}
	}
	return nil
}
