package mqtt

import (
	"context"
	"fmt"
	"math/rand"
	"mqttcli/logger"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type config struct {
	Server   string `env:"MQTT_SERVER"`
	Username string `env:"MQTT_USERNAME" envDefault:""`
	Password string `env:"MQTT_PASSWORD" envDefault:""`
	Topic    string `env:"MQTT_TOPIC"`
	ClientID string `env:"MQTT_CLIENT_ID" envDefault:""`
	QoS      byte   `env:"MQTT_QOS" envDefault:"1"`
}

type Connection struct {
	ctx context.Context
	cfg config
	cm  *autopaho.ConnectionManager
}

func NewConnection() (*Connection, context.CancelFunc) {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	cfg.QoS = ensureValidQoS(cfg.QoS)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	return &Connection{ctx: ctx, cfg: cfg}, stop
}

func (m *Connection) Publish(payload []byte) {
	m.connect()

	if _, err := m.cm.Publish(m.ctx, &paho.Publish{
		QoS:     m.cfg.QoS,
		Topic:   m.cfg.Topic,
		Payload: payload,
	}); err != nil {
		if m.ctx.Err() == nil {
			panic(err)
		}
	}

	m.disconnect()
	<-m.cm.Done()
}

func (m *Connection) Subscribe() {
	m.connect()

	if _, err := m.cm.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{Topic: m.cfg.Topic, QoS: m.cfg.QoS},
		},
	}); err != nil {
		logger.FailF("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
	}

	ch := make(chan bool)
	m.cm.AddOnPublishReceived(func(received autopaho.PublishReceived) (bool, error) {
		_, err := os.Stdout.Write(received.Packet.Payload)
		if err != nil {
			logger.Fail("Error reading:", err)
		}

		ch <- true
		return true, nil
	})

	for {
		select {
		case <-ch:
			m.disconnect()
		case <-m.ctx.Done():
			/* noop */
		}
		break
	}

	<-m.cm.Done()
}

func (m *Connection) connect() {
	u, err := url.Parse(m.cfg.Server)
	if err != nil {
		panic(err)
	}

	username := m.cfg.Username
	password := m.cfg.Password
	if u.User != nil {
		username = u.User.Username()
		urlPassword, passwordIsSet := u.User.Password()
		if passwordIsSet {
			password = urlPassword
		}
		u.User = nil
	}

	if m.cfg.ClientID == "" {
		m.cfg.ClientID = fmt.Sprintf("mqttcli-%08x", rand.Uint32())
	}

	cliCfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     20,
		CleanStartOnInitialConnection: true,
		OnConnectError: func(err error) {
			logger.Fail("error whilst attempting connection", err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID:      m.cfg.ClientID,
			OnClientError: func(err error) { logger.Fail("client error:", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					logger.FailF("server requested disconnect: %s", d.Properties.ReasonString)
				} else {
					logger.FailF("server requested disconnect; reason code: %d", d.ReasonCode)
				}
			},
		},
	}

	if username != "" && password != "" {
		cliCfg.ConnectUsername = username
		cliCfg.ConnectPassword = []byte(password)
	}

	cm, err := autopaho.NewConnection(m.ctx, cliCfg)
	if err != nil {
		panic(err)
	}
	// Wait for the connection to come up
	if err = cm.AwaitConnection(m.ctx); err != nil {
		panic(err)
	}

	m.cm = cm
}

func (m *Connection) disconnect() {
	if err := m.cm.Disconnect(m.ctx); err != nil {
		panic(err)
	}
}

func ensureValidQoS(qos byte) byte {
	if qos < 0 || qos > 2 {
		return 1
	}
	return qos
}
