import { onchainTable } from "ponder";

export const message = onchainTable("message", (t) => ({
  index: t.integer().primaryKey(),
  id: t.hex().notNull()
}));
