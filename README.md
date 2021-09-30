# ruth-m-api

The API backend for ruth-m-level-1.

I spent a lot of time trying to get an API Gateway working for this, but it did not work.
The spec file is still here (spec.yml) which Neil claims worked with his setup, but it did
not work for me.  The terraform to make it work is also here, it is just commented out. It
did take like ten minutes to spin up, I'm not sure if it's the best solution.

The deploy the API:

Download the key for the level-1-service-account@roi-takeoff-user47.iam.gserviceaccount.com
service account in the user47 account to ~/roi-takeoff-user47.json and run ./deploy.sh.  To
tear down, run ./tear-down.sh.  The API will be available at https://events-api-ex7otr565q-uc.a.run.app.
