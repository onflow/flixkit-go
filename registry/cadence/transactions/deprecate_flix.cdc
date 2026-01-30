import "FLIXRegistry"

transaction(idOrAlias: String, registryName: String) {

    prepare(signer: AuthAccount) {
        let registry = signer.borrow<&FLIXRegistry.Registry>(from: FLIXRegistry.StoragePath(name: registryName))
                            ?? panic("Could not borrow a reference to the Registry")

        registry.deprecate(idOrAlias: idOrAlias)
    }
}
