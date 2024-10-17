package test

func stCopy(dest *NginxBody, src *NginxDomain) {
	stCopyCopy{}.Copy(dest, src)
	return
}

type stCopyCopy struct{}

func (d stCopyCopy) Copy(dest *NginxBody, src *NginxDomain) {
	if src == nil {
		return
	}
	// basic map
	dest.Cluster = src.Cluster
	dest.Conf.ClientMaxBodySize = src.Conf.ClientMaxBodySize
	dest.Conf.Domain = src.Conf.Domain
	dest.Conf.Enable = src.Conf.Enable
	dest.Conf.EnableLimitConnZone = src.Conf.EnableLimitConnZone
	dest.Conf.ErrorPage = src.Conf.ErrorPage
	dest.Conf.LimitConnPerServer = src.Conf.LimitConnPerServer
	dest.Conf.Rewrite = src.Conf.Rewrite
	dest.Conf.SSLCertificate = src.Conf.SSLCertificate
	dest.Conf.SSLCertificateKey = src.Conf.SSLCertificateKey
	dest.Conf.SSLOn = src.Conf.SSLOn
	dest.Conf.Server = src.Conf.Server
	dest.Conf.ServiceType = src.Conf.ServiceType
	dest.ConfDomain = src.Conf.Domain
	dest.Domain = src.Domain
	dest.Product = src.Product
	dest.Project = src.Project
	dest.ProjectCname = src.GetName()
	dest.ProjectID = src.ProjectID
	// slice map
	dest.Conf.Listens = src.Conf.Listens
	dest.Conf.ListensSSL = src.Conf.ListensSSL
	dest.Conf.Locations = make([]LocationItem, 0, len(src.Conf.Locations))
	for i := 0; i < len(src.Conf.Locations); i++ {
		dest.Conf.Locations[i] = d.testNginxDomainLocationItemToTestLocationItem(src.Conf.Locations[i])
	}
	dest.Conf.UpstreamItems = make([]UpstreamItem, 0, len(src.Conf.UpstreamItems))
	for i := 0; i < len(src.Conf.UpstreamItems); i++ {
		dest.Conf.UpstreamItems[i] = d.testNginxDomainUpstreamItemToTestUpstreamItem(src.Conf.UpstreamItems[i])
	}
	dest.SliceStruct = make([]Password2, 0, len(src.SliceStruct))
	for i := 0; i < len(src.SliceStruct); i++ {
		dest.SliceStruct[i] = d.pTestPasswordToTestPassword2(src.SliceStruct[i])
	}
	dest.SliceStruct2 = make([]*Password2, 0, len(src.SliceStruct2))
	for i := 0; i < len(src.SliceStruct2); i++ {
		dest.SliceStruct2[i] = d.testPasswordToPTestPassword2(src.SliceStruct2[i])
	}
	// map map
	dest.Map = src.Map
	dest.MapStruct100 = make(map[string]*Password2, len(src.MapStruct))
	for key, value := range src.MapStruct {
		dest.MapStruct100[key] = d.testPasswordToPTestPassword2(value)
	}
	dest.MapStruct2 = make(map[string]Password2, len(src.MapStruct2))
	for key, value := range src.MapStruct2 {
		dest.MapStruct2[key] = d.pTestPasswordToTestPassword2(value)
	}
	// pointer map
	// method map
	return
}
func (d stCopyCopy) testNginxDomainLocationItemToTestLocationItem(src NginxDomainLocationItem) (dest LocationItem) {
	// basic map
	dest.ConfID = src.ConfID
	dest.HeaderHost = src.HeaderHost
	dest.Huidu.Enable = src.Huidu.Enable
	dest.Huidu.HuiduKey = src.Huidu.HuiduKey
	dest.Huidu.Upstream = src.Huidu.Upstream
	dest.Huidu.Upstreamhuidu = src.Huidu.Upstreamhuidu
	dest.Key = src.Key
	dest.LimitConnZone.Enable = src.LimitConnZone.Enable
	dest.LimitConnZone.PerServer = src.LimitConnZone.PerServer
	dest.LimitReqZone.Burst = src.LimitReqZone.Burst
	dest.LimitReqZone.Enable = src.LimitReqZone.Enable
	dest.LimitReqZone.Zone = src.LimitReqZone.Zone
	dest.Rewrite = src.Rewrite
	dest.SubDirectoryPath = src.SubDirectoryPath
	dest.UpstreamName = src.UpstreamName
	// slice map
	dest.Huidu.ArgsHuidu.Args = src.Huidu.ArgsHuidu.Args
	dest.Huidu.Content = src.Huidu.Content
	dest.Huidu.HeaderHuidu.Header = src.Huidu.HeaderHuidu.Header
	dest.Huidu.IPHuidu.Ips = src.Huidu.IPHuidu.Ips
	// map map
	// pointer map
	// method map
	return
}
func (d stCopyCopy) testNginxDomainUpstreamItemToTestUpstreamItem(src NginxDomainUpstreamItem) (dest UpstreamItem) {
	// basic map
	dest.CheckFall = src.CheckFall
	dest.CheckHTTPExpectAlive = src.CheckHTTPExpectAlive
	dest.CheckHTTPSend = src.CheckHTTPSend
	dest.CheckInterval = src.CheckInterval
	dest.CheckRise = src.CheckRise
	dest.CheckTimeout = src.CheckTimeout
	dest.ConfID = src.ConfID
	dest.LoadbalanceType = src.LoadbalanceType
	dest.LoadbalanceValue = src.LoadbalanceValue
	dest.Name = src.Name
	dest.PodServers.Port = src.PodServers.Port
	dest.PodServers.SyncFromPod = src.PodServers.SyncFromPod
	// slice map
	dest.Servers = make([]UpstreamServerItem, 0, len(src.Servers))
	for i := 0; i < len(src.Servers); i++ {
		dest.Servers[i] = d.testNginxDomainUpstreamServerItemToTestUpstreamServerItem(src.Servers[i])
	}
	// map map
	// pointer map
	// method map
	return
}
func (d stCopyCopy) testNginxDomainUpstreamServerItemToTestUpstreamServerItem(src NginxDomainUpstreamServerItem) (dest UpstreamServerItem) {
	// basic map
	dest.Flag = src.Flag
	dest.FromPod = src.FromPod
	dest.HP = src.HP
	dest.Healthy = src.Healthy
	dest.Weight = src.Weight
	// slice map
	// map map
	// pointer map
	// method map
	return
}
func (d stCopyCopy) pTestPasswordToTestPassword2(src *Password) (dest Password2) {
	if src == nil {
		return
	}
	// basic map
	dest.Nest.Ipone1 = src.Nest.Ipone
	dest.PasswordName = src.PasswordName
	// slice map
	dest.NestSlice = make([]Password2Nest, 0, len(src.NestSlice))
	for i := 0; i < len(src.NestSlice); i++ {
		dest.NestSlice[i] = d.testPasswordNestToTestPassword2Nest(src.NestSlice[i])
	}
	// map map
	dest.NestMap = make(map[string]Password2Nest, len(src.NestMap))
	for key, value := range src.NestMap {
		dest.NestMap[key] = d.testPasswordNestToTestPassword2Nest(value)
	}
	// pointer map
	// method map
	return
}
func (d stCopyCopy) testPasswordNestToTestPassword2Nest(src PasswordNest) (dest Password2Nest) {
	// basic map
	dest.Ipone1 = src.Ipone
	// slice map
	// map map
	// pointer map
	// method map
	return
}
func (d stCopyCopy) testPasswordToPTestPassword2(src Password) (dest *Password2) {
	dest = new(Password2)
	// basic map
	dest.Nest.Ipone1 = src.Nest.Ipone
	dest.PasswordName = src.PasswordName
	// slice map
	dest.NestSlice = make([]Password2Nest, 0, len(src.NestSlice))
	for i := 0; i < len(src.NestSlice); i++ {
		dest.NestSlice[i] = d.testPasswordNestToTestPassword2Nest(src.NestSlice[i])
	}
	// map map
	dest.NestMap = make(map[string]Password2Nest, len(src.NestMap))
	for key, value := range src.NestMap {
		dest.NestMap[key] = d.testPasswordNestToTestPassword2Nest(value)
	}
	// pointer map
	// method map
	return
}
func nginxDomain2NginxBodyDTO(dest *NginxBody, src *NginxDomain) {
	nginxDomain2NginxBodyDTOCopy{}.Copy(dest, src)
	return
}

