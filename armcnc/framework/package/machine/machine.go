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
				item.Time = item.Time.Add(-525600 * time.Minute)
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
			iniFile.Section("IO").Key("ESTOP_PIN").SetValue(data.Io.EstopPin)
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
			data = machine.DefaultIni(data)
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

func (machine *Machine) DefaultIni(data INI) INI {
	if data.Emc.Machine == "" {
		data.Emc.Machine = "machine"
	}
	if data.Emc.Debug == "" {
		data.Emc.Debug = "0"
	}
	if data.Emc.Version == "" {
		data.Emc.Version = "1.1"
	}
	if data.Display.Display == "" {
		data.Display.Display = "halui"
	}
	if data.Display.CycleTime == "" {
		data.Display.CycleTime = "0.1000"
	}
	if data.Display.PositionOffset == "" {
		data.Display.PositionOffset = "RELATIVE"
	}
	if data.Display.PositionFeedback == "" {
		data.Display.PositionFeedback = "ACTUAL"
	}
	if data.Display.Arcdivision == "" {
		data.Display.Arcdivision = "64"
	}
	if data.Display.MaxFeedOverride == "" {
		data.Display.MaxFeedOverride = "1.5"
	}
	if data.Display.MinSpindleOverride == "" {
		data.Display.MinSpindleOverride = "0.3"
	}
	if data.Display.MaxSpindleOverride == "" {
		data.Display.MaxSpindleOverride = "1.5"
	}
	if data.Display.DefaultLinearVelocity == "" {
		data.Display.DefaultLinearVelocity = "30.000"
	}
	if data.Display.MinLinearVelocity == "" {
		data.Display.MinLinearVelocity = "0.000"
	}
	if data.Display.MaxLinearVelocity == "" {
		data.Display.MaxLinearVelocity = "100.000"
	}
	if data.Display.DefaultAngularVelocity == "" {
		data.Display.DefaultAngularVelocity = "13.333"
	}
	if data.Display.MinAngularVelocity == "" {
		data.Display.MinAngularVelocity = "0.000"
	}
	if data.Display.MaxAngularVelocity == "" {
		data.Display.MaxAngularVelocity = "16.667"
	}
	if data.Display.ProgramPrefix == "" {
		data.Display.ProgramPrefix = "../../files"
	}
	if data.Display.OpenFile == "" {
		data.Display.OpenFile = "../../programs/armcnc.ngc"
	}
	if data.Display.Increments == "" {
		data.Display.Increments = "10mm,5mm,1mm,.5mm,.1mm,.05mm,.01mm,.005mm"
	}
	if data.Python.PathAppend == "" {
		data.Python.PathAppend = "../../scripts/python"
	}
	if data.Python.Toplevel == "" {
		data.Python.Toplevel = "../../scripts/python/main.py"
	}
	if data.Filter.ProgramExtension == "" {
		data.Filter.ProgramExtension = ".py Python Script"
	}
	if data.Filter.Py == "" {
		data.Filter.Py = "python"
	}
	if data.Rs274ngc.Features == "" {
		data.Rs274ngc.Features = "30"
	}
	if data.Rs274ngc.SubroutinePath == "" {
		data.Rs274ngc.SubroutinePath = "../../scripts/ngc"
	}
	if data.Rs274ngc.ParameterFile == "" {
		data.Rs274ngc.ParameterFile = "machine.var"
	}
	if data.Emcmot.Emcmot == "" {
		data.Emcmot.Emcmot = "motmod"
	}
	if data.Emcmot.CommTimeout == "" {
		data.Emcmot.CommTimeout = "1.000"
	}
	if data.Emcmot.BasePeriod == "" {
		data.Emcmot.BasePeriod = "200000"
	}
	if data.Emcmot.ServoPeriod == "" {
		data.Emcmot.ServoPeriod = "1000000"
	}
	if data.Emcio.Emcio == "" {
		data.Emcio.Emcio = "io"
	}
	if data.Emcio.CycleTime == "" {
		data.Emcio.CycleTime = "0.100"
	}
	if data.Emcio.ToolTable == "" {
		data.Emcio.ToolTable = "machine.tbl"
	}
	if data.Task.Task == "" {
		data.Task.Task = "milltask"
	}
	if data.Task.CycleTime == "" {
		data.Task.CycleTime = "0.010"
	}
	if data.Hal.HalFile == "" {
		data.Hal.HalFile = "machine.hal"
	}
	if data.Traj.Spindles == "" {
		data.Traj.Spindles = "1"
	}
	if data.Traj.Coordinates == "" {
		data.Traj.Coordinates = "XYZ"
	}
	if data.Traj.LinearUnits == "" {
		data.Traj.LinearUnits = "mm"
	}
	if data.Traj.AngularUnits == "" {
		data.Traj.AngularUnits = "deg"
	}
	if data.Traj.PositionFile == "" {
		data.Traj.PositionFile = "machine.position"
	}
	if data.Spindle0.MinForwardVelocity == "" {
		data.Spindle0.MinForwardVelocity = "0"
	}
	if data.Spindle0.MaxForwardVelocity == "" {
		data.Spindle0.MaxForwardVelocity = "12000"
	}
	if data.Kins.Joints == "" {
		data.Kins.Joints = "3"
	}
	if data.Kins.Kinematics == "" {
		data.Kins.Kinematics = "trivkins coordinates=XYZ"
	}
	if data.AxisX.MaxVelocity == "" {
		data.AxisX.MaxVelocity = "0.000"
	}
	if data.AxisX.MaxAcceleration == "" {
		data.AxisX.MaxAcceleration = "0.000"
	}
	if data.AxisX.MinLimit == "" {
		data.AxisX.MinLimit = "0.000"
	}
	if data.AxisX.MaxLimit == "" {
		data.AxisX.MaxLimit = "0.000"
	}
	if data.Joint0.Type == "" {
		data.Joint0.Type = "LINEAR"
	}
	if data.Joint0.Home == "" {
		data.Joint0.Home = "0.000"
	}
	if data.Joint0.MaxVelocity == "" {
		data.Joint0.MaxVelocity = "0.000"
	}
	if data.Joint0.MaxAcceleration == "" {
		data.Joint0.MaxAcceleration = "0.000"
	}
	if data.Joint0.StepgenMaxaccel == "" {
		data.Joint0.StepgenMaxaccel = "0.000"
	}
	if data.Joint0.Scale == "" {
		data.Joint0.Scale = "1600"
	}
	if data.Joint0.Ferror == "" {
		data.Joint0.Ferror = "1.000"
	}
	if data.Joint0.MinLimit == "" {
		data.Joint0.MinLimit = "0.000"
	}
	if data.Joint0.MaxLimit == "" {
		data.Joint0.MaxLimit = "0.000"
	}
	if data.Joint0.HomeOffset == "" {
		data.Joint0.HomeOffset = "0.000"
	}
	if data.Joint0.HomeSearchVel == "" {
		data.Joint0.HomeSearchVel = "0.000"
	}
	if data.Joint0.HomeLarchVel == "" {
		data.Joint0.HomeLarchVel = "0.000"
	}
	if data.Joint0.HomeFinalVel == "" {
		data.Joint0.HomeFinalVel = "0.000"
	}
	if data.Joint0.VolatileHome == "" {
		data.Joint0.VolatileHome = "1"
	}
	if data.Joint0.HomeIgnoreLimits == "" {
		data.Joint0.HomeIgnoreLimits = "NO"
	}
	if data.Joint0.HomeUseIndex == "" {
		data.Joint0.HomeUseIndex = "NO"
	}
	if data.Joint0.HomeSequence == "" {
		data.Joint0.HomeSequence = "0"
	}
	if data.Joint0.Backlash == "" {
		data.Joint0.Backlash = "0.00"
	}
	if data.AxisY.MaxVelocity == "" {
		data.AxisY.MaxVelocity = "0.000"
	}
	if data.AxisY.MaxAcceleration == "" {
		data.AxisY.MaxAcceleration = "0.000"
	}
	if data.AxisY.MinLimit == "" {
		data.AxisY.MinLimit = "0.000"
	}
	if data.AxisY.MaxLimit == "" {
		data.AxisY.MaxLimit = "0.000"
	}
	if data.Joint1.Type == "" {
		data.Joint1.Type = "LINEAR"
	}
	if data.Joint1.Home == "" {
		data.Joint1.Home = "0.000"
	}
	if data.Joint1.MaxVelocity == "" {
		data.Joint1.MaxVelocity = "0.000"
	}
	if data.Joint1.MaxAcceleration == "" {
		data.Joint1.MaxAcceleration = "0.000"
	}
	if data.Joint1.StepgenMaxaccel == "" {
		data.Joint1.StepgenMaxaccel = "0.000"
	}
	if data.Joint1.Scale == "" {
		data.Joint1.Scale = "1600"
	}
	if data.Joint1.Ferror == "" {
		data.Joint1.Ferror = "1.000"
	}
	if data.Joint1.MinLimit == "" {
		data.Joint1.MinLimit = "0.000"
	}
	if data.Joint1.MaxLimit == "" {
		data.Joint1.MaxLimit = "0.000"
	}
	if data.Joint1.HomeOffset == "" {
		data.Joint1.HomeOffset = "0.000"
	}
	if data.Joint1.HomeSearchVel == "" {
		data.Joint1.HomeSearchVel = "0.000"
	}
	if data.Joint1.HomeLarchVel == "" {
		data.Joint1.HomeLarchVel = "0.000"
	}
	if data.Joint1.HomeFinalVel == "" {
		data.Joint1.HomeFinalVel = "0.000"
	}
	if data.Joint1.VolatileHome == "" {
		data.Joint1.VolatileHome = "1"
	}
	if data.Joint1.HomeIgnoreLimits == "" {
		data.Joint1.HomeIgnoreLimits = "NO"
	}
	if data.Joint1.HomeUseIndex == "" {
		data.Joint1.HomeUseIndex = "NO"
	}
	if data.Joint1.HomeSequence == "" {
		data.Joint1.HomeSequence = "0"
	}
	if data.Joint1.Backlash == "" {
		data.Joint1.Backlash = "0.00"
	}
	if data.AxisZ.MaxVelocity == "" {
		data.AxisZ.MaxVelocity = "0.000"
	}
	if data.AxisZ.MaxAcceleration == "" {
		data.AxisZ.MaxAcceleration = "0.000"
	}
	if data.AxisZ.MinLimit == "" {
		data.AxisZ.MinLimit = "0.000"
	}
	if data.AxisZ.MaxLimit == "" {
		data.AxisZ.MaxLimit = "0.000"
	}
	if data.Joint2.Type == "" {
		data.Joint2.Type = "LINEAR"
	}
	if data.Joint2.Home == "" {
		data.Joint2.Home = "0.000"
	}
	if data.Joint2.MaxVelocity == "" {
		data.Joint2.MaxVelocity = "0.000"
	}
	if data.Joint2.MaxAcceleration == "" {
		data.Joint2.MaxAcceleration = "0.000"
	}
	if data.Joint2.StepgenMaxaccel == "" {
		data.Joint2.StepgenMaxaccel = "0.000"
	}
	if data.Joint2.Scale == "" {
		data.Joint2.Scale = "1600"
	}
	if data.Joint2.Ferror == "" {
		data.Joint2.Ferror = "1.000"
	}
	if data.Joint2.MinLimit == "" {
		data.Joint2.MinLimit = "0.000"
	}
	if data.Joint2.MaxLimit == "" {
		data.Joint2.MaxLimit = "0.000"
	}
	if data.Joint2.HomeOffset == "" {
		data.Joint2.HomeOffset = "0.000"
	}
	if data.Joint2.HomeSearchVel == "" {
		data.Joint2.HomeSearchVel = "0.000"
	}
	if data.Joint2.HomeLarchVel == "" {
		data.Joint2.HomeLarchVel = "0.000"
	}
	if data.Joint2.HomeFinalVel == "" {
		data.Joint2.HomeFinalVel = "0.000"
	}
	if data.Joint2.VolatileHome == "" {
		data.Joint2.VolatileHome = "1"
	}
	if data.Joint2.HomeIgnoreLimits == "" {
		data.Joint2.HomeIgnoreLimits = "NO"
	}
	if data.Joint2.HomeUseIndex == "" {
		data.Joint2.HomeUseIndex = "NO"
	}
	if data.Joint2.HomeSequence == "" {
		data.Joint2.HomeSequence = "0"
	}
	if data.Joint2.Backlash == "" {
		data.Joint2.Backlash = "0.00"
	}
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
		data.Joint3.Type = "ANGULAR"
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
		data.Joint4.Type = "ANGULAR"
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
		data.Joint5.Type = "ANGULAR"
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
	return data
}

