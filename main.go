package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	. "github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

//Users crea una matriz de string para listar los usuarios
var Users []string

// PowerShell struct
type PowerShell struct {
	powerShell string
}

//New Crea una nueva sesion de ps
func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

// ejecuta el comando y el argunmento
func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

//GetUserListLinux Trae una lista de usuarios compatible con sistemas unix/linux
func GetUserListLinux() {
	file, err := os.Open("/etc/passwd")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		// skip all line starting with #
		if equal := strings.Index(line, "#"); equal < 0 {
			// get the username and description
			lineSlice := strings.FieldsFunc(line, func(divide rune) bool {
				return divide == ':' // we divide at colon
			})

			if len(lineSlice) > 0 {
				Users = append(Users, lineSlice[0])
			}

		}

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	// now we have a list of users
	// iterate(cycle) each of them to
	// print out HomeDir, GroupID, description, etc

	for _, name := range Users {

		usr, err := user.Lookup(name)
		if err != nil {
			panic(err)
		}

		// see https://golang.org/pkg/os/user/#User
		fmt.Printf("username:%s\n", usr.Username)
		fmt.Printf("homedir:%s\n", usr.HomeDir)
		fmt.Printf("groupID:%s\n", usr.Gid)
		fmt.Printf("DisplayName:%s\n", usr.Name)
		fmt.Println("*********************************")

	}
}

//ProcessList Trae una lista de procesos corriendo en el sistema, Usuarios y Version del SO
func ProcessList() {
	infoStat, _ := host.Info()
	fmt.Printf("Total processes: %d\n", infoStat.Procs)

	miscStat, _ := load.Misc()
	fmt.Printf("Running processes: %v\n", miscStat.ProcsRunning)
	fmt.Printf("Sistema Operativo: %v\n", infoStat.OS)
	fmt.Printf("Versión de SO: %v\n", infoStat.PlatformVersion)

	// Valido SO y lista usuarios

	if infoStat.OS == "windows" {
		New()
	} else {
		GetUserListLinux()
	}
}

//PrintCPUData Imprime datos del CPU
func PrintCPUData() {
	// Muentra informacion Basica del CPU:
	fmt.Println("Name:", CPU.BrandName)
	fmt.Println("PhysicalCores:", CPU.PhysicalCores)
	fmt.Println("ThreadsPerCore:", CPU.ThreadsPerCore)
	fmt.Println("LogicalCores:", CPU.LogicalCores)
	fmt.Println("Family", CPU.Family, "Model:", CPU.Model, "Vendor ID:", CPU.VendorID)
	fmt.Println("Features:", fmt.Sprintf(strings.Join(CPU.FeatureSet(), ",")))
	fmt.Println("Cacheline bytes:", CPU.CacheLine)
	fmt.Println("L1 Data Cache:", CPU.Cache.L1D, "bytes")
	fmt.Println("L1 Instruction Cache:", CPU.Cache.L1D, "bytes")
	fmt.Println("L2 Cache:", CPU.Cache.L2, "bytes")
	fmt.Println("L3 Cache:", CPU.Cache.L3, "bytes")
	fmt.Println("Frequency", CPU.Hz, "hz")

	// Muentra si tenemos instrucciones de CPU especificas:
	if CPU.Supports(SSE, SSE2) {
		fmt.Println("We have Streaming SIMD 2 Extensions")
	}
}
func main() {
	PrintCPUData()
	ProcessList()
}
