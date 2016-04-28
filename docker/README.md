Docker file for Factorio Server Manager

This dockerfile builds an image for running Factorio Server Manager.

UI is exposed by default on port 8080.


To Build:
docker build -t factorio-server-manager .

To Run:
docker run -d -p 8080:8080 -v factorio-directory:/factorio factorio-server-manager


