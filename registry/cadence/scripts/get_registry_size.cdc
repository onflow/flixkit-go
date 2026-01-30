import "FLIXRegistry"
import "FLIXRegistryInterface"

pub fun main(address: Address, registryName: String): Int {
    let account = getAccount(address)
    let registry = account.getCapability(FLIXRegistry.PublicPath(name: registryName))
                            .borrow<&FLIXRegistry.Registry{FLIXRegistryInterface.Queryable}>()
                            ?? panic("Could not borrow a reference to the Registry")

    assert(registryName == registry.getName())

    return registry.getIds().length
}