package player

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/cmd"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/effect"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/form"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/hud"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/permission"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/position"
	"github.com/gonethernet/legacy-ghopertunnel/legacy/player/session"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

// Player represents a connected user in the proxy session.
type Player struct {
	mu         *sync.Mutex
	client     *minecraft.Conn
	server     *minecraft.Conn
	formsMu    sync.Mutex
	forms      map[uint32]interface{}
	lastFormID uint32
	se         *session.Session
}

// NewPlayer creates and initializes a new Player instance.
func NewPlayer(client, server *minecraft.Conn, mu *sync.Mutex, sess *session.Session) *Player {
	return &Player{
		client: client,
		server: server,
		mu:     mu,
		se:     sess,
		forms:  make(map[uint32]interface{}),
	}
}

// Name returns the display name of the player.
func (p *Player) Name() string { return p.client.IdentityData().DisplayName }

// XUID returns the Xbox User ID of the player.
func (p *Player) XUID() string { return p.client.IdentityData().XUID }

// Message sends a chat message to the player's client.
func (p *Player) Message(msg string) error {
	return p.client.WritePacket(&packet.Text{Message: msg, TextType: packet.TextTypeChat})
}

// Messagef sends a formatted chat message to the player's client.
func (p *Player) Messagef(msg string, args ...interface{}) error {
	return p.Message(fmt.Sprintf(msg, args...))
}

// SendMessage sends a chat message to the destination server.
func (p *Player) SendMessage(msg string) error {
	return p.server.WritePacket(&packet.Text{Message: msg, TextType: packet.TextTypeChat})
}

// SendMessagef sends a formatted chat message to the destination server.
func (p *Player) SendMessagef(msg string, args ...interface{}) error {
	return p.SendMessage(fmt.Sprintf(msg, args...))
}

// Health returns the current health value from the session.
func (p *Player) Health() int32 {
	return p.se.Health()
}

// HeldSlot returns the current inventory slot held by the player.
func (p *Player) HeldSlot() int32 {
	return p.se.HeldSlot()
}

// GameMode returns the current game mode of the player.
func (p *Player) GameMode() int32 {
	return p.se.GameMode()
}

// PermissionLevel returns the current permission level from the session.
func (p *Player) PermissionLevel() permission.Permission {
	return p.se.PermissionLevel()
}

// Position returns the current spatial position of the player.
func (p *Player) Position() position.Position {
	inputs := p.se.Inputs()
	if inputs == nil {
		return position.Position{}
	}
	pos := inputs.Position
	return position.NewPosition(float64(pos.X()), float64(pos.Y()), float64(pos.Z()))
}

// Toast sends a toast notification request to the client.
func (p *Player) Toast(title, msg string) error {
	return p.client.WritePacket(&packet.ToastRequest{Title: title, Message: msg})
}

// Popup sends a popup text message to the client.
func (p *Player) Popup(text string) error {
	return p.client.WritePacket(&packet.Text{Message: text, TextType: packet.TextTypePopup})
}

// Popupf sends a formatted popup text message to the client.
func (p *Player) Popupf(msg string, args ...interface{}) error {
	return p.Popup(fmt.Sprintf(msg, args...))
}

// SetHeldSlot updates the current held item slot for the client.
func (p *Player) SetHeldSlot(slot int32) error {
	return p.client.WritePacket(&packet.MobEquipment{InventorySlot: byte(slot)})
}

// Title sends a main title text to the player's screen.
func (p *Player) Title(t Title) error {
	return p.client.WritePacket(&packet.SetTitle{
		Text:            t.Text(),
		ActionType:      packet.TitleActionSetTitle,
		FadeOutDuration: int32(t.FadeOutDuration()),
		FadeInDuration:  int32(t.FadeInDuration()),
	})
}

// SubTitle sends a subtitle text to the player's screen.
func (p *Player) SubTitle(t Title) error {
	return p.client.WritePacket(&packet.SetTitle{
		Text:            t.Text(),
		ActionType:      packet.TitleActionSetSubtitle,
		FadeOutDuration: int32(t.FadeOutDuration()),
		FadeInDuration:  int32(t.FadeInDuration()),
	})
}

