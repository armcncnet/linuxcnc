#!/usr/bin/env python
# -*- coding: utf-8 -*-
# Mar 2017 Code by Michel Trahan, Bob Bevins and Sylvain Deschene
# Remap of M6 in Pure Python, remap.py includes change_epilog
# Machine is Biesse 346 1995, with 3 position rack toolchanger
#-----------------------------------------------------------------------------------

#from stdglue import *
import sys
import linuxcnc
import hal
from interpreter import *
import emccanon
from util import lineno
from linuxcnc import ini
import ConfigParser
#import websocket
import subprocess
#import json
import os
sys.path.append('/var/www/html/python/')
import cncConfig
import toolChange
import cncTools

cnc_stat = linuxcnc.stat()
cnc_cmd = linuxcnc.command()


def toolChangeM6(self, **words):
    #基本信息
    if not self.task:
        yield INTERP_OK

    disableChangeSig(self);
    try:

        #self.toolSetInfo = toolChange.getToolChangeSet()
        self.toolchangetype = toolChange.getToolChangeSet()['toolchangetype']
        self.ensorData = toolChange.getToolensorPosition()
        historyData = cncTools.get_history_data()
        taskId = historyData["id"]

        self.params[4990] = int(taskId)

        isNewTask = False
        if self.params[4990] != self.params[4991]:
            isNewTask = True

        self.params[4991] = self.params[4990]


        executeGcode(self, "M70")
        executeGcode(self, "G21")

        executeGcode(self, "G49")
        

        moveZtoSafePos(self, self.ensorData['searchfeed'])

        #停止主轴
        stopSpindle(self)

        gcode = "G4 P2"
        self.execute(gcode)

        noProbeFlag = False
        if self.selected_tool == 0 or self.ensorData['flag'] == '0':
            noProbeFlag = True

        cnc_stat.poll()
        #安全位置
        maxLimit = cnc_stat.axis[2]["max_position_limit"]-5
        #手动换刀
        if self.toolchangetype == '1':
            
            
            #移动到换刀位置
            moveToToolChange(self)
            
            #通知用户换刀
            #数字输出信号1:输出换刀标记
            cncTools.set_change_sig(1)
            #cnc_cmd.set_digital_output(1, 1)
            #subprocess.Popen('halcmd setp hal_linktoolchange.changed_sig 1', stdout = subprocess.PIPE, shell = True)

            #数字输入信号0:等待输入，
            gcode = "M66 P0 L3 Q10000"
            #print gcode + "wait change tool"
            executeGcode(self, gcode)
            #print gcode
            #emccanon.WAIT(0,1,3,10.0)
            #int WAIT(int index, /* index of the motion exported input */
            #    int input_type, /*DIGITAL_INPUT or ANALOG_INPUT */
            #    int wait_type,  /* 0 - immediate, 1 - rise, 2 - fall, 3 - be high, 4 - be low */
            #    double timeout) /* time to wait [in seconds], if the input didn't change the value -1 is returned */
            yield INTERP_EXECUTE_FINISH
            #print("last5399--:{}".format(self.params[5399]))
            if self.params[5399] == -1.0:
                #msg = "5399 -1"
                #print msg
                self.set_errormsg("换刀时发生错误")
                disableChangeSig(self);
                yield INTERP_ERROR
            
        #自动换刀
        elif self.toolchangetype =='2':
            print "toolchangetype2"
            spindleToolNum = cnc_stat.tool_in_spindle;
            print spindleToolNum;
            #主轴中没刀
            if spindleToolNum == 0:
                print "no tool in pocket"
                #数字输出信号1:输出换刀标记
                cncTools.set_notool_sig(1)
                cncTools.set_tool_sig(0)
                #数字输入信号0:等待输入，
                gcode = "M66 P2 L3 Q10000"
                #print gcode + "wait change tool"
                executeGcode(self, gcode)
                yield INTERP_EXECUTE_FINISH
                spindleToolNum = cncTools.get_spindle_tool() 

            if spindleToolNum == "":
                #self.set_errormsg("换刀时发生错误")
                disableChangeSig(self);
                yield INTERP_ERROR

            #主轴中有刀
            if spindleToolNum > 0:
                """
                if self.selected_tool == spindleToolNum:
                    self.set_errormsg("换刀时发生错误,相同刀具")
                    disableChangeSig(self);
                    yield INTERP_ERROR
                """

                #获取刀具位置
                posData = toolChange.getToolPosition(spindleToolNum)
                print posData
                gcode = "G53 G0 X{} Y{}".format(posData[0],posData[1])
                self.execute(gcode)
                gcode = "G53 G1 Z{} F200".format(float(posData[2])+1)
                print gcode
                self.execute(gcode)
                #打开拉爪
                gcode = "G91 G1 W{} F500".format(-15)
                print gcode
                self.execute(gcode)
                gcode = "G4 P1"
                self.execute(gcode)

                gcode = "G1 Z6 F200"
                self.execute(gcode)

                #绝对距离模式
                gcode = "G90"
                self.execute(gcode)
                #抬刀
                gcode = "G53 G0 Z{}".format(maxLimit)
                self.execute(gcode)

            #移动到新刀具位置 
            posData = toolChange.getToolPosition(self.selected_tool)
            #print posData
            gcode = "G53 G0 X{} Y{}".format(posData[0],posData[1])
            self.execute(gcode)
            gcode = "G53 G1 Z{} F200".format(float(posData[2]))
            print gcode
            self.execute(gcode)
            #收紧拉爪
            gcode = "G91 G1 W{} F500".format(15)
            print gcode
            self.execute(gcode)
            gcode = "G4 P1"
            self.execute(gcode)

            gcode = "G1 Z6 F200"
            self.execute(gcode)

            #绝对距离模式
            gcode = "G90"
            self.execute(gcode)
            #抬刀
            gcode = "G53 G0 Z{}".format(maxLimit)
            self.execute(gcode)
            print gcode
            

        else :
            self.set_errormsg("tool_probe_m6 remap error:change type not found")
            disableChangeSig(self);
        disableChangeSig(self); 
        #如果是卸刀或没有对刀台，不需要执行以下代码
        if not noProbeFlag:
            gcode = "G4 P1"
            self.execute(gcode)

            #移动到对刀位置
            gcode = "G49"
            self.execute(gcode)
            gcode = "G53 G0 X{} Y{} Z{}".format(self.ensorData['x'], self.ensorData['y'], self.ensorData['z'])
            self.execute(gcode)
            print("#移动到对刀位置");
            print gcode

            #self.params[1000] = self.params[5063]
            self.params[1000] = cnc_stat.probed_position[2]
            
            #增量距离模式
            gcode = "G91"
            self.execute(gcode)

            gcode = "G38.3 Z{} F{}".format(self.ensorData['maxprobe'], self.ensorData['searchfeed'])
            self.execute(gcode)
            print gcode
            yield INTERP_EXECUTE_FINISH

            if self.params[5070] == 0:
                self.execute("G90")
                self.set_errormsg("换刀时发生错误")
                yield INTERP_ERROR

            #print("last5063:{}".format(self.params[1000]))

            gcode = "G1 Z{} F{}".format(self.ensorData['latchdist'], self.ensorData['searchfeed'])
            self.execute(gcode)
            #print gcode

            gcode = "G4 P0.5"
            self.execute(gcode)
            #print gcode

            gcode = "G38.3 Z{} F{}".format((0-float(self.ensorData['latchdist'])*2), self.ensorData['probefeed'])
            self.execute(gcode)
            #print gcode
            yield INTERP_EXECUTE_FINISH
            if self.params[5070] == 0:
                self.execute("G90")
                self.set_errormsg("换刀时发生错误")
                yield INTERP_ERROR

            gcode = "G0 Z2"
            self.execute(gcode)
            #print gcode

            gcode = "G0 X1.1"
            self.execute(gcode)
            #print gcode

            #print("current5063:{}".format(self.params[5063]))

            #绝对距离模式
            gcode = "G90"
            self.execute(gcode)
            gcode = "G53 G0 Z{}".format(maxLimit)
            self.execute(gcode)
            #print gcode

            gcode = "G4 P1"
            self.execute(gcode)
            #print gcode

        #if not noProbeFlag:
            cnc_stat.poll()

            #print("current5063:{}".format(self.params[5063]))
            #print("last5063:{}".format(self.params[1000]))
            #self.ensorData = toolChange.getToolensorPosition()

            #toolOffset =  abs(float(self.ensorData['height'])) + abs(float(self.ensorData['zOffset']))-abs(float(self.params[5063]))
            #toolOffset =  118-abs(float(self.ensorData['height']))-abs(float(self.params[5063]))
            
            #toolOffset =  abs(float(self.params[5063]))-abs(float(self.ensorData['height']))
            probed_position_z = cnc_stat.probed_position[2]
            #print("probed_position_z:{}".format(probed_position_z))
            #五轴换刀后对刀使用下面
            #toolOffset =  abs(float(probed_position_z))-abs(float(self.ensorData['height']))
            #3轴换刀后使用下面
            
            toolOffset =  abs(float(probed_position_z))-abs(float(self.params[1000]))
            
            #toolOffset =  abs(float(probed_position_z))
            print("toolOffset2:")
            print(cnc_stat.probed_position[2])
            print("toolOffset1:")
            print(self.params[1000])
            print("toolOffset:")
            print(toolOffset)
            
            #if isNewTask:
            #    toolOffset = 0
            #else :
            #    toolOffset = self.params[5063] - self.params[1000]
            #toolOffset = 0-toolOffset
            #动态刀具长度偏移
            #gcode = "G43.1 Z{}".format(toolOffset)
            #print gcode
            #self.execute(gcode)

        #M61设置当前刀具号
        gcode = "m61 q{}".format(self.selected_tool)
        self.execute(gcode)

        gcode = "G0 X0 Y0"
        #print gcode
        self.execute(gcode)

        #self.execute("M3")
        gcode = "G4 P1"
        self.execute(gcode)

        self.execute("M72")

        if float(self.params[1000]) >0 or float(self.params[1000])<0 :
            print(232323232323)
            gcode = "G43.1 Z{}".format(toolOffset)
            print gcode
            self.execute(gcode)

        
    except InterpreterException,e:
        disableChangeSig(self); 
        msg = "换刀时发生错误end"
        self.set_errormsg(msg)
        yield INTERP_ERROR
    yield INTERP_OK

