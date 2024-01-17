package sftpprovider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/sftp"
	"go-kit/src/common/configs"
	"go-kit/src/common/logger"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type SftpProvider struct {
	config *configs.Sftp
	client *sftp.Client
}

func NewSftpProvider(cf *configs.Config) *SftpProvider {
	return &SftpProvider{
		config: cf.Sftp,
	}
}

func (s *SftpProvider) connect(ctx context.Context) error {
	auths := make([]ssh.AuthMethod, 0)
	auths = append(auths, ssh.Password(s.config.Pass))
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	conn, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            s.config.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	})
	if err != nil {
		logger.Error(ctx, err, "Failed to connect to host [%s]", addr)
		return fmt.Errorf("[SFTP] failed to connect to host %s, %w, ", addr, err)
	}

	// Create new SFTP client
	s.client, err = sftp.NewClient(conn)
	if err != nil {
		logger.Error(ctx, err, "Unable to start SFTP subsystem")
		return fmt.Errorf("[SFTP] unable to start SFTP subsystem %w", err)
	}
	return nil
}

func (s *SftpProvider) disconnect(ctx context.Context) {
	err := s.client.Close()
	if err != nil {
		logger.Error(ctx, err, "Failed to close sftp connection")
	}
}

func (s *SftpProvider) ReadDir(ctx context.Context, path string) ([]os.FileInfo, error) {
	logger.Debug(ctx, "Get all file path:[%v]", path)
	err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer s.disconnect(ctx)

	files, err := s.client.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *SftpProvider) CreateFile(ctx context.Context, path string, buff *bytes.Buffer) error {
	logger.Debug(ctx, "Create file path:[%v]", path)
	err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer s.disconnect(ctx)

	file, err := s.client.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}

	_, err = file.Write(buff.Bytes())
	if err != nil {
		return err
	}
	return nil
}
