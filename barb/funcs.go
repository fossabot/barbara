package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func reply(ctx *cli.Context) {
	client := getClient()
	if len(ctx.Args()) != 1 {
		exitError(errors.New("invalid arguments"))
	}

	args := ctx.Args()
	myRepo, err := repo()
	if err != nil {
		exitError(err)
	}

	f, err := ioutil.TempFile("", "barbara-edit")
	if err != nil {
		exitError(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	if err := runProgram(os.Getenv("EDITOR"), f.Name()); err != nil {
		exitError(err)
	}

	content, err := ioutil.ReadFile(f.Name())
	if err != nil {
		exitError(err)
	}

	_, err = client.AddComment(myRepo, args[0], string(content))
	if err != nil {
		exitError(err)
	}

	fmt.Printf("Comment on ticket %s posted!\n", args[0])
}

func get(ctx *cli.Context) {
	client := getClient()

	if len(ctx.Args()) != 1 {
		exitError(errors.New("invalid arguments"))
	}

	args := ctx.Args()

	myRepo, err := repo()
	if err != nil {
		exitError(err)
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		exitError(err)
	}

	issue, err := client.Issue(myRepo, i, nil)
	if err != nil {
		exitError(err)
	}

	comments, err := client.Comments(myRepo, args[0], nil)
	if err != nil {
		exitError(err)
	}

	f, err := ioutil.TempFile("", "barbara-edit")
	if err != nil {
		exitError(err)
	}

	color.Output = f

	line()
	color.New(color.FgBlue).Printf("From: %s\n", issue.User.Login)
	color.New(color.FgBlue).Printf("Title: %s\n", issue.Title)
	color.New(color.FgBlue).Printf("Number: %d\n", issue.Number)
	color.New(color.FgBlue).Printf("State: %s\n", issue.State)
	color.New(color.FgBlue).Printf("URL: %s\n", issue.URL)
	line()
	fmt.Fprintln(f, issue.Body)

	for _, comment := range comments {
		fmt.Fprintln(f)
		line()
		color.New(color.FgWhite).Printf("From: %s\n", comment.User.Login)
		color.New(color.FgWhite).Printf("Date: %s\n", comment.CreatedAt.Local())
		line()
		fmt.Fprintln(f)
		fmt.Fprintln(f, comment.Body)
	}

	fmt.Fprintln(f)

	f.Close()
	defer os.Remove(f.Name())

	if err := runProgram("less", "-R", f.Name()); err != nil {
		exitError(err)
	}
}
