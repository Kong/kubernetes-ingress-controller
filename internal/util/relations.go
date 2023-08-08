package util

type ForeignRelations struct {
	Consumer, ConsumerGroup, Route, Service []string
}

type Rel struct {
	Consumer, ConsumerGroup, Route, Service string
}

func (relations *ForeignRelations) GetCombinations() []Rel {
	var cartesianProduct []Rel

	if len(relations.Consumer) > 0 {
		consumers := relations.Consumer
		if len(relations.Route)+len(relations.Service) > 0 {
			for _, service := range relations.Service {
				for _, consumer := range consumers {
					cartesianProduct = append(cartesianProduct, Rel{
						Service:  service,
						Consumer: consumer,
					})
				}
			}
			for _, route := range relations.Route {
				for _, consumer := range consumers {
					cartesianProduct = append(cartesianProduct, Rel{
						Route:    route,
						Consumer: consumer,
					})
				}
			}
		} else {
			for _, consumer := range relations.Consumer {
				cartesianProduct = append(cartesianProduct, Rel{Consumer: consumer})
			}
		}
	} else {
		for _, consumerGroup := range relations.ConsumerGroup {
			cartesianProduct = append(cartesianProduct, Rel{ConsumerGroup: consumerGroup})
		}
		for _, service := range relations.Service {
			cartesianProduct = append(cartesianProduct, Rel{Service: service})
		}
		for _, route := range relations.Route {
			cartesianProduct = append(cartesianProduct, Rel{Route: route})
		}
	}

	return cartesianProduct
}
