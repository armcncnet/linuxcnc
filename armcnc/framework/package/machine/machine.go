/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachinePackage

import (
	"armcnc/framework/config"
	"armcnc/framework/utils/file"
	"armcnc/framework/utils/ini"
	"github.com/djherbis/times"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Machine struct {
	Path string `json:"path"`
}

type Data struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Describe     string    `json:"describe"`
	Version      string    `json:"version"`
	Control      int       `json:"control"`
	Coordinate   string    `json:"coordinate"`
	Increments   string    `json:"increments"`
	LinearUnits  string    `json:"linear_units"`
	AngularUnits string    `json:"angular_units"`
	Time         time.Time `json:"-"`
}

func Init() *Machine {
	return &Machine{
		Path: Config.Get.Basic.Workspace + "/configs/",
	}
}

func (machine *Machine) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(machine.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		if file.IsDir() {
			item.Path = file.Name()
			timeData, _ := times.Stat(machine.Path + file.Name())
			item.Time = timeData.BirthTime()
			if strings.Contains(file.Name(), "default_") {
				item.Time = item.Time.Add(-10 * time.Minute)
			}
			ini := machine.GetIni(file.Name())
			if ini.Emc.Version != "" {
				item.Version = ini.Emc.Version
				item.Coordinate = ini.Traj.Coordinates
				item.Increments = ini.Display.Increments
				item.LinearUnits = ini.Traj.LinearUnits
				item.AngularUnits = ini.Traj.AngularUnits

				user := machine.GetUser(file.Name())
				item.Name = user.Base.Name
				item.Describe = user.Base.Describe
				item.Control = user.Base.Control
				data = append(data, item)
			}
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Time.After(data[j].Time)
	})
	return data
}

func (machine *Machine) GetUser(path string) USER {
	data := USER{}
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.user")
	if exists {
		userFile, err := IniUtils.Load(machine.Path + path + "/machine.user")
		if err == nil {
			err = IniUtils.MapTo(userFile, &data)
		}
	}
	return data
}

func (machine *Machine) UpdateUser(path string, data UserJson) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.user")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + path + "/machine.user")
		if err == nil {
			iniFile.Section("BASE").Key("NAME").SetValue(data.Base.Name)
			iniFile.Section("BASE").Key("DESCRIBE").SetValue(data.Base.Describe)
			iniFile.Section("BASE").Key("CONTROL").SetValue(strconv.Itoa(data.Base.Control))
			iniFile.Section("HANDWHEEL").Key("X_VELOCITY").SetValue(data.HandWheel.XVelocity)
			iniFile.Section("HANDWHEEL").Key("Y_VELOCITY").SetValue(data.HandWheel.YVelocity)
			iniFile.Section("HANDWHEEL").Key("Z_VELOCITY").SetValue(data.HandWheel.ZVelocity)
			iniFile.Section("HANDWHEEL").Key("A_VELOCITY").SetValue(data.HandWheel.AVelocity)
			iniFile.Section("HANDWHEEL").Key("B_VELOCITY").SetValue(data.HandWheel.BVelocity)
			iniFile.Section("HANDWHEEL").Key("C_VELOCITY").SetValue(data.HandWheel.CVelocity)
			iniFile.Section("TOOL").Key("METHOD").SetValue(data.Tool.Method)
			iniFile.Section("TOOL").Key("X_POSITION").SetValue(data.Tool.XPosition)
			iniFile.Section("TOOL").Key("Y_POSITION").SetValue(data.Tool.YPosition)
			iniFile.Section("TOOL").Key("Z_POSITION").SetValue(data.Tool.ZPosition)
			iniFile.Section("TOOL").Key("Z_HEIGHT").SetValue(data.Tool.ZHeight)
			iniFile.Section("TOOL").Key("MAX_SEARCH_DISTANCE").SetValue(data.Tool.MaxSearchDistance)
			iniFile.Section("TOOL").Key("LATCH_SEARCH_DISTANCE").SetValue(data.Tool.LatchSearchDistance)
			iniFile.Section("TOOL").Key("SEARCH_VELOCITY").SetValue(data.Tool.SearchVelocity)
			iniFile.Section("TOOL").Key("LATCH_SEARCH_VELOCITY").SetValue(data.Tool.LatchSearchVelocity)
			iniFile.Section("TOOL").Key("POCKETS").SetValue(data.Tool.Pockets)
			err = IniUtils.SaveTo(iniFile, machine.Path+path+"/machine.user")
			if err == nil {
				status = true
			}
		}
	}
	return status
}

