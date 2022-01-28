package analyzer

import (
	"fmt"
	"os/exec"
	_ "srm/lib/ptree"
	"strconv"
	"strings"
)

func AnalyzeContainer() {

}

func GetCpuInfo() int {
	cmd := exec.Command("./cpuinfo.sh")
	output, err := cmd.Output()

	if output != nil && err == nil {
		//fmt.Println("cpuinfo: ", num)
		num, err := strconv.Atoi(strings.Split(string(output), "\n")[0])
		if err != nil {
			fmt.Println(err)
		}
		return num
	}
	return 0
}
