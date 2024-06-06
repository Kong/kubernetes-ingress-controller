package util

type ForeignRelations struct {
	Consumer, ConsumerGroup, Route, Service []string
}

type Rel struct {
	Consumer, ConsumerGroup, Route, Service string
}

func (relations *ForeignRelations) GetCombinations() []Rel {
	var (
		lConsumer      = len(relations.Consumer)
		lConsumerGroup = len(relations.ConsumerGroup)
		lRoutes        = len(relations.Route)
		lServices      = len(relations.Service)
		l              = lRoutes + lServices
	)

	var cartesianProduct []Rel

	// gocritic I don't care that you think switch statements are the one true god of readability, the language offers
	// multiple options for a reason. go away, gocritic.
	if lConsumer > 0 { //nolint:gocritic
		if l > 0 {
			cartesianProduct = make([]Rel, 0, l*lConsumer)
			for _, consumer := range relations.Consumer {
				for _, service := range relations.Service {
					cartesianProduct = append(cartesianProduct, Rel{
						Service:  service,
						Consumer: consumer,
					})
				}
				for _, route := range relations.Route {
					cartesianProduct = append(cartesianProduct, Rel{
						Route:    route,
						Consumer: consumer,
					})
				}
			}

		} else {
			cartesianProduct = make([]Rel, 0, len(relations.Consumer))
			for _, consumer := range relations.Consumer {
				cartesianProduct = append(cartesianProduct, Rel{Consumer: consumer})
			}
		}
	} else if lConsumerGroup > 0 {
		if l > 0 {
			cartesianProduct = make([]Rel, 0, l*lConsumerGroup)
			for _, group := range relations.ConsumerGroup {
				for _, service := range relations.Service {
					cartesianProduct = append(cartesianProduct, Rel{
						Service:       service,
						ConsumerGroup: group,
					})
				}
				for _, route := range relations.Route {
					cartesianProduct = append(cartesianProduct, Rel{
						Route:         route,
						ConsumerGroup: group,
					})
				}
			}
		} else {
			cartesianProduct = make([]Rel, 0, lConsumerGroup)
			for _, group := range relations.ConsumerGroup {
				cartesianProduct = append(cartesianProduct, Rel{ConsumerGroup: group})
			}
		}
	} else if l > 0 {
		cartesianProduct = make([]Rel, 0, l)
		for _, service := range relations.Service {
			cartesianProduct = append(cartesianProduct, Rel{Service: service})
		}
		for _, route := range relations.Route {
			cartesianProduct = append(cartesianProduct, Rel{Route: route})
		}
	}

	return cartesianProduct
}
