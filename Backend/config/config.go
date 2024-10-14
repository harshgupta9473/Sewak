package config





func IfRoleExistINOURSYSTEM(role string)bool{
	
    roles:=[4]string{"customer","admin","Provider","worker"}
	for _,v:=range roles{
		if v==role{
			return true
		}
	}
	return false
}