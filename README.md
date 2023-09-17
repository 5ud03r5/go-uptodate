# go-uptodate
Keeps track of security updates of used applications, still in development

# Backend
Provide a way of being notified on every spotted application vulnerability along with provided by application owner security fix.
This part will be done in microservices structure and each microservice will be responsible for notifying main go-uptodate backend about such update.

go-uptodate consist of 3 main components - Microservice, Application, User. They are described below in respective sections.

## Microservice
Microservice is a software which is build on top of any preffered language capable of handling JSON communication
Such microservice needs to be first registered in go-uptodate backend to perform requests (auth mechanism not yet implemented, will be done with service accounts approach)

Microservice structure should be capable of:
1. Getting information of found vulnerabilities (which can be found on offciial CVE website) and update version status if it is vulnerable or not
2. Getting information from application release notes about security fixes for specific application version
3. Sending such informations to main go-uptodate backend to developed endpoint

## Application
Application is a go-uptodate abstraction of an application which is then used to notify user about any changes on the application. Additionally it helps track information about previous versions.
Application update request (comming from the endpoint) is an upsert operation, ID is defined in a way to prevent duplicates and store history of the application.

## User
User is anyone who is subscribing to an application within go-uptodate backend.
User can register and then create binding between himself and application
User provides an endpoint and email which will serve as communication layer about any updates for subscribed applications.
Endpoint will work as a microservice which should be capable of handling JSON (not yet defined, might be protobuf as well) communication.

# Frontend
Not yet defined, probably will be implemented.
