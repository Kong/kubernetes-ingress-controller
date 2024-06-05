package util

type ForeignRelations struct {
	Consumer, ConsumerGroup, Route, Service []string
}

type Rel struct {
	Consumer, ConsumerGroup, Route, Service string
}

func (relations *ForeignRelations) GetCombinations() []Rel {
	var cartesianProduct []Rel

	// gocritic I don't care that you think switch statements are the one true god of readability, the language offers
	// multiple options for a reason. go away, gocritic.
	if len(relations.Consumer) > 0 { //nolint:gocritic
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
	} else if len(relations.ConsumerGroup) > 0 {
		groups := relations.ConsumerGroup
		if len(relations.Route)+len(relations.Service) > 0 {
			for _, service := range relations.Service {
				for _, group := range groups {
					cartesianProduct = append(cartesianProduct, Rel{
						Service:       service,
						ConsumerGroup: group,
					})
				}
			}
			for _, route := range relations.Route {
				for _, group := range groups {
					cartesianProduct = append(cartesianProduct, Rel{
						Route:         route,
						ConsumerGroup: group,
					})
				}
			}
		} else {
			for _, group := range relations.ConsumerGroup {
				cartesianProduct = append(cartesianProduct, Rel{ConsumerGroup: group})
			}
		}
	} else {
		for _, service := range relations.Service {
			cartesianProduct = append(cartesianProduct, Rel{Service: service})
		}
		for _, route := range relations.Route {
			cartesianProduct = append(cartesianProduct, Rel{Route: route})
		}
	}

	return cartesianProduct
}
