import "FLIXRegistryInterface"
import "FLIXRegistry"

access(all) 
contract FLIXSchema_draft {

    access(all)
    struct FLIX: FLIXRegistryInterface.InteractionTemplate {
        access(all)
        let id: String

        access(all)
        let data: AnyStruct

        access(all)
        let cadenceBodyHash: String

        access(all)
        fun getId(): String {
            return self.id
        }

        access(all)
        fun getData(): AnyStruct {
            return self.data
        }

        access(all)
        fun getVersion(): String {
            return "draft"
        }

        access(all)
        fun getCadenceBodyHash(): String {
            return self.cadenceBodyHash
        }

        pub(set)
        var status: FLIXRegistryInterface.FLIXStatus

        init(id: String, data: AnyStruct, cadenceBodyHash: String) {
            self.id = id
            self.data = data
            self.cadenceBodyHash = cadenceBodyHash
            self.status = FLIXRegistry.FLIXStatus.active
        }
    }

}