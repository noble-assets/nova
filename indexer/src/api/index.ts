import { IMT, IMTNode } from "@zk-kit/imt";
import { Hono } from "hono";
import { client, graphql } from "ponder";
import { encodePacked, keccak256 } from "viem";

import { db } from "ponder:api";
import schema from "ponder:schema";

const app = new Hono();

app.use("/sql/*", client({ db, schema }));

app.use("/", graphql({ db, schema }));
app.use("/graphql", graphql({ db, schema }));

app.use("/prove/:messageId", async (ctx) => {
    const messageId = ctx.req.param("messageId") as `0x${string}`;
  
    const hash = (values: IMTNode[]) => {
      return keccak256(encodePacked(
        ["bytes32", "bytes32"],
        [values[0] as `0x${string}`, values[1] as `0x${string}`]
      ));
    };
    const tree = new IMT(hash, 32, "0x0000000000000000000000000000000000000000000000000000000000000000");
  
    const messages = await db.query.message.findMany();
    for (const message of messages) {
      tree.insert(message.id);
    }
  
    const index = messages.find((m) => m.id === messageId)?.index;
    if (!index) {
      return ctx.json({ error: "unknown message id" }, 404);
    }
  
    return ctx.json(tree.createProof(index));
  });

export default app;
