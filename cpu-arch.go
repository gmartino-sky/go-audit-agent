package main

import ( "exec";
		"io/ioutil";
		)
		
func main () {
	if cmd, e := exec.Run("/usr/bin/lscpu", nil, nil, exec.DevNull, exec.Pipe, exec.MergeWithStdout); e == nil {
        b, _ := ioutil.ReadAll(cmd.Stdout)
        println("output: " + string(b))
    }
}