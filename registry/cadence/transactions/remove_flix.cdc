import "FLIXRegistry"
import "FLIXRegistryInterface"

transaction(id: String, registryName: String) {

    prepare(signer: AuthAccount) {
        // Borrow a reference to the Registry with removal capability
        let registry = signer.borrow<&FLIXRegistry.Registry{FLIXRegistryInterface.Removable}>(from: FLIXRegistry.StoragePath(name: registryName))
                            ?? panic("Could not borrow a reference to the Registry")

        // Remove the FLIX item using the provided id
        registry.remove(id: id)
    }
}
