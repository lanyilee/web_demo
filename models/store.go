package models

//Store 数据存储
type Store struct {
	User          UserStore

	App           AppStore
	UserBind      UserBindStore
}

//StoreWatchEvent watch event
type StoreWatchEvent struct {
	Delete   bool
	Put      bool
	Data     interface{}
	PrevData interface{}
}

//StoreWatchFiled watch 字段
type StoreWatchFiled interface {
	Key() string
	Version() string
}
