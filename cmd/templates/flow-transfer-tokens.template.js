/**
    This binding file was auto generated based on FLIX template. 
    Changes to this file might get overwritten 
**/

import * as fcl from "@onflow/fcl"
import flixTemplate from "./flow-transfer-tokens.template.json"

const parameterNames = ["to", "amount"];

export async function transferTokens({to, amount}) {
  const transactionId = await fcl.mutate({
    template: flixTemplate,
    args: (arg, t) => [arg(to, t.Address), arg(amount, t.UFix64)]
  });

  return transactionId
}





