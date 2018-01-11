package http

import (
	"net/http"

	"github.com/romana/rlog"

	"github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/acl"
	"github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/distributor"
	"github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/sofia"
)

var (
	configuration configurationHandler
)

type configurationHandler struct{}

// requestForm is the equivalent of a parsed POST form
type requestForm map[string][]string

// get the `key` from the request
func (c requestForm) Get(key string) string {
	if vs := c[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (configurationHandler) Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cr := requestForm(r.PostForm)
	ctx := r.Context()

	var err error
	switch cr.Get("key_value") {
	case "acl.conf":
		err = acl.Handler(ctx, cr.Get("hostname"), w)
	case "distributor.conf":
		err = distributor.Handler(ctx, cr.Get("hostname"), w)
	case "sofia.conf":
		err = sofia.Handler(ctx, cr.Get("hostname"), w)
	default:
		rlog.Infof("configuration request not supported [%s]", cr.Get("key_value"))
		notFound(w)
		return
	}

	// check for error
	if err != nil {
		rlog.Errorf("could not load module configuration [%s]", err.Error())
		notFound(w)
	}
	return
}