func (machine *Machine) GetIni(path string) INI {
	data := INI{}
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + path + "/machine.ini")
		if err == nil {
			err = IniUtils.MapTo(iniFile, &data)
			if data.AxisA.MaxVelocity == "" {
				data.AxisA.MaxVelocity = "0.000"
			}
			if data.AxisA.MaxAcceleration == "" {
				data.AxisA.MaxAcceleration = "0.000"
			}
			if data.AxisA.MinLimit == "" {
				data.AxisA.MinLimit = "0.000"
			}
			if data.AxisA.MaxLimit == "" {
				data.AxisA.MaxLimit = "0.000"
			}
			if data.Joint3.Type == "" {
				data.Joint3.Type = "LINEAR"
			}
			if data.Joint3.Home == "" {
				data.Joint3.Home = "0.000"
			}
			if data.Joint3.MaxVelocity == "" {
				data.Joint3.MaxVelocity = "0.000"
			}
			if data.Joint3.MaxAcceleration == "" {
				data.Joint3.MaxAcceleration = "0.000"
			}
			if data.Joint3.StepgenMaxaccel == "" {
				data.Joint3.StepgenMaxaccel = "0.000"
			}
			if data.Joint3.Scale == "" {
				data.Joint3.Scale = "1600"
			}
			if data.Joint3.Ferror == "" {
				data.Joint3.Ferror = "1.000"
			}
			if data.Joint3.MinLimit == "" {
				data.Joint3.MinLimit = "0.000"
			}
			if data.Joint3.MaxLimit == "" {
				data.Joint3.MaxLimit = "0.000"
			}
			if data.Joint3.HomeOffset == "" {
				data.Joint3.HomeOffset = "0.000"
			}
			if data.Joint3.HomeSearchVel == "" {
				data.Joint3.HomeSearchVel = "0.000"
			}
			if data.Joint3.HomeLarchVel == "" {
				data.Joint3.HomeLarchVel = "0.000"
			}
			if data.Joint3.HomeFinalVel == "" {
				data.Joint3.HomeFinalVel = "0.000"
			}
			if data.Joint3.VolatileHome == "" {
				data.Joint3.VolatileHome = "1"
			}
			if data.Joint3.HomeIgnoreLimits == "" {
				data.Joint3.HomeIgnoreLimits = "NO"
			}
			if data.Joint3.HomeUseIndex == "" {
				data.Joint3.HomeUseIndex = "NO"
			}
			if data.Joint3.HomeSequence == "" {
				data.Joint3.HomeSequence = "0"
			}
			if data.Joint3.Backlash == "" {
				data.Joint3.Backlash = "0.00"
			}
			if data.AxisB.MaxVelocity == "" {
				data.AxisB.MaxVelocity = "0.000"
			}
			if data.AxisB.MaxAcceleration == "" {
				data.AxisB.MaxAcceleration = "0.000"
			}
			if data.AxisB.MinLimit == "" {
				data.AxisB.MinLimit = "0.000"
			}
			if data.AxisB.MaxLimit == "" {
				data.AxisB.MaxLimit = "0.000"
			}
			if data.Joint4.Type == "" {
				data.Joint4.Type = "LINEAR"
			}
			if data.Joint4.Home == "" {
				data.Joint4.Home = "0.000"
			}
			if data.Joint4.MaxVelocity == "" {
				data.Joint4.MaxVelocity = "0.000"
			}
			if data.Joint4.MaxAcceleration == "" {
				data.Joint4.MaxAcceleration = "0.000"
			}
			if data.Joint4.StepgenMaxaccel == "" {
				data.Joint4.StepgenMaxaccel = "0.000"
			}
			if data.Joint4.Scale == "" {
				data.Joint4.Scale = "1600"
			}
			if data.Joint4.Ferror == "" {
				data.Joint4.Ferror = "1.000"
			}
			if data.Joint4.MinLimit == "" {
				data.Joint4.MinLimit = "0.000"
			}
			if data.Joint4.MaxLimit == "" {
				data.Joint4.MaxLimit = "0.000"
			}
			if data.Joint4.HomeOffset == "" {
				data.Joint4.HomeOffset = "0.000"
			}
			if data.Joint4.HomeSearchVel == "" {
				data.Joint4.HomeSearchVel = "0.000"
			}
			if data.Joint4.HomeLarchVel == "" {
				data.Joint4.HomeLarchVel = "0.000"
			}
			if data.Joint4.HomeFinalVel == "" {
				data.Joint4.HomeFinalVel = "0.000"
			}
			if data.Joint4.VolatileHome == "" {
				data.Joint4.VolatileHome = "1"
			}
			if data.Joint4.HomeIgnoreLimits == "" {
				data.Joint4.HomeIgnoreLimits = "NO"
			}
			if data.Joint4.HomeUseIndex == "" {
				data.Joint4.HomeUseIndex = "NO"
			}
			if data.Joint4.HomeSequence == "" {
				data.Joint4.HomeSequence = "0"
			}
			if data.Joint4.Backlash == "" {
				data.Joint4.Backlash = "0.00"
			}
			if data.AxisC.MaxVelocity == "" {
				data.AxisC.MaxVelocity = "0.000"
			}
			if data.AxisC.MaxAcceleration == "" {
				data.AxisC.MaxAcceleration = "0.000"
			}
			if data.AxisC.MinLimit == "" {
				data.AxisC.MinLimit = "0.000"
			}
			if data.AxisC.MaxLimit == "" {
				data.AxisC.MaxLimit = "0.000"
			}
			if data.Joint5.Type == "" {
				data.Joint5.Type = "LINEAR"
			}
			if data.Joint5.Home == "" {
				data.Joint5.Home = "0.000"
			}
			if data.Joint5.MaxVelocity == "" {
				data.Joint5.MaxVelocity = "0.000"
			}
			if data.Joint5.MaxAcceleration == "" {
				data.Joint5.MaxAcceleration = "0.000"
			}
			if data.Joint5.StepgenMaxaccel == "" {
				data.Joint5.StepgenMaxaccel = "0.000"
			}
			if data.Joint5.Scale == "" {
				data.Joint5.Scale = "1600"
			}
			if data.Joint5.Ferror == "" {
				data.Joint5.Ferror = "1.000"
			}
			if data.Joint5.MinLimit == "" {
				data.Joint5.MinLimit = "0.000"
			}
			if data.Joint5.MaxLimit == "" {
				data.Joint5.MaxLimit = "0.000"
			}
			if data.Joint5.HomeOffset == "" {
				data.Joint5.HomeOffset = "0.000"
			}
			if data.Joint5.HomeSearchVel == "" {
				data.Joint5.HomeSearchVel = "0.000"
			}
			if data.Joint5.HomeLarchVel == "" {
				data.Joint5.HomeLarchVel = "0.000"
			}
			if data.Joint5.HomeFinalVel == "" {
				data.Joint5.HomeFinalVel = "0.000"
			}
			if data.Joint5.VolatileHome == "" {
				data.Joint5.VolatileHome = "1"
			}
			if data.Joint5.HomeIgnoreLimits == "" {
				data.Joint5.HomeIgnoreLimits = "NO"
			}
			if data.Joint5.HomeUseIndex == "" {
				data.Joint5.HomeUseIndex = "NO"
			}
			if data.Joint5.HomeSequence == "" {
				data.Joint5.HomeSequence = "0"
			}
			if data.Joint5.Backlash == "" {
				data.Joint5.Backlash = "0.00"
			}
		}
	}
	return data
}

