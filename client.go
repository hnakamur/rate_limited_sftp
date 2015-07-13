package rate_limited_sftp

import (
	"github.com/mxk/go-flowrate/flowrate"
	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

func NewClient(conn *ssh.Client, readLimitBytesPerSecond, writeLimitBytesPerSecond int64, opts ...func(*sftp.Client) error) (*sftp.Client, error) {
	s, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	if err := s.RequestSubsystem("sftp"); err != nil {
		return nil, err
	}
	pw, err := s.StdinPipe()
	if err != nil {
		return nil, err
	}
	pr, err := s.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if readLimitBytesPerSecond > 0 {
		pr = flowrate.NewReader(pr, readLimitBytesPerSecond)
	}
	if writeLimitBytesPerSecond > 0 {
		pw = flowrate.NewWriter(pw, writeLimitBytesPerSecond)
	}
	return sftp.NewClientPipe(pr, pw, opts...)
}
