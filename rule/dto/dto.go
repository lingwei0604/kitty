package dto

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lingwei0604/kitty/pkg/config"
	jwt2 "github.com/lingwei0604/kitty/pkg/kjwt"
	pb "github.com/lingwei0604/kitty/proto"
	"github.com/spf13/cast"
)

type Payload struct {
	RuleNames   []string               `json:"rule_names" schema:"rule_names"`
	Channel     string                 `json:"channel" schema:"channel"`
	VersionCode int                    `json:"version_code" schema:"version_code"`
	Os          uint8                  `json:"os" schema:"os"`
	UserId      uint64                 `json:"user_id" schema:"user_id"`
	Imei        string                 `json:"imei" schema:"imei"`
	Idfa        string                 `json:"idfa" schema:"idfa"`
	Oaid        string                 `json:"oaid" schema:"oaid"`
	Suuid       string                 `json:"suuid" schema:"suuid"`
	Mac         string                 `json:"mac" schema:"mac"`
	AndroidId   string                 `json:"android_id" schema:"android_id"`
	PackageName string                 `json:"package_name" schema:"package_name"`
	Ip          string                 `json:"ip" schema:"ip"`
	Q           map[string][]string    `json:"-" schema:"-"`
	B           map[string]interface{} `json:"-" schema:"-"`
	DMP         Dmp                    `json:"-" schema:"-"`
	Context     context.Context        `json:"-" schema:"-"`
	Redis       redis.UniversalClient  `json:"-" schema:"-"`
}

type Dmp struct {
	pb.DmpResp
}

func FromClaim(claim jwt2.Claim) *Payload {
	return FromClaimWithRedis(claim, nil)
}

func FromClaimWithRedis(claim jwt2.Claim, client redis.UniversalClient) *Payload {
	versionCode, _ := strconv.Atoi(claim.VersionCode)
	return &Payload{
		Channel:     claim.Channel,
		VersionCode: versionCode,
		Suuid:       claim.Suuid,
		UserId:      claim.UserId,
		Redis:       client,
	}
}

func FromTenant(tenant *config.Tenant) *Payload {
	return FromTenantWithRedis(tenant, nil)
}

func FromTenantWithRedis(tenant *config.Tenant, client redis.UniversalClient) *Payload {
	versionCode, _ := strconv.Atoi(tenant.VersionCode)
	return &Payload{
		Channel:     tenant.Channel,
		VersionCode: versionCode,
		Os:          tenant.Os,
		UserId:      tenant.UserId,
		Imei:        tenant.Imei,
		Idfa:        tenant.Idfa,
		Oaid:        tenant.Oaid,
		Suuid:       tenant.Suuid,
		Mac:         tenant.Mac,
		AndroidId:   tenant.AndroidId,
		PackageName: tenant.PackageName,
		Ip:          tenant.Ip,
		Context:     tenant.Context,
		Redis:       client,
	}
}

func (p *Payload) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (p Payload) Now() time.Time {
	return time.Now()
}

func (p Payload) Date(s string) time.Time {
	date, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		panic(err)
	}
	return date
}

func (p Payload) DaysAgo(s string) int {
	if s == "" {
		return 0
	}
	return int(time.Now().Sub(p.DateTime(s)).Hours() / 24)
}

func (p Payload) HoursAgo(s string) int {
	if s == "" {
		return 0
	}
	return int(time.Now().Sub(p.DateTime(s)).Hours())
}

func (p Payload) MinutesAgo(s string) int {
	if s == "" {
		return 0
	}
	return int(time.Now().Sub(p.DateTime(s)).Minutes())
}

func (p Payload) DateTime(s string) time.Time {
	date, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		panic(err)
	}
	return date
}

func (p Payload) IsBefore(s string) bool {
	var (
		t   time.Time
		err error
	)
	if len(s) == 10 {
		t, err = time.ParseInLocation("2006-01-02", s, time.Local)
	} else {
		t, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	}
	if err != nil {
		panic(err)
	}
	return time.Now().Before(t)
}

func (p Payload) IsAfter(s string) bool {
	var (
		t   time.Time
		err error
	)
	if len(s) == 10 {
		t, err = time.ParseInLocation("2006-01-02", s, time.Local)
	} else {
		t, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	}
	if err != nil {
		panic(err)
	}
	return time.Now().After(t)
}

func (p Payload) IsBetween(begin string, end string) bool {
	return p.IsAfter(begin) && p.IsBefore(end)
}

func (p Payload) IsWeekday(day int) bool {
	return time.Now().Weekday() == time.Weekday(day)
}

func (p Payload) IsWeekend() bool {
	if weekday := time.Now().Weekday(); weekday == 0 || weekday == 6 {
		return true
	}
	return false
}

func (p Payload) IsToday(s string) bool {
	return time.Now().Format("2006-01-02") == s
}

func (p Payload) IsHourRange(begin int, end int) bool {
	now := time.Now().Hour()
	return now >= begin && now <= end
}

func (p Payload) IsBlackListed() bool {
	return p.DMP.BlackType == pb.DmpResp_BLACK
}

func (p Payload) SIsMember(key string, needle string) bool {
	if p.Redis == nil {
		return false
	}
	if p.Context == nil {
		p.Context = context.Background()
	}
	ok, _ := p.Redis.SIsMember(p.Context, key, needle).Result()
	return ok
}

func (p Payload) Get(key string) string {
	if p.Redis == nil {
		return ""
	}
	if p.Context == nil {
		p.Context = context.Background()
	}
	result, _ := p.Redis.Get(p.Context, key).Result()
	return result
}

func (p Payload) ToString(str interface{}) string {
	return cast.ToString(str)
}

func (p Payload) ToInt(int interface{}) int {
	return cast.ToInt(int)
}

func (p Payload) Hash() int {
	h := fnv.New32a()
	if p.Suuid != "" {
		h.Write([]byte(p.Suuid))
	} else {
		h.Write([]byte(strconv.Itoa(int(p.UserId))))
	}

	return int(h.Sum32())
}

func (p Payload) Percent(percent int) bool {
	if p.Hash()%100 < percent {
		return true
	}
	return false
}

type Data map[string]interface{}

type Response struct {
	Code    uint `json:"code"`
	Message uint `json:"message"`
	Data    Data `json:"data"`
}

func (p Response) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (d Dmp) IsBlackListed() bool {
	return d.BlackType == pb.DmpResp_BLACK
}

func (d Dmp) RegisterRisk() int {
	if d.Skynet == nil {
		return 0
	}
	return int(d.Skynet.Register)
}

func (d Dmp) RegisteredDays() int {
	date, err := time.ParseInLocation("2006-01-02", d.Register, time.Local)
	if err != nil {
		panic(err)
	}
	return int(time.Now().Sub(date).Hours() / 24)
}

func (d Dmp) BrowseRisk() int {
	if d.Skynet == nil {
		return 0
	}
	return int(d.Skynet.Browse)
}

func (d Dmp) FissionRisk() int {
	if d.Skynet == nil {
		return 0
	}
	return int(d.Skynet.Fission)
}

func (d Dmp) OverallRisk() int {
	if d.Skynet == nil {
		return 0
	}
	return int(d.Skynet.Level)
}

func (d Dmp) LoginRisk() int {
	if d.Skynet == nil {
		return 0
	}
	return int(d.Skynet.Login)
}

func (d Dmp) TaskRisk() int {
	if d.Skynet == nil {
		return 0
	}
	return int(d.Skynet.Task)
}

func (d Dmp) WithdrawRisk() int {
	if d.Skynet == nil {
		return 0
	}
	return int(d.Skynet.Withdraw)
}
