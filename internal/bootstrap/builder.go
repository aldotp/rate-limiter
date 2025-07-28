package bootstrap

func (b *Bootstrap) BuildDependencies() (*Bootstrap, error) {
	if err := b.setConfig(); err != nil {
		return nil, err
	}
	if err := b.setLogger(); err != nil {
		return nil, err
	}
	if err := b.setRedisClient(); err != nil {
		return nil, err
	}
	return b, nil
}
