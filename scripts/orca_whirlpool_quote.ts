import NodeWallet from "@project-serum/anchor/dist/cjs/nodewallet";
import {AnchorProvider, BN} from "@project-serum/anchor";
import {
    AccountFetcher,
    buildWhirlpoolClient,
    ORCA_WHIRLPOOL_PROGRAM_ID,
    PDAUtil, swapQuoteByInputToken,
    WhirlpoolContext
} from "@orca-so/whirlpools-sdk";
import {Percentage} from "@orca-so/common-sdk";
import { PublicKey, Keypair, Connection } from "@solana/web3.js";

async function getQuote() {
    const args = process.argv.slice(2);
    const config = new PublicKey(args[0]);
    const tokenAMint = new PublicKey(args[1]);
    const tokenBMint = new PublicKey(args[2]);
    const inputToken = new PublicKey(args[3]);
    let endpoint = args[4];
    if (!endpoint) {
        endpoint = "https://api.devnet.solana.com"
    }

    const wallet = new NodeWallet(Keypair.generate());
    const provider = new AnchorProvider(
        new Connection(endpoint, "confirmed"),
        wallet,
        AnchorProvider.defaultOptions()
    );

    const fetcher = new AccountFetcher(provider.connection);
    // @ts-ignore (incompatible provider types because orca uses an older version of anchor
    const ctx = WhirlpoolContext.withProvider(provider, ORCA_WHIRLPOOL_PROGRAM_ID);

    const whirlpoolPda = PDAUtil.getWhirlpool(
        ORCA_WHIRLPOOL_PROGRAM_ID,
        config,
        tokenAMint,
        tokenBMint,
        64,
    );

    const whirlpoolClient = buildWhirlpoolClient(ctx);
    const whirlpool = await whirlpoolClient.getPool(whirlpoolPda.publicKey, true);

    const swapQuote =  await swapQuoteByInputToken(
        whirlpool,
        inputToken,
        new BN(100),
        Percentage.fromFraction(0, 100),
        ORCA_WHIRLPOOL_PROGRAM_ID,
        fetcher,
        true,
    );
    const swapQuoteString =  {
        estimatedAmountIn: swapQuote.estimatedAmountIn.toString(),
        estimatedAmountOut: swapQuote.estimatedAmountOut.toString(),
        estimatedEndTickIndex: swapQuote.estimatedEndTickIndex,
        estimatedEndSqrtPrice: swapQuote.estimatedEndSqrtPrice.toString(),
        estimatedFeeAmount: swapQuote.estimatedFeeAmount.toString(),
        amount: swapQuote.amount.toString(),
        amountSpecifiedIsInput: swapQuote.amountSpecifiedIsInput,
        aToB: swapQuote.aToB,
        otherAmountThreshold: swapQuote.otherAmountThreshold.toString(),
        sqrtPriceLimit: swapQuote.sqrtPriceLimit.toString(),
        tickArray0: swapQuote.tickArray0.toString(),
        tickArray1: swapQuote.tickArray1.toString(),
        tickArray2: swapQuote.tickArray2.toString(),
    };
    console.log(JSON.stringify(swapQuoteString));
}

try {
    getQuote();
} catch (e) {
    console.log(JSON.stringify({
        error: e.toString()
    }))
}

