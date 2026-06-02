package iot

type BootstrapStatus int

const (
	BootstrapStatusACCESSxDENIED      BootstrapStatus = -1
	BootstrapStatusAUTHENTICATEDxONLY BootstrapStatus = 1
	BootstrapStatusPENDING            BootstrapStatus = 2
	BootstrapStatusREADY              BootstrapStatus = 3
)

var bootstrapStatuses = map[int]string{
	int(BootstrapStatusACCESSxDENIED):      "denied",
	int(BootstrapStatusAUTHENTICATEDxONLY): "authenticated only",
	int(BootstrapStatusPENDING):            "pending",
	int(BootstrapStatusREADY):              "ready",
}

func (thisBootstrapStatus BootstrapStatus) String() string {
	return bootstrapStatuses[int(thisBootstrapStatus)]
}

// Val helps implement the IEnum interface
func (thisBootstrapStatus BootstrapStatus) Val() int {
	return int(thisBootstrapStatus)
}

// Values helps implement the IEnum interface
func (thisBootstrapStatus BootstrapStatus) Values() map[int]string {
	return bootstrapStatuses
}
