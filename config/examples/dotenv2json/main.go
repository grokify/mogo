package main

import (
	"fmt"

	"github.com/grokify/mogo/config"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	InputFile  string `short:"i" long:"input" description:"Input file in .env format" required:"true"`
	Outputfile string `short:"o" long:"output" description:"Output file in JSON format" required:"true"`
}

type AppConfig struct {
	RingCentralTokenJSON  string `env:"RINGCENTRAL_TOKEN_JSON" json:"ringcentralTokenJson"`
	RingCentralServerURL  string `env:"RINGCENTRAL_SERVER_URL" json:"ringcentralServerUrl"`
	RingCentralWebhookURL string `env:"RINGCENTRAL_WEBHOOK_URL" json:"ringcentralWebhookUrl"`
	RingCentralBotID      string `env:"RINGCENTRAL_BOT_ID" json:"ringcentralBotId"`
	GoogleSvcAccountJWT   string `env:"GOOGLE_SERVICE_ACCOUNT_JWT" json:"googleServiceAccountJwt"`
	GoogleSpreadsheetID   string `env:"GOOGLE_SPREADSHEET_ID" json:"googleSpreadsheetId"`
	GoogleSheetIndex      int64  `env:"GOOGLE_SHEET_INDEX" json:"googleSheetIndex"`
}

func main() {
	opts := &Options{}
	_, err := flags.Parse(opts)
	if err != nil {
		panic(err)
	}

	cfg := &AppConfig{}

	err = config.EnvFileToJSONFile(cfg, opts.InputFile, opts.Outputfile, 0600, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("WROTE: %v\n", opts.Outputfile)
	fmt.Println("DONE")
}
