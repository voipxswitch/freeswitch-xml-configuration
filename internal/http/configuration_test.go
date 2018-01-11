package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

    "github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/acl"
    "github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/distributor"
    "github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/sofia"
)

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	moduleData := filepath.Join(wd, "../../moduledata")
	templatePath := filepath.Join(wd, "../../templates")

	// init each module for testing
	acl.New(moduleData, templatePath)
	distributor.New(moduleData, templatePath)
	sofia.New(moduleData, templatePath)

	notFoundTemplatePath = filepath.Join(templatePath, notFoundTemplate)

	os.Exit(m.Run())
}

func TestConfigHandlerAcl(t *testing.T) {
	expect := `<document type="freeswitch/xml">
    <section name="configuration" description="FreeSWITCH Configuration">
        <configuration name="acl.conf" description="Network Lists">
            <network-lists>
                <list name="lan" default="allow">
                    <node type="deny" cidr="192.168.42.0/24"/>
                </list>
                <list name="proxy" default="deny">
                    <node type="allow" cidr="10.10.10.11/32"/>
                    <node type="allow" cidr="10.10.10.12/32"/>
                </list>
            </network-lists>
        </configuration>
    </section>
</document>
`

	//create fake request
	form := url.Values{} // Create fake form (as if it was posted)
	form.Add("hostname", "fs-01")
	form.Add("section", "configuration")
	form.Add("tag_name", "configuration")
	form.Add("key_name", "name")
	form.Add("key_value", "acl.conf")
	r, _ := http.NewRequest("POST", "http://nowhere.local", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	configuration.Handler(w, r)
	if w.Body.String() != expect {
		t.Errorf("\n\nExpected:\n%s\n\nGot:\n%s\n", expect, w.Body.String())
	}
}

func TestConfigHandlerAclNotFound(t *testing.T) {
	expect := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
    <section name="result">
        <result status="not found" />
    </section>
</document>
`

	//create fake request
	form := url.Values{} // Create fake form (as if it was posted)
	form.Add("hostname", "fs-02")
	form.Add("section", "configuration")
	form.Add("tag_name", "configuration")
	form.Add("key_name", "name")
	form.Add("key_value", "acl.conf")
	r, _ := http.NewRequest("POST", "http://nowhere.local", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	configuration.Handler(w, r)
	if w.Body.String() != expect {
		t.Errorf("\n\nExpected:\n%s\n\nGot:\n%s\n", expect, w.Body.String())
	}
}

func TestConfigHandlerDistributor(t *testing.T) {
	expect := `<document type="freeswitch/xml">
    <section name="configuration" description="FreeSWITCH Configuration">
        <configuration name="distributor.conf" description="Distributor Configuration">
            <lists>
                <list name="proxy" default="2">
                    <node name="proxy-01.local" weight="1"/>
                    <node name="proxy-02.local" weight="1"/>
                </list>
            </lists>
        </configuration>
    </section>
</document>
`

	//create fake request
	form := url.Values{} // Create fake form (as if it was posted)
	form.Add("hostname", "fs-01")
	form.Add("section", "configuration")
	form.Add("tag_name", "configuration")
	form.Add("key_name", "name")
	form.Add("key_value", "distributor.conf")
	r, _ := http.NewRequest("POST", "http://nowhere.local", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	configuration.Handler(w, r)
	if w.Body.String() != expect {
		t.Errorf("\n\nExpected:\n%s\n\nGot:\n%s\n", expect, w.Body.String())
	}
}

func TestConfigHandlerDistributorNotFound(t *testing.T) {
	expect := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
    <section name="result">
        <result status="not found" />
    </section>
</document>
`

	//create fake request
	form := url.Values{} // Create fake form (as if it was posted)
	form.Add("hostname", "fs-02")
	form.Add("section", "configuration")
	form.Add("tag_name", "configuration")
	form.Add("key_name", "name")
	form.Add("key_value", "distributor.conf")
	r, _ := http.NewRequest("POST", "http://nowhere.local", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	configuration.Handler(w, r)
	if w.Body.String() != expect {
		t.Errorf("\n\nExpected:\n%s\n\nGot:\n%s\n", expect, w.Body.String())
	}
}

func TestConfigHandlerSofia(t *testing.T) {
	expect := `<document type="freeswitch/xml">
    <section name="configuration" description="FreeSWITCH Configuration">
        <configuration name="sofia.conf" description="sofia Endpoint">
            <global_settings>
                <param name="log-level" value="4"/>
                <param name="debug-presence" value="0"/>
                <param name="capture-server" value="udp:127.0.0.1:9060;hep=3"/>
            </global_settings>
            <profiles>
                <profile name="internal">
                    <aliases></aliases>
                    <domains>
                        <domain name="all" alias="true" parse="false"/>
                    </domains>
                    <gateways>
                        <gateway name="proxy-01.local">
                            <param name="register" value="false"/>
                            <param name="username" value="$${hostname}"/>
                            <param name="ping" value="20"/>
                        </gateway>
                        <gateway name="proxy-02.local">
                            <param name="register" value="false"/>
                            <param name="username" value="$${hostname}"/>
                            <param name="ping" value="20"/>
                        </gateway>
                    </gateways>
                    <settings>
                        <param name="debug" value="0"/>
                        <param name="sip-trace" value="no"/>
                        <param name="sip-capture" value="yes"/>
                        <param name="track-calls" value="true"/>
                        <param name="enable-timer" value="true"/>
                        <param name="session-timeout" value="600"/>
                        <param name="enable-compact-headers" value="true"/>
                        <param name="caller-id-type" value="pid"/>
                        <param name="context" value="internal"/>
                        <param name="sip-port" value="5060"/>
                        <param name="sip-ip" value="$${public_sip_ip}"/>
                        <param name="rtp-ip" value="$${public_rtp_ip}"/>
                        <param name="rtp-timeout-sec" value="6000"/>
                        <param name="rtp-hold-timeout-sec" value="1800"/>
                        <param name="rtp-timer-name" value="soft"/>
                        <param name="dialplan" value="XML"/>
                        <param name="dtmf-duration" value="2000"/>
                        <param name="rfc2833-pt" value="101"/>
                        <param name="inbound-codec-prefs" value="$${global_codec_prefs}"/>
                        <param name="outbound-codec-prefs" value="$${global_codec_prefs}"/>
                        <param name="inbound-codec-negotiation" value="generous"/>
                        <param name="inbound-codec-negotiation" value="generous"/>
                        <param name="log-auth-failures" value="true"/>
                        <param name="forward-unsolicited-mwi-notify" value="false"/>
                        <param name="hold-music" value="$${hold_music}"/>
                        <param name="apply-inbound-acl" value="proxy"/>
                        <param name="local-network-acl" value="localnet.auto"/>
                        <param name="manage-presence" value="true"/>
                        <param name="sip-options-respond-503-on-busy" value="true"/>
                        <param name="sip-messages-respond-200-ok" value="true"/>
                        <param name="send-display-update" value="false"/>
                        <param name="record-path" value="/var/lib/freeswitch/recordings"/>
                        <param name="record-template" value="${caller_id_number}.${target_domain}.${strftime(%Y-%m-%d-%H-%M-%S)}.wav"/>
                        <param name="inbound-late-negotiation" value="true"/>
                        <param name="inbound-zrtp-passthru" value="true"/>
                        <param name="auth-calls" value="false"/>
                        <param name="auth-all-packets" value="false"/>
                        <param name="challenge-realm" value="auto_from"/>
                        <param name="inbound-reg-force-matching-username" value="true"/>
                        <param name="tls" value="false"/>
                    </settings>
                </profile>
            </profiles>
        </configuration>
    </section>
</document>
`

	//create fake request
	form := url.Values{} // Create fake form (as if it was posted)
	form.Add("hostname", "fs-01")
	form.Add("section", "configuration")
	form.Add("tag_name", "configuration")
	form.Add("key_name", "name")
	form.Add("key_value", "sofia.conf")
	r, _ := http.NewRequest("POST", "http://nowhere.local", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	configuration.Handler(w, r)
	if w.Body.String() != expect {
		t.Errorf("\n\nExpected:\n%s\n\nGot:\n%s\n", expect, w.Body.String())
	}
}

func TestConfigHandlerSofiaNotFound(t *testing.T) {
	expect := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
    <section name="result">
        <result status="not found" />
    </section>
</document>
`

	//create fake request
	form := url.Values{} // Create fake form (as if it was posted)
	form.Add("hostname", "fs-02")
	form.Add("section", "configuration")
	form.Add("tag_name", "configuration")
	form.Add("key_name", "name")
	form.Add("key_value", "sofia.conf")
	r, _ := http.NewRequest("POST", "http://nowhere.local", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	configuration.Handler(w, r)
	if w.Body.String() != expect {
		t.Errorf("\n\nExpected:\n%s\n\nGot:\n%s\n", expect, w.Body.String())
	}
}
