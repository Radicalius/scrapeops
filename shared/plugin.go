package shared

type CronConfiguration struct {
	Schedule string
	Item     Item
}

type PluginConfiguration struct {
	Providers    []Provider
	CronConfigs  []CronConfiguration
	Schematizers []Schematizer[any]
}

func NewPluginConfiguration() *PluginConfiguration {
	return &PluginConfiguration{
		Providers:   make([]Provider, 0),
		CronConfigs: make([]CronConfiguration, 0),
	}
}

func (pc *PluginConfiguration) RegisterProvider(p Provider) {
	pc.Providers = append(pc.Providers, p)
}

func (pc *PluginConfiguration) RegisterCron(expr string, i Item) {
	pc.CronConfigs = append(pc.CronConfigs, CronConfiguration{
		Schedule: expr,
		Item:     i,
	})
}

func (pc *PluginConfiguration) RegisterSchematizer(s Schematizer[any]) {
	pc.Schematizers = append(pc.Schematizers, s)
}
