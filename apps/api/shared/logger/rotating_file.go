package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type RotatingFileConfig struct {
	Path       string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
}

type RotatingFileWriter struct {
	mu         sync.Mutex
	path       string
	maxSize    int64
	maxBackups int
	maxAge     time.Duration
	compress   bool
	file       *os.File
}

func NewRotatingFileWriter(cfg RotatingFileConfig) (*RotatingFileWriter, error) {
	if cfg.MaxSizeMB <= 0 {
		cfg.MaxSizeMB = 100
	}
	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = 7
	}
	if cfg.MaxAgeDays <= 0 {
		cfg.MaxAgeDays = 14
	}
	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil {
		return nil, err
	}
	writer := &RotatingFileWriter{
		path:       cfg.Path,
		maxSize:    int64(cfg.MaxSizeMB) * 1024 * 1024,
		maxBackups: cfg.MaxBackups,
		maxAge:     time.Duration(cfg.MaxAgeDays) * 24 * time.Hour,
		compress:   cfg.Compress,
	}
	if err := writer.open(); err != nil {
		return nil, err
	}
	return writer, nil
}

func (w *RotatingFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		if err := w.open(); err != nil {
			return 0, err
		}
	}
	if info, err := w.file.Stat(); err == nil && info.Size()+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	return w.file.Write(p)
}

func (w *RotatingFileWriter) open() error {
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

func (w *RotatingFileWriter) rotate() error {
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return err
		}
		w.file = nil
	}
	if _, err := os.Stat(w.path); err == nil {
		rotatedPath := w.rotatedPath()
		if err := os.Rename(w.path, rotatedPath); err != nil {
			return err
		}
		if w.compress {
			_ = gzipFile(rotatedPath)
		}
	}
	if err := w.cleanup(); err != nil {
		return err
	}
	return w.open()
}

func (w *RotatingFileWriter) rotatedPath() string {
	ext := filepath.Ext(w.path)
	base := strings.TrimSuffix(w.path, ext)
	return fmt.Sprintf("%s-%s%s", base, time.Now().Format("20060102-150405.000000000"), ext)
}

func (w *RotatingFileWriter) cleanup() error {
	files, err := filepath.Glob(w.backupGlob())
	if err != nil {
		return err
	}
	type backup struct {
		path    string
		modTime time.Time
	}
	backups := make([]backup, 0, len(files))
	cutoff := time.Now().Add(-w.maxAge)
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(path)
			continue
		}
		backups = append(backups, backup{path: path, modTime: info.ModTime()})
	}
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].modTime.After(backups[j].modTime)
	})
	for i := w.maxBackups; i < len(backups); i++ {
		_ = os.Remove(backups[i].path)
	}
	return nil
}

func (w *RotatingFileWriter) backupGlob() string {
	ext := filepath.Ext(w.path)
	base := strings.TrimSuffix(w.path, ext)
	return base + "-*" + ext + "*"
}

func gzipFile(path string) error {
	source, err := os.Open(path)
	if err != nil {
		return err
	}
	defer source.Close()

	targetPath := path + ".gz"
	target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer target.Close()

	gzipWriter := gzip.NewWriter(target)
	if _, err := io.Copy(gzipWriter, source); err != nil {
		_ = gzipWriter.Close()
		return err
	}
	if err := gzipWriter.Close(); err != nil {
		return err
	}
	return os.Remove(path)
}