func (machine *Machine) UpdateIni(path string, data IniJson) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + path + "/machine.ini")
		if err == nil {
			iniFile.Section("EMC").Key("MACHINE").SetValue(data.Emc.Machine)
			iniFile.Section("EMC").Key("DEBUG").SetValue(data.Emc.Debug)
			iniFile.Section("EMC").Key("VERSION").SetValue(data.Emc.Version)
			iniFile.Section("DISPLAY").Key("DISPLAY").SetValue(data.Display.Display)
			iniFile.Section("DISPLAY").Key("CYCLE_TIME").SetValue(data.Display.CycleTime)
			iniFile.Section("DISPLAY").Key("POSITION_OFFSET").SetValue(data.Display.PositionOffset)
			iniFile.Section("DISPLAY").Key("POSITION_FEEDBACK").SetValue(data.Display.PositionFeedback)
			iniFile.Section("DISPLAY").Key("ARCDIVISION").SetValue(data.Display.Arcdivision)
			iniFile.Section("DISPLAY").Key("MAX_FEED_OVERRIDE").SetValue(data.Display.MaxFeedOverride)
			iniFile.Section("DISPLAY").Key("MIN_SPINDLE_OVERRIDE").SetValue(data.Display.MinSpindleOverride)
			iniFile.Section("DISPLAY").Key("MAX_SPINDLE_OVERRIDE").SetValue(data.Display.MaxSpindleOverride)
			iniFile.Section("DISPLAY").Key("DEFAULT_LINEAR_VELOCITY").SetValue(data.Display.DefaultLinearVelocity)
			iniFile.Section("DISPLAY").Key("MIN_LINEAR_VELOCITY").SetValue(data.Display.MinLinearVelocity)
			iniFile.Section("DISPLAY").Key("DEFAULT_ANGULAR_VELOCITY").SetValue(data.Display.DefaultAngularVelocity)
			iniFile.Section("DISPLAY").Key("MIN_ANGULAR_VELOCITY").SetValue(data.Display.MinAngularVelocity)
			iniFile.Section("DISPLAY").Key("MAX_ANGULAR_VELOCITY").SetValue(data.Display.MaxAngularVelocity)
			iniFile.Section("DISPLAY").Key("PROGRAM_PREFIX").SetValue(data.Display.ProgramPrefix)
			iniFile.Section("DISPLAY").Key("OPEN_FILE").SetValue(data.Display.OpenFile)
			iniFile.Section("DISPLAY").Key("INCREMENTS").SetValue(data.Display.Increments)
			iniFile.Section("PYTHON").Key("PATH_APPEND").SetValue(data.Python.PathAppend)
			iniFile.Section("PYTHON").Key("TOPLEVEL").SetValue(data.Python.Toplevel)
			iniFile.Section("FILTER").Key("PROGRAM_EXTENSION").SetValue(data.Filter.ProgramExtension)
			iniFile.Section("FILTER").Key("py").SetValue(data.Filter.Py)
			iniFile.Section("RS274NGC").Key("FEATURES").SetValue(data.Rs274ngc.Features)
			iniFile.Section("RS274NGC").Key("SUBROUTINE_PATH").SetValue(data.Rs274ngc.SubroutinePath)
			iniFile.Section("RS274NGC").Key("PARAMETER_FILE").SetValue(data.Rs274ngc.ParameterFile)
			iniFile.Section("EMCMOT").Key("EMCMOT").SetValue(data.Emcmot.Emcmot)
			iniFile.Section("EMCMOT").Key("COMM_TIMEOUT").SetValue(data.Emcmot.CommTimeout)
			iniFile.Section("EMCMOT").Key("BASE_PERIOD").SetValue(data.Emcmot.BasePeriod)
			iniFile.Section("EMCMOT").Key("SERVO_PERIOD").SetValue(data.Emcmot.ServoPeriod)
			iniFile.Section("EMCIO").Key("EMCIO").SetValue(data.Emcio.Emcio)
			iniFile.Section("EMCIO").Key("CYCLE_TIME").SetValue(data.Emcio.CycleTime)
			iniFile.Section("EMCIO").Key("TOOL_TABLE").SetValue(data.Emcio.ToolTable)
			iniFile.Section("TASK").Key("TASK").SetValue(data.Task.Task)
			iniFile.Section("TASK").Key("CYCLE_TIME").SetValue(data.Task.CycleTime)
			iniFile.Section("HAL").Key("HALFILE").SetValue(data.Hal.HalFile)
			iniFile.Section("TRAJ").Key("SPINDLES").SetValue(data.Traj.Spindles)
			iniFile.Section("TRAJ").Key("COORDINATES").SetValue(data.Traj.Coordinates)
			iniFile.Section("TRAJ").Key("LINEAR_UNITS").SetValue(data.Traj.LinearUnits)
			iniFile.Section("TRAJ").Key("ANGULAR_UNITS").SetValue(data.Traj.AngularUnits)
			iniFile.Section("TRAJ").Key("POSITION_FILE").SetValue(data.Traj.PositionFile)
			iniFile.Section("SPINDLE_0").Key("MAX_FORWARD_VELOCITY").SetValue(data.Spindle0.MaxForwardVelocity)
			iniFile.Section("SPINDLE_0").Key("MIN_FORWARD_VELOCITY").SetValue(data.Spindle0.MinForwardVelocity)
			iniFile.Section("KINS").Key("JOINTS").SetValue(data.Kins.Joints)
			iniFile.Section("KINS").Key("KINEMATICS").SetValue(data.Kins.Kinematics)
			iniFile.Section("AXIS_X").Key("MAX_VELOCITY").SetValue(data.AxisX.MaxVelocity)
			iniFile.Section("AXIS_X").Key("MAX_ACCELERATION").SetValue(data.AxisX.MaxAcceleration)
			iniFile.Section("AXIS_X").Key("MIN_LIMIT").SetValue(data.AxisX.MinLimit)
			iniFile.Section("AXIS_X").Key("MAX_LIMIT").SetValue(data.AxisX.MaxLimit)
			iniFile.Section("Joint0").Key("TYPE").SetValue(data.Joint0.Type)
			iniFile.Section("Joint0").Key("HOME").SetValue(data.Joint0.Home)
			iniFile.Section("Joint0").Key("MAX_VELOCITY").SetValue(data.Joint0.MaxVelocity)
			iniFile.Section("Joint0").Key("MAX_ACCELERATION").SetValue(data.Joint0.MaxAcceleration)
			iniFile.Section("Joint0").Key("STEPGEN_MAXACCEL").SetValue(data.Joint0.StepgenMaxaccel)
			iniFile.Section("Joint0").Key("SCALE").SetValue(data.Joint0.Scale)
			iniFile.Section("Joint0").Key("FERROR").SetValue(data.Joint0.Ferror)
			iniFile.Section("Joint0").Key("MIN_LIMIT").SetValue(data.Joint0.MinLimit)
			iniFile.Section("Joint0").Key("MAX_LIMIT").SetValue(data.Joint0.MaxLimit)
			iniFile.Section("Joint0").Key("HOME_OFFSET").SetValue(data.Joint0.HomeOffset)
			iniFile.Section("Joint0").Key("HOME_SEARCH_VEL").SetValue(data.Joint0.HomeSearchVel)
			iniFile.Section("Joint0").Key("HOME_LATCH_VEL").SetValue(data.Joint0.HomeLarchVel)
			iniFile.Section("Joint0").Key("HOME_FINAL_VEL").SetValue(data.Joint0.HomeFinalVel)
			iniFile.Section("Joint0").Key("VOLATILE_HOME").SetValue(data.Joint0.VolatileHome)
			iniFile.Section("Joint0").Key("HOME_IGNORE_LIMITS").SetValue(data.Joint0.HomeIgnoreLimits)
			iniFile.Section("Joint0").Key("HOME_USE_INDEX").SetValue(data.Joint0.HomeUseIndex)
			iniFile.Section("Joint0").Key("HOME_SEQUENCE").SetValue(data.Joint0.HomeSequence)
			iniFile.Section("Joint0").Key("BACKLASH").SetValue(data.Joint0.Backlash)
			iniFile.Section("AXIS_Y").Key("MAX_VELOCITY").SetValue(data.AxisY.MaxVelocity)
			iniFile.Section("AXIS_Y").Key("MAX_ACCELERATION").SetValue(data.AxisY.MaxAcceleration)
			iniFile.Section("AXIS_Y").Key("MIN_LIMIT").SetValue(data.AxisY.MinLimit)
			iniFile.Section("AXIS_Y").Key("MAX_LIMIT").SetValue(data.AxisY.MaxLimit)
			iniFile.Section("Joint1").Key("TYPE").SetValue(data.Joint1.Type)
			iniFile.Section("Joint1").Key("HOME").SetValue(data.Joint1.Home)
			iniFile.Section("Joint1").Key("MAX_VELOCITY").SetValue(data.Joint1.MaxVelocity)
			iniFile.Section("Joint1").Key("MAX_ACCELERATION").SetValue(data.Joint1.MaxAcceleration)
			iniFile.Section("Joint1").Key("STEPGEN_MAXACCEL").SetValue(data.Joint1.StepgenMaxaccel)
			iniFile.Section("Joint1").Key("SCALE").SetValue(data.Joint1.Scale)
			iniFile.Section("Joint1").Key("FERROR").SetValue(data.Joint1.Ferror)
			iniFile.Section("Joint1").Key("MIN_LIMIT").SetValue(data.Joint1.MinLimit)
			iniFile.Section("Joint1").Key("MAX_LIMIT").SetValue(data.Joint1.MaxLimit)
			iniFile.Section("Joint1").Key("HOME_OFFSET").SetValue(data.Joint1.HomeOffset)
			iniFile.Section("Joint1").Key("HOME_SEARCH_VEL").SetValue(data.Joint1.HomeSearchVel)
			iniFile.Section("Joint1").Key("HOME_LATCH_VEL").SetValue(data.Joint1.HomeLarchVel)
			iniFile.Section("Joint1").Key("HOME_FINAL_VEL").SetValue(data.Joint1.HomeFinalVel)
			iniFile.Section("Joint1").Key("VOLATILE_HOME").SetValue(data.Joint1.VolatileHome)
			iniFile.Section("Joint1").Key("HOME_IGNORE_LIMITS").SetValue(data.Joint1.HomeIgnoreLimits)
			iniFile.Section("Joint1").Key("HOME_USE_INDEX").SetValue(data.Joint1.HomeUseIndex)
			iniFile.Section("Joint1").Key("HOME_SEQUENCE").SetValue(data.Joint1.HomeSequence)
			iniFile.Section("Joint1").Key("BACKLASH").SetValue(data.Joint1.Backlash)
			iniFile.Section("AXIS_Z").Key("MAX_VELOCITY").SetValue(data.AxisZ.MaxVelocity)
			iniFile.Section("AXIS_Z").Key("MAX_ACCELERATION").SetValue(data.AxisZ.MaxAcceleration)
			iniFile.Section("AXIS_Z").Key("MIN_LIMIT").SetValue(data.AxisZ.MinLimit)
			iniFile.Section("AXIS_Z").Key("MAX_LIMIT").SetValue(data.AxisZ.MaxLimit)
			iniFile.Section("Joint2").Key("TYPE").SetValue(data.Joint2.Type)
			iniFile.Section("Joint2").Key("HOME").SetValue(data.Joint2.Home)
			iniFile.Section("Joint2").Key("MAX_VELOCITY").SetValue(data.Joint2.MaxVelocity)
			iniFile.Section("Joint2").Key("MAX_ACCELERATION").SetValue(data.Joint2.MaxAcceleration)
			iniFile.Section("Joint2").Key("STEPGEN_MAXACCEL").SetValue(data.Joint2.StepgenMaxaccel)
			iniFile.Section("Joint2").Key("SCALE").SetValue(data.Joint2.Scale)
			iniFile.Section("Joint2").Key("FERROR").SetValue(data.Joint2.Ferror)
			iniFile.Section("Joint2").Key("MIN_LIMIT").SetValue(data.Joint2.MinLimit)
			iniFile.Section("Joint2").Key("MAX_LIMIT").SetValue(data.Joint2.MaxLimit)
			iniFile.Section("Joint2").Key("HOME_OFFSET").SetValue(data.Joint2.HomeOffset)
			iniFile.Section("Joint2").Key("HOME_SEARCH_VEL").SetValue(data.Joint2.HomeSearchVel)
			iniFile.Section("Joint2").Key("HOME_LATCH_VEL").SetValue(data.Joint2.HomeLarchVel)
			iniFile.Section("Joint2").Key("HOME_FINAL_VEL").SetValue(data.Joint2.HomeFinalVel)
			iniFile.Section("Joint2").Key("VOLATILE_HOME").SetValue(data.Joint2.VolatileHome)
			iniFile.Section("Joint2").Key("HOME_IGNORE_LIMITS").SetValue(data.Joint2.HomeIgnoreLimits)
			iniFile.Section("Joint2").Key("HOME_USE_INDEX").SetValue(data.Joint2.HomeUseIndex)
			iniFile.Section("Joint2").Key("HOME_SEQUENCE").SetValue(data.Joint2.HomeSequence)
			iniFile.Section("Joint2").Key("BACKLASH").SetValue(data.Joint2.Backlash)
			iniFile.Section("AXIS_A").Key("MAX_VELOCITY").SetValue(data.AxisA.MaxVelocity)
			iniFile.Section("AXIS_A").Key("MAX_ACCELERATION").SetValue(data.AxisA.MaxAcceleration)
			iniFile.Section("AXIS_A").Key("MIN_LIMIT").SetValue(data.AxisA.MinLimit)
			iniFile.Section("AXIS_A").Key("MAX_LIMIT").SetValue(data.AxisA.MaxLimit)
			iniFile.Section("Joint3").Key("TYPE").SetValue(data.Joint3.Type)
			iniFile.Section("Joint3").Key("HOME").SetValue(data.Joint3.Home)
			iniFile.Section("Joint3").Key("MAX_VELOCITY").SetValue(data.Joint3.MaxVelocity)
			iniFile.Section("Joint3").Key("MAX_ACCELERATION").SetValue(data.Joint3.MaxAcceleration)
			iniFile.Section("Joint3").Key("STEPGEN_MAXACCEL").SetValue(data.Joint3.StepgenMaxaccel)
			iniFile.Section("Joint3").Key("SCALE").SetValue(data.Joint3.Scale)
			iniFile.Section("Joint3").Key("FERROR").SetValue(data.Joint3.Ferror)
			iniFile.Section("Joint3").Key("MIN_LIMIT").SetValue(data.Joint3.MinLimit)
			iniFile.Section("Joint3").Key("MAX_LIMIT").SetValue(data.Joint3.MaxLimit)
			iniFile.Section("Joint3").Key("HOME_OFFSET").SetValue(data.Joint3.HomeOffset)
			iniFile.Section("Joint3").Key("HOME_SEARCH_VEL").SetValue(data.Joint3.HomeSearchVel)
			iniFile.Section("Joint3").Key("HOME_LATCH_VEL").SetValue(data.Joint3.HomeLarchVel)
			iniFile.Section("Joint3").Key("HOME_FINAL_VEL").SetValue(data.Joint3.HomeFinalVel)
			iniFile.Section("Joint3").Key("VOLATILE_HOME").SetValue(data.Joint3.VolatileHome)
			iniFile.Section("Joint3").Key("HOME_IGNORE_LIMITS").SetValue(data.Joint3.HomeIgnoreLimits)
			iniFile.Section("Joint3").Key("HOME_USE_INDEX").SetValue(data.Joint3.HomeUseIndex)
			iniFile.Section("Joint3").Key("HOME_SEQUENCE").SetValue(data.Joint3.HomeSequence)
			iniFile.Section("Joint3").Key("BACKLASH").SetValue(data.Joint3.Backlash)
			iniFile.Section("AXIS_B").Key("MAX_VELOCITY").SetValue(data.AxisB.MaxVelocity)
			iniFile.Section("AXIS_B").Key("MAX_ACCELERATION").SetValue(data.AxisB.MaxAcceleration)
			iniFile.Section("AXIS_B").Key("MIN_LIMIT").SetValue(data.AxisB.MinLimit)
			iniFile.Section("AXIS_B").Key("MAX_LIMIT").SetValue(data.AxisB.MaxLimit)
			iniFile.Section("Joint4").Key("TYPE").SetValue(data.Joint4.Type)
			iniFile.Section("Joint4").Key("HOME").SetValue(data.Joint4.Home)
			iniFile.Section("Joint4").Key("MAX_VELOCITY").SetValue(data.Joint4.MaxVelocity)
			iniFile.Section("Joint4").Key("MAX_ACCELERATION").SetValue(data.Joint4.MaxAcceleration)
			iniFile.Section("Joint4").Key("STEPGEN_MAXACCEL").SetValue(data.Joint4.StepgenMaxaccel)
			iniFile.Section("Joint4").Key("SCALE").SetValue(data.Joint4.Scale)
			iniFile.Section("Joint4").Key("FERROR").SetValue(data.Joint4.Ferror)
			iniFile.Section("Joint4").Key("MIN_LIMIT").SetValue(data.Joint4.MinLimit)
			iniFile.Section("Joint4").Key("MAX_LIMIT").SetValue(data.Joint4.MaxLimit)
			iniFile.Section("Joint4").Key("HOME_OFFSET").SetValue(data.Joint4.HomeOffset)
			iniFile.Section("Joint4").Key("HOME_SEARCH_VEL").SetValue(data.Joint4.HomeSearchVel)
			iniFile.Section("Joint4").Key("HOME_LATCH_VEL").SetValue(data.Joint4.HomeLarchVel)
			iniFile.Section("Joint4").Key("HOME_FINAL_VEL").SetValue(data.Joint4.HomeFinalVel)
			iniFile.Section("Joint4").Key("VOLATILE_HOME").SetValue(data.Joint4.VolatileHome)
			iniFile.Section("Joint4").Key("HOME_IGNORE_LIMITS").SetValue(data.Joint4.HomeIgnoreLimits)
			iniFile.Section("Joint4").Key("HOME_USE_INDEX").SetValue(data.Joint4.HomeUseIndex)
			iniFile.Section("Joint4").Key("HOME_SEQUENCE").SetValue(data.Joint4.HomeSequence)
			iniFile.Section("Joint4").Key("BACKLASH").SetValue(data.Joint4.Backlash)
			iniFile.Section("AXIS_C").Key("MAX_VELOCITY").SetValue(data.AxisC.MaxVelocity)
			iniFile.Section("AXIS_C").Key("MAX_ACCELERATION").SetValue(data.AxisC.MaxAcceleration)
			iniFile.Section("AXIS_C").Key("MIN_LIMIT").SetValue(data.AxisC.MinLimit)
			iniFile.Section("AXIS_C").Key("MAX_LIMIT").SetValue(data.AxisC.MaxLimit)
			iniFile.Section("Joint5").Key("TYPE").SetValue(data.Joint5.Type)
			iniFile.Section("Joint5").Key("HOME").SetValue(data.Joint5.Home)
			iniFile.Section("Joint5").Key("MAX_VELOCITY").SetValue(data.Joint5.MaxVelocity)
			iniFile.Section("Joint5").Key("MAX_ACCELERATION").SetValue(data.Joint5.MaxAcceleration)
			iniFile.Section("Joint5").Key("STEPGEN_MAXACCEL").SetValue(data.Joint5.StepgenMaxaccel)
			iniFile.Section("Joint5").Key("SCALE").SetValue(data.Joint5.Scale)
			iniFile.Section("Joint5").Key("FERROR").SetValue(data.Joint5.Ferror)
			iniFile.Section("Joint5").Key("MIN_LIMIT").SetValue(data.Joint5.MinLimit)
			iniFile.Section("Joint5").Key("MAX_LIMIT").SetValue(data.Joint5.MaxLimit)
			iniFile.Section("Joint5").Key("HOME_OFFSET").SetValue(data.Joint5.HomeOffset)
			iniFile.Section("Joint5").Key("HOME_SEARCH_VEL").SetValue(data.Joint5.HomeSearchVel)
			iniFile.Section("Joint5").Key("HOME_LATCH_VEL").SetValue(data.Joint5.HomeLarchVel)
			iniFile.Section("Joint5").Key("HOME_FINAL_VEL").SetValue(data.Joint5.HomeFinalVel)
			iniFile.Section("Joint5").Key("VOLATILE_HOME").SetValue(data.Joint5.VolatileHome)
			iniFile.Section("Joint5").Key("HOME_IGNORE_LIMITS").SetValue(data.Joint5.HomeIgnoreLimits)
			iniFile.Section("Joint5").Key("HOME_USE_INDEX").SetValue(data.Joint5.HomeUseIndex)
			iniFile.Section("Joint5").Key("HOME_SEQUENCE").SetValue(data.Joint5.HomeSequence)
			iniFile.Section("Joint5").Key("BACKLASH").SetValue(data.Joint5.Backlash)
			err = IniUtils.SaveTo(iniFile, machine.Path+path+"/machine.ini")
			if err == nil {
				status = true
			}
		}
	}
	return status
}

