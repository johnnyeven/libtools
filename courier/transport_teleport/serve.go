package transport_teleport

import (
	"github.com/henrylee2cn/teleport"
	"github.com/henrylee2cn/teleport/proto/pbproto"
	"github.com/johnnyeven/libtools/conf"
	"time"
)

type ServeTeleport struct {
	Network            string        `conf:"env"`
	IP                 string        `conf:"env"`
	Port               uint16        `conf:"env"`
	DefaultDialTimeout time.Duration `conf:"env"`
	RedialTimes        int32         `conf:"env"`
	RedialInterval     time.Duration `conf:"env"`
	DefaultBodyCodec   string        `conf:"env"`
	DefaultSessionAge  time.Duration `conf:"env"`
	DefaultContextAge  time.Duration `conf:"env"`
	SlowCometDuration  time.Duration `conf:"env"`
	PrintDetail        bool          `conf:"env"`
	CountTime          bool          `conf:"env"`
	Plugins            []tp.Plugin   `conf:"-"`

	server tp.Peer
}

func (s *ServeTeleport) Init() {
	tp.SetDefaultProtoFunc(pbproto.NewPbProtoFunc())
	s.server = tp.NewPeer(tp.PeerConfig{
		Network:            s.Network,
		LocalIP:            s.IP,
		ListenPort:         s.Port,
		DefaultDialTimeout: s.DefaultDialTimeout,
		RedialTimes:        s.RedialTimes,
		RedialInterval:     s.RedialInterval,
		DefaultBodyCodec:   s.DefaultBodyCodec,
		DefaultSessionAge:  s.DefaultSessionAge,
		DefaultContextAge:  s.DefaultContextAge,
		SlowCometDuration:  s.SlowCometDuration,
		PrintDetail:        s.PrintDetail,
		CountTime:          s.CountTime,
	}, s.Plugins...)
}

func (s ServeTeleport) MarshalDefaults(v interface{}) {
	if h, ok := v.(*ServeTeleport); ok {
		if h.Port == 0 {
			h.Port = 9090
		}
	}
}

func (s *ServeTeleport) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Port": 9090,
	}
}

func (s *ServeTeleport) Start() error {
	return s.server.ListenAndServe()
}

func (s *ServeTeleport) Stop() error {
	return s.server.Close()
}

func (s *ServeTeleport) RegisterCallRouter(route interface{}, plugins ...tp.Plugin) []string {
	return s.server.RouteCall(route, plugins...)
}

func (s *ServeTeleport) RegisterPushRouter(route interface{}, plugins ...tp.Plugin) []string {
	return s.server.RoutePush(route, plugins...)
}
