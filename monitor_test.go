package monitor

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }
type FileMonitorSuite struct {}
var _ = Suite(&FileMonitorSuite{})

func (a *FileMonitorSuite) TestX(c *C) {
	c.Check(true, Equals, true)
}
