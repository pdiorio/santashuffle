<p align="center"><img src="gopher.png" alt="santagopher" width="200"/></p>

# Santa Shuffle

A Go project to create and notify Secret Santa assignments. 

Santa Shuffle helps you generate matches for Secret Santa participants. 
It allows you to specify participants, their email address for notification, and also other participants they should never be matched with as a gift-giver (e.g. spouses, enemies, etc).

While a dry-run mode option is available to generate matches locally and display all matches openly, the default mode aims to keep the matches secret to the facilitator by quietly emailing the assignments using standard smtp configuration with basic authentication (e.g. gmail with 2fa enabled can be used with application passwords). 

See example configurations in `/conf` for both specifying participants and also providing email templates + tokens.
