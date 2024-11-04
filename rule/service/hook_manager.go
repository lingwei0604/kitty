package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/pkg/errors"
)

type hooks struct {
	OnChange   string `json:"onChange" yaml:"onChange"`
	PreUpdate  string `json:"preUpdate" yaml:"preUpdate"`
	PostUpdate string `json:"postUpdate" yaml:"postUpdate"`
}

type hookResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type HookManager struct {
	doer  contract.HttpDoer
	hooks hooks
	data  []byte
}

func NewHookManager(doer contract.HttpDoer, data []byte) (*HookManager, error) {
	var hooks hooks
	c := koanf.New(".")

	if err := c.Load(rawbytes.Provider(data), kyaml.Parser()); err != nil {
		return nil, errors.Wrap(err, "error creating HookManager")
	}
	if err := c.Unmarshal("hooks", &hooks); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling hooks")
	}

	return &HookManager{doer: doer, hooks: hooks, data: data}, nil
}

func (h *HookManager) OnChange() error {
	if h.hooks.OnChange == "" {
		return nil
	}
	return h.callHook(h.hooks.OnChange)
}

func (h *HookManager) PreUpdate() error {
	if h.hooks.PreUpdate == "" {
		return nil
	}
	return h.callHook(h.hooks.PreUpdate)
}

func (h *HookManager) PostUpdate() error {
	if h.hooks.PostUpdate == "" {
		return nil
	}
	return h.callHook(h.hooks.PostUpdate)
}

func (h *HookManager) callHook(url string) error {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(h.data))
	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/yaml; charset=utf-8")
	resp, err := h.doer.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to contact hook server")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("hook server returns non-200 code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response body from hook server")
	}
	var result hookResult
	if err := json.Unmarshal(body, &result); err != nil {
		return errors.Wrap(err, "unable to unmarshal hook results")
	}
	if result.Code != 0 {
		return fmt.Errorf("hook server returns \"%s\" with code %d", result.Msg, result.Code)
	}
	return nil
}
