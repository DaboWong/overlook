package ds

type IData interface {
	GetUnique() string
	Container() IDataContainer
}


type Identify struct {
	Unique string
	Con    IDataContainer
}

func (self *Identify) GetUnique() string {
	return self.Unique
}

func (self *Identify) Container() IDataContainer {
	return self.Con
}