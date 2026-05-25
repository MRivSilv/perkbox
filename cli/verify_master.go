package cli

func verifyMasterPwd(masterPwd string) bool {
	verifiedPwd := readPassword("Confirm your Master Password: ")
	if verifiedPwd == masterPwd {
		return true
	}
	return false
}
