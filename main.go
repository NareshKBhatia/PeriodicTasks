package main

import (
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {

	freqencies := []string{"secondly", "minutely", "hourly", "daily", "weekly", "monthly", "yearly"}

	log.Println("Periodic Task Scheduler started")

	PeriodicTasks := GetPeriodicTasks()

	for _, task := range PeriodicTasks {

		fmt.Println(task.Name, task.Start, task.Command, task.Repeat)
		rpt, _ := strconv.Atoi(task.Repeat)

		t := time.Duration(0)
		if slices.Contains(freqencies, task.Frequency) {
			switch task.Frequency {
			case "secondly":
				t = time.Second
			case "minutely":
				t = time.Minute
			case "hourly":
				t = time.Hour
			case "daily":
				t = 24 * time.Hour 
			default:
				t = time.Minute
			}
		} else {
			log.Println("frequency should be one of ", freqencies)
			continue
		}

        if task.Start != "" {
           hhmm := strings.Split(task.Start,":")
           hh,_ := strconv.Atoi(hhmm[0])
           mm,_ := strconv.Atoi(hhmm[1])
           rpt,_ := strconv.Atoi(task.Repeat)

	       go myAlarm(hh,mm, task.Command, task.WorkingDir,time.Duration( time.Duration(rpt) * t ))
        } else {
		   go myTicker(time.Duration(time.Duration(rpt)*t), task.Command, task.WorkingDir)
        }

	}

	select {}

}

func myTicker(interval time.Duration, myCmd string, myDir string) {

	ticker := time.NewTicker(interval)

	done := make(chan bool)

	go func() {

		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				log.Println("runs every ", interval)
				go exec_it(myCmd, myDir)
			}

		}
	}()

}

func myAlarm(hh int, mm int, myCmd string, myDir string, frequency time.Duration) {

	TargetDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hh, mm, 0, 0, time.Local)

	if TargetDate.Before(time.Now()) {
		TargetDate = TargetDate.Add(24 * time.Hour)
	}
	durationToTarget := time.Until(TargetDate)
	log.Println("TargetDate", TargetDate, "until", durationToTarget)

	time.AfterFunc(durationToTarget, func() {
		go exec_it(myCmd, myDir)
	})

	c := time.Tick(frequency)
	for _ = range c {
		TargetDate = TargetDate.Add(frequency)
		durationToTarget := time.Until(TargetDate)
		log.Println("TargetDate", TargetDate, "until", durationToTarget)
		time.AfterFunc(durationToTarget, func() {
		    go exec_it(myCmd, myDir)
		})
	}

}

func exec_it(myCmd string, myDir string) {

	Flds := strings.Fields(myCmd)

	commandFullPath, err := exec.LookPath(Flds[0])
	if err != nil {
		log.Println(err)
	}

	cmd := exec.Command(commandFullPath)

	for _, x := range Flds[1:] {
		cmd.Args = append(cmd.Args, x)
	}

	cmd.Dir = myDir

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
	}

	log.Println(cmd, string(output))

}