// TitleAndSubTitle sends both a title and a subtitle simultaneously.
func (p *Player) TitleAndSubTitle(title, sub Title) error {
	err := p.Title(title)
	if err != nil {
		return err
	}
	return p.SubTitle(sub)
}

// SendScoreboard displays a custom scoreboard on the player's sidebar.
func (p *Player) SendScoreboard(sb Scoreboard) error {
	_ = p.client.WritePacket(&packet.SetDisplayObjective{
		DisplaySlot:   "sidebar",
		ObjectiveName: "scoreboard",
		DisplayName:   sb.Title(),
		CriteriaName:  "dummy",
	})
	entries := make([]protocol.ScoreboardEntry, 0, len(sb.Lines()))
	for i, line := range sb.Lines() {
		entries = append(entries, protocol.ScoreboardEntry{
			EntryID:       int64(i),
			ObjectiveName: "scoreboard",
			Score:         int32(i),
			IdentityType:  protocol.ScoreboardIdentityFakePlayer,
			DisplayName:   line,
		})
	}
	return p.client.WritePacket(&packet.SetScore{
		ActionType: packet.ScoreboardActionModify,
		Entries:    entries,
	})
}

// RemoveScoreboard removes the current scoreboard from the player's sidebar.
func (p *Player) RemoveScoreboard() error {
	return p.client.WritePacket(&packet.RemoveObjective{ObjectiveName: "scoreboard"})
}

// JukeboxPopup displays a jukebox-style popup message.
func (p *Player) JukeboxPopup(text string) error {
	return p.client.WritePacket(&packet.Text{
		Message:  text,
		TextType: packet.TextTypeJukeboxPopup,
	})
}

// JukeboxPopupf displays a formatted jukebox-style popup message.
func (p *Player) JukeboxPopupf(text string, args ...interface{}) error {
	return p.client.WritePacket(&packet.Text{
		Message:  fmt.Sprintf(text, args...),
		TextType: packet.TextTypeJukeboxPopup,
	})
}

// Rotation returns the current look direction of the player.
func (p *Player) Rotation() *position.Rotation {
	p.mu.Lock()
	defer p.mu.Unlock()
	return position.NewRotation(p.client, p.se.Inputs())
}

// HandleFormResponse processes responses from UI forms sent to the client.
func (p *Player) HandleFormResponse(id uint32, data []byte) {
	p.formsMu.Lock()
	f, ok := p.forms[id]
	delete(p.forms, id)
	p.formsMu.Unlock()
	if !ok || data == nil || string(data) == "null" {
		return
	}
	switch formType := f.(type) {
	case form.CustomForm:
		var res []json.RawMessage
		if err := json.Unmarshal(data, &res); err == nil {
			elems := formType.Elements()
			for i, raw := range res {
				if i >= len(elems) {
					break
				}
				p.handleCustomElement(elems[i], raw)
			}
			formType.Submit(p)
		}
	case form.ModalForm:
		var res bool
		if err := json.Unmarshal(data, &res); err == nil {
			btns := formType.Buttons()
			if res && btns[0].Submit != nil {
				btns[0].Submit(p)
			} else if !res && btns[1].Submit != nil {
				btns[1].Submit(p)
			}
		}
	case form.SimpleForm:
		var index uint32
		if err := json.Unmarshal(data, &index); err == nil {
			elems := formType.Elements()
			btns := form.ElementsToButtons(elems)
			if int(index) < len(btns) {
				if btns[index].Submit != nil {
					btns[index].Submit(p)
				}
			}
		}
	default:
		fmt.Printf("unknown form type: %v", f)
	}
}

