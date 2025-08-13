import { createConfig } from "ponder";

import MerkleTreeHookAbi from "./abis/MerkleTreeHookAbi";

export default createConfig({
  chains: {
    anvil: {
      id: 31337,
      rpc: "http://localhost:8545",
    },
  },
  contracts: {
    MerkleTreeHook: {
      chain: "anvil",
      abi: MerkleTreeHookAbi,
      address: "0x0000000000000000000000000000000000000000",
      startBlock: 0,
    },
  },
});
