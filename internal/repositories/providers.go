package repositories

type Providers interface {
	Convert(interface{}) (interface{}, error)

	PreProcess(interface{}, interface{}, interface{}, interface{}) (interface{}, error)

	InstallComponent(interface{}) error
	UpdateComponent(interface{}) error
	UninstallComponent(interface{}) error

	GetAllName() ([]string, error)
	Add(interface{}) error
	Remove(interface{}) error
	Update(interface{}) error
	GetDetail(string) (interface{}, error)
	GetDetailFromComponent(interface{}) (interface{}, error)
	GetFromModuleReleaseID(uint) ([]interface{}, error)
}
