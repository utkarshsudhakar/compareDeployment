package utils



func CheckEnv(env1,env2 string,envList []string) bool {
	count:=0
    for _, b := range envList {
        if b == env1 || b == env2 {
            count++
        }
    }
	if count==2{
		return true
	}
    return false
}