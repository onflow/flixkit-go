/**
    This binding file was auto generated based on FLIX template. 
    Changes to this file might get overwritten 
**/

import * as fcl from "@onflow/fcl"
import flixTemplate from "./flow-transfer-tokens.template.json"

const parameterNames = ["amount", "to"];

export async function transferTokens({amount, to}) {
  const transactionId = await fcl.mutate({
    template: flixTemplate,
    args: (arg, t) => [arg(amount, t.UFix64), arg(to, t.Address)]
  });

  return transactionId
}





