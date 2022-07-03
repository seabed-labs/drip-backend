CREATE TABLE "token_account_balance" (
    "pubkey" varchar(255) PRIMARY KEY,
    "mint" varchar(255) NOT NULL,
    "owner" varchar(255) NOT NULL,
    "amount" numeric NOT NULL,
    "state" varchar(255) NOT NULL
);

ALTER TABLE "token_account_balance" ADD FOREIGN KEY ("mint") REFERENCES "token" ("pubkey");
