# Language Agnostic TEsting

Web service that allows to run tests for programms written in these languages:

* C
* Python
* Planned: Go, C++, C#, Pascal

# Contents
- [How it works](#how-it-works)
- [Requirements](#requirements)
- [Quick start](#quick-start)
- [Architecture](#architecture)
- [Tests structure](#tests-structure)
- [Service start](#service-start)
- [Requests](#requests)

# How it works

* ✉️ Web service receives solution source code for specific task
* 🔨 Source code is built inside separate docker container
* 🧪 If build succeeded, then solution is tested with various test cases
* 📊 User receives test result

This testing system is "language agnostic" because:

* All input parameters is passed via standard input
* Result is received in specific format via standard output

This is example of program for this testing system:

```python
n = int(input()) # Receive count
s = 0
for _ in range(n):
    s += int(input()) # Receive numbers n times
print(s) # Output addition result
```

> Yes I know about command line arguments, but I've built it this way, so programms still can be executed and tested manually as usual.
> 
> I don't consider that it's good idea to teach begginers in programming about command line arguments at first lesson.

# Requirements

* docker-compose

# Quick start

```bash
# Run all containers in detached mode for dev environment
./run-docker-compose.sh dev up -d

# Get id of "manage" container and open interactive bash shell inside of it
sudo docker exec -it $(sudo docker ps | grep late_manage | cut -d' ' -f1) bash
```

Inside opened bash shell:

```bash
# Fill database with sample project, create user and run all available requests
./test_service.sh
```

# Architecture

Service have 4 containers:

* 🕸 web - web service written in Go, that:
	* Receives requests from clients
	* Communicates with **db**
	* Sends solutions into runner container
	* Responds with test result
* 🏃 runner - internal web service written in Python, that:
	* Receives solutions from **web** service
	* Builds solutions (if it's not written with interpreted language)
	* Tests solutions
	* Responds with test result to **web** service
* 🏗 manage - container with Bash and Python scripts, that could be used for:
	* Filling database with tests
	* Creating users
	* Giving tokens to users, that's required to send any solutions for testing
* 🗄 db - PostgreSQL container (postgres:latest)

# Tests structure

Main purpose of this web service is testing of specific programms, so let's figure out how you need to set them up.

Tests is organized this way:

`"tests"` -> `project` -> `unit` -> `task`

* `"tests"` - folder in project root directory, that contains projects
* `project` - folder with arbitrary name, that contains units
* `unit` - folder with arbitrary name, that contains tasks
* `task` - folder with arbitrary name, that contains actual test data

`project`, `unit` and `task` folders contains file `desc.json`, that contains descripton for according folder. Here are neccessary fields for every folder type:

* `project`
	* "name" - human readable name of project
* `unit`
	* "name" - human readable name of unit
* `task`
	* "name" - human readable name of project
	* "position" - position inside unit when it will be presented to user
	* "desc" - text description that will help user to understant given task
	* "input" - format of input data for program
		* "name" - name of input parameter
		* "type" - type of passed values
		* "range" - range (from min to max) in which values for random tests will be generated
		* "dimensions" - if this field is not specified, then it is just single value, if value specified - it is specification of array size and dimensions. Each value specifies size of dimension. Examples:
			* [50] - programm could receive array from 1 up to 50 elements (size will be randomly generated in each test)
			* [10, 5] - matrix 10x5 (each size for each dimension will be generated randomly in range [1,10]x[1,5])
			* [3, 0] - if zero is specified, size will be the same as previous one (in given example, possible sizes for matrix are 1x1, 2x2, 3x3)
	* "output" - text description of output format

This is example of `desc.json` file for some `task`:

```json
{
	"name": "Add to array",
	"desc": "Add a number to all values in an array",
	"input": [
		{"name": "A", "type": "int", "range": ["-1000", "1000"]}, 
		{"name": "B", "type": "int", "range": ["-1000", "1000"], "dimensions": [50]}
	],
	"output": "Result of adding A to B"
}
```

Apart from `desc.json` file, task folder also must contain other files:

* `complete_solution.*` - file with source code of reference solution. Output of this file will be compared with incoming solutions - if output differs, than test of incoming solution fails
* `fixed_tests.txt` - file with tests for solution. It contains values that will be passed into both reference and incoming solutions
* `template.*` - file with template for solution. Contents of this file could be used on UI side, to provide user with sample code for easy start

I have [repository](https://github.com/kee-reel/late-sample-project) with example project - you can use it for for reference.

# Service start

You can easily start web service with docker-compose:

```bash
./run-docker-compose.sh dev up -d # Run all containers in detached mode for dev environment
```

> You can check file [/manage/test\_service.sh](/manage/test_service.sh) - it sets up environment and makes sample requests.

After that you can manage web server via **manage** container. To open interactive bash shell inside **manage**:

```bash
# Get id of manage container and open bash inside "manage" of it
sudo docker exec -it $(sudo docker ps | grep late_manage | cut -d' ' -f1) bash
```

Then you need to prepare tests - you can use mine for this time:

```bash
cd tests # Go inside (tests folder is created during container build process)
git clone https://github.com/kee-reel/late-sample-project # Clone sample project
cd .. # Go back
```

Tests are ready, lets insert them into database:

```bash
python3 fill_db.py # Fill database with sample project
```

All set, now we can try to send requests to web server.

# Requests

API responds with error codes that ould be found in [web/api/errors.go](/web/api/errors.go)

If you know what Swagger is -- you can check specs for all requests in [swagger.json](/web/docs/swagger.json).

If you don't -- check out hosted Swagger docs [here](https://kee-reel.com/cyber-api/).
