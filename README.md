# INFO 344 In-Class Coding

[Fork this repo](https://help.github.com/articles/fork-a-repo/) and use your fork for all the code you write in-class.

I will add more directories and starter code to this repo as we go along. To pull these new files into your fork, execute the commands below.

## On a Lab Machine?

The lab machines completely reset when you log out, so execute all of these commands each time you need to pull updates.

First, tell git who you are:

```bash
git config --global user.name "Your Name"
git config --global user.email your-netid@uw.edu
```

Then clone your **forked repo** to your lab machine, and execute these commands **from within the repo directory** to pull the updates:

```bash
git remote add upstream https://github.com/info344-a17/info344-in-class.git
git pull upstream master
```

If you end up in `vim` to confirm a merge message, just hit `Esc` and type `:wq` (for "write and quit") to accept the default message and get back to your command prompt.

## On Your Own Laptop?

After cloning your **forked repo** for the first time, execute this command once **from within the repo directory** to setup the upstream remote:

```bash
git remote add upstream https://github.com/info344-a17/info344-in-class.git
```

When I ask you to pull updates, execute this command **from within the repo directory**:

```bash
git pull upstream master
```

If you end up in `vim` to confirm a merge message, just hit `Esc` and type `:wq` (for "write and quit") to accept the default message and get back to your command prompt.
