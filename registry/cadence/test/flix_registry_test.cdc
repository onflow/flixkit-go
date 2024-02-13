import Test
import BlockchainHelpers
import "FLIXSchema_draft"
import "FLIXRegistry"

access(all) let REGISTRY_OWNER = Test.createAccount()
access(all) let TEMPLATE_ID = "aTestId"
access(all) let ALIAS = "anAlias"
access(all) let NEW_ALIAS = "aNewAlias"
access(all) let REGISTRY_NAME = "someName"
access(all) let FLIX_DATA: AnyStruct = "same data"
access(all) let CADENCE_BODY_HASH = "someHash"

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

    let flixSchema_draftErr = Test.deployContract(
        name: "FLIXSchema_draft",
        path: "../contracts/flix_schema_draft.cdc",
        arguments: []
    )
    Test.expect(flixSchema_draftErr, Test.beNil())
    Test.createSnapshot(name: CONTRACT_DEPLOYED_SNAPSHOT)

    let createTxResult = executeTransaction(
        "../transactions/create_registry.cdc",
        [REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(createTxResult, Test.beSucceeded())
    Test.createSnapshot(name: REGISTRY_CREATED_SNAPSHOT)

    let publishTxResult = executeTransaction(
        "../transactions/publish_flix.cdc",
        [ALIAS, TEMPLATE_ID, FLIX_DATA, CADENCE_BODY_HASH, REGISTRY_NAME],
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
fun testShouldCreateRegistry() {
    Test.loadSnapshot(name: REGISTRY_CREATED_SNAPSHOT)

    let scriptResult = executeScript(
        "../scripts/get_registry_size.cdc",
        [REGISTRY_OWNER.address, REGISTRY_NAME]
    )
    Test.expect(scriptResult, Test.beSucceeded())

    let registrySize = scriptResult.returnValue! as! Int
    Test.assertEqual(0, registrySize)
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

    let flix = lookupScriptResult.returnValue! as! FLIXSchema_draft.FLIX
    Test.assertEqual(TEMPLATE_ID, flix.id)
    Test.assertEqual(FLIX_DATA, flix.getData())
    Test.assertEqual(CADENCE_BODY_HASH, flix.cadenceBodyHash)
    Test.assertEqual(FLIXRegistry.FLIXStatus.active, flix.status)
    Test.assertEqual("draft", flix.getVersion())
}

access(all)
fun testShouldResolveFlixAfterPublished() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let resolveScriptResult = executeScript(
        "../scripts/resolve.cdc",
        [REGISTRY_OWNER.address, CADENCE_BODY_HASH, REGISTRY_NAME]
    )

    Test.expect(resolveScriptResult, Test.beSucceeded())

    let resolvedFlix = resolveScriptResult.returnValue! as! FLIXSchema_draft.FLIX
    Test.assertEqual(TEMPLATE_ID, resolvedFlix.id)
    Test.assertEqual(FLIX_DATA, resolvedFlix.getData())
    Test.assertEqual(CADENCE_BODY_HASH, resolvedFlix.cadenceBodyHash)
    Test.assertEqual(FLIXRegistry.FLIXStatus.active, resolvedFlix.status)
}

access(all)
fun testShouldLinkAlias() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let txResult = executeTransaction(
        "../transactions/link_alias.cdc",
        [NEW_ALIAS, TEMPLATE_ID, REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(txResult, Test.beSucceeded())

    let typ = Type<FLIXRegistry.AliasLinked>()
    let events = Test.eventsOfType(typ)
    let event = events[0] as! FLIXRegistry.AliasLinked
    Test.assertEqual(1, events.length)
    Test.assertEqual(REGISTRY_NAME, event.registryName)
    Test.assertEqual(REGISTRY_OWNER.address, event.registryOwner)
    Test.assertEqual(NEW_ALIAS, event.alias)
    Test.assertEqual(TEMPLATE_ID, event.id)

    let scriptResult = executeScript(
        "../scripts/get_all_alias.cdc",
        [REGISTRY_OWNER.address, REGISTRY_NAME]
    )
    Test.expect(scriptResult, Test.beSucceeded())

    let aliases = scriptResult.returnValue! as! {String: String}

    Test.assertEqual(TEMPLATE_ID, aliases[ALIAS]!)
    Test.assertEqual(TEMPLATE_ID, aliases[NEW_ALIAS]!)
}

access(all)
fun testShouldUnlinkAlias() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let txResult = executeTransaction(
        "../transactions/unlink_alias.cdc",
        [ALIAS, REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(txResult, Test.beSucceeded())

    let typ = Type<FLIXRegistry.AliasUnlinked>()
    let events = Test.eventsOfType(typ)
    let event = events[0] as! FLIXRegistry.AliasUnlinked
    Test.assertEqual(1, events.length)
    Test.assertEqual(REGISTRY_NAME, event.registryName)
    Test.assertEqual(REGISTRY_OWNER.address, event.registryOwner)
    Test.assertEqual(ALIAS, event.alias)

    let scriptResult = executeScript(
        "../scripts/get_all_alias.cdc",
        [REGISTRY_OWNER.address, REGISTRY_NAME]
    )
    Test.expect(scriptResult, Test.beSucceeded())

    let aliases = scriptResult.returnValue! as! {String: String}

    Test.assertEqual(nil, aliases[ALIAS])
    Test.assertEqual(TEMPLATE_ID, aliases[NEW_ALIAS]!)
}

access(all)
fun testShouldDeprecateFlixWithTemplateId() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let lookupScriptResultBefore = executeScript(
        "../scripts/lookup.cdc",
        [REGISTRY_OWNER.address, TEMPLATE_ID, REGISTRY_NAME]
    )
    let flixBefore = lookupScriptResultBefore.returnValue! as! FLIXSchema_draft.FLIX
    Test.assertEqual(FLIXRegistry.FLIXStatus.active, flixBefore.status)

    let txResult = executeTransaction(
        "../transactions/deprecate_flix.cdc",
        [TEMPLATE_ID, REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(txResult, Test.beSucceeded())

    let typ = Type<FLIXRegistry.Deprecated>()
    let events = Test.eventsOfType(typ)
    let event = events[0] as! FLIXRegistry.Deprecated
    Test.assertEqual(1, events.length)
    Test.assertEqual(REGISTRY_NAME, event.registryName)
    Test.assertEqual(REGISTRY_OWNER.address, event.registryOwner)
    Test.assertEqual(TEMPLATE_ID, event.id)

    let lookupScriptResultAfter = executeScript(
        "../scripts/lookup.cdc",
        [REGISTRY_OWNER.address, TEMPLATE_ID, REGISTRY_NAME]
    )

    Test.expect(lookupScriptResultAfter, Test.beSucceeded())

    let flixAfter = lookupScriptResultAfter.returnValue! as! FLIXSchema_draft.FLIX
    Test.assertEqual(FLIXRegistry.FLIXStatus.deprecated, flixAfter.status)
}

access(all)
fun testShouldRemoveFlix() {
    Test.loadSnapshot(name: FLIX_PUBLISHED_SNAPSHOT)

    let removeTxResult = executeTransaction(
        "../transactions/remove_flix.cdc",
        [TEMPLATE_ID, REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(removeTxResult, Test.beSucceeded())

    let typ = Type<FLIXRegistry.Removed>()
    let events = Test.eventsOfType(typ)
    let event = events[0] as! FLIXRegistry.Removed
    Test.assertEqual(1, events.length)
    Test.assertEqual(REGISTRY_NAME, event.registryName)
    Test.assertEqual(REGISTRY_OWNER.address, event.registryOwner)
    Test.assertEqual(TEMPLATE_ID, event.id)

    let scriptResult = executeScript(
        "../scripts/get_registry_size.cdc",
        [REGISTRY_OWNER.address, REGISTRY_NAME]
    )
    Test.expect(scriptResult, Test.beSucceeded())

    let registrySize = scriptResult.returnValue! as! Int
    Test.assertEqual(0, registrySize)
}

access(all)
fun testShouldNotEmitEventWhenRemovingNonexistentFlix() {
    Test.loadSnapshot(name: REGISTRY_CREATED_SNAPSHOT)

    let removeTxResult = executeTransaction(
        "../transactions/remove_flix.cdc",
        [TEMPLATE_ID, REGISTRY_NAME],
        REGISTRY_OWNER
    )
    Test.expect(removeTxResult, Test.beSucceeded())

    let typ = Type<FLIXRegistry.Removed>()
    Test.assertEqual(0, Test.eventsOfType(typ).length)

    let scriptResult = executeScript(
        "../scripts/get_registry_size.cdc",
        [REGISTRY_OWNER.address, REGISTRY_NAME]
    )
    Test.expect(scriptResult, Test.beSucceeded())

    let registrySize = scriptResult.returnValue! as! Int
    Test.assertEqual(0, registrySize)
}

access(all)
fun testShouldThrowExeptionWhenDeprecatingNonexistentFlix() {
    Test.loadSnapshot(name: REGISTRY_CREATED_SNAPSHOT)

    Test.expectFailure(fun(): Void {
        let txResult = executeTransaction(
            "../transactions/deprecate_flix.cdc",
            [TEMPLATE_ID, REGISTRY_NAME],
            REGISTRY_OWNER
        )
        Test.expect(txResult, Test.beFailed())
        let err: Test.Error? = txResult.error
        panic(err!.message) // to trick expectFaliure, so we can match on the substring of the actual error
    }, errorMessageSubstring: "FLIX does not exist with the given id or alias: aTestId")

}

access(all)
fun testShouldCreatePublicPath() {
    Test.assertEqual(/public/flix_test, FLIXRegistry.PublicPath(name: "test"))
}

access(all)
fun testShouldCreateStoragePath() {
    Test.assertEqual(/storage/flix_test, FLIXRegistry.StoragePath(name: "test"))
}