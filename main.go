package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
	"time"

	wapi "github.com/iamacarpet/go-win64api"
	. "github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

//Users crea una matriz de string para listar los usuarios
var Users []string

//GetUserListWindows lista los usuarios locales para windows
func GetUserListWindows() {
	users, err := wapi.ListLocalUsers()
	if err != nil {
		fmt.Printf("Error fetching user list, %s.\r\n", err.Error())
		return
	}

	for _, u := range users {
		fmt.Printf("%s (%s)\r\n", u.Username, u.FullName)
		fmt.Printf("\tIs Enabled:                   %t\r\n", u.IsEnabled)
		fmt.Printf("\tIs Locked:                    %t\r\n", u.IsLocked)
		fmt.Printf("\tIs Admin:                     %t\r\n", u.IsAdmin)
		fmt.Printf("\tPassword Never Expires:       %t\r\n", u.PasswordNeverExpires)
		fmt.Printf("\tUser can't change password:   %t\r\n", u.NoChangePassword)
		fmt.Printf("\tPassword Age:                 %.0f days\r\n", (u.PasswordAge.Hours() / 24))
		fmt.Printf("\tLast Logon Time:              %s\r\n", u.LastLogon.Format(time.RFC850))
		fmt.Printf("\tBad Password Count:           %d\r\n", u.BadPasswordCount)
		fmt.Printf("\tNumber Of Logons:             %d\r\n", u.NumberOfLogons)
	}
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

		// skip todas las lineas que empiezen con #
		if equal := strings.Index(line, "#"); equal < 0 {
			// Obtengo el username y la descripcion
			lineSlice := strings.FieldsFunc(line, func(divide rune) bool {
				return divide == ':' // divido por :
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

	// Ahora que tenemos la lista de usuarios iteramos cada uno para imprimir Homedir,GroupID, description,etc
	for _, name := range Users {

		usr, err := user.Lookup(name)
		if err != nil {
			panic(err)
		}

		// ver https://golang.org/pkg/os/user/#User
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

		GetUserListWindows()
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