func (machine *Machine) DefaultUser(data USER) USER {
	if data.Base.Name == "" {
		data.Base.Name = "机床名称"
	}
	if data.Base.Describe == "" {
		data.Base.Describe = "机床的描述信息"
	}
	if data.HandWheel.Status == "" {
		data.HandWheel.Status = "NO"
	}
	if data.HandWheel.XVelocity == "" {
		data.HandWheel.XVelocity = "10.000"
	}
	if data.HandWheel.YVelocity == "" {
		data.HandWheel.YVelocity = "10.000"
	}
	if data.HandWheel.ZVelocity == "" {
		data.HandWheel.ZVelocity = "6.000"
	}
	if data.HandWheel.AVelocity == "" {
		data.HandWheel.AVelocity = "10.000"
	}
	if data.HandWheel.BVelocity == "" {
		data.HandWheel.BVelocity = "10.000"
	}
	if data.HandWheel.CVelocity == "" {
		data.HandWheel.CVelocity = "10.000"
	}
	if data.Tool.Method == "" {
		data.Tool.Method = "MANUAL"
	}
	if data.Tool.XPosition == "" {
		data.Tool.XPosition = "0.000"
	}
	if data.Tool.YPosition == "" {
		data.Tool.YPosition = "0.000"
	}
	if data.Tool.ZPosition == "" {
		data.Tool.ZPosition = "0.000"
	}
	if data.Tool.ZHeight == "" {
		data.Tool.ZHeight = "0.000"
	}
	if data.Tool.MaxSearchDistance == "" {
		data.Tool.MaxSearchDistance = "0.000"
	}
	if data.Tool.LatchSearchDistance == "" {
		data.Tool.LatchSearchDistance = "0.000"
	}
	if data.Tool.SearchVelocity == "" {
		data.Tool.SearchVelocity = "0.000"
	}
	if data.Tool.LatchSearchVelocity == "" {
		data.Tool.LatchSearchVelocity = "0.000"
	}
	if data.Tool.Pockets == "" {
		data.Tool.Pockets = "[]"
	}
	if data.IO.EstopPin == "" {
		data.IO.EstopPin = ""
	}
	return data
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
