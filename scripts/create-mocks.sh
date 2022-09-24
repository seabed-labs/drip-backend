cd pkg/service/clients/solana || exit
mockgen -source=client.go -destination=mock.go -package=solana
cd ../../repository || exit
mockgen -source=repository.go -destination=mock.go -package=repository