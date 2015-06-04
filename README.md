Go Get Android SDK Tools
========================

Installing the [Android SDK Tools](https://developer.android.com/sdk/index.html)
and managing its archives in an automated fashion is a PITA (Probably Isn't That Amusing). And although others have written wrappers before I didn't find anyone
yet who wrote one in go.

### Goals:

 * Self contained binaries, simplest for CI systems
 * Cross platform
 * Idempotent - install from scratch or just keep an install up-to-date
 * Flexible enough to handle changes as the tool changes (i.e. download urls not being stable, archive installation not being idempotent in tools, )

The usual disclaimer applies. This is a work in progress and has some rough edges.

### A few alternatives worth looking at:

 * [android-sdk-installer](https://github.com/embarkmobile/android-sdk-installer) - Bash Scripts
 * [sdk-manager-plugin](https://github.com/JakeWharton/sdk-manager-plugin) - Gradle Plugin
