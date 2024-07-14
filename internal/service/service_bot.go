package service

// Run to run discord bot.
func (s *service) Run() error {
	return s.discord.Run()
}

// Stop to stop discord bot.
func (s *service) Stop() error {
	return s.discord.Close()
}

// InitPlayer to init player.
func (s *service) InitPlayer(guildID string) {
	s.discord.InitPlayer(guildID)
}
