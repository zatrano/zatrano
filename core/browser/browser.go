package browser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zatrano/framework/core"
	testkit "github.com/zatrano/framework/core/testing"
)

// Browser is a lightweight browser-testing toolkit stub.
type Browser struct {
	tc      *testkit.TestCase
	last    *testkit.TestResponse
	path    string
	form    map[string]string
	headers map[string]string
}

// New creates a browser session around an application.
func New(app *core.Application) (*Browser, error) {
	tc, err := testkit.New(app)
	if err != nil {
		return nil, err
	}
	return &Browser{
		tc:      tc,
		form:    map[string]string{},
		headers: map[string]string{},
	}, nil
}

// Visit performs a GET request.
func (b *Browser) Visit(path string) *Browser {
	b.path = path
	b.last = b.tc.Get(path)
	return b
}

// AssertOK asserts the last response was 200.
func (b *Browser) AssertOK() *Browser {
	b.ensure()
	b.last.AssertOK()
	return b
}

// AssertStatus asserts the last response status.
func (b *Browser) AssertStatus(status int) *Browser {
	b.ensure()
	b.last.AssertStatus(status)
	return b
}

// AssertSee asserts the body contains text.
func (b *Browser) AssertSee(text string) *Browser {
	b.ensure()
	if !strings.Contains(b.last.String(), text) {
		panic(fmt.Sprintf("browser: expected to see %q in body", text))
	}
	return b
}

// AssertDontSee asserts the body does not contain text.
func (b *Browser) AssertDontSee(text string) *Browser {
	b.ensure()
	if strings.Contains(b.last.String(), text) {
		panic(fmt.Sprintf("browser: did not expect to see %q", text))
	}
	return b
}

// AssertPathIs asserts the current path.
func (b *Browser) AssertPathIs(path string) *Browser {
	if b.path != path {
		panic(fmt.Sprintf("browser: expected path %q, got %q", path, b.path))
	}
	return b
}

// Type queues a form field value.
func (b *Browser) Type(name, value string) *Browser {
	b.form[name] = value
	return b
}

// Press submits the queued form via POST to the current path (or action if provided).
func (b *Browser) Press(button string, action ...string) *Browser {
	_ = button
	target := b.path
	if len(action) > 0 && action[0] != "" {
		target = action[0]
	}
	b.path = target
	b.last = b.tc.Post(target, b.form)
	b.form = map[string]string{}
	return b
}

// ClickLink follows the first anchor whose text matches.
func (b *Browser) ClickLink(text string) *Browser {
	b.ensure()
	re := regexp.MustCompile(`(?is)<a[^>]+href=["']([^"']+)["'][^>]*>\s*` + regexp.QuoteMeta(text) + `\s*</a>`)
	m := re.FindStringSubmatch(b.last.String())
	if m == nil {
		// looser: any href near the text
		re2 := regexp.MustCompile(`(?is)href=["']([^"']+)["'][^>]*>[^<]*` + regexp.QuoteMeta(text))
		m = re2.FindStringSubmatch(b.last.String())
	}
	if m == nil {
		panic(fmt.Sprintf("browser: link %q not found", text))
	}
	return b.Visit(m[1])
}

// Status returns the last status code.
func (b *Browser) Status() int {
	b.ensure()
	return b.last.StatusCode
}

// Body returns the last response body.
func (b *Browser) Body() string {
	b.ensure()
	return b.last.String()
}

// Response returns the underlying test response.
func (b *Browser) Response() *testkit.TestResponse {
	return b.last
}

func (b *Browser) ensure() {
	if b.last == nil {
		panic("browser: no response yet; call Visit first")
	}
}
