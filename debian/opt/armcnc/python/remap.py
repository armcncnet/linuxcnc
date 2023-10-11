#!/usr/bin/env python
# -*- coding: utf-8 -*-
# Mar 2017 Code by Michel Trahan, Bob Bevins and Sylvain Deschene
# Remap of M6 in Pure Python, remap.py includes change_epilog
# Machine is Biesse 346 1995, with 3 position rack toolchanger
#-----------------------------------------------------------------------------------

#from stdglue import *
import sys
import time
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
    #初始化换刀相关信号
    disableChangeSig(self);
    try:

        #获取换刀类型
        self.toolChangeSet = toolChange.getToolChangeSet()
        #获取对刀台信息
        self.ensorData = toolChange.getToolensorPosition()

        if not self.toolChangeSet or not self.ensorData:
            self.set_errormsg("获取刀库信息失败")
            disableChangeSig(self);
            yield INTERP_ERROR

        toolNotDidFirst = 0
        if self.params[4990] == 1:
            toolNotDidFirst = 1

        self.execute("M70", lineno()) #保存模态
        self.execute("M9")  #关闭冷却液
        self.execute("M5")  #关主轴
        self.execute("G21") #设置单位
        self.execute("G30.1")  #保存当前位置到in #5181-#5183
        self.execute("G49")  #关闭补偿
        self.execute("G94") #设置速度单位units/min
        #self.execute("G40") #关闭刀具补偿

        newX = '{:.3f}'.format(self.params[5181])
        newY = '{:.3f}'.format(self.params[5182])
        newZ = '{:.3f}'.format(self.params[5183])
        
        #移动到安全高度
        gcode = "G90 G53 G1 F{} Z{}".format(self.toolChangeSet['fastfeed'],self.toolChangeSet['safez'])
        self.execute(gcode)

        cnc_stat.poll()
        #手动换刀
        if self.toolChangeSet['toolchangetype'] == '1':
            
            #移动到换刀位置
            gcode = "G53 G0 X{} Y{} Z{}".format(self.toolChangeSet['x'], self.toolChangeSet['y'], self.toolChangeSet['z'])
            self.execute(gcode)
            
            #通知用户换刀
            #数字输出信号1:输出换刀标记
            cncTools.set_change_sig(1)
            
            #数字输入信号0:等待输入，
            gcode = "M66 P0 L3 Q10000000"

            self.execute(gcode)

            #emccanon.WAIT(0,1,3,10.0)
            #int WAIT(int index, /* index of the motion exported input */
            #    int input_type, /*DIGITAL_INPUT or ANALOG_INPUT */
            #    int wait_type,  /* 0 - immediate, 1 - rise, 2 - fall, 3 - be high, 4 - be low */
            #    double timeout) /* time to wait [in seconds], if the input didn't change the value -1 is returned */
            yield INTERP_EXECUTE_FINISH
            #print("last5399--:{}".format(self.params[5399]))

            if self.params[5399] == -1.0:
                self.set_errormsg("换刀时发生错误")
                disableChangeSig(self);
                yield INTERP_ERROR
            
        #自动换刀
        elif self.toolChangeSet['toolchangetype'] =='2':
            cnc_stat.poll()
            spindleToolNum = cnc_stat.tool_in_spindle; #获取主轴中的刀具号
            #绝对距离模式
            gcode = "G90"
            self.execute(gcode)

            #主轴中没刀
            if spindleToolNum == 0 and (not toolNotDidFirst):
                
                #数字输出信号1:输出换刀标记
                cncTools.set_notool_sig(1)
                cncTools.set_tool_sig(0)
                #数字输入信号0:等待输入，
                gcode = "M66 P2 L3 Q100000000"
                self.execute(gcode)

                yield INTERP_EXECUTE_FINISH
                spindleToolNum = cncTools.get_spindle_tool() 

            if spindleToolNum == "":

                self.set_errormsg("换刀时发生错误,spindleToolNum=null")
                
                disableChangeSig(self);
                yield INTERP_ERROR
            sameTool = False
            if self.selected_tool == spindleToolNum:
                sameTool = True

            #收起护罩
            gcode = "M64 P7"
            self.execute(gcode)
            gcode = "G4 P{}".format(self.toolChangeSet['waittime'])
            self.execute(gcode)

            #主轴中有刀且和需要换的刀号不同
            if spindleToolNum > 0 and not sameTool:

                #获取刀具位置
                posData = toolChange.getToolPosition(spindleToolNum)
                
                if not posData:
                    self.set_errormsg("换刀时发生错误,刀库不存在")
                    disableChangeSig(self);
                    
                    yield INTERP_ERROR

                #绝对距离模式
                gcode = "G90"
                self.execute(gcode)
                
                #移动到主轴中刀具的刀库号
                if self.toolChangeSet['toolchangealign'] =='2':
                    putToolOffset = float(posData[1]) - float(self.toolChangeSet['puttooloffset'])

                gcode = "G53 G0 X{} Y{} F{}".format(posData[0],putToolOffset,self.toolChangeSet['fastfeed'])
                self.execute(gcode)

                #gcode = "G53 G1 Z{} F{}".format(self.toolChangeSet['pocketsafedeep'], self.toolChangeSet['fastfeed'])
                #self.execute(gcode)

                

                gcode = "G53 G0 Z{} F{}".format(float(posData[2]), self.toolChangeSet['fastfeed'])
                self.execute(gcode)

                gcode = "G53 G1 Y{} F{}".format(float(posData[1]), self.toolChangeSet['slowfeed'])
                self.execute(gcode)

                #打开拉爪
                gcode = "M64 P6"
                self.execute(gcode)
                gcode = "G4 P{}".format(self.toolChangeSet['waittime'])
                self.execute(gcode)
                #提刀
                #upToolDeep = float(posData[2]) + float(self.toolChangeSet['uptooldeep'])
                gcode = "G91 G1 Z{} F{}".format(self.toolChangeSet['uptooldeep'], self.toolChangeSet['slowfeed'])
                self.execute(gcode)

                #收起拉爪
                gcode = "M65 P6"
                self.execute(gcode)
                gcode = "G4 P{}".format(self.toolChangeSet['waittime'])
                self.execute(gcode)
                
                #绝对距离模式
                gcode = "G90"
                self.execute(gcode)
                #抬主轴至安全高度
                gcode = "G53 G1 F{} Z{}".format(self.toolChangeSet['fastfeed'],self.toolChangeSet['safez'])
                #print(gcode)
                self.execute(gcode)


            #移动到新刀具位置 
            posData = toolChange.getToolPosition(self.selected_tool)
            if not sameTool:
                if not posData:
                    self.set_errormsg("换刀时发生错误,目标刀库不存在")
                    disableChangeSig(self);
                    
                    yield INTERP_ERROR
                #移动到主轴中刀具的刀库号
                gcode = "G53 G0 X{} Y{} F{}".format(posData[0],posData[1] , self.toolChangeSet['fastfeed'])
                self.execute(gcode)

                gcode = "G53 G1 Z{} F{}".format(self.toolChangeSet['pocketsafedeep'], self.toolChangeSet['fastfeed'])
                self.execute(gcode)

                #打开拉爪
                gcode = "M64 P6"
                self.execute(gcode)
                gcode = "G4 P{}".format(self.toolChangeSet['waittime'])
                self.execute(gcode)

                gcode = "G53 G1 Z{} F{}".format(float(posData[2]), self.toolChangeSet['slowfeed'])
                self.execute(gcode)

                #收起拉爪
                gcode = "M65 P6"
                self.execute(gcode)
                gcode = "G4 P{}".format(self.toolChangeSet['waittime'])
                self.execute(gcode)

                

                if self.toolChangeSet['toolchangealign'] =='2':
                    putToolOffset = float(posData[1]) - float(self.toolChangeSet['puttooloffset'])

                gcode = "G53 G1 Y{} F{}".format(putToolOffset, self.toolChangeSet['slowfeed'])
                self.execute(gcode)

                #绝对距离模式
                gcode = "G90"
                self.execute(gcode)
                #抬主轴至安全高度
                gcode = "G53 G1 F{} Z{}".format(self.toolChangeSet['fastfeed'],self.toolChangeSet['safez'])
                self.execute(gcode)

                
            
        #所有换刀信号恢复初始值
        disableChangeSig(self); 

        #移动到对刀位置
        gcode = "G49"
        self.execute(gcode)
        gcode = "G90"
        self.execute(gcode)


        if self.toolChangeSet['toolchangetype'] =='2' and self.toolChangeSet['isdodge'] =='1':
            
            safeXY = float(self.toolChangeSet['safexy'])
            #开启避让模式
            posData[0] = float(posData[0])
            posData[1] = float(posData[1])
            #刀库位置：
            #0：平行于X轴，Y负极
            #2：平行于X轴，Y正极
            #1：平行于Y轴，X正极
            #3：平行于Y轴，X负极
            #绕开已有刀具位置，回到对刀器位置
            if self.toolChangeSet['toolchangealign'] == "0":
                safeXYLine = posData[1]+safeXY
                #向Y方向移动安全距离
                gcode = "G53 G1 X{} Y{} Z{} F{}".format(posData[0], safeXYLine, self.toolChangeSet['safez'], self.toolChangeSet['fastfeed'])
                self.execute(gcode)

                #移动到对刀器X方向位置
                gcode = "G53 G1 X{} F{}".format(self.ensorData['x'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)
                #移动到对刀器Y方向位置
                gcode = "G53 G1 Y{} F{}".format(self.ensorData['y'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)

            elif self.toolChangeSet['toolchangealign'] == "2":
                safeXYLine = posData[1]-safeXY
                gcode = "G53 G1 X{} Y{} Z{} F{}".format(posData[0], safeXYLine, self.toolChangeSet['safez'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)

                gcode = "G53 G1 X{} F{}".format(self.ensorData['x'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)

                gcode = "G53 G1 Y{} F{}".format(self.ensorData['y'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)

            elif self.toolChangeSet['toolchangealign'] == "3":
                safeXYLine = posData[0] + safeXY
                gcode = "G53 G1 X{} Y{} Z{} F{}".format(safeXYLine,posData[1], self.toolChangeSet['safez'],self.toolChangeSet['fastfeed'])

                self.execute(gcode)

                gcode = "G53 G1 Y{} F{}".format(self.ensorData['y'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)

                gcode = "G53 G1 X{} F{} ".format(self.ensorData['x'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)

            elif self.toolChangeSet['toolchangealign'] == "1":
                safeXYLine = posData[0]-safeXY
                gcode = "G53 G1 X{} Y{} Z{} F{}".format(safeXYLine,posData[1],self.toolChangeSet['safez'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)
                
                gcode = "G53 G1 Y{} F{}".format(self.ensorData['y'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)

                gcode = "G53 G1 X{} F{}".format(self.ensorData['x'], self.toolChangeSet['fastfeed'])

                self.execute(gcode)
        else:

            gcode = "G53 G1 X{} Y{} F{}".format(self.ensorData['x'], self.ensorData['y'], self.toolChangeSet['fastfeed'])
            self.execute(gcode)

        #移动到对刀器安全高度
        gcode = "G53 G1 Z{} F{} ".format(self.ensorData['z'], self.toolChangeSet['fastfeed'])

        self.execute(gcode)
        #增量距离模式
        gcode = "G91"
        self.execute(gcode)
        #第一次对刀
        gcode = "G38.2 Z{} F{}".format(self.ensorData['maxprobe'], self.ensorData['searchfeed'])

        self.execute(gcode)
        yield INTERP_EXECUTE_FINISH

        if self.params[5070] == 0:
            self.execute("G90")
            self.set_errormsg("换刀时发生错误g38.3-1")
            yield INTERP_ERROR

        #抬刀
        gcode = "G0 Z{}".format(self.ensorData['latchdist'])

        self.execute(gcode)

        #M61设置当前刀具号
        gcode = "m61 q{}".format(self.selected_tool)
        self.execute(gcode)
        

        gcode = "G38.2 Z{} F{}".format((0-float(self.ensorData['latchdist'])*2), self.ensorData['probefeed'])

        self.execute(gcode)
        #print gcode
        yield INTERP_EXECUTE_FINISH
        if self.params[5070] == 0:
            self.execute("G90")
            self.set_errormsg("换刀时发生错误g38.3-2")
            yield INTERP_ERROR

        self.params[4990] = 1

        #抬刀
        gcode = "G0 Z{}".format(self.ensorData['latchdist'])

        self.execute(gcode)

        #绝对距离模式
        gcode = "G90"

        self.execute(gcode)
        #抬主轴至安全高度
        gcode = "G53 G1 F{} Z{}".format(self.toolChangeSet['fastfeed'],self.toolChangeSet['safez'])
        self.execute(gcode)

        #开启避让
        if self.toolChangeSet['toolchangetype'] =='2' and self.toolChangeSet['isdodge'] =='1':
            if self.toolChangeSet['toolchangealign'] == "0":
                safeXYLine = float(posData[1]) + float(safeXY)
                gcode = "G53 G1 Y{}".format(safeXYLine)

                self.execute(gcode)

            elif self.toolChangeSet['toolchangealign'] == "1":
                safeXYLine = float(posData[0]) - float(safeXY)
                gcode = "G53 G1 X{}".format(safeXYLine)

                self.execute(gcode)

            elif self.toolChangeSet['toolchangealign'] == "2":
                safeXYLine = float(posData[1]) - float(safeXY)
                gcode = "G53 G1 Y{}".format(safeXYLine)

                self.execute(gcode)

            elif self.toolChangeSet['toolchangealign'] == "3":
                safeXYLine = float(posData[0]) + float(safeXY)
                gcode = "G53 G1 X{}".format(safeXYLine)

                self.execute(gcode)

        gcode = "G4 P0.5"
        self.execute(gcode)

        #收起护罩
        gcode = "M65 P7"
        self.execute(gcode)
        gcode = "G4 P{}".format(self.toolChangeSet['waittime'])
        self.execute(gcode)
        
        #绝对距离模式
        gcode = "G90"
        self.execute(gcode)

        cnc_stat.poll()
        probed_position_z = cnc_stat.probed_position[2]
        #5轴对刀器使用下面这3行
        #toolOffset =  abs(float(probed_position_z))-abs(float(self.ensorData['height']))
        #print(self.params[4991])
        toolOffset = 0
        if toolNotDidFirst :
            toolOffset = probed_position_z - self.params[1000] + self.params[4991]
            gcode = "G43.1 Z{}".format(toolOffset)

            self.execute(gcode)

        self.params[4991] = toolOffset

        newX = '{:.3f}'.format(self.params[5181])
        newY = '{:.3f}'.format(self.params[5182])
        newZ = '{:.3f}'.format(self.params[5183])

        #移动到原始位置
        gcode = "G53 G0 X{} Y{} Z{}".format(newX, newY, self.toolChangeSet['safez'])

        self.execute(gcode)

        #绝对距离模式
        gcode = "G90"
        self.execute(gcode)

        
        
        #if toolNotDidFirst :
        #    newZ = float(newZ) - float(self.params[1000]) + float(self.params[5063])

        self.params[1000] = probed_position_z

        self.execute("M72")

        if toolNotDidFirst :
            gcode = "G43.1 Z{}".format(toolOffset)

            self.execute(gcode)

    except InterpreterException,e:
        disableChangeSig(self); 
        msg = "换刀时发生错误end"
        self.set_errormsg(msg)
        yield INTERP_ERROR
    yield INTERP_OK

def disableChangeSig(self):
    print("----disableChangeSig-----")
    cncTools.set_change_sig(0)
    cncTools.set_changed_sig(0)
    cncTools.set_notool_sig(0)
    cncTools.set_tool_sig(0)




