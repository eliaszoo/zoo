package registry

type Service struct {
	Name 		string
	Version 	string
	IP			string
	Port 		int
	MetaData 	map[string]string
}
