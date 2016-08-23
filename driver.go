package main

import (
	"os"
	"path/filepath"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/ehazlett/libsecret/store"
)

type SecretDriver struct {
	root         string
	fs           map[string]*FS
	storeAddr    string
	storeBackend store.Backend
	storeOpts    map[string]interface{}
}

func NewSecretDriver(root string, backend store.Backend, addr string, opts map[string]interface{}) (*SecretDriver, error) {
	log.Debugf("opts: %v", opts)

	log.Debugf("backend: %s", backend)

	return &SecretDriver{
		root:         root,
		storeAddr:    addr,
		storeBackend: backend,
		storeOpts:    opts,
		fs:           map[string]*FS{},
	}, nil
}

func (d *SecretDriver) resolvePath(name string) string {
	return filepath.Join(d.root, name)
}

func (d *SecretDriver) Create(r volume.Request) volume.Response {
	log.Debugf("create: %v", r)

	p := d.resolvePath(r.Name)

	errStr := ""
	if err := os.MkdirAll(p, 0755); err != nil {
		errStr = err.Error()
	}

	return volume.Response{Err: errStr}
}

func (d *SecretDriver) Remove(r volume.Request) volume.Response {
	log.Debugf("remove: %v", r)

	p := d.resolvePath(r.Name)

	errStr := ""
	if err := os.RemoveAll(p); err != nil {
		errStr = err.Error()
	}

	return volume.Response{Err: errStr}
}

func (d *SecretDriver) Path(r volume.Request) volume.Response {
	log.Debugf("path: %v", r)

	p := d.resolvePath(r.Name)
	return volume.Response{Mountpoint: p}
}

func (d *SecretDriver) Mount(r volume.MountRequest) volume.Response {
	log.Debugf("mount: %v", r)

	p := d.resolvePath(r.Name)

	errStr := ""

	fs, err := NewFS(p, d.storeBackend, d.storeAddr, d.storeOpts)
	if err != nil {
		errStr = err.Error()
	}

	if err := fs.Mount(r.Name); err != nil {
		errStr = err.Error()
	}

	d.fs[r.Name] = fs

	return volume.Response{
		Mountpoint: filepath.Join(d.root, r.Name),
		Err:        errStr,
	}
}

func (d *SecretDriver) Capabilities(r volume.Request) volume.Response {
	var response volume.Response
	log.Debugf("capabilities: %#v", r)
	return response
}

func (d *SecretDriver) Get(r volume.Request) volume.Response {
	var response volume.Response
	log.Debugf("get: %#v", r)
	return response
}

func (d *SecretDriver) List(r volume.Request) volume.Response {
	var response volume.Response
	log.Debugf("list: %#v", r)
	return response
}

func (d *SecretDriver) Unmount(r volume.UnmountRequest) volume.Response {
	log.Debugf("unmount: %v", r)

	p := d.resolvePath(r.Name)
	if err := syscall.Unmount(p, 0); err != nil {
		log.Fatal(err)
	}

	return volume.Response{}
}
