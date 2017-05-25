package main

import (
	"log"
	"syscall"
)

func sendCtrlCToPid(pid int) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		log.Fatalf("LoadDLL: %v\n", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		log.Fatalf("FindProc: %v\n", e)
	}
	r, _, e := p.Call(uintptr(syscall.CTRL_C_EVENT), uintptr(pid))
	if r == 0 {
		log.Fatalf("GenerateConsoleCtrlEvent: %v\n", e)
	}
}

func setCtrlHandlingIsDisabledForThisProcess(disabled bool) {
	disabledInt := 0
	if(disabled){
		disabledInt = 1
	}
	
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		log.Fatalf("LoadDLL: %v\n", e)
	}
	p, e := d.FindProc("SetConsoleCtrlHandler")
	if e != nil {
		log.Fatalf("FindProc: %v\n", e)
	}
	r, _, e := p.Call(uintptr(0), uintptr(disabledInt))
	if r == 0 {
		log.Fatalf("SetConsoleCtrlHandler: %v\n", e)
	}
}
