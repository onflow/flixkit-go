import "FLIXRegistry"
import "FLIXSchema_v1_1_0"

transaction(alias: String, flix: FLIXSchema_v1_1_0.FLIX, registryName: String) {

    prepare(signer: AuthAccount) {
        let registry = signer.borrow<&FLIXRegistry.Registry>(from: FLIXRegistry.StoragePath(name: registryName))
                            ?? panic("Could not borrow a reference to the Registry")
        registry.publish(alias: alias, flix: flix)
    }
}

