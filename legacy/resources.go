package legacy

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"golang.org/x/oauth2"
)

// Resources ...
type Resources struct {
	cfg   *Config
	src   oauth2.TokenSource
	raddr string
}

// Load returns all the resources pack that will be applied to the proxy.
func (r *Resources) Load() ([]*resource.Pack, error) {
	var packs []*resource.Pack
	identity := login.IdentityData{
		DisplayName: "",
	}
	dialer := minecraft.Dialer{
		TokenSource:  r.src,
		IdentityData: identity,
	}
	conn, err := dialer.Dial("raknet", r.raddr)
	if err == nil {
		packs = conn.ResourcePacks()
		_ = conn.Close()
	} else {
		logger.Printf("resources lookup: %v", err)
	}
	dir := r.cfg.Resources.Path
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return packs, nil
		}
		return nil, fmt.Errorf("read resources folder: %v", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == "cache" {
			continue
		}
		path := filepath.Join(dir, name)
		pack, err := resource.ReadPath(path)
		if err != nil {
			logger.Printf("load resource pack '%s': %v", name, err)
			continue
		}
		packs = append(packs, pack)
		logger.Printf("pack loaded '%s'", pack.Manifest().Header.Name)
	}
	return packs, nil
}

// clean cleans up unsupported characters of resources pack names for windows.
func clean(filename string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	filename = re.ReplaceAllString(filename, "_")
	filename = strings.TrimSpace(filename)
	re = regexp.MustCompile(`\s+`)
	filename = re.ReplaceAllString(filename, "_")
	return filename
}
