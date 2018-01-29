package ds

type IData interface {
	Unique() string
	Container() IDataContainer
}
