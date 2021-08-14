package service

func Init() {
	go MonitorReset()
	go monitorInput()
}
