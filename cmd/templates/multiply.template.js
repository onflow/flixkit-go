/**
    This binding file was auto generated based on FLIX template. 
    Changes to this file might get overwritten 
**/

import * as fcl from "@onflow/fcl"
import flixTemplate from "./multiply.template.json"

const parameterNames = ["x", "y"];

export async function multiplyTwoIntegers({x, y}) {
  const info = await fcl.query({
    template: flixTemplate,
    args: (arg, t) => [arg(x, t.Int), arg(y, t.Int)]
  });

  return info
}