def disableChangeSig(self):
    print(disableChangeSig)
    cncTools.set_change_sig(0)
    cncTools.set_changed_sig(0)
    cncTools.set_notool_sig(0)
    cncTools.set_tool_sig(0)

def executeGcode(self, gcode):
    self.execute(gcode, lineno())
    print(gcode)

#停止主轴
def stopSpindle(self):

    executeGcode(self, "M5")

#打开主轴
def enableSpindle(self):
    executeGcode(self, "M3")

#Z轴移动到安全位置
def moveZtoSafePos(self, speed):
    
    cnc_stat.poll()
    maxLimit = cnc_stat.axis[2]["max_position_limit"]-5
    gcode = "G53 G0 Z{:.3f} F{}".format(maxLimit, speed)
    
    executeGcode(self, gcode)

def getToolPos(self, toolNumber):
    global cnc_stat
    cnc_stat.poll()
    return cnc_stat.tool_table[toolNumber]

def moveToCurrentToolPocket(self):
    print "moveToCurrentToolPocket:"
    #self.current_tool = 5
    toolNumber = self.current_tool
    data = toolChange.getToolPosition(toolNumber)
    #print data
    #print data[0]
    gcode = "G53 G0 X{} Y{}".format(data[0],data[1])
    executeGcode(self, gcode)

    gcode = "G53 G0 Z{} F100".format(data[2])
    executeGcode(self, gcode)


