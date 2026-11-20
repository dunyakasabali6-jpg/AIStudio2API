package app

import (
	"context"
	"errors"

	"github.com/Mag1cFall/AIStudio2API/internal/aistudio"
	"github.com/Mag1cFall/AIStudio2API/internal/api"
	"github.com/Mag1cFall/AIStudio2API/internal/camoufoxnative"
	"github.com/Mag1cFall/AIStudio2API/internal/config"
)

var errBrowserLoginDisabled = errors.New("browser-based login is disabled in the server deployment; import an AI Studio Cookie instead")

type disabledLoginDriver struct{}

func (disabledLoginDriver) Login(context.Context, aistudio.IsolatedLoginRequest) (aistudio.IsolatedLoginResult, error) {
	return aistudio.IsolatedLoginResult{}, errBrowserLoginDisabled
}

// prepareWAABackend selects the in-process SpiderMonkey proof backend for the server build.
// Camoufox is deliberately not prepared or downloaded in this deployment profile.
func prepareWAABackend(_ config.Config, requests *requestRegistry, accounts []*aistudio.Account) aistudio.IsolatedLoginDriver {
	requests.logStage("service", "INFO", api.LogStage{Name: "运行时装配", Step: 2, Steps: 3}, "使用 SpiderMonkey WAA 后端", "accounts", len(accounts))
	return disabledLoginDriver{}
}

// newWAAWorker always starts the in-process SpiderMonkey worker; browser executables are not supported.
func newWAAWorker(ctx context.Context, accountID string, options camoufoxnative.Options) (*aistudio.NativeWorker, error) {
	options.ExecutablePath = ""
	return aistudio.NewSpiderMonkeyWorker(ctx, accountID, options)
}
