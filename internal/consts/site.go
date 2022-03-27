package consts

import "time"

const (
	SiteDefaultHost = "penguin-stats.io"

	SiteGlobalMirrorHost        = "penguin-stats.io"
	SiteChinaMainlandMirrorHost = "penguin-stats.cn"

	ShimCompatibilityHeaderKey   = "X-Penguin-Compatible"
	ShimCompatibilityHeaderValue = "frontend-v2@v3.4.0"
)

var SiteHosts = []string{
	SiteGlobalMirrorHost,
	SiteChinaMainlandMirrorHost,
}

var Servers = []string{
	"CN",
	"US",
	"JP",
	"KR",
}

var ServerNameMapping = map[string]string{
	"CN": "国服",
	"US": "美服",
	"JP": "日服",
	"KR": "韩服",
}

var LocMap = map[string]*time.Location{
	"CN": time.FixedZone("UTC+8", +8*60*60),
	"US": time.FixedZone("UTC-7", -7*60*60),
	"JP": time.FixedZone("UTC+9", +9*60*60),
	"KR": time.FixedZone("UTC+9", +9*60*60),
}
