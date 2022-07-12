# Quiz problem

Create a program that will read in a quiz provided via a CSV file (more details below) and will then give the quiz to a user keeping track of how many questions they get right and how many they get incorrect. Regardless of whether the answer is correct or wrong the next question should be asked immediately afterwards.

Show all flag support
```shell
./quiz -h [--help]
```

Run
```shell
go build .; ./quiz -csv="problems.csv" -limit=10
```

Ref:
- Problem: https://courses.calhoun.io/lessons/les_goph_01
