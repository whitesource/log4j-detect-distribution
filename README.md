# Log4jDetect

WhiteSource Log4j Detect is a free CLI tool that quickly scans your projects to find vulnerable Log4j versions
containing the following known CVEs:

* CVE-2021-45046
* CVE-2021-44228

It provides the exact path — both to direct or indirect dependencies, along with the fixed version for speedy
remediation

The supported packages managers are:

* gradle
* maven

In addition, the tool will search for vulnerable files with the `.jar` extension.

### Prerequisites:

* Download the log4j-detect binary based on your OS platform (see installation steps below)

---
**NOTE**

The relevant binaries must be installed for the scan to work, i.e:

* `gradle` if the scanned project is a gradle project (contains a `settings.gradle` or a `build.gradle` file)
* `mvn` if the scanned project is a maven project (contains a `pom.xml` file)

---

## Usage

In order to scan your project, simply run the following command:

```shell
log4j-detect scan -d PROJECT_DIR
```

## Installation

### linux

```shell
VERSION=1.0.0
ARCH=x64 # or ARCH=arm64
wget "https://github.com/whitesource/icu-log4j-distribution/releases/download/latest/log4j-detect--$VERSION-linux-$ARCH.tar.gz"
tar -xzvf log4j-detect-linux-$ARCH.tar.gz
./log4j-detect -h
```

### mac

```shell
VERSION=1.0.0
ARCH=x64 # or ARCH=arm64 
wget "https://github.com/whitesource/icu-log4j-distribution/releases/download/latest/log4j-detect-$VERSION-darwin-$ARCH.tar.gz"
tar -xzvf log4j-detect-darwin-$ARCH.tar.gz
./log4j-detect -h
```

### windows

Download and
extract https://github.com/whitesource/icu-log4j-distribution/releases/download/latest/log4j-detect-VERSION-windows.zip,
and the run:

```shell
.\log4j-detect.exe -h
```