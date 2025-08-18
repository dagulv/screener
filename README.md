# Stock screener

A system that fetches financial information from the stock market and presenting the data concisely.

---

## Table of Contents

* [Overview](#overview)
* [Architecture](#architecture)
* [Features](#features)
* [License](#license)

---

## Overview

Stock screener is a free stock screener system that uses data from the stock market to filter, sort and present financial data in a concise manner. It is inspired by stock terminals and other stock screener systems. The system provides an clean and easy to interface as well as advanced tools such as the modular data table and the Magic Formula. 

**Motivation:** Stock screener systems can be overwhelming and expensive, this project aims at providing an intuitive system while also being completely free.

**Problem it solves:** It simplifies the user's process to find potentially great stocks with an easier learning curve than other stock screening systems.

**What I learned:** Building this system helped me improve my data proccessing skills, how to build reliable web scrapers, process and clean data efficiently, and integrate it with a backend API for use in a full-stack application.

---

## Architecture

- **Frontend**: Svelte SPA
- **Backend API**: Go
- **Database**: PostgreSQL
- **Communication**: HTTP REST API between frontend and backend

---

## Features

* View financial data over 650 Nordic companies
* Sort companies by any metric
* Frontend SPA with dynamic routing
* Realtime updates / WebSockets (if applicable)

---

## TODO

- Improved screening functionality

---

## License

This project is licensed under the [MIT License](LICENSE).

---
