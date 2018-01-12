# Foundry Fred's Fire Drills

<pre>
         .'.---.'.
        //   ,   \\
       ||   '|    ||
       ||    |    ||
       ||   -'-   ||
  .-"''-.,_     _,.-''"-.
 / .'--,___'"""'___,--'. \
 |  /:////_'---'_\\\\:\  |
  \|:|// "_     _" \\|:|/
   '-/| (◕)     (◕) |\-'
     \\     | |     //
      '|   (._.)   |'
       |           |
       \     ‿     /
        '--.___.--'
       --------------
      | FOUNDRY FRED |
       --------------
</pre>

Howdy y'all! I'm Foundry Fred and I'm here to help teach you how to put out production fires before they happen.
I can run you through some fire drills to simulate real life production issues.

## Getting Started with the fire-drills CLI

```
git clone https://github.com/ljfranklin/fire-drills.git fire-drills
cd fire-drills
wget -O ./fire-drills https://github.com/ljfranklin/fire-drills/releases/download/v0.0.1/fire-drills_v0.0.1_osx
chmod +x ./fire-drills
./fire-drills -h
```

## Writing Your Own Fire Drills
In the `/drills` directory you'll notice several existing drills.
Each subdirectory in this folder contains a `drill.yml` as well as `setup` and `teardown` scripts for creating and cleaning up the testing environment.
Check out the `example-drill` for what a minimal drill might consist of.

Example `drill.yml` file:
```
setup_cmd: ./setup
teardown_cmd: ./teardown

required_env_vars:
  - MY_ENV_VAR

summary: |
  This is an example drill. More coming soon!

prompt: |
  This is where you would introduce the prompt to the user.

solution: |
  This is where you would show the solution and ways to check your solution is correct.

hints:
  - This is a vague hint.
  - This is a little more specific.
  - This tells the user exactly where to look.
```