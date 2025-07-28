# Nexpull

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

![Go Version](https://img.shields.io/badge/GO-1.23.0-red)
[![Go Report Card](https://goreportcard.com/badge/github.com/zachdehooge/radar_database)](https://goreportcard.com/report/github.com/zachdehooge/radar_database)
[![CI/CD Pipeline](https://github.com/Zachdehooge/Radar-Database/actions/workflows/CICD-Pipeline.yml/badge.svg)](https://github.com/Zachdehooge/Radar-Database/actions/workflows/CICD-Pipeline.yml)
<!--[![GO Snyk Check Master](https://github.com/Zachdehooge/Radar_Database/actions/workflows/snyk.yml/badge.svg?branch=main)](https://github.com/Zachdehooge/Radar_Database/actions/workflows/snyk.yml)-->
![Alt](https://repobeats.axiom.co/api/embed/593c1f4cb3a0b113e91f1675f5330631fe2926ea.svg "Repobeats analytics image")

## About
* An application for downloading Level II data from the NEXRAD Level II archive
* Takes user input for the month, day, year, radar site code (KHTX - Huntsville AL for example), along with beginning and ending time frames in Zulu
* (Format for times goes as follows: HHMMSS) with all times needing to be in Zulu

<!--## Installation

1. Download Github repo
2. Run `.exe` file
3. If smartscreen comes up, click more info -> run anyway // this warning is entirely harmless and only shows because the app is not signed-->

## Docker Installation
1. Make sure Docker is installed
2. Sign into Docker
3. Open a command line prompt
4. Run `docker login`
5. Run `docker pull zachdehooge/nexpull`
6. Run `docker run -it --rm -v .:/app zachdehooge/nexpull`

## Issues
* Be sure to open an issue and I will be more than happy to fix it!

## Roadmap
1. More to come...
