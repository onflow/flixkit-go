import Test
import BlockchainHelpers
import "FLIXSchema_v1_1_0"
import "FLIXRegistry"

access(all) let REGISTRY_OWNER = Test.createAccount()
access(all) let TEMPLATE_ID = "aTestId"
access(all) let ALIAS = "anAlias"
access(all) let NEW_ALIAS = "aNewAlias"
access(all) let REGISTRY_NAME = "someName"
access(all) let CADENCE_BODY_HASH = "someHash"
access(all) let SCHEMA_VERSION = "1.1.0"

access(all) let CONTRACT_DEPLOYED_SNAPSHOT = "contract deployed"
access(all) let REGISTRY_CREATED_SNAPSHOT = "registry created"
access(all) let FLIX_PUBLISHED_SNAPSHOT = "flix published"

access(all)
fun setup() {

    let flixRegistryInterfaceErr = Test.deployContract(
        name: "FLIXRegistryInterface",
        path: "../contracts/flix_registry_interface.cdc",
        arguments: []
    )
    Test.expect(flixRegistryInterfaceErr, Test.beNil())

    let flixRegistryErr = Test.deployContract(
        name: "FLIXRegistry",
        path: "../contracts/flix_registry.cdc",
        arguments: []
    )
    Test.expect(flixRegistryErr, Test.beNil())

    let flix_schema_v1_1_0 = Test.deployContract(
        name: "FLIXSchema_v1_1_0",
        path: "../contracts/flix_schema_v1_1_0.cdc",
        arguments: []
    )
    Test.expect(flix_schema_v1_1_0, Test.beNil())
    Test.createSnapshot(name: CONTRACT_DEPLOYED_SNAPSHOT)

    let createTxResult = executeTransaction(
        "../transactions/create_registry.cdc",
        [REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(createTxResult, Test.beSucceeded())
    Test.createSnapshot(name: REGISTRY_CREATED_SNAPSHOT)

    let flixData = FLIXSchema_v1_1_0.Data(
        type: "transaction", 
        interface: "asadf23234...fas234234", 
        messages: [
            FLIXSchema_v1_1_0.Message(
                key: "title",
                i18n: [
                    FLIXSchema_v1_1_0.I18n(
                        tag: "en-US",
                        translation: "Transfer FLOW",
                    ),
                    FLIXSchema_v1_1_0.I18n(
                        tag: "fr-FR",
                        translation: "FLOW de transfert",
                    ),
                    FLIXSchema_v1_1_0.I18n(
                        tag: "zh-CN",
                        translation: "转移流程",
                    )
                ]
            ),
            FLIXSchema_v1_1_0.Message(
                key: "description",
                i18n: [
                    FLIXSchema_v1_1_0.I18n(
                        tag: "en-US",
                        translation: "Transfer {amount} FLOW to {to}",
                    ),
                    FLIXSchema_v1_1_0.I18n(
                        tag: "fr-FR",
                        translation: "Transférez {amount} FLOW à {to}",
                    ),
                    FLIXSchema_v1_1_0.I18n(
                        tag: "zh-CN",
                        translation: "将 {amount} FLOW 转移到 {to}",
                    )
                ]
            ),
            FLIXSchema_v1_1_0.Message(
                key: "signer",
                i18n: [
                    FLIXSchema_v1_1_0.I18n(
                        tag: "en-US",
                        translation: "Sign this message to transfer FLOW",
                    ),
                    FLIXSchema_v1_1_0.I18n(
                        tag: "fr-FR",
                        translation: "Signez ce message pour transférer FLOW.",
                    ),
                    FLIXSchema_v1_1_0.I18n(
                        tag: "zh-CN",
                        translation: "签署此消息以转移FLOW。",
                    )
                ]
            )
        ],
        cadence: FLIXSchema_v1_1_0.Cadence(
            body: "import \"FlowToken\"\ntransaction(amount: UFix64, to: Address) {\n    let vault: @FungibleToken.Vault\n    prepare(signer: AuthAccount) {\n        self.vault <- signer\n        .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n        .withdraw(amount: amount)\n\n        self.vault <- FungibleToken.getVault(signer)\n    }\n    execute {\n        getAccount(to)\n        .getCapability(/public/flowTokenReceiver)!\n        .borrow<&{FungibleToken.Receiver}>()!\n        .deposit(from: <-self.vault)\n    }\n}", 
            networkPins: [
                FLIXSchema_v1_1_0.NetworkPin(
                    network: "mainnet",
                    pinSelf: "186e262ce6fe06b5075ec6569a0e5482a79c471881182612d8e4a665c2977f3e"
                ),
                FLIXSchema_v1_1_0.NetworkPin(
                    network: "testnet",
                    pinSelf: "f93977d7a297f559e97259cb2a95fed0f87cfeec46c5257a26adc26a260d6c4c"
                )
            ]
        ),
        dependencies: [ 
            FLIXSchema_v1_1_0.Dependency(
                contracts: [
                    FLIXSchema_v1_1_0.Contract(
                        contract: "FlowToken",
                        networks: [
                            FLIXSchema_v1_1_0.Network(
                                network: "mainnet",
                                address: "0x1654653399040a61",
                                dependencyPinBlockHeight: 10123123123,
                                dependencyPin: FLIXSchema_v1_1_0.DependencyPin(
                                    pin: "c8cb7cc7a1c2a329de65d83455016bc3a9b53f9668c74ef555032804bac0b25b",
                                    pinSelf: "38d0cca4b74c4e88213df636b4cfc2eb6e86fd8b2b84579d3b9bffab3e0b1fcb",
                                    pinContractName: "FlowToken",
                                    pinContractAddress: "0x1654653399040a61",
                                    imports: [
                                        FLIXSchema_v1_1_0.Import(
                                            pin: "b8a3ed26c222ed67016a28021d8fee5603b948533cbc992b3c90f71a61b2b312",
                                            pinSelf: "7bc3056ba5d39d130f45411c2c05bb549db8ce727c11a1cb821136a621be27fb",
                                            pinContractName: "FungibleToken",
                                            pinContractAddress: "0xf233dcee88fe0abe",
                                            imports: []
                                        )
                                    ]
                                )
                            ), 
                            FLIXSchema_v1_1_0.Network(
                                network: "testnet",
                                address: "0x7e60df042a9c0868",
                                dependencyPinBlockHeight: 10123123123,
                                dependencyPin: FLIXSchema_v1_1_0.DependencyPin(
                                    pin: "c8cb7cc7a1c2a329de65d83455016bc3a9b53f9668c74ef555032804bac0b25b",
                                    pinSelf: "38d0cca4b74c4e88213df636b4cfc2eb6e86fd8b2b84579d3b9bffab3e0b1fcb",
                                    pinContractName: "FlowToken",
                                    pinContractAddress: "0x7e60df042a9c0868",
                                    imports: [
                                        FLIXSchema_v1_1_0.Import(
                                            pin: "b8a3ed26c222ed67016a28021d8fee5603b948533cbc992b3c90f71a61b2b312",
                                            pinSelf: "7bc3056ba5d39d130f45411c2c05bb549db8ce727c11a1cb821136a621be27fb",
                                            pinContractName: "FungibleToken",
                                            pinContractAddress: "0x9a0766d93b6608b7",
                                            imports: []
                                        )
                                    ]
                                )
                            )
                        ]
                    )
                ]
            )
        ],
        parameters: [
            FLIXSchema_v1_1_0.Parameter(
                label: "amount",
                index: 0,
                type: "UFix64",
                messages: [
                    FLIXSchema_v1_1_0.Message(
                        key: "title",
                        i18n: [
                            FLIXSchema_v1_1_0.I18n(
                                tag: "en-US",
                                translation: "Amount"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "fr-FR",
                                translation: "Montant"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "zh-CN",
                                translation: "数量"
                            )
                        ]
                    ),
                    FLIXSchema_v1_1_0.Message(
                        key: "description",
                        i18n: [
                            FLIXSchema_v1_1_0.I18n(
                                tag: "en-US",
                                translation: "Amount of FLOW token to transfer"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "fr-FR",
                                translation: "Quantité de token FLOW à transférer"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "zh-CN",
                                translation: "要转移的 FLOW 代币数量"
                            )
                        ])
                ],
                balance: [
                    FLIXSchema_v1_1_0.Balance(
                        network: "mainnet",
                        pin: "A.0xABC123DEF456.FlowToken"
                    ),
                    FLIXSchema_v1_1_0.Balance(
                        network: "testnet",
                        pin: "A.0xXYZ678DEF123.FlowToken"
                    )
                ]
            ),
            FLIXSchema_v1_1_0.Parameter(
                label: "to",
                index: 1,
                type: "Address",
                messages: [
                    FLIXSchema_v1_1_0.Message(
                        key: "title", 
                        i18n: [
                            FLIXSchema_v1_1_0.I18n(
                                tag: "en-US",
                                translation: "To"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "fr-FR",
                                translation: "Pour"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "zh-CN",
                                translation: "到"
                            )
                        ]
                    ),
                    FLIXSchema_v1_1_0.Message(
                        key: "description", 
                        i18n: [
                            FLIXSchema_v1_1_0.I18n(
                                tag: "en-US",
                                translation: "Recipient of the FLOW token transfer"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "fr-FR",
                                translation: "Le compte vers lequel transférer les jetons FLOW"
                            ),
                            FLIXSchema_v1_1_0.I18n(
                                tag: "zh-CN",
                                translation: "将 FLOW 代币转移到的帐户"
                            )
                        ]
                    )
                ],
                balance: []
            )
        ]
    )

    let flix = FLIXSchema_v1_1_0.FLIX(
        id: TEMPLATE_ID,
        data: flixData,
        cadenceBodyHash: CADENCE_BODY_HASH
    )

    let publishTxResult = executeTransaction(
        "../transactions/publish_flix_v1_1_0.cdc",
        [ALIAS, flix, REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(publishTxResult, Test.beSucceeded())
    Test.createSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)
}

access(all)
fun testShouldEmitContractInitializedEvent() {
    Test.loadSnapshot(name: REGISTRY_CREATED_SNAPSHOT)

    let typ = Type<FLIXRegistry.ContractInitialized>()
    Test.assertEqual(1, Test.eventsOfType(typ).length)
}

access(all)
fun testShouldEmitRegistryCreatedEvent() {
    Test.loadSnapshot(name: REGISTRY_CREATED_SNAPSHOT)

    let typ = Type<FLIXRegistry.RegistryCreated>()
    let events = Test.eventsOfType(typ)
    let event = events[0] as! FLIXRegistry.RegistryCreated
    Test.assertEqual(1, events.length)
    Test.assertEqual(REGISTRY_NAME, event.name)
}

access(all)
fun testShouldEmitEventAfterFlixPublished() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let typ = Type<FLIXRegistry.Published>()
    let events = Test.eventsOfType(typ)
    let event = events[0] as! FLIXRegistry.Published
    Test.assertEqual(1, events.length)
    Test.assertEqual(REGISTRY_NAME, event.registryName)
    Test.assertEqual(REGISTRY_OWNER.address, event.registryOwner)
    Test.assertEqual(ALIAS, event.alias)
    Test.assertEqual(TEMPLATE_ID, event.id)
    Test.assertEqual(CADENCE_BODY_HASH, event.cadenceBodyHash)
}

access(all)
fun testShouldContainFlixAfterPublished() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let scriptResult = executeScript(
        "../scripts/get_registry_size.cdc",
        [REGISTRY_OWNER.address, REGISTRY_NAME]
    )
    Test.expect(scriptResult, Test.beSucceeded())

    let registrySize = scriptResult.returnValue! as! Int
    Test.assertEqual(1, registrySize)
}

access(all)
fun testShouldLookupFlixAfterPublished() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let lookupScriptResult = executeScript(
        "../scripts/lookup.cdc",
        [REGISTRY_OWNER.address, ALIAS, REGISTRY_NAME]
    )

    Test.expect(lookupScriptResult, Test.beSucceeded())

    let flix = lookupScriptResult.returnValue! as! FLIXSchema_v1_1_0.FLIX
    Test.assertEqual(TEMPLATE_ID, flix.id)
    Test.assertEqual("transaction", flix.getData().type)
    Test.assertEqual(CADENCE_BODY_HASH, flix.cadenceBodyHash)
    Test.assertEqual("active", flix.status)
    Test.assertEqual(SCHEMA_VERSION, flix.getVersion())
}

access(all)
fun testShouldResolveFlixAfterPublished() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let resolveScriptResult = executeScript(
        "../scripts/resolve.cdc",
        [REGISTRY_OWNER.address, CADENCE_BODY_HASH, REGISTRY_NAME]
    )

    Test.expect(resolveScriptResult, Test.beSucceeded())

    let resolvedFlix = resolveScriptResult.returnValue! as! FLIXSchema_v1_1_0.FLIX
    Test.assertEqual(TEMPLATE_ID, resolvedFlix.id)
    Test.assertEqual("transaction", resolvedFlix.getData().type)
    Test.assertEqual(CADENCE_BODY_HASH, resolvedFlix.cadenceBodyHash)
    Test.assertEqual("active", resolvedFlix.status)
}
