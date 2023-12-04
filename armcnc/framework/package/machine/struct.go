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
	Joint5 JOINT `ini:"JOINT_5"`
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

type IniJson struct {
	Emc struct {
		Machine string `json:"Machine"`
		Debug   string `json:"Debug"`
		Version string `json:"Version"`
	} `json:"Emc"`
	Display struct {
		Display                string `json:"Display"`
		CycleTime              string `json:"CycleTime"`
		PositionOffset         string `json:"PositionOffset"`
		PositionFeedback       string `json:"PositionFeedback"`
		Arcdivision            string `json:"Arcdivision"`
		MaxFeedOverride        string `json:"MaxFeedOverride"`
		MinSpindleOverride     string `json:"MinSpindleOverride"`
		MaxSpindleOverride     string `json:"MaxSpindleOverride"`
		DefaultLinearVelocity  string `json:"DefaultLinearVelocity"`
		MinLinearVelocity      string `json:"MinLinearVelocity"`
		MaxLinearVelocity      string `json:"MaxLinearVelocity"`
		DefaultAngularVelocity string `json:"DefaultAngularVelocity"`
		MinAngularVelocity     string `json:"MinAngularVelocity"`
		MaxAngularVelocity     string `json:"MaxAngularVelocity"`
		ProgramPrefix          string `json:"ProgramPrefix"`
		OpenFile               string `json:"OpenFile"`
		Increments             string `json:"Increments"`
	} `json:"Display"`
	Python struct {
		PathAppend string `json:"PathAppend"`
		Toplevel   string `json:"Toplevel"`
	} `json:"Python"`
	Filter struct {
		ProgramExtension string `json:"ProgramExtension"`
		Py               string `json:"Py"`
	} `json:"Filter"`
	Rs274ngc struct {
		Features       string `json:"Features"`
		SubroutinePath string `json:"SubroutinePath"`
		ParameterFile  string `json:"ParameterFile"`
	} `json:"Rs274ngc"`
	Emcmot struct {
		Emcmot      string `json:"Emcmot"`
		CommTimeout string `json:"CommTimeout"`
		BasePeriod  string `json:"BasePeriod"`
		ServoPeriod string `json:"ServoPeriod"`
	} `json:"Emcmot"`
	Emcio struct {
		Emcio     string `json:"Emcio"`
		CycleTime string `json:"CycleTime"`
		ToolTable string `json:"ToolTable"`
	} `json:"Emcio"`
	Task struct {
		Task      string `json:"Task"`
		CycleTime string `json:"CycleTime"`
	} `json:"Task"`
	Hal struct {
		HalFile string `json:"HalFile"`
	} `json:"Hal"`
	Traj struct {
		Spindles     string `json:"Spindles"`
		Coordinates  string `json:"Coordinates"`
		LinearUnits  string `json:"LinearUnits"`
		AngularUnits string `json:"AngularUnits"`
		PositionFile string `json:"PositionFile"`
	} `json:"Traj"`
	Spindle0 SpindleJson `json:"Spindle0"`
	Kins     struct {
		Joints     string `json:"Joints"`
		Kinematics string `json:"Kinematics"`
	} `json:"Kins"`
	AxisX  AxisJson  `json:"AxisX"`
	Joint0 JointJson `json:"Joint0"`
	AxisY  AxisJson  `json:"AxisY"`
	Joint1 JointJson `json:"Joint1"`
	AxisZ  AxisJson  `json:"AxisZ"`
	Joint2 JointJson `json:"Joint2"`
	AxisA  AxisJson  `json:"AxisA"`
	Joint3 JointJson `json:"Joint3"`
	AxisB  AxisJson  `json:"AxisB"`
	Joint4 JointJson `json:"Joint4"`
	AxisC  AxisJson  `json:"AxisC"`
	Joint5 JointJson `json:"Joint5"`
}

type SpindleJson struct {
	MaxForwardVelocity string `json:"MaxForwardVelocity"`
	MinForwardVelocity string `json:"MinForwardVelocity"`
}

type AxisJson struct {
	MaxVelocity     string `json:"MaxVelocity"`
	MaxAcceleration string `json:"MaxAcceleration"`
	MinLimit        string `json:"MinLimit"`
	MaxLimit        string `json:"MaxLimit"`
}

type JointJson struct {
	Type             string `json:"Type"`
	Home             string `json:"Home"`
	MaxVelocity      string `json:"MaxVelocity"`
	MaxAcceleration  string `json:"MaxAcceleration"`
	StepgenMaxaccel  string `json:"StepgenMaxaccel"`
	Scale            string `json:"Scale"`
	Ferror           string `json:"Ferror"`
	MinLimit         string `json:"MinLimit"`
	MaxLimit         string `json:"MaxLimit"`
	HomeOffset       string `json:"HomeOffset"`
	HomeSearchVel    string `json:"HomeSearchVel"`
	HomeLarchVel     string `json:"HomeLarchVel"`
	HomeFinalVel     string `json:"HomeFinalVel"`
	VolatileHome     string `json:"VolatileHome"`
	HomeIgnoreLimits string `json:"HomeIgnoreLimits"`
	HomeUseIndex     string `json:"HomeUseIndex"`
	HomeSequence     string `json:"HomeSequence"`
	Backlash         string `json:"Backlash"`
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
	Io struct {
		EstopPin         string `ini:"ESTOP_PIN"`
		SpindleEnablePin string `ini:"SPINDLE_ENABLE_PIN"`
		SpindlePwmPin    string `ini:"SPINDLE_PWM_PIN"`
		XHomePin         string `ini:"X_HOME_PIN"`
		YHomePin         string `ini:"Y_HOME_PIN"`
		ZHomePin         string `ini:"Z_HOME_PIN"`
		AHomePin         string `ini:"A_HOME_PIN"`
		BHomePin         string `ini:"B_HOME_PIN"`
		CHomePin         string `ini:"C_HOME_PIN"`
	} `ini:"IO"`
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
	Io struct {
		EstopPin         string `json:"EstopPin"`
		SpindleEnablePin string `json:"SpindleEnablePin"`
		SpindlePwmPin    string `json:"SpindlePwmPin"`
		XHomePin         string `json:"XHomePin"`
		YHomePin         string `json:"YHomePin"`
		ZHomePin         string `json:"ZHomePin"`
		AHomePin         string `json:"AHomePin"`
		BHomePin         string `json:"BHomePin"`
		CHomePin         string `json:"CHomePin"`
	} `json:"Io"`
}
