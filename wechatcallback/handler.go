package wechatcallback

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.yingzhongshare.com/mkt/kitty/app/entity"
	"git.yingzhongshare.com/mkt/kitty/app/msg"
	"git.yingzhongshare.com/mkt/kitty/app/repository"
	"git.yingzhongshare.com/mkt/kitty/pkg/kerr"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type UserRepository interface {
	GetUserByOpenID(ctx context.Context, packageName string, openID string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}

// Handler handles wechat callback
type Handler struct {
	ur     UserRepository
	token  string
	logger log.Logger
}

func NewHandler(ur UserRepository, logger log.Logger) *Handler {
	return &Handler{ur: ur, token: "donews123", logger: logger}
}

func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
	if h.signatureValidate(r) {
		io.WriteString(w, r.FormValue("echostr"))
	}
}

// UnbindWechatUser unbind wechat user
func (h *Handler) UnbindWechatUser(w http.ResponseWriter, r *http.Request) {
	if !h.signatureValidate(r) {
		h.logger.Log("msg", "signature validation failed")
		return
	}
	packageName := mux.Vars(r)["packageName"]
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Log("err", err, "msg", "read body failed")
		return
	}
	defer r.Body.Close()
	var payload WeChatUserInfoChangedPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		h.logger.Log("err", err, "msg", "UnbindWechatUser json marshal error", "body", string(data))
		return
	}
	if payload.Event != "user_authorization_revoke" {
		h.logger.Log("msg", "user authorization not revoked", "event", payload.Event)
		return
	}

	if err := h.unbindOpenID(r.Context(), packageName, payload.OpenID); err != nil {
		h.logger.Log("err", err, "msg", "unable to unbind openID", "packageName", packageName, "open_id", payload.OpenID)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func (h *Handler) unbindOpenID(ctx context.Context, packageName string, openID string) error {
	user, err := h.ur.GetUserByOpenID(ctx, packageName, openID)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return kerr.NotFoundErr(err, msg.ErrorRecordNotFound)
	}
	if err != nil {
		return err
	}

	user.WechatUnionId = sql.NullString{}
	user.WechatOpenId = sql.NullString{}
	user.WechatExtra = nil

	err = h.ur.Save(ctx, user)
	if err != nil {
		return err
	}
	// TODO: dispatch user update event

	return nil
}

func (h *Handler) signatureValidate(r *http.Request) bool {
	r.ParseForm()
	signature := r.Form.Get("signature")
	timestamp := r.Form.Get("timestamp")
	nonce := r.Form.Get("nonce")
	echostr := r.Form.Get("echostr")
	if signature != "" && timestamp != "" && nonce != "" && echostr != "" {
		if h.token != "" {
			if h.checkSignature(signature, timestamp, nonce) {
				return true
			}
		}
	}
	return false
}

// check Signature checks wechat signature
func (h *Handler) checkSignature(signature string, timestamp string, nonce string) bool {
	str := fmt.Sprintf("%s%s%s", h.token, timestamp, nonce)
	return signature == GetSha1(str)
}

func GetSha1(str string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(str))
	return fmt.Sprintf("%x", sha1.Sum(nil))
}
