#!bin/bash

sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'thisIsSparta'"
sudo -u postgres psql -c "CREATE DATABASE youth2k"
sudo -u postgres psql -c "ALTER DATABASE youth2k SET TIMEZONE TO 'UTC';"