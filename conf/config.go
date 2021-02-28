package conf

type AppConf struct{
	KafkaConf  `ini:"kafka"`
	EtcdConf `ini:"etcd"`

}

type EtcdConf struct{
	Address string `ini:"address"`
	Timeout int `ini:"timeout"`
	Key string `ini:"collect_log_key"`
}


type KafkaConf struct{
	Address string `ini:"address"`
	ChanMaxSize int `ini:"chan_max_size"`
}

// -----最初版本-----------
type TaillogConf struct{
	FileName string
}