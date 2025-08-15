![Promo Image 1](https://raw.githubusercontent.com/zarkones/ControlSTUDIO/production/promo/promo1.png)

# INTRODUCTION
ControlSTUDIO is an adversary simulation framework made fully in Go, with support for malleable command and control (C2) profiles.

# SETUP
If you wish you can modify the default C2 profile located at "agent/profile.json".

Run: sh build.sh

All files would be located in the "export" directory. To start the C2 server just run one the one inside of "export/c2". It would automatically setup a database for you and everything you need. Run it with "-h" argument to see available commands.

To run the UI double-click on "export/ControlSTUDIO", then go to settings and enter the address of the C2 server. Then go back and click on the menu button of the C2 node in the diagram, there you'd be able to import a C2 profile. Default C2 profile is located at "agent/profile.json".

IMPORTANT!

The current C2 doesn't have authentication itself. Either connect to it locally or via VPN. However, the listener services (not the same as C2 API for management) would run on different interface and port depending on your C2 profile.
