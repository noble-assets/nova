const MerkleTreeHookAbi = [
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "bytes32",
        name: "messageId",
        type: "bytes32",
      },
      { indexed: false, internalType: "uint32", name: "index", type: "uint32" },
    ],
    name: "InsertedIntoTree",
    type: "event",
  },
] as const;

export default MerkleTreeHookAbi;
