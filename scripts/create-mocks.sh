#!/bin/bash

cd pkg/service/clients/solana || exit
mockgen -source=client.go -destination=mock.go -package=solana

cd ../tokenregistry || exit
mockgen -source=client.go -destination=mock.go -package=tokenregistry

cd ../coingecko || exit
mockgen -source=client.go -destination=mock.go -package=coingecko

cd ../orcawhirlpool || exit
mockgen -source=client.go -destination=mock.go -package=orcawhirlpool

cd ../../repository || exit
mockgen -source=repository.go -destination=mock.go -package=repository

cd ../base || exit
mockgen -source=base.go -destination=mock.go -package=base

cd ../config || exit
mockgen -source=config.go -destination=mock.go -package=config
