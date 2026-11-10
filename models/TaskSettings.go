package m

type TaskSettings struct {
	Common TaskCommonSettings
	Custom TaskCustomSettings
}

func NewTaskSettingsFromFile(filePath string) (ts *TaskSettings, err error) {
	return ReadObjectFromDisk[TaskSettings](filePath)
}
