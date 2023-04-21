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

cd analytics || exit
mockgen -source=repository.go -destination=mock.go -package=repository

cd ../queue || exit
mockgen -source=repository.go -destination=mock.go -package=repository

cd ../transactioncheckpoint || exit
mockgen -source=repository.go -destination=mock.go -package=repository

cd ../

cd ../base || exit
mockgen -source=base.go -destination=mock.go -package=base

cd ../config || exit
mockgen -source=config.go -destination=mock.go -package=config

cd ../../job/token || exit
mockgen -source=token.go -destination=mock.go -package=token

cd ../tokenaccount || exit
mockgen -source=tokenaccount.go -destination=mock.go -package=tokenaccount