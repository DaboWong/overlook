package ds

import "log"

type IDataContainer interface {
	GetData(unique string) (IData, bool)
	AddData(v ...IData)
}

type DataContainer struct {
	data map[string]IData
}

func NewDataContainer() IDataContainer {
	return &DataContainer{
		make(map[string]IData, 0),
	}
}

func (self *DataContainer) GetData(unique string) (IData, bool) {
	v, ok := self.data[unique]
	return v, ok
}

func (self *DataContainer) addData(data IData) {
	if _, ok := self.data[data.Unique()]; ok {
		log.Println("duplicated data add, name: ", data.Unique())
	} else {
		self.data[data.Unique()] = data
	}
}

func (self *DataContainer) AddData(v ...IData) {
	for i := 0; i < len(v); i++ {
		self.addData(v[i])
	}
}
