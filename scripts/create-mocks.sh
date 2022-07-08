cd pkg/clients/solana
mockgen -source=client.go -destination=mock.go -package=solana
cd ../../repository
mockgen -source=repository.go -destination=mock.go -package=repository