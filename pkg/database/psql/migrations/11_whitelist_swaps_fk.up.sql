ALTER TABLE "vault_whitelist" ADD FOREIGN KEY ("vault_pubkey") REFERENCES "vault" ("pubkey");

