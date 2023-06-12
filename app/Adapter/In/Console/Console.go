package Console

import (
	"fmt"
	"github.com/Enrikerf/pfm/commandExecutor/app/Application/Port/In/ManageEngine"
	"os"
	"os/signal"
)

type Console interface {
	Run()
}
type console struct {
	manageEngineUseCase ManageEngine.UseCase
}

func NewConsole(engine ManageEngine.UseCase) Console {

	return &console{manageEngineUseCase: engine}
}

func (console *console) Run() {

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	go console.exitWatchdog(channel)
	var exit bool = false
	for {
		i := console.getOption()
		switch i {
		case "1":
			console.manageEngineUseCase.Turnaround()
		case "2":
			console.manageEngineUseCase.RpmControl(60.0 * 2)
		case "3":
			console.manageEngineUseCase.StopRpmControl()
		case "4":
			console.setGas()
		case "5":
			console.manageEngineUseCase.StepResponse()
		case "6":
			exit = true
			console.manageEngineUseCase.TearDown()
		default:
			fmt.Println("- command not found")
		}
		if exit {
			break
		}
	}
}

func (console *console) getOption() string {
	var option string
	fmt.Println("commands:")
	fmt.Println("\t1) MakeLap")
	fmt.Println("\t2) up rpm control")
	fmt.Println("\t3) down rpm control")
	fmt.Println("\t4) set gas")
	fmt.Println("\t5) step response")
	fmt.Println("\t6) exit")
	fmt.Printf("choose: ")
	_, err := fmt.Scanf("%s", &option)
	if err != nil {
		fmt.Println(err)
		console.manageEngineUseCase.TearDown()
		panic(err)
	}
	return option
}

func (console *console) setGas() {
	console.manageEngineUseCase.UnBrake()
	var i int
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		fmt.Println(err)
		return
	}
	console.manageEngineUseCase.SetGas(i)
}

func (console *console) exitWatchdog(channel chan os.Signal) {
	for range channel {
		var i string
		fmt.Println("\nDo you want to exit? (y/n)")
		_, err := fmt.Scanf("%s", &i)
		if err != nil {
			fmt.Println(err)
			console.manageEngineUseCase.TearDown()
			panic(err)
		}
		if i == "y" {
			fmt.Println("tearing down engine")
			console.manageEngineUseCase.TearDown()
			fmt.Println("exit succeed")
			os.Exit(1)
		}
	}
}
