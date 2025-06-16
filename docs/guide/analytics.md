---
title: Anonymous Analytics
---

# Data Collection of Analytics

Limascope collects anonymous user configurations using a simple beacon written in Go. _Why?_ Limascope is an open source project with no funding. As a result, there is no time to do user studies of Limascope. Analytics are collected to prioritize features and fixes based on how people use Limascope.

## Where is Data Stored

Limascope sends anonymous data to DigitalOcean, where it is written to a flat file for processing.

## Opting Out

Limascope analytics helps to prioritize features and spend time on the most important improvements. If you do not want to be tracked, use the `--no-analytics` flag or `DOZZLE_NO_ANALYTICS` environment variable.