// handleCustomElement processes individual elements within a custom form.
func (p *Player) handleCustomElement(e form.Element, data []byte) {
	switch t := e.(type) {
	case form.Toggle:
		if t.Value != nil {
			var v bool
			_ = json.Unmarshal(data, &v)
			t.Value(v)
		}
	case *form.Toggle:
		if t.Value != nil {
			var v bool
			_ = json.Unmarshal(data, &v)
			t.Value(v)
		}
	case form.Slider:
		if t.Selected != nil {
			var v float32
			_ = json.Unmarshal(data, &v)
			t.Selected(v)
		}
	case *form.Slider:
		if t.Selected != nil {
			var v float32
			_ = json.Unmarshal(data, &v)
			t.Selected(v)
		}
	case form.Input:
		if t.Final != nil {
			var v string
			_ = json.Unmarshal(data, &v)
			t.Final(v)
		}
	case *form.Input:
		if t.Final != nil {
			var v string
			_ = json.Unmarshal(data, &v)
			t.Final(v)
		}
	case form.Dropdown:
		if t.Selected != nil {
			var i int
			if err := json.Unmarshal(data, &i); err == nil && i < len(t.Options) {
				t.Selected(t.Options[i])
			}
		}
	case *form.Dropdown:
		if t.Selected != nil {
			var i int
			if err := json.Unmarshal(data, &i); err == nil && i < len(t.Options) {
				t.Selected(t.Options[i])
			}
		}
	}
}

// CommandsData returns the map of all registered custom commands.
func (p *Player) CommandsData() map[string]cmd.RegisteredCommand {
	return cmd.CustomCommands
}

// IsCustom checks if a specific command name is a registered custom command.
func (p *Player) IsCustom(name string) bool {
	name = strings.TrimPrefix(name, "/")
	_, ok := cmd.CustomCommands[name]
	return ok
}

// Session returns the internal proxy session associated with the player.
func (p *Player) Session() *session.Session {
	return p.se
}

// SendModalForm sends a modal (two-button) form to the player's client.
func (p *Player) SendModalForm(f form.ModalForm) error {
	id := p.nextID(f)
	btns := f.Buttons()
	data, _ := json.Marshal(struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Button1 string `json:"button1"`
		Button2 string `json:"button2"`
	}{
		Type:    "modal",
		Title:   f.Title(),
		Content: f.Message(),
		Button1: btns[0].Text,
		Button2: btns[1].Text,
	})
	return p.client.WritePacket(&packet.ModalFormRequest{FormID: id, FormData: data})
}

// RegisterCommand registers a new custom command and syncs it with the client.
func (p *Player) RegisterCommand(c cmd.Command) {
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	cmd.CustomCommands[c.Name()] = cmd.RegisteredCommand{Type: t}
	p.SendCommands()
}

// SendCommands sends the available custom commands.
func (p *Player) SendCommands() {
	pk := p.se.AvailableCommands()
	if pk == nil {
		pk = &packet.AvailableCommands{}
	}
	for name, rc := range cmd.CustomCommands {
		val := reflect.New(rc.Type)
		if val.Kind() == reflect.Ptr && val.Elem().CanAddr() {
		}

		cmdInstance := val.Interface().(cmd.Command)
		aliases := cmdInstance.Aliases()
		baseOverloads := cmd.NewCommand(rc.Type, &pk.Enums, &pk.EnumValues)

		aliasOffset := uint32(0xFFFFFFFF)
		if len(aliases) > 0 {
			aliasEnum := protocol.CommandEnum{Type: strings.ToLower(name) + "Aliases"}
			for _, alias := range aliases {
				valIndex := uint32(len(pk.EnumValues))
				pk.EnumValues = append(pk.EnumValues, strings.ToLower(alias))
				aliasEnum.ValueIndices = append(aliasEnum.ValueIndices, valIndex)
			}
			aliasOffset = uint32(len(pk.Enums))
			pk.Enums = append(pk.Enums, aliasEnum)
		}

		found := false
		for i, existing := range pk.Commands {
			if strings.EqualFold(existing.Name, name) {
				pk.Commands[i].Description = cmdInstance.Description()
				pk.Commands[i].PermissionLevel = byte(p.PermissionLevel().Level())
				pk.Commands[i].Overloads = baseOverloads
				pk.Commands[i].AliasesOffset = aliasOffset
				found = true
				break
			}
		}
		if !found {
			pk.Commands = append(pk.Commands, protocol.Command{
				Name:            strings.ToLower(name),
				Description:     cmdInstance.Description(),
				PermissionLevel: byte(p.PermissionLevel().Level()),
				Overloads:       baseOverloads,
				AliasesOffset:   aliasOffset,
			})
		}
	}
	_ = p.client.WritePacket(pk)
}