func (machine *Machine) GetLaunch(path string) string {
	content := ""
	exists, _ := FileUtils.PathExists(machine.Path + path + "/launch/launch.py")
	if exists {
		contentByte, err := FileUtils.ReadFile(machine.Path + path + "/launch/launch.py")
		if err == nil {
			content = string(contentByte)
		}
	}
	return content
}

func (machine *Machine) UpdateLaunch(path string, content string) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/launch/launch.py")
	if exists {
		write := FileUtils.WriteFile(content, machine.Path+path+"/launch/launch.py")
		if write == nil {
			status = true
		}
	}
	return status
}

func (machine *Machine) GetTable(path string) string {
	content := ""
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.tbl")
	if exists {
		contentByte, err := FileUtils.ReadFile(machine.Path + path + "/machine.tbl")
		if err == nil {
			content = string(contentByte)
		}
	}
	return content
}

func (machine *Machine) UpdateTable(path string, content string) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.tbl")
	if exists {
		write := FileUtils.WriteFile(content, machine.Path+path+"/machine.tbl")
		if write == nil {
			status = true
		}
	}
	return status
}

func (machine *Machine) GetHal(path string) string {
	content := ""
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.hal")
	if exists {
		contentByte, err := FileUtils.ReadFile(machine.Path + path + "/machine.hal")
		if err == nil {
			content = string(contentByte)
		}
	}
	return content
}

func (machine *Machine) UpdateHal(path string, content string) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.hal")
	if exists {
		write := FileUtils.WriteFile(content, machine.Path+path+"/machine.hal")
		if write == nil {
			status = true
		}
	}
	return status
}

func (machine *Machine) GetXml(path string) string {
	content := ""
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.xml")
	if exists {
		contentByte, err := FileUtils.ReadFile(machine.Path + path + "/machine.xml")
		if err == nil {
			content = string(contentByte)
		}
	}
	return content
}

func (machine *Machine) UpdateXml(path string, content string) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.xml")
	if exists {
		write := FileUtils.WriteFile(content, machine.Path+path+"/machine.xml")
		if write == nil {
			status = true
		}
	}
	return status
}

func (machine *Machine) Delete(path string) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path)
	if exists {
		err := os.RemoveAll(machine.Path + path)
		if err == nil {
			status = true
		}
	}
	return status
}
