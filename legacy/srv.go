package legacy

import (
	"context"
	"errors"
	"fmt"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/cmd"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/session"
	"iter"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/realms"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"golang.org/x/oauth2"
)

// Server ...
type Server struct {
	statusProvider minecraft.ServerStatusProvider
	packs          []*resource.Pack
	conn           *minecraft.Conn
	listener       *minecraft.Listener
	src            oauth2.TokenSource
	raddr          string
	isRealm        bool
	realm          realms.Realm
	cfg            *Config
	listenConfig   *minecraft.ListenConfig
}

// ResourcesPacks returns all applied resources packs.
func (s *Server) ResourcesPacks() []*resource.Pack {
	return s.packs
}

// RemoteAddr ...
func (s *Server) RemoteAddr() string {
	return s.raddr
}

// LocalAddr ...
func (s *Server) LocalAddr() string {
	return s.listener.Addr().String()
}

// Close closes the connection.
func (s *Server) Close() error {
	return s.conn.Close()
}

// StatusProvider ...
func (s *Server) StatusProvider() minecraft.ServerStatusProvider {
	return s.statusProvider
}

// IsFull ...
func (s *Server) IsFull() bool {
	if s.listenConfig.MaximumPlayers == s.listener.PlayerCount() {
		return true
	}
	return false
}

// Realm returns the realm if the remote direct is a realm, otherwise nil is returned.
func (s *Server) Realm() (*realms.Realm, bool) {
	if !s.isRealm {
		return nil, false
	}
	return &s.realm, true
}

