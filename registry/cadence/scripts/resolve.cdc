import "FLIXRegistry"
import "FLIXRegistryInterface"

pub fun main(accountAddress: Address, cadenceBodyHash: String, registryName: String): AnyStruct{FLIXRegistryInterface.InteractionTemplate}? {
    let account = getAccount(accountAddress)
    let registry = account.getCapability(FLIXRegistry.PublicPath(name: registryName))
                            .borrow<&FLIXRegistry.Registry{FLIXRegistryInterface.Queryable}>()
                            ?? panic("Could not borrow a reference to the Registry")

    return registry.resolve(cadenceBodyHash: cadenceBodyHash)
}