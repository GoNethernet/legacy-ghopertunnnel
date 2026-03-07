package legacy

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/realms"
	"golang.org/x/oauth2"
)

// logger ...
var logger = log.New(os.Stdout, "", log.LstdFlags)

// Start starts the listening on the chosen server, this reads the config and applies values in the proxy, any error
// (config, listening etc.) is returned if non nil.
func Start() (*Server, error) {
	var (
		raddr   string
		rInfo   realms.Realm
		isRealm bool
	)
	_ = createFiles()
	token := tokenSrc()
	cfg := readConfig()
	raddr = cfg.Network.RemoteDirect
	if !strings.Contains(cfg.Network.RemoteDirect, ":") {
		isRealm = true
		c, cancel := context.WithTimeout(context.Background(), time.Second*20)
		client := realms.NewClient(token, http.DefaultClient)
		realm, err := client.Realm(c, cfg.Network.RemoteDirect)
		cancel()
		if err != nil {
			errR := err.Error()
			if strings.Contains(errR, "403") {
				logger.Fatalf("realm: realm '%s' not found in the default minecraft list, add the realm normally from minecraft first", realm.Name)
			} else if strings.Contains(errR, "404") {
				logger.Fatalf("realm: realm is unknown or unreachable")
			} else if strings.Contains(errR, "401") {
				logger.Fatalf("realm: realm rejected your connection request")
			} else if strings.Contains(errR, "502") {
				logger.Fatalf("realm: context deadline data, bad gateway, code: %s", err)
			} else {
				logger.Fatalf("unhandled error: %v", err)
			}
		}
		rInfo = realm
		raddr, err = realm.Address(context.Background())
		if err != nil {
			if !strings.Contains(raddr, ":") {
				logger.Fatalf("realm address: could not resolve address, excpected: non-char port, got: %s", raddr)
			}
			logger.Printf("realm address: %v", err)
			return nil, err
		}
	}
	resources := Resources{
		raddr: raddr,
		src:   token,
		cfg:   cfg,
	}
	packs, err := resources.Load()
	if err != nil {
		logger.Fatalf("resources packs: %v", err)
	}
	p, _ := minecraft.NewForeignStatusProvider(raddr)
	listenConfig := &minecraft.ListenConfig{
		ResourcePacks:          packs,
		AllowUnknownPackets:    true,
		AllowInvalidPackets:    false,
		AuthenticationDisabled: true,
		HTTPClient:             http.DefaultClient,
		StatusProvider:         p,
		TexturePacksRequired:   cfg.Resources.Required,
	}
	listener, err := listenConfig.Listen("raknet", cfg.Network.LocalAddress)
	if err != nil {
		logger.Printf("listen: %s", err)
		return nil, err
	}
	if isRealm {
		logger.Printf("listening... local address: %s, remote address: %s (realm: %s)", cfg.Network.LocalAddress, raddr, rInfo.Name)
	} else {
		logger.Printf("listening... local address: %s, remote address: %s", cfg.Network.LocalAddress, raddr)
	}
	if cfg.Whitelist.Enabled {
		logger.Printf("whitelist enabled: %v", cfg.Whitelist.Names)
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			logger.Printf("listener accept: %v", err)
			continue
		}
		return &Server{
			raddr:          raddr,
			cfg:            cfg,
			conn:           c.(*minecraft.Conn),
			listener:       listener,
			src:            token,
			packs:          packs,
			statusProvider: p,
			realm:          rInfo,
			listenConfig:   listenConfig,
			isRealm:        isRealm,
		}, nil
	}
}

// Config represent a customizable section for the proxy.
type Config struct {
	// Network ...
	Network struct {
		// LocalAddress represent the address where the proxy will start listening to, the default is 0.0.0.0:19132.
		LocalAddress string `json:"local_address"`
		// RemoteDirect is the address/the realm code where the client will be forwarded to.
		RemoteDirect string `json:"remote_direct"`
	} `json:"network"`
	// Resources ...
	Resources struct {
		// CachePath represent the path where the server resources packs will be cached, the default is "resources/cache".
		CachePath string `json:"cache_path"`
		// Path is the path where the proxy will read resources pack all appends them to the server resources
		// packs, the default is "resources".
		Path string `json:"path"`
		// Required defines if the resources packs are mandatory or not.
		Required bool `json:"required"`
	}
	// Whitelist ...
	Whitelist struct {
		// Enabled defines if the whitelist is enabled, the default value is 'false'.
		Enabled bool `json:"enabled"`
		// Names represent a slice of all minecraft game nametags that will be able to join the proxy if whitelist is enabled.
		Names []string `json:"names"`
	}
}

// readConfig creates the config if it is not created, then reads and returns it.
func readConfig() *Config {
	c := &Config{}
	path := "config.json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if c.Network.LocalAddress == "" {
			c.Network.LocalAddress = "0.0.0.0:19132"
		}
		if c.Resources.Path == "" {
			c.Resources.Path = "resources"
		}
		if c.Resources.CachePath == "" {
			c.Resources.CachePath = "resources/cache"
		}
		file, err := os.Create(path)
		if err != nil {
			logger.Fatalf("config create error: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logger.Fatalf("config close: %v", err)
			}
		}(file)
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(c)
		logger.Printf("default configuration created. edit the file and restart.")
		os.Exit(0)
	}
	file, err := os.Open(path)
	if err != nil {
		logger.Fatalf("config open error: %v", err)
	}
	defer file.Close()
	_ = json.NewDecoder(file).Decode(&c)
	return c
}

// tokenPath is the path where the token.json is located.
var tokenPath = "token.json"

// tokenSrc returns a token source for using with a gophertunnel client. It either reads it from the
// token.json file if cached or requests logging in with a device code.
func tokenSrc() oauth2.TokenSource {
	token := new(oauth2.Token)
	data, err := os.ReadFile(tokenPath)
	if err == nil {
		_ = json.Unmarshal(data, token)
	}
	oldClient := http.DefaultClient
	defer func() { http.DefaultClient = oldClient }()
	http.DefaultClient = &http.Client{
		Timeout: time.Second * 35,
	}
	src := auth.RefreshTokenSource(token)
	tok, err := src.Token()
	if err != nil || tok.AccessToken == "" {
		_ = os.Remove(tokenPath)
		token, err = auth.RequestLiveToken()
		if err != nil {
			logger.Fatalf("token login failed: %v.", err)
		}
		src = auth.RefreshTokenSource(token)
	}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		<-c
		tok, _ := src.Token()
		b, _ := json.Marshal(tok)
		_ = os.WriteFile(tokenPath, b, 0644)
		os.Exit(0)
	}()
	return src
}

// createFiles create runtime files if not created yet.
func createFiles() error {
	for _, file := range []string{"resources", "resources/cache"} {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			_ = os.Mkdir(file, 0755)
		}
	}
	return nil
}
