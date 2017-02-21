package main

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	desc  string
	flags []*Flag
	set   *flag.FlagSet
}

type Flag struct {
	FullName  string
	ShortName string
	Default   interface{}
	Usage     string
}

func NewFlags(desc string) *Flags {
	f := &Flags{
		desc: desc,
		set:  flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
	f.set.Usage = f.Usage
	return f
}

func (fs *Flags) String(p *string, full, short, value, usage string) {
	fs.addFlag(full, short, value, usage)
	if full != "" {
		fs.set.StringVar(p, full, value, usage)
	}
	if short != "" {
		fs.set.StringVar(p, short, value, usage)
	}
}

func (fs *Flags) Int64(p *int64, full, short string, value int64, usage string) {
	fs.addFlag(full, short, value, usage)
	if full != "" {
		fs.set.Int64Var(p, full, value, usage)
	}
	if short != "" {
		fs.set.Int64Var(p, short, value, usage)
	}
}

func (fs *Flags) Bool(p *bool, full, short string, value bool, usage string) {
	fs.addFlag(full, short, value, usage)
	if full != "" {
		fs.set.BoolVar(p, full, value, usage)
	}
	if short != "" {
		fs.set.BoolVar(p, short, value, usage)
	}
}

func (fs *Flags) addFlag(full, short string, value interface{}, usage string) {
	f := &Flag{
		FullName:  full,
		ShortName: short,
		Default:   value,
		Usage:     usage,
	}
	fs.flags = append(fs.flags, f)
}

func (fs *Flags) Parse() {
	fs.set.Parse(os.Args[1:])
}

func (fs *Flags) Usage() {
	if fs.desc != "" {
		fmt.Println(fs.desc)
		fmt.Println()
	}

	fmt.Println("Options:")
	fmt.Println()

	for i, f := range fs.flags {
		if f.FullName != "" && f.ShortName != "" {
			fmt.Printf("  -%s, -%s", f.ShortName, f.FullName)
		} else if f.FullName != "" {
			fmt.Printf("  -%s", f.FullName)
		} else {
			fmt.Printf("  -%s", f.ShortName)
		}

		if v := fmt.Sprint(f.Default); v != "" {
			fmt.Printf("[=%v]", v)
		}

		fmt.Println()
		fmt.Println("     ", f.Usage)
		if i != len(fs.flags)-1 {
			fmt.Println()
		}
	}
}
