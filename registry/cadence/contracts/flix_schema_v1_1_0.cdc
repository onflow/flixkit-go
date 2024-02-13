import "FLIXRegistryInterface"
import "FLIXRegistry"

access(all) 
contract FLIXSchema_v1_1_0 {

    access(all) 
    struct FLIX: FLIXRegistryInterface.InteractionTemplate {
        access(all)
        let id: String

        access(all)
        let data: Data

        access(all)
        let cadenceBodyHash: String

        access(all)
        fun getId(): String {
            return self.id
        }

        access(all)
        fun getData(): Data {
            return self.data
        }

        access(all)
        fun getVersion(): String {
            return "1.1.0"
        }

        access(all)
        fun getCadenceBodyHash(): String {
            return self.cadenceBodyHash
        }

        access(all)
        fun getStatus(): FLIXRegistryInterface.FLIXStatus {
            return self.status
        }

        pub(set)
        var status: FLIXRegistryInterface.FLIXStatus

        init(id: String, data: Data, cadenceBodyHash: String) {
            self.id = id
            self.data = data
            self.cadenceBodyHash = cadenceBodyHash
            self.status = FLIXRegistry.FLIXStatus.active
        }
    }

    access(all) struct Data {
        access(all) var type: String
        access(all) var interface: String
        access(all) var messages: [Message]
        access(all) var cadence: Cadence
        access(all) var dependencies: [Dependency]
        access(all) var parameters: [Parameter]

        init(type: String, interface: String, messages: [Message], cadence: Cadence, dependencies: [Dependency], parameters: [Parameter]) {
            self.type = type
            self.interface = interface
            self.messages = messages
            self.cadence = cadence
            self.dependencies = dependencies
            self.parameters = parameters
        }
    }

    access(all) struct Message {
        access(all) var key: String
        access(all) var i18n: [I18n]

        init(key: String, i18n: [I18n]) {
            self.key = key
            self.i18n = i18n
        }
    }

    access(all) struct I18n {
        access(all) var tag: String
        access(all) var translation: String

        init(tag: String, translation: String) {
            self.tag = tag
            self.translation = translation
        }
    }

    access(all) struct Cadence {
        access(all) var body: String
        access(all) var networkPins: [NetworkPin]

        init(body: String, networkPins: [NetworkPin]) {
            self.body = body
            self.networkPins = networkPins
        }
    }

     access(all) struct NetworkPin {
        access(all) var network: String
        access(all) var pinSelf: String

        init(network: String, pinSelf: String) {
            self.network = network
            self.pinSelf = pinSelf
        }
    }

    access(all) struct Dependency {
        access(all) var contracts: [Contract]

        init(contracts: [Contract]) {
            self.contracts = contracts
        }
    }

    access(all) struct Contract {
        access(all) var contract: String
        access(all) var networks: [Network]

        init(contract: String, networks: [Network]) {
            self.contract = contract
            self.networks = networks
        }
    }

    access(all) struct Network {
        access(all) var network: String
        access(all) var address: String
        access(all) var dependencyPinBlockHeight: UInt64
        access(all) var dependencyPin: DependencyPin

        init(network: String, address: String, dependencyPinBlockHeight: UInt64, dependencyPin: DependencyPin) {
            self.network = network
            self.address = address
            self.dependencyPinBlockHeight = dependencyPinBlockHeight
            self.dependencyPin = dependencyPin
        }
    }

    access(all) struct DependencyPin {
        access(all) var pin: String
        access(all) var pinSelf: String
        access(all) var pinContractName: String
        access(all) var pinContractAddress: String
        access(all) var imports: [Import]

        init(pin: String, pinSelf: String, pinContractName: String, pinContractAddress: String, imports: [Import]) {
            self.pin = pin
            self.pinSelf = pinSelf
            self.pinContractName = pinContractName
            self.pinContractAddress = pinContractAddress
            self.imports = imports
        }
    }

    access(all) struct Import {
        access(all) var pin: String
        access(all) var pinSelf: String
        access(all) var pinContractName: String
        access(all) var pinContractAddress: String
        access(all) var imports: [Import]

        init(pin: String, pinSelf: String, pinContractName: String, pinContractAddress: String, imports: [Import]) {
            self.pin = pin
            self.pinSelf = pinSelf
            self.pinContractName = pinContractName
            self.pinContractAddress = pinContractAddress
            self.imports = imports
        }
    }

    access(all) struct Parameter {
        access(all) var label: String
        access(all) var index: Int
        access(all) var type: String
        access(all) var messages: [Message]
        access(all) var balance: [Balance]

        init(label: String, index: Int, type: String, messages: [Message], balance: [Balance]) {
            self.label = label
            self.index = index
            self.type = type
            self.messages = messages
            self.balance = balance
        }
    }

    access(all) struct Balance {
        access(all) var network: String
        access(all) var pin: String

        init(network: String, pin: String) {
            self.network = network
            self.pin = pin
        }
    }

}