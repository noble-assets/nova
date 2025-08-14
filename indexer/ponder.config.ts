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
      address: "0xb19b36b1456E65E3A6D514D3F715f204BD59f431",
      startBlock: 0,
    },
  },
});
