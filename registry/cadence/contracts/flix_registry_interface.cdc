access(all)
contract interface FLIXRegistryInterface {

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
    enum FLIXStatus: UInt8 {

        access(all)
        case active

        access(all)
        case deprecated
    }

    access(all)
    struct interface InteractionTemplate {
        access(all)
        fun getId(): String

        access(all)
        fun getVersion(): String

        access(all)
        fun getCadenceBodyHash(): String

        pub(set)
        var status: FLIXStatus

        access(all)
        fun getData(): AnyStruct
    }

    access(all)
    resource interface Queryable{
        access(all)
        fun getName(): String

        access(all)
        fun resolve(cadenceBodyHash: String): {FLIXRegistryInterface.InteractionTemplate}?

        access(all)
        fun lookup(idOrAlias: String): {FLIXRegistryInterface.InteractionTemplate}?

        access(all)
        fun getIds(): [String]

        access(all)
        fun getAllAlias(): {String: String}
    }

    access(all)
    resource interface Admin {

        // Add the flix to the registry and add or update the alias to point to this flix
        access(all) 
        fun publish(alias: String, flix: {FLIXRegistryInterface.InteractionTemplate})

        access(all)
        fun link(alias: String, id: String)

        access(all)
        fun unlink(alias: String)

        access(all)
        fun deprecate(idOrAlias: String)
    }

    access(all)
    resource interface Removable {
        access(all)
        fun remove(id: String): AnyStruct?
    }

    access(all)
    resource Registry: Queryable, Removable, Admin {
        access(all)
        fun getName(): String
    }

    access(all)
    fun createRegistry(name: String): @Registry

}