package nested

// Origin represents the type of environment we're running the app in
type Origin int

const (
	OriginUNDEFINED Origin = iota
	OriginENGLISH   Origin = 1
	OriginFRENCH    Origin = 2
	OriginGERMAN    Origin = 3
	OriginSPANISH   Origin = 4
	OriginITALIAN   Origin = 5
)

var Origins = map[int]string{
	int(OriginUNDEFINED): "- undefined -",
	int(OriginENGLISH):   "en",
	int(OriginFRENCH):    "fr",
	int(OriginGERMAN):    "de",
	int(OriginSPANISH):   "es",
	int(OriginITALIAN):   "it",
}

func (thisOrigin Origin) String() string {
	return Origins[int(thisOrigin)]
}

// Val helps implement the IEnum interface
func (thisOrigin Origin) Val() int {
	return int(thisOrigin)
}

// Values helps implement the IEnum interface
func (thisOrigin Origin) Values() map[int]string {
	return Origins
}
