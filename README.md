<h1 align="center">
   <a href="#"> Temperature by PostCode with Golang </a>
</h1>

<h3 align="center">
   This program allows users to provide a ZIP code and receive the location's temperature in Celsius, Fahrenheit and Kelvin.
</h3>

<h4 align="center"> 
	 Status: Finished
</h4>

# About
In this project we use Golang, http Rest, Docker and Google Cloud Platform to Deploy (Cloud Run)

## Usage
1. Clone the repository.

## Dependencies
This program uses all dependencies in go.mod file
1. Run "go mod tidy" to download all dependencies.

#### Running Project

```bash

# Clone this repository
$ git clone https://github.com/renatafborges/cloud-run.git

# Access folder contaning docker-compose.yaml and let's up all containers:
$ docker-compose up

# The server will start at
Starting web server on port :8080

# After running docker-compose, the automated tests will start automatically
If succeed: Automated test executed with sucess

# Server:
https://cloudrun-goexpert-qqjgvrgoma-uc.a.run.app/temperature/{zipcode}

Valid ZipCode
Example: https://cloudrun-goexpert-qqjgvrgoma-uc.a.run.app/temperature/87020260

The expected result:
{"temp_C":"22.0","temp_F":"71.6","temp_K":"295.0"}

Invalid ZipCode
Example: https://cloudrun-goexpert-qqjgvrgoma-uc.a.run.app/temperature/abc00000

The expected result:
invalid zipcode

Not found ZipCode
https://cloudrun-goexpert-qqjgvrgoma-uc.a.run.app/temperature/99999999

The expected result:
cannot find zipCode

## Author
Made with love by Renata Borges 👋🏽 [Get in Touch!](Https://www.linkedin.com/in/renataborgestech)
