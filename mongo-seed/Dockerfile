# NOTHING WORKS FOR NOW
FROM mongo:3.0.15

# Seed data file
COPY init.json /init.json

# --drop the collection, then recreate with by importing
CMD mongoimport --host mongo --user who --password dat --db devdb --collection subscribers --type json --drop --file /init.json --jsonArray