// Accept ...
func (s *Server) Accept() iter.Seq[*player.Player] {
	return func(yield func(*player.Player) bool) {
		var (
			mu sync.Mutex
		)
		httpClient := &http.Client{
			Timeout: time.Second * 30,
		}
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
		ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
		defer cancel()
		dialer := minecraft.Dialer{
			TokenSource:         s.src,
			ClientData:          s.conn.ClientData(),
			IdentityData:        s.conn.IdentityData(),
			HTTPClient:          httpClient,
			KeepXBLIdentityData: true,
			DownloadResourcePack: func(id uuid.UUID, version string, current, total int) bool {
				name := fmt.Sprintf("%s_%s", id, version)
				path := filepath.Join(s.cfg.Resources.CachePath, name+".mcpack")
				_, err := os.Stat(path)
				return os.IsNotExist(err)
			},
		}
		serverConn, err := dialer.DialContext(ctx, "raknet", s.raddr)
		_, _ = s.src.Token()
		if err != nil {
			_ = s.listener.Disconnect(s.conn, "failed to connect to remote server")
			logger.Fatalf("dial failed, name %s : %v", s.conn.IdentityData().DisplayName, err)
			return
		}
		if s.cfg.Whitelist.Enabled {
			whitelisted := false
			for _, name := range s.cfg.Whitelist.Names {
				if name == s.conn.IdentityData().DisplayName {
					whitelisted = true
					break
				}
			}
			if !whitelisted {
				_ = s.listener.Disconnect(s.conn, "unwhitelisted client")
				logger.Printf("unwhitelisted client, exepected: %v, got: %s", s.cfg.Whitelist.Names, s.conn.IdentityData().DisplayName)
				return
			}
		}
		logger.Printf("connected, server: %s, player: %s", s.raddr, s.conn.ClientData().ThirdPartyName)
		_ = os.MkdirAll(s.cfg.Resources.CachePath, 0755)
		for _, pack := range s.packs {
			packName := clean(pack.Name())
			cachePath := filepath.Join(s.cfg.Resources.CachePath, packName+".mcpack")
			if _, err := os.Stat(cachePath); os.IsNotExist(err) {
				packData := make([]byte, pack.Len())
				_, _ = pack.ReadAt(packData, 0)
				_ = os.WriteFile(cachePath, packData, 0644)
				logger.Printf("pack cached: %s", packName)
			}
		}
		var (
			g       sync.WaitGroup
			once    sync.Once
			players []protocol.PlayerListEntry
			ready   = make(chan struct{})
		)
		g.Add(2)
		go func() {
			defer g.Done()
			if err := s.conn.StartGame(serverConn.GameData()); err != nil {
				_ = serverConn.Close()
				logger.Fatalf("start game: %v", err)
			}
		}()
		go func() {
			defer g.Done()
			if err := serverConn.DoSpawn(); err != nil {
				_ = s.conn.Close()
				logger.Fatalf("do spawn: %v", err)
			}
		}()
		g.Wait()
		done := make(chan struct{})
		se := session.New(&mu)
		p := player.NewPlayer(s.conn, serverConn, &mu, se)
		go func() {
			loginHandler := session.LoginPacket{}
			defer func() {
				_ = serverConn.Close()
				select {
				case <-done:
				default:
					close(done)
				}
				logger.Fatalf("player disconnected")
			}()
			for {
				pk, err := s.conn.ReadPacket()
				if err != nil {
					return
				}
				se.UpdateFromClient(pk)
				if formPk, ok := pk.(*packet.ModalFormResponse); ok {
					data, hasData := formPk.ResponseData.Value()
					if hasData {
						p.HandleFormResponse(formPk.FormID, data)
					} else {
						p.HandleFormResponse(formPk.FormID, nil)
					}
					continue
				}

				if _, ok := pk.(*packet.Login); ok {
					if err := loginHandler.Handle(pk, s.conn); err != nil {
						logger.Printf("%v", err)
						continue
					}
				}

				if _, ok := pk.(*packet.PlayerAuthInput); ok {
					once.Do(func() { close(ready) })
				}

				if cmdReq, ok := pk.(*packet.CommandRequest); ok {
					lineStr := strings.TrimPrefix(cmdReq.CommandLine, "/")
					parts := strings.Split(lineStr, " ")
					if len(parts) == 0 || parts[0] == "" {
						continue
					}
					name := parts[0]

					if rc, ok := cmd.CustomCommands[name]; ok {
						line := cmd.NewLine(parts[1:], p, players)
						cmdValue := reflect.New(rc.Type)
						cmdInstance := cmdValue.Interface().(cmd.Command)

						if p.PermissionLevel().Level() < cmdInstance.PermissionLevel().Level() {
							_ = p.Message("§cUnknown command: " + name + ". Please check that the command exists and that you have permission to use it.")
						} else {
							val := cmdValue.Elem()
							parser := cmd.Parser{}
							failed := false

							for i := 0; i < val.NumField(); i++ {
								field := val.Field(i)
								fieldType := val.Type().Field(i)
								fieldName := strings.ToLower(fieldType.Name)
								optional := strings.HasPrefix(field.Type().Name(), "Optional[")

								if err := parser.ParseArgument(line, field, optional, fieldName); err != nil {
									_ = p.Message("§c" + err.Error())
									failed = true
									break
								}
							}

							if !failed {
								if arg, ok := line.Next(); ok {
									_ = p.Message(fmt.Sprintf("§cSyntax error: Unexpected \"%s\": at \"%s\"", arg, arg))
								} else {
									cmdInstance.Run(p)
								}
							}
						}
						p.Session().ClearCommandData()
						continue
					}
				}

				if err := serverConn.WritePacket(pk); err != nil {
					return
				}
			}
		}()
		go func() {
			textHandler := session.TextPacket{}
			disconnectHandler := session.DisconnectPacket{}
			respawnHandler := session.RespawnPacket{}
			deathHandler := session.DeathPacket{}
			defer func() {
				_ = s.conn.Close()
				select {
				case <-done:
				default:
					close(done)
				}
			}()
			for {
				pk, err := serverConn.ReadPacket()
				if err != nil {
					var disc minecraft.DisconnectError
					if errors.As(err, &disc) {
						_ = s.listener.Disconnect(s.conn, disc.Error())
					}
					return
				}
				se.UpdateFromServer(pk)
				if _, ok := pk.(*packet.SetHealth); ok {
					once.Do(func() { close(ready) })
				}
				if _, ok := pk.(*packet.Text); ok {
					if err := textHandler.Handle(pk, s.conn); err != nil {
						logger.Printf("%v", err)
						continue
					}
				}
				if t, ok := pk.(*packet.PlayerList); ok {
					if t.ActionType == packet.PlayerListActionAdd {
						players = append(players, t.Entries...)
					} else {
						for _, entry := range t.Entries {
							for i, pEntry := range players {
								if pEntry.UUID == entry.UUID {
									players = append(players[:i], players[i+1:]...)
									break
								}
							}
						}
					}
				}
				if _, ok := pk.(*packet.Disconnect); ok {
					if err := disconnectHandler.Handle(pk, s.conn); err != nil {
						logger.Printf("%v", err)
						continue
					}
				}
				if _, ok := pk.(*packet.DeathInfo); ok {
					if err := deathHandler.Handle(pk, s.conn); err != nil {
						logger.Printf("%v", err)
						continue
					}
				}
				if _, ok := pk.(*packet.Respawn); ok {
					if err := respawnHandler.Handle(pk, s.conn); err != nil {
						logger.Printf("%v", err)
						continue
					}
				}
				if cpk, ok := pk.(*packet.AvailableCommands); ok {
					se.UpdateFromServer(cpk)

					for name, rc := range cmd.CustomCommands {
						cmdInstance := reflect.New(rc.Type).Interface().(cmd.Command)

						if p.PermissionLevel().Level() < cmdInstance.PermissionLevel().Level() {
							continue
						}

						ovl := cmd.BuildCommand(rc.Type, &cpk.Enums, &cpk.EnumValues)
						found := false
						for i, existing := range cpk.Commands {
							if strings.EqualFold(existing.Name, name) {
								logger.Printf("command '%s' is overwriting an existing server command", name)
								cpk.Commands[i].Description = cmdInstance.Description()
								cpk.Commands[i].PermissionLevel = byte(cmdInstance.PermissionLevel().Level())
								cpk.Commands[i].Overloads = ovl
								found = true
								break
							}
						}
						if !found {
							cpk.Commands = append(cpk.Commands, protocol.Command{
								Name:            name,
								Description:     cmdInstance.Description(),
								PermissionLevel: byte(cmdInstance.PermissionLevel().Level()),
								Overloads:       ovl,
								AliasesOffset:   0xFFFFFFFF,
							})
						}
					}
					_ = s.conn.WritePacket(cpk)
					continue
				}
				if err := s.conn.WritePacket(pk); err != nil {
					return
				}

			}
		}()
		<-ready
		if !yield(p) {
			return
		}
		<-done
	}
}
