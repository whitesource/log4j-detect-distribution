# Log4jDetect

WhiteSource Log4j Detect is a free CLI tool that quickly scans your projects to find vulnerable Log4j versions
containing the following known CVEs:

* CVE-2021-45046
* CVE-2021-44228
* CVE-2021-4104
* CVE-2021-45105
* CVE-2021-44832
* CVE-2020-9488
* CVE-2020-9493
* CVE-2022-23302
* CVE-2022-23305
* CVE-2022-23307

It provides the exact path to direct and indirect dependencies, along with the fixed version for speedy remediation.

The supported packages managers are:

* gradle
* maven
* bundler

In addition, the tool will search for vulnerable files with the `.jar`,`.gem` extensions.

### Prerequisites:

* Download the log4j-detect binary based on your OS platform (see installation steps below)

---
**NOTE**

1. For mac users, if the following message appears:
   "log4j-detect can't be opened because Apple cannot check it for malicious software", please follow the steps
   [described here](https://support.apple.com/en-il/guide/mac-help/mchleab3a043/mac)


2. The relevant binaries must be installed for the scan to work, i.e:
    * `gradle` if the scanned project is a gradle project (contains a `settings.gradle` or a `build.gradle` file)
    * `mvn` if the scanned project is a maven project (contains a `pom.xml` file)
    * `ruby`/`jruby` and `gem`/`jgem` if the scanned project is a bundler project (contains a `Gemfile.lock`/`gems.locked` file)


3. Building the projects before scanning will improve scan time and reduce potential scan errors

    * maven projects __must__ be built prior to scanning, e.g. with the following command:
       ```shell
       mvn install
       ```

    * bundler projects __must__ be built prior to scanning, e.g. with the following command:
       ```shell
       jbundler install
       ```

    * It is not necessary to run `gradle build` prior to scanning a `gradle` project, but that will greatly decrease the
      scan time

---

## Usage

In order to scan your project, simply run the following command:

```shell
log4j-detect scan -d PROJECT_DIR
```

The folder can include source code that uses supported package managers in the project, as well binaries with the
supported extensions mentioned above.
It may error if it's run in a location which has protected folders it cannot access, such as Windows system folders.

## Installation

### Linux

```shell
ARCH=amd64 # or ARCH=arm64
wget "https://github.com/whitesource/log4j-detect-distribution/releases/latest/download/log4j-detect-1.5.0-linux-$ARCH.tar.gz"
tar -xzvf log4j-detect-1.5.0-linux-$ARCH.tar.gz
chmod +x log4j-detect
./log4j-detect -h
```

### Mac

```shell
ARCH=amd64 # or ARCH=arm64 
wget "https://github.com/whitesource/log4j-detect-distribution/releases/latest/download/log4j-detect-1.5.0-darwin-$ARCH.tar.gz"
tar -xzvf log4j-detect-1.5.0-darwin-$ARCH.tar.gz
chmod +x log4j-detect
./log4j-detect -h
```

### Windows

```powershell
Invoke-WebRequest -Uri "https://github.com/whitesource/log4j-detect-distribution/releases/latest/download/log4j-detect-1.5.0-windows-amd64.zip" -OutFile "log4j-detect.zip"
Expand-Archive -LiteralPath 'log4j-detect.zip'
cd log4j-detect
.\log4j-detect.exe -h
```
