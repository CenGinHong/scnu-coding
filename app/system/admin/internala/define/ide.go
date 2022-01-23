package define

type ListContainerResp struct {
	ContainerId string `json:"containerId"`
	UserId      int    `json:"userId"`
	UserDetail  *struct {
		UserId   int    `json:"-"`
		Username string `json:"username"`
	} `json:"userDetail"`
	LabId     int `json:"labId"`
	LabDetail *struct {
		LabId    int    `json:"labId"`
		Title    string `json:"title"`
		CourseId int    `json:"courseId"`
	} `json:"labDetail"`
	State       string `json:"state"`
	Status      string `json:"status"`
	Memory      int64  `json:"memory"`
	MemoryLimit int64  `json:"memoryLimit"`
}
type ServerInfo struct {
	ID                string `json:"ID"`
	Name              string `json:"name"`
	OperatingSystem   string `json:"operatingSystem"`
	KernelVersion     string `json:"kernelVersion"`
	Architecture      string `json:"architecture"`
	ContainersRunning int    `json:"containersRunning"`
	ContainersPaused  int    `json:"containersPaused"`
	ContainersStopped int    `json:"containersStopped"`
	Images            int    `json:"images"`
	NCPU              string `json:"NCPU"`
	MemTotal          string `json:"memTotal"`
}
