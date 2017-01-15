# datafs
A de-duplicated filesystem that syncs to S3 and works on OSX, Windows and Linux

## Anatomy of a Go -> C -> Dokany

mounter := kbfs/libdokan.Mounter interface

```Go
type Mounter interface {
	Dir() string
	Mount(*dokan.Config, logger.Logger) error
	Unmount() error
}
```

options := kbfs/libdokan.StartOptions struct
```Go
type StartOptions struct {
	KbfsParams  libkbfs.InitParams
	RuntimeDir  string
	Label       string
	DokanConfig dokan.Config
}
```

ctx := kbfs/libkbfs.Context interface
```Go
type Context interface {
	GetRunMode() libkb.RunMode
	GetLogDir() string
	GetDataDir() string
	ConfigureSocketInfo() (err error)
	GetSocket(clearError bool) (net.Conn, rpc.Transporter, bool, error)
	NewRPCLogFactory() *libkb.RPCLogFactory
}
```


then to start:
```
err = libdokan.Start(mounter, options, ctx)
```
