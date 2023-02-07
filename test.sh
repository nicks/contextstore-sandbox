#!/bin/bash

go build -o main.out ./main.go
for ((i=1;i<=100;i++)); 
do 
   ./main.out
done