def moveToSelectedToolPocket(self, toolNumber):
    print "moveToSelectedToolPocket:"
    print "selectedTool:"
    print toolNumber
    data = toolChange.getToolPosition(toolNumber)

    gcode = "G53 G0 X{} Y{}".format(data[0],data[1])
    executeGcode(self, gcode)

    gcode = "G53 G0 Z{} F100".format(data[2])
    executeGcode(self, gcode)
#移动到换刀位置
def moveToToolChange(self):
    print("moveToToolChange")
    data = toolChange.getToolChangeSet()
    #gcode = "G49"
    #executeGcode(self, gcode)

    gcode = "G53 G0 X{} Y{} Z{}".format(data['x'], data['y'], data['z'])
    executeGcode(self, gcode)
#移动到对刀位置
def moveToToolensor(self):
    print("moveToToolensor")
    data = toolChange.getToolensorPosition()

    #gcode = "G49"
    #executeGcode(self, gcode)

    gcode = "G53 G0 X{} Y{} Z{}".format(data['x'], data['y'], data['z'])
    executeGcode(self, gcode)
    
#对完刀后重新设置
def finishToolChange(self):
    #绝对距离模式
    gcode = "G90"
    self.execute(gcode)
    #print gcode
    #set tool 1 Z offset from the machine origin to 1.5)
    gcode = "G10 L1 P{} Z{}".format(self.selected_tool, self.params[5063])
    self.execute(gcode)
    #print gcode

    currentOffset = self.params[5063] - self.params[1000]
    #动态刀具长度偏移
    gcode = "G43.1 Z{}".format(currentOffset)
    self.execute(gcode)
    #print gcode
    #G53移入机器坐标
    moveZtoSafePos(self)
    #gcode = "G53 G0 Z10 F100"
    #self.execute(gcode)
    #print gcode
    #M61设置当前刀具号
    gcode = "m61 q{}".format(self.selected_tool)
    self.execute(gcode)

"""    
def moveToWorkPos(self):
    s.poll()
    x = s.actual_position[0]
    y = s.actual_position[1]
    z = s.actual_position[2]
    gcode = "G53 G0 X{:.3f} Y{:.3f}".format(x, y)
    print(gcode)
    self.execute(gcode)
    gcode = "G53 G0 Z{:.3f}".format(z)
    print(gcode)
    self.execute(gcode)

"""




