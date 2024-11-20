package contracts

func GetCoreContracts() map[string]map[string]string {
	// TODO: this should be from core contracts module
	return map[string]map[string]string{}
}

func GetCoreContractForNetwork(contractName string, network string) string {
	// TODO: this should be from core contracts module
	c := GetCoreContracts()
	return c[contractName][network]
}
