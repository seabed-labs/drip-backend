import NodeWallet from "@project-serum/anchor/dist/cjs/nodewallet";
import {AnchorProvider, BN} from "@project-serum/anchor";
import {
    AccountFetcher,
    buildWhirlpoolClient,
    ORCA_WHIRLPOOL_PROGRAM_ID,
    swapQuoteByInputToken,
    WhirlpoolContext
} from "@orca-so/whirlpools-sdk";
import {Percentage} from "@orca-so/common-sdk";
import { PublicKey, Keypair, Connection } from "@solana/web3.js";

// example usage
// npx ts-node ./pkg/service/orcawhirlpool/orcaWhirlpoolQuoteEstimate.ts GSFnjnJ7TdSsGWb6JgFhWakWrv8VGZUAghnY3EA8Xj46 7ihthG4cFydyDnuA3zmJrX13ePGpLcANf3tHLmKLPN7M 100000 https://api.devnet.solana.com
async function getQuote() {
    const args = process.argv.slice(2);
    if (args.length != 4) {
        console.log(JSON.stringify({
            error: `invalid number of arguments ${args.length}, expected 6`
        }));
        return;
    }
    const whirlpoolPubkey = new PublicKey(args[0]);
    const inputToken = new PublicKey(args[1]);
    const inputAmount =  new BN(args[2]);
    const connectionUrl = args[3];

    // Don't need to sign anything, so a random keypair is fine
    const wallet = new NodeWallet(Keypair.generate());
    const provider = new AnchorProvider(
        new Connection(connectionUrl, "confirmed"),
        wallet,
        AnchorProvider.defaultOptions()
    );
    const fetcher = new AccountFetcher(provider.connection);
    // @ts-ignore - orca uses an older anchor version, so the provider is incompatible
    const ctx = WhirlpoolContext.withProvider(provider, ORCA_WHIRLPOOL_PROGRAM_ID);
    const whirlpoolClient = buildWhirlpoolClient(ctx);

    const whirlpool = await whirlpoolClient.getPool(whirlpoolPubkey, true);

    const swapQuote =  await swapQuoteByInputToken(
        whirlpool,
        inputToken,
        inputAmount,
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

async function main() {
    try {
        await getQuote();
    } catch(e) {
        console.log(JSON.stringify({
            error: JSON.stringify(e)
        }));
    }
}

main();


