package rights

type ShortRights struct {
	Name  string
	Rules []string
}

type OutRights struct {
	Id    string
	Name  string
	Rules []string
}

type RightsId struct {
	Value string
}

type RightsList struct {
	List []OutRights
}
