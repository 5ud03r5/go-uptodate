# go-uptodate
Keeps track of security updates of used applications

Still in development...

# Description
Provide a way o being notified on every spotted application vulnerability along with provided by application owner security fix.
This part will be done in microservices structure and each microservice will be responsible for notifying main go-uptodate backend about such update.

Microservice structure should be capable of:
1. Getting information of found vulnerabilities (which can be found on offciial CVE website) and update version status if it is vulnerable or not
2. Getting information from application release notes about security fixes for specific application version
3. Sending such informations to main go-uptodate backend to developed endpoint

Application update request is an upsert operation, ID is defined in a way to prevent duplicates and store history of the application
