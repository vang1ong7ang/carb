package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	var cmd T
	val := reflect.ValueOf(&cmd)
	if len(os.Args) < 2 {
		usage(val)
		return
	}

	method := strings.Title(os.Args[1])

	fun := val.MethodByName(method)
	if fun.IsValid() == false {
		usage(val)
		return
	}

	args := reflect.New(fun.Type().In(0)).Elem()
	set := flag.NewFlagSet(method, flag.ExitOnError)

	for i := 0; i < args.NumField(); i++ {
		switch args.Field(i).Kind() {
		case reflect.String:
			set.StringVar(
				args.Field(i).Addr().Interface().(*string),
				strings.ToLower(args.Type().Field(i).Name),
				args.Type().Field(i).Tag.Get("default"),
				args.Type().Field(i).Tag.Get("help"),
			)
		case reflect.Int:
			set.IntVar(
				args.Field(i).Addr().Interface().(*int),
				strings.ToLower(args.Type().Field(i).Name),
				func() int { ret, _ := strconv.Atoi(args.Type().Field(i).Tag.Get("default")); return ret }(),
				args.Type().Field(i).Tag.Get("help"),
			)
		}
	}

	set.Parse(os.Args[2:])

	fun.Call([]reflect.Value{args})

}

func usage(val reflect.Value) {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "%s <COMMAND> [-k=v]...\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "to see help message for each command, try:\n")
	fmt.Fprintf(os.Stderr, "%s <COMMAND> -help\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "availiable commmands are:\n")
	for i := 0; i < val.NumMethod(); i++ {
		fmt.Fprintf(os.Stderr, "\t%s\n", strings.ToLower(val.Type().Method(i).Name))
	}
}
