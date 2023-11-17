/**
 ******************************************************************************
 * @file    struct.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachinePackage

type INI struct {
	Emc struct {
		Machine string `ini:"MACHINE"`
		Debug   string `ini:"DEBUG"`
		Version string `ini:"VERSION"`
	} `ini:"EMC"`
	Display struct {
		Display                string `ini:"DISPLAY"`
		CycleTime              string `ini:"CYCLE_TIME"`
		PositionOffset         string `ini:"POSITION_OFFSET"`
		PositionFeedback       string `ini:"POSITION_FEEDBACK"`
		Arcdivision            string `ini:"ARCDIVISION"`
		MaxFeedOverride        string `ini:"MAX_FEED_OVERRIDE"`
		MinSpindleOverride     string `ini:"MIN_SPINDLE_OVERRIDE"`
		MaxSpindleOverride     string `ini:"MAX_SPINDLE_OVERRIDE"`
		DefaultLinearVelocity  string `ini:"DEFAULT_LINEAR_VELOCITY"`
		MinLinearVelocity      string `ini:"MIN_LINEAR_VELOCITY"`
		MaxLinearVelocity      string `ini:"MAX_LINEAR_VELOCITY"`
		DefaultAngularVelocity string `ini:"DEFAULT_ANGULAR_VELOCITY"`
		MinAngularVelocity     string `ini:"MIN_ANGULAR_VELOCITY"`
		MaxAngularVelocity     string `ini:"MAX_ANGULAR_VELOCITY"`
		ProgramPrefix          string `ini:"PROGRAM_PREFIX"`
		OpenFile               string `ini:"OPEN_FILE"`
		Increments             string `ini:"INCREMENTS"`
	} `ini:"DISPLAY"`
	Python struct {
		PathAppend string `ini:"PATH_APPEND"`
		Toplevel   string `ini:"TOPLEVEL"`
	} `ini:"PYTHON"`
	Filter struct {
		ProgramExtension string `ini:"PROGRAM_EXTENSION"`
		Py               string `ini:"py"`
	} `ini:"FILTER"`
	Rs274ngc struct {
		Features       string `ini:"FEATURES"`
		SubroutinePath string `ini:"SUBROUTINE_PATH"`
		ParameterFile  string `ini:"PARAMETER_FILE"`
	} `ini:"RS274NGC"`
	Emcmot struct {
		Emcmot      string `ini:"EMCMOT"`
		CommTimeout string `ini:"COMM_TIMEOUT"`
		BasePeriod  string `ini:"BASE_PERIOD"`
		ServoPeriod string `ini:"SERVO_PERIOD"`
	} `ini:"EMCMOT"`
	Emcio struct {
		Emcio     string `ini:"EMCIO"`
		CycleTime string `ini:"CYCLE_TIME"`
		ToolTable string `ini:"TOOL_TABLE"`
	} `ini:"EMCIO"`
	Task struct {
		Task      string `ini:"TASK"`
		CycleTime string `ini:"CYCLE_TIME"`
	} `ini:"TASK"`
	Hal struct {
		HalFile string `ini:"HALFILE"`
	} `ini:"HAL"`
	Traj struct {
		Spindles     string `ini:"SPINDLES"`
		Coordinates  string `ini:"COORDINATES"`
		LinearUnits  string `ini:"LINEAR_UNITS"`
		AngularUnits string `ini:"ANGULAR_UNITS"`
		PositionFile string `ini:"POSITION_FILE"`
	} `ini:"TRAJ"`
	Spindle0 SPINDLE `ini:"SPINDLE_0"`
	Kins     struct {
		Joints     string `ini:"JOINTS"`
		Kinematics string `ini:"KINEMATICS"`
	} `ini:"KINS"`
	AxisX  AXIS  `ini:"AXIS_X"`
	Joint0 JOINT `ini:"JOINT_0"`
	AxisY  AXIS  `ini:"AXIS_Y"`
	Joint1 JOINT `ini:"JOINT_1"`
	AxisZ  AXIS  `ini:"AXIS_Z"`
	Joint2 JOINT `ini:"JOINT_2"`
	AxisA  AXIS  `ini:"AXIS_A"`
	Joint3 JOINT `ini:"JOINT_3"`
	AxisB  AXIS  `ini:"AXIS_B"`
	Joint4 JOINT `ini:"JOINT_4"`
	AxisC  AXIS  `ini:"AXIS_C"`
}

type SPINDLE struct {
	MaxForwardVelocity string `ini:"MAX_FORWARD_VELOCITY"`
	MinForwardVelocity string `ini:"MIN_FORWARD_VELOCITY"`
}

type AXIS struct {
	MaxVelocity     string `ini:"MAX_VELOCITY"`
	MaxAcceleration string `ini:"MAX_ACCELERATION"`
	MinLimit        string `ini:"MIN_LIMIT"`
	MaxLimit        string `ini:"MAX_LIMIT"`
}

type JOINT struct {
	Type             string `ini:"TYPE"`
	Home             string `ini:"HOME"`
	MaxVelocity      string `ini:"MAX_VELOCITY"`
	MaxAcceleration  string `ini:"MAX_ACCELERATION"`
	StepgenMaxaccel  string `ini:"STEPGEN_MAXACCEL"`
	Scale            string `ini:"SCALE"`
	Ferror           string `ini:"FERROR"`
	MinLimit         string `ini:"MIN_LIMIT"`
	MaxLimit         string `ini:"MAX_LIMIT"`
	HomeOffset       string `ini:"HOME_OFFSET"`
	HomeSearchVel    string `ini:"HOME_SEARCH_VEL"`
	HomeLarchVel     string `ini:"HOME_LATCH_VEL"`
	HomeFinalVel     string `ini:"HOME_FINAL_VEL"`
	VolatileHome     string `ini:"VOLATILE_HOME"`
	HomeIgnoreLimits string `ini:"HOME_IGNORE_LIMITS"`
	HomeUseIndex     string `ini:"HOME_USE_INDEX"`
	HomeSequence     string `ini:"HOME_SEQUENCE"`
	Backlash         string `ini:"BACKLASH"`
}

type USER struct {
	Base struct {
		Name     string `ini:"NAME"`
		Describe string `ini:"DESCRIBE"`
		Control  int    `ini:"CONTROL"`
	} `ini:"BASE"`
	HandWheel struct {
		Status    string `ini:"STATUS"`
		XVelocity string `ini:"X_VELOCITY"`
		YVelocity string `ini:"Y_VELOCITY"`
		ZVelocity string `ini:"Z_VELOCITY"`
		AVelocity string `ini:"A_VELOCITY"`
		BVelocity string `ini:"B_VELOCITY"`
		CVelocity string `ini:"C_VELOCITY"`
	} `ini:"HANDWHEEL"`
	Tool struct {
		Method              string `ini:"METHOD"`
		XPosition           string `ini:"X_POSITION"`
		YPosition           string `ini:"Y_POSITION"`
		ZPosition           string `ini:"Z_POSITION"`
		ZHeight             string `ini:"Z_HEIGHT"`
		MaxSearchDistance   string `ini:"MAX_SEARCH_DISTANCE"`
		LatchSearchDistance string `ini:"LATCH_SEARCH_DISTANCE"`
		SearchVelocity      string `ini:"SEARCH_VELOCITY"`
		LatchSearchVelocity string `ini:"LATCH_SEARCH_VELOCITY"`
		Pockets             string `ini:"POCKETS"`
	} `ini:"TOOL"`
}

type UserJson struct {
	Base struct {
		Name     string `json:"Name"`
		Describe string `json:"Describe"`
		Control  int    `json:"Control"`
	} `json:"Base"`
	HandWheel struct {
		Status    string `json:"Status"`
		XVelocity string `json:"XVelocity"`
		YVelocity string `json:"YVelocity"`
		ZVelocity string `json:"ZVelocity"`
		AVelocity string `json:"AVelocity"`
		BVelocity string `json:"BVelocity"`
		CVelocity string `json:"CVelocity"`
	} `json:"HandWheel"`
	Tool struct {
		Method              string `json:"Method"`
		XPosition           string `json:"XPosition"`
		YPosition           string `json:"YPosition"`
		ZPosition           string `json:"ZPosition"`
		ZHeight             string `json:"ZHeight"`
		MaxSearchDistance   string `json:"MaxSearchDistance"`
		LatchSearchDistance string `json:"LatchSearchDistance"`
		SearchVelocity      string `json:"SearchVelocity"`
		LatchSearchVelocity string `json:"LatchSearchVelocity"`
		Pockets             string `json:"Pockets"`
	} `json:"Tool"`
}
