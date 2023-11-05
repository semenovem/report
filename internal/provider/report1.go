package provider

func (p *Provider) Report1() (any, error) {
	ll := p.logger.Named("Report1")

	ll.Debug("got")

	return "", nil
}
