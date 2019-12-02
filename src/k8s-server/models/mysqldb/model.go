package mysqldb

import (
	"github.com/go-gorp/gorp"
	"github.com/astaxie/beego"
	"k8s-server/conf"
)

var (
	cfg       = beego.AppConfig
	model     Model
	modelLock = new(sync.Mutex)
)

// Model is an public struct but should be initialized only once.
type Model struct {
	db     *gorp.DbMap
	inited bool
	config conf.MySQLConfig
}

// GetModel returns initiated model.
func GetModel() (*Model, error) {
	modelLock.Lock()
	defer modelLock.Unlock()
	var err error
	if !model.inited {
		if err = model.initDB(); err == nil {
			model.inited = true
		}
	}
	return &model, err
}

func (m *Model) initDB() error {
	m.config = conf.HPCMySQLConfig()
	dbMap, err := m.newDBMap()
	if err != nil {
		return err
	}
	m.db = dbMap
	return nil
}