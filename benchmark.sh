#!/bin/bash

go run main.go -broker=kafka -count=1000000 -message_size=1 -batch_size=1000000
go run main.go -broker=ripple -count=50000 -message_size=1 -batch_size=61440