import { ponder } from "ponder:registry";
import { message } from "ponder:schema";

ponder.on("MerkleTreeHook:InsertedIntoTree", async ({ event, context }) => {
  await context.db.insert(message).values({
    index: event.args.index,
    id: event.args.messageId,
  });
});