// SendSimpleForm sends a simple (button list) form to the player's client.
func (p *Player) SendSimpleForm(f form.SimpleForm) error {
	if err := form.Validate(f); err != nil {
		logger.Printf(err.Error())
		return err
	}
	id := p.nextID(f)
	data, _ := json.Marshal(struct {
		Type    string        `json:"type"`
		Title   string        `json:"title"`
		Content string        `json:"content"`
		Buttons []form.Button `json:"buttons"`
	}{
		Type:    "form",
		Title:   f.Title(),
		Content: "",
		Buttons: form.ElementsToButtons(f.Elements()),
	})
	return p.client.WritePacket(&packet.ModalFormRequest{FormID: id, FormData: data})
}

// SendCustomForm sends a complex custom form to the player's client.
func (p *Player) SendCustomForm(f form.CustomForm) error {
	if err := form.Validate(f); err != nil {
		logger.Printf(err.Error())
		return err
	}
	id := p.nextID(f)
	data, _ := json.Marshal(struct {
		Type    string         `json:"type"`
		Title   string         `json:"title"`
		Content []form.Element `json:"content"`
	}{
		Type:    "custom_form",
		Title:   f.Title(),
		Content: f.Elements(),
	})
	return p.client.WritePacket(&packet.ModalFormRequest{FormID: id, FormData: data})
}

// nextID generates and stores a new unique ID for a UI form.
func (p *Player) nextID(f interface{}) uint32 {
	p.formsMu.Lock()
	defer p.formsMu.Unlock()
	p.lastFormID++
	p.forms[p.lastFormID] = f
	return p.lastFormID
}

// Disconnect forcefully closes the connection with an optional message.
func (p *Player) Disconnect(msg string) error {
	return p.client.WritePacket(&packet.Disconnect{
		Message:                 msg,
		Reason:                  packet.DisconnectReasonDisconnected,
		HideDisconnectionScreen: msg == "",
	})
}

// SendBossbar displays a boss health bar to the player.
func (p *Player) SendBossbar(b Bossbar) error {
	r, g, bCol, _ := b.Colour().RGBA()

	var colorID uint32
	switch {
	case r > 0 && g == 0 && bCol == 0:
		colorID = 1
	case r == 0 && g > 0 && bCol == 0:
		colorID = 3
	case r == 0 && g == 0 && bCol > 0:
		colorID = 4
	case r > 0 && g > 0 && bCol == 0:
		colorID = 2
	case r > 0 && g == 0 && bCol > 0:
		colorID = 5
	case r == 0 && g > 0 && bCol > 0:
		colorID = 6
	default:
		colorID = 0
	}
	return p.client.WritePacket(&packet.BossEvent{
		BossEntityUniqueID: p.client.GameData().EntityUniqueID,
		EventType:          packet.BossEventShow,
		BossBarTitle:       b.Text(),
		HealthPercentage:   float32(b.Health()) / 100,
		Colour:             colorID,
		Overlay:            uint32(len(b.Text())),
	})
}

// RemoveBossbar removes any active boss bar from the player's screen.
func (p *Player) RemoveBossbar() error {
	return p.client.WritePacket(&packet.BossEvent{
		BossEntityUniqueID: p.client.GameData().EntityUniqueID,
		EventType:          packet.BossEventHide,
	})
}

// HideHud triggers a texture animation to hide specific HUD elements.
func (p *Player) HideHud(h hud.Hud) error {
	return p.client.WritePacket(&packet.OnScreenTextureAnimation{
		AnimationType: uint32(h.Type()),
	})
}

// ShowHud triggers a texture animation to show specific HUD elements.
func (p *Player) ShowHud(h hud.Hud) error {
	return p.client.WritePacket(&packet.OnScreenTextureAnimation{
		AnimationType: uint32(h.Type()),
	})
}

// Transfer initiates a server transfer for the player.
func (p *Player) Transfer(addr string, port uint16) error {
	if port == 0 {
		port = 19132
	}
	return p.client.WritePacket(&packet.Transfer{
		Address: addr,
		Port:    port,
	})
}

// AddEffect applies a status effect to the player.
func (p *Player) AddEffect(e effect.Effect) error {
	return p.client.WritePacket(&packet.MobEffect{
		EntityRuntimeID: p.client.GameData().EntityRuntimeID,
		Operation:       packet.MobEffectAdd,
		EffectType:      e.Type(),
		Amplifier:       int32(e.Force()),
		Duration:        int32(e.Duration()),
		Particles:       e.Particles(),
	})
}

