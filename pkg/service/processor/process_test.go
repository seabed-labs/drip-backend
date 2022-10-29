package processor

//func TestImpl_UpsertTokenSwapByAddress(t *testing.T) {
//	privKey := "[95,189,40,215,74,154,138,123,245,115,184,90,2,187,104,25,241,164,79,247,14,69,207,235,40,245,13,157,149,20,13,227,252,155,201,43,89,96,76,119,162,241,148,53,80,193,126,159,80,213,140,166,144,139,205,143,160,238,11,34,192,249,59,31]"
//	appConfig := config.AppConfig{
//		Environment: config.DevnetEnv,
//		Wallet:      privKey,
//	}
//	solanaClient, err := solana.NewSolanaClient(&appConfig)
//	assert.NoError(t, err)
//	psqlConfig, err := config.NewPSQLConfig()
//	assert.NoError(t, err)
//	gormDB, err := psql.NewGORMDatabase(psqlConfig)
//	assert.NoError(t, err)
//	repo := query.Use(gormDB)
//	processor := impl{
//		repo:   repository.NewRepository(solanaClient, repo),
//		solanaClient: solanaClient,
//	}
//	err = processor.UpsertTokenSwapByAddress(context.Background(), "8Vx5D6dKfv1vyLzL6hdxB1KNnRr5yiHSaw8bDh5LVNAE")
//	assert.NoError(t, err)
//}
