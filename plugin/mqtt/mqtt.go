package mqtt

import (
	"fmt"
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"github.com/glycerine/zygomys/zygo"
	"github.com/goiiot/libmqtt"
	"go.uber.org/zap"
	"net/url"

	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

// set up a pluginConfig struct
type pluginConfig struct {
	Prefix          string
	Interval        int
	Address         string
	LispFilter      string `yaml:"lisp_filter"`
	Template        string `yaml:"template"`
	TemplateOnClick string `yaml:"template_on_click"`
	Subscribe       string
}

type plugin struct {
	l              *zap.SugaredLogger
	cfg            pluginConfig
	cnt            int
	ev             int
	nextTs         time.Time
	m              libmqtt.Client
	zl             *zygo.Zlisp
	lastMQTTUpdate time.Time
	lastMessage    string
}

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{
		l: cfg.Logger,
	}
	p.cfg, err = loadConfig(cfg.Config)
	return p, err
}

type Event struct {
	Msg string
	TS  time.Time
}

func (p *plugin) Init() (err error) {
	//p.zl = zygo.NewZlispSandbox()
	p.zl = zygo.NewZlisp()
	if len(p.cfg.LispFilter) > 0 {
		err := p.zl.LoadString(p.cfg.LispFilter)
		if err != nil {
			return err
		}
	}
	AddZygoFuncs(p.zl)

	p.m, err = libmqtt.NewClient(
		//libmqtt.WithVersion(libmqtt.V5,true),
		libmqtt.WithKeepalive(30, 1.2),
		// enable auto reconnect and set backoff strategy
		libmqtt.WithAutoReconnect(true),
		libmqtt.WithBackoffStrategy(time.Second, time.Minute, 1.2),
		// use RegexRouter for topic routing, if not specified
		// will use TextRouter, which will match full text
		libmqtt.WithRouter(libmqtt.NewRegexRouter()),
	)
	var mqttOpts = []libmqtt.Option{}
	mqttUrl, err := url.Parse(p.cfg.Address)
	if err != nil {
		return err
	}
	switch mqttUrl.Scheme {
	case "tcp":
		mqttOpts = append(mqttOpts, libmqtt.WithCustomTLS(nil))
	case "tls":
	default:
		return fmt.Errorf("MQTT protocol [%s] not supported", mqttUrl)
	}
	if mqttUrl.User != nil {
		pass, _ := mqttUrl.User.Password()
		mqttOpts = append(mqttOpts, libmqtt.WithIdentity(mqttUrl.User.Username(), pass))
	}

	mqttOpts = append(mqttOpts, libmqtt.WithConnHandleFunc(func(client libmqtt.Client, server string, code byte, err error) {
		if err != nil {
			// failed
			panic(fmt.Sprintf("failed to connect: %s", err))
		}

		if code != libmqtt.CodeSuccess {
			// server rejected or in error
			panic(fmt.Sprintf("server rejected with %s", ReasonString(code)))
		}

		// success
		// you are now connected to the `server`
		// (the `server` is one of your provided `servers` when create the client)
		// start your business logic here or send a signal to your logic to start

		// subscribe some topic(s)
		p.m.Subscribe([]*libmqtt.Topic{
			//{Name: "events/#"},
			//{Name: "metrics/mpower/power/socket2"},
			{Name: p.cfg.Subscribe},
		}...)

		// publish some topic message(s)
		// go func() {
		// 	for {
		// 		p.m.Publish([]*libmqtt.PublishPacket{
		// 			{
		// 				TopicName: "foo",
		// 				Payload:   []byte("bar"),
		// 				Qos:       libmqtt.Qos0,
		// 				Props: &libmqtt.PublishProps{
		// 					RespTopic: "bar",
		// 					UserProps: map[string][]string{
		// 						"testprop1": []string{"testprop1-val"},
		// 					},
		// 					SubIDs:      nil,
		// 					ContentType: "application/json",
		// 				},
		// 			},
		// 		}...)
		// 		time.Sleep(time.Second)
		// 	}
		// } ()
	}))
	mqttOpts = append(mqttOpts, libmqtt.WithSubHandleFunc(func(client libmqtt.Client, topics []*libmqtt.Topic, err error) {
		p.l.Infof("[mqtt] subscription:")
		for _, t := range topics {
			p.l.Infof("[mqtt] subscribed to %s", t.Name)
		}
	}))

	err = p.m.ConnectServer(mqttUrl.Host, mqttOpts...)
	if len(p.cfg.LispFilter) == 0 {
		p.m.HandleTopic(".*", func(client libmqtt.Client, topic string, qos libmqtt.QosLevel, msg []byte) {
			p.lastMessage = string(msg)
			p.lastMQTTUpdate = time.Now()
		})
	} else {
		p.m.HandleTopic(".*", func(client libmqtt.Client, topic string, qos libmqtt.QosLevel, msg []byte) {
			p.UpdateFromMQTT(string(msg))
		})
	}

	return err
}

func (p *plugin) GetUpdateInterval() int {
	return p.cfg.Interval
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	// TODO precompile and preallcate
	tpl, _ := util.NewTemplate("uberEvent", p.cfg.Template)
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function
	util.WaitForTs(&p.nextTs)
	update.FullText = tpl.ExecuteString(Event{
		Msg: p.lastMessage,
		TS:  p.lastMQTTUpdate,
	})
	update.Markup = "pango"
	update.Color = util.GetColorPct(
		int(
			(time.Now().Sub(p.lastMQTTUpdate) / (time.Minute / 3))))
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	tpl, _ := util.NewTemplate("uberEvent", p.cfg.TemplateOnClick)
	update.FullText = tpl.ExecuteString(Event{
		Msg: p.lastMessage,
		TS:  p.lastMQTTUpdate,
	})
	update.Markup = "pango"
	update.Color = "pango"
	update.Color = util.GetColorPct(
		int(
			(time.Now().Sub(p.lastMQTTUpdate) / (time.Minute / 3))))
	p.ev++
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Interval = 10000
	cfg.Prefix = "ex: "
	// optionally, check for pluginConfig validity after GetConfig call
	return cfg, c.GetConfig(&cfg)
}

func (p *plugin) UpdateFromMQTT(v string) {
	p.zl.AddGlobal("x", &zygo.SexpStr{S: v})
	err := p.zl.LoadString(p.cfg.LispFilter)
	if err != nil {
		// this should not happen as we check same string in init but oh well
		// maybe live editing will be there someday
		p.l.Errorf("error parsing expression: %s", err)
		return
	}

	sexp, err := p.zl.Run()
	defer p.zl.Clear()
	if err != nil {
		p.l.Errorf("error: %s", err)
		p.lastMessage = fmt.Sprintf("LISP run err: %s on %s", err, v)
		return
	}

	switch out := sexp.(type) {
	case *zygo.SexpStr:
		p.lastMessage = out.S
		p.lastMQTTUpdate = time.Now()
	default:
		p.lastMessage = fmt.Sprintf("LISP retval should be string: %+v %s", out, err)
	}
}
