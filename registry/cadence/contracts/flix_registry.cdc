import "FLIXRegistryInterface"

access(all) 
contract FLIXRegistry: FLIXRegistryInterface {

    access(all)
    event ContractInitialized()

    access(all)
    event RegistryCreated(name: String)

    access(all)
    event Published(registryOwner: Address, registryName: String, registryUuid: UInt64, alias: String, id: String, cadenceBodyHash: String)

    access(all)
    event Removed(registryOwner: Address, registryName: String, registryUuid: UInt64, id: String)

    access(all)
    event Deprecated(registryOwner: Address, registryName: String, registryUuid: UInt64, id: String)

    access(all)
    event AliasLinked(registryOwner: Address, registryName: String, registryUuid: UInt64, alias: String, id: String)

    access(all)
    event AliasUnlinked(registryOwner: Address, registryName: String, registryUuid: UInt64, alias: String)

    access(all)
    resource Registry: FLIXRegistryInterface.Queryable, FLIXRegistryInterface.Removable, FLIXRegistryInterface.Admin {

        access(account) var flixes: {String: {FLIXRegistryInterface.InteractionTemplate}}
        access(account) var aliases: {String: String}
        access(account) var cadenceBodyHashes: {String: {FLIXRegistryInterface.InteractionTemplate}}
        access(account) let name: String

        access(all)
        fun getName(): String {
            return self.name
        }

        access(all) 
        fun publish(alias: String, flix: {FLIXRegistryInterface.InteractionTemplate}) {
            self.flixes[flix.getId()] = flix
            self.aliases[alias] = flix.getId()
            self.cadenceBodyHashes[flix.getCadenceBodyHash()] = flix
            emit Published(registryOwner: self.owner!.address, registryName: self.name, registryUuid: self.uuid, alias: alias, id: flix.getId(), cadenceBodyHash: flix.getCadenceBodyHash())
        }

        access(all)
        fun link(alias: String, id: String) {
            self.aliases[alias] = id
            emit AliasLinked(registryOwner: self.owner!.address, registryName: self.name, registryUuid: self.uuid, alias: alias, id: id)
        }

        access(all)
        fun unlink(alias: String) {
            self.aliases.remove(key: alias)
            emit AliasUnlinked(registryOwner: self.owner!.address, registryName: self.name, registryUuid: self.uuid, alias: alias)
        }

        access(all)
        fun lookup(idOrAlias: String): {FLIXRegistryInterface.InteractionTemplate}? {
            if self.aliases.containsKey(idOrAlias) {
                return self.flixes[self.aliases[idOrAlias]!]
            }
            return self.flixes[idOrAlias]
        }

        access(all)
        fun resolve(cadenceBodyHash: String): {FLIXRegistryInterface.InteractionTemplate}? {
            return self.cadenceBodyHashes[cadenceBodyHash]
        }

        access(all)
        fun deprecate(idOrAlias: String) {
            var flix = self.lookup(idOrAlias: idOrAlias) ?? panic("FLIX does not exist with the given id or alias: ".concat(idOrAlias))
            flix.status = "deprecated"
            self.flixes[flix.getId()] = flix
            emit Deprecated(registryOwner: self.owner!.address, registryName: self.name, registryUuid: self.uuid, id: flix.getId())
        }

        access(all)
        fun remove(id: String): {FLIXRegistryInterface.InteractionTemplate}? {
            let removed = self.flixes.remove(key: id)
            if(removed != nil) { emit Removed(registryOwner: self.owner!.address, registryName: self.name, registryUuid: self.uuid, id: id) }
            return removed
        }

        access(all)
        fun getIds(): [String] {
            return self.flixes.keys
        }

        access(all)
        fun getAllAlias(): {String: String} {
            return self.aliases
        }

        init(name: String) {
            self.name = name
            self.flixes = {}
            self.aliases = {}
            self.cadenceBodyHashes = {}
            emit RegistryCreated(name: name)
        }
    }

    access(all)
    fun createRegistry(name: String): @Registry {
        return <- create Registry(name: name)
    }

    access(all)
    fun PublicPath(name: String): PublicPath {
        return PublicPath(identifier: "flix_".concat(name))!
    }

        access(all)
    fun StoragePath(name: String): StoragePath {
        return StoragePath(identifier: "flix_".concat(name))!
    }

    init() {
        emit ContractInitialized()
    }

}