// RemoveEffect removes a status effect from the player.
func (p *Player) RemoveEffect(e effect.Effect) error {
	return p.client.WritePacket(&packet.MobEffect{
		EntityRuntimeID: p.client.GameData().EntityRuntimeID,
		Operation:       packet.MobEffectRemove,
		EffectType:      e.Type(),
	})
}

// Execute triggers a command request from the client-side.
func (p *Player) Execute(cmd string) error {
	return p.client.WritePacket(&packet.CommandRequest{
		CommandLine: cmd,
		Internal:    true,
		Version:     protocol.CurrentVersion,
	})
}

// Sprinting checks if the player is currently sprinting.
func (p *Player) Sprinting() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	inputs := p.se.Inputs()
	if inputs == nil {
		return false
	}
	return inputs.InputData.Load(packet.InputFlagSprinting)
}

// Sneaking checks if the player is currently sneaking.
func (p *Player) Sneaking() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	inputs := p.se.Inputs()
	if inputs == nil {
		return false
	}
	return inputs.InputData.Load(packet.InputFlagSneaking)
}

// Walking checks if the player is currently moving forward.
func (p *Player) Walking() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	inputs := p.se.Inputs()
	if inputs == nil {
		return false
	}
	return inputs.InputData.Load(packet.InputFlagUp)
}

// Jumping checks if the player is currently attempting to jump.
func (p *Player) Jumping() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	inputs := p.se.Inputs()
	if inputs == nil {
		return false
	}
	return inputs.InputData.Load(packet.InputFlagWantUp)
}

// JumpSprinting checks if the player is both jumping and sprinting.
func (p *Player) JumpSprinting() bool {
	return p.Jumping() && p.Sprinting()
}

// doAction simulates player input flags for a specific duration.
func (p *Player) doAction(elapsed time.Duration, flags ...uint32) error {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	endTime := time.Now().Add(elapsed)
	for time.Now().Before(endTime) {
		p.mu.Lock()
		ptr := p.se.Inputs()
		p.mu.Unlock()
		if ptr == nil {
			return fmt.Errorf("empty input bitset")
		}
		current := *ptr
		bitset := protocol.NewBitset(packet.PlayerAuthInputBitsetSize)
		for _, flag := range flags {
			bitset.Set(int(flag))
		}
		pk := &packet.PlayerAuthInput{
			Pitch:      current.Pitch,
			Yaw:        current.Yaw,
			Position:   current.Position,
			MoveVector: current.MoveVector,
			HeadYaw:    current.HeadYaw,
			InputData:  bitset,
			InputMode:  current.InputMode,
			Tick:       current.Tick + 1,
			Delta:      current.Delta,
		}
		if err := p.server.WritePacket(pk); err != nil {
			return err
		}
		<-ticker.C
	}
	return nil
}

// DoWalk simulates the walk input for the specified duration.
func (p *Player) DoWalk(elapsed time.Duration) error {
	return p.doAction(elapsed, packet.InputFlagUp)
}

// DoJump simulates the jump input for the specified duration.
func (p *Player) DoJump(elapsed time.Duration) error {
	return p.doAction(elapsed, packet.InputFlagWantUp)
}

// DoSneak simulates the sneak input for the specified duration.
func (p *Player) DoSneak(elapsed time.Duration) error {
	return p.doAction(elapsed, packet.InputFlagUp, packet.InputFlagSneaking)
}

// DoSprint simulates the sprint input for the specified duration.
func (p *Player) DoSprint(elapsed time.Duration) error {
	return p.doAction(elapsed, packet.InputFlagUp, packet.InputFlagSprinting)
}

// DoSprintJump simulates a sprint-jump action for the specified duration.
func (p *Player) DoSprintJump(elapsed time.Duration) error {
	return p.doAction(elapsed, packet.InputFlagUp, packet.InputFlagSprinting, packet.InputFlagWantUp)
}

// ClientLatency returns the measured network latency for the client.
func (p *Player) ClientLatency() time.Duration {
	return p.client.Latency()
}

// ServerLatency returns the measured network latency for the destination server.
func (p *Player) ServerLatency() time.Duration {
	return p.server.Latency()
}