type nginxDomain2NginxBodyDTOCopy struct{}

func (d nginxDomain2NginxBodyDTOCopy) Copy(dest *NginxBody, src *NginxDomain) {
	if src == nil {
		return
	}
	// basic map
	dest.Cluster = src.Cluster
	dest.Conf.ClientMaxBodySize = src.Conf.ClientMaxBodySize
	dest.Conf.Domain = src.Conf.Domain
	dest.Conf.Enable = src.Conf.Enable
	dest.Conf.EnableLimitConnZone = src.Conf.EnableLimitConnZone
	dest.Conf.ErrorPage = src.Conf.ErrorPage
	dest.Conf.LimitConnPerServer = src.Conf.LimitConnPerServer
	dest.Conf.Rewrite = src.Conf.Rewrite
	dest.Conf.SSLCertificate = src.Conf.SSLCertificate
	dest.Conf.SSLCertificateKey = src.Conf.SSLCertificateKey
	dest.Conf.SSLOn = src.Conf.SSLOn
	dest.Conf.Server = src.Conf.Server
	dest.Conf.ServiceType = src.Conf.ServiceType
	dest.ConfDomain = src.Conf.Domain
	dest.Domain = src.Domain
	dest.Product = src.Product
	dest.Project = src.Project
	dest.ProjectCname = src.GetName()
	dest.ProjectID = src.ProjectID
	// slice map
	dest.Conf.Listens = src.Conf.Listens
	dest.Conf.ListensSSL = src.Conf.ListensSSL
	dest.Conf.Locations = make([]LocationItem, 0, len(src.Conf.Locations))
	for i := 0; i < len(src.Conf.Locations); i++ {
		dest.Conf.Locations[i] = d.testNginxDomainLocationItemToTestLocationItem(src.Conf.Locations[i])
	}
	dest.Conf.UpstreamItems = make([]UpstreamItem, 0, len(src.Conf.UpstreamItems))
	for i := 0; i < len(src.Conf.UpstreamItems); i++ {
		dest.Conf.UpstreamItems[i] = d.testNginxDomainUpstreamItemToTestUpstreamItem(src.Conf.UpstreamItems[i])
	}
	dest.SliceStruct = make([]Password2, 0, len(src.SliceStruct))
	for i := 0; i < len(src.SliceStruct); i++ {
		dest.SliceStruct[i] = d.pTestPasswordToTestPassword2(src.SliceStruct[i])
	}
	dest.SliceStruct2 = make([]*Password2, 0, len(src.SliceStruct2))
	for i := 0; i < len(src.SliceStruct2); i++ {
		dest.SliceStruct2[i] = d.testPasswordToPTestPassword2(src.SliceStruct2[i])
	}
	// map map
	dest.Map = src.Map
	dest.MapStruct100 = make(map[string]*Password2, len(src.MapStruct))
	for key, value := range src.MapStruct {
		dest.MapStruct100[key] = d.testPasswordToPTestPassword2(value)
	}
	dest.MapStruct2 = make(map[string]Password2, len(src.MapStruct2))
	for key, value := range src.MapStruct2 {
		dest.MapStruct2[key] = d.pTestPasswordToTestPassword2(value)
	}
	// pointer map
	// method map
	return
}
func (d nginxDomain2NginxBodyDTOCopy) testNginxDomainLocationItemToTestLocationItem(src NginxDomainLocationItem) (dest LocationItem) {
	// basic map
	dest.ConfID = src.ConfID
	dest.HeaderHost = src.HeaderHost
	dest.Huidu.Enable = src.Huidu.Enable
	dest.Huidu.HuiduKey = src.Huidu.HuiduKey
	dest.Huidu.Upstream = src.Huidu.Upstream
	dest.Huidu.Upstreamhuidu = src.Huidu.Upstreamhuidu
	dest.Key = src.Key
	dest.LimitConnZone.Enable = src.LimitConnZone.Enable
	dest.LimitConnZone.PerServer = src.LimitConnZone.PerServer
	dest.LimitReqZone.Burst = src.LimitReqZone.Burst
	dest.LimitReqZone.Enable = src.LimitReqZone.Enable
	dest.LimitReqZone.Zone = src.LimitReqZone.Zone
	dest.Rewrite = src.Rewrite
	dest.SubDirectoryPath = src.SubDirectoryPath
	dest.UpstreamName = src.UpstreamName
	// slice map
	dest.Huidu.ArgsHuidu.Args = src.Huidu.ArgsHuidu.Args
	dest.Huidu.Content = src.Huidu.Content
	dest.Huidu.HeaderHuidu.Header = src.Huidu.HeaderHuidu.Header
	dest.Huidu.IPHuidu.Ips = src.Huidu.IPHuidu.Ips
	// map map
	// pointer map
	// method map
	return
}
func (d nginxDomain2NginxBodyDTOCopy) testNginxDomainUpstreamItemToTestUpstreamItem(src NginxDomainUpstreamItem) (dest UpstreamItem) {
	// basic map
	dest.CheckFall = src.CheckFall
	dest.CheckHTTPExpectAlive = src.CheckHTTPExpectAlive
	dest.CheckHTTPSend = src.CheckHTTPSend
	dest.CheckInterval = src.CheckInterval
	dest.CheckRise = src.CheckRise
	dest.CheckTimeout = src.CheckTimeout
	dest.ConfID = src.ConfID
	dest.LoadbalanceType = src.LoadbalanceType
	dest.LoadbalanceValue = src.LoadbalanceValue
	dest.Name = src.Name
	dest.PodServers.Port = src.PodServers.Port
	dest.PodServers.SyncFromPod = src.PodServers.SyncFromPod
	// slice map
	dest.Servers = make([]UpstreamServerItem, 0, len(src.Servers))
	for i := 0; i < len(src.Servers); i++ {
		dest.Servers[i] = d.testNginxDomainUpstreamServerItemToTestUpstreamServerItem(src.Servers[i])
	}
	// map map
	// pointer map
	// method map
	return
}
func (d nginxDomain2NginxBodyDTOCopy) testNginxDomainUpstreamServerItemToTestUpstreamServerItem(src NginxDomainUpstreamServerItem) (dest UpstreamServerItem) {
	// basic map
	dest.Flag = src.Flag
	dest.FromPod = src.FromPod
	dest.HP = src.HP
	dest.Healthy = src.Healthy
	dest.Weight = src.Weight
	// slice map
	// map map
	// pointer map
	// method map
	return
}
func (d nginxDomain2NginxBodyDTOCopy) pTestPasswordToTestPassword2(src *Password) (dest Password2) {
	if src == nil {
		return
	}
	// basic map
	dest.Nest.Ipone1 = src.Nest.Ipone
	dest.PasswordName = src.PasswordName
	// slice map
	dest.NestSlice = make([]Password2Nest, 0, len(src.NestSlice))
	for i := 0; i < len(src.NestSlice); i++ {
		dest.NestSlice[i] = d.testPasswordNestToTestPassword2Nest(src.NestSlice[i])
	}
	// map map
	dest.NestMap = make(map[string]Password2Nest, len(src.NestMap))
	for key, value := range src.NestMap {
		dest.NestMap[key] = d.testPasswordNestToTestPassword2Nest(value)
	}
	// pointer map
	// method map
	return
}
func (d nginxDomain2NginxBodyDTOCopy) testPasswordNestToTestPassword2Nest(src PasswordNest) (dest Password2Nest) {
	// basic map
	dest.Ipone1 = src.Ipone
	// slice map
	// map map
	// pointer map
	// method map
	return
}
func (d nginxDomain2NginxBodyDTOCopy) testPasswordToPTestPassword2(src Password) (dest *Password2) {
	dest = new(Password2)
	// basic map
	dest.Nest.Ipone1 = src.Nest.Ipone
	dest.PasswordName = src.PasswordName
	// slice map
	dest.NestSlice = make([]Password2Nest, 0, len(src.NestSlice))
	for i := 0; i < len(src.NestSlice); i++ {
		dest.NestSlice[i] = d.testPasswordNestToTestPassword2Nest(src.NestSlice[i])
	}
	// map map
	dest.NestMap = make(map[string]Password2Nest, len(src.NestMap))
	for key, value := range src.NestMap {
		dest.NestMap[key] = d.testPasswordNestToTestPassword2Nest(value)
	}
	// pointer map
	// method map
	return
}
