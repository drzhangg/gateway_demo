package load_balance

//配置主题
type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

type LoadBalanceZkConf struct {
	observers    []Observer
	path         string
	zkHosts      []string
	confIpWeight map[string]string
	activeList   []string
	format       string
}

type Observer interface {
	Update()
}

type LoadBalanceObserver struct {
	ModuleConf *LoadBalanceZkConf
}

//func (l *LoadBalanceObserver) Update() {
//	fmt.Println("Update get conf:", l.ModuleConf.GetConf())
//}

func NewLoadBalanceObserver(conf *LoadBalanceZkConf) *LoadBalanceObserver {
	return &LoadBalanceObserver{
		ModuleConf: conf,
	}
}
