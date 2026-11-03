package chunk

type Config struct {
	maxWords     int
	overlapWords int
	maxChunks    int
}

func DefaultConfig() Config {
	return Config{
		maxWords:     350,
		overlapWords: 50,
	}
}

func (c Config) MaxWords(n int) Config {
	c.maxWords = n
	return c
}

func (c Config) OverlapWords(n int) Config {
	c.overlapWords = n
	return c
}

func (c Config) MaxChunks(n int) Config {
	c.maxChunks = n
	return c
}
