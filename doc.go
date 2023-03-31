/*
go-scaffold is an opinionated scaffolding tool. It is intended for internal use at AboveSoftware,
but shared as OSS for anyone to use or modify freely.

If there is go.mod file present in the desired folder, go-scaffold will fail.
This is a tool for creating a project, not for modifying existing ones.

Usage:

  - go-scaffold [folder]

If no folder is passed to go-scaffold, the project will be created on pwd.

If a folder is passed to go-scaffold, the project will be created on such folder.

If the folder passed is "help", no project will be created, a help message will be
shown instead. If you want to create a project in a help folder, first create the
folder, cd to it, and then run "go-scaffold" with no arguments.

This tool will ask for some basic options, for example:
  - Web library
  - DB library
  - DBMS
  - etc

You can cancel the execution at any time before the final stage of accepting the configuration, and nothing will be created
*/
package main
