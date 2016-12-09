package mtty_test

import (
	"bytes"

	gc "gopkg.in/check.v1"
	stdtesting "testing"

	"github.com/mattyw/mtty"
)

type mttySuite struct{}

func TestAll(t *stdtesting.T) {
	gc.TestingT(t)
}

var _ = gc.Suite(&mttySuite{})

func (s *mttySuite) TestOpen(c *gc.C) {
	out := new(bytes.Buffer)
	m := mtty.Mtty{
		Stdout: out,
		Stderr: new(bytes.Buffer),
		Stdin:  new(bytes.Buffer),
	}
	m.SetLastOut([]byte(`
    hello there
    foo bar foo
    foo_bar.go:17
    foo_bar.go:17:37
    foo
`))
	m.Open()
	c.Assert(string(out.Bytes()), gc.Equals, `0) foo_bar.go:17
1) foo_bar.go:17:37
`)
}
