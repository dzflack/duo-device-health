## Duo Device Health

Bypass client-side checks for [Duo Device Health Application](https://duo.com/docs/device-health).

**WARNING**: Multiple failed attempts will lock out you account. If you fail to auth with 2 consecutive attempts, wait for awhile before you try again. Alternatively reset your failed auth attempts by using another device (i.e. a phone) to successfully auth, and then try again.

**WARNING**: Tested and working on MacOS with Firefox only. 

### How To

* First make sure actual Duo Device Health app is not running

* Then clone repo

```
git clone git@github.com:dzflack/duo-device-health.git
```

* Modify the following fields `HealthCheckStartTimestamp, HealthCheckEndTimestamp, DeviceID, DeviceName, HealthCheckLengthMillis` in [main.go](https://github.com/dzflack/duo-device-health/blob/master/main.go#L63-L79)

* Run `main.go`

```
go run main.go
```

* Make sure we are listening, i.e. the following is displayed in terminal:

```
2020/06/09 10:47:29 Listening on port 53106
```

* If you are not using the TLS environment variables specified under "Options," use firefox to login to your DUO protected account. (With TLS, you can alternatively use Chrome.)

### Options

With the use of environment variables, you can modify this program's behaviour:

* `DUO_LOCAL_TLS_CERT` and `DUO_LOCAL_TLS_KEY` - Specify these as filepaths to a cert and key file. _For use with Chrome, the cert should be signed by a CA (even if the CA is just you), and that CA's cert should be imported as an Authority in your Chrome settings._
* `DUO_LOCAL_PORT` - Override port on which to listen.
* `DEVICE_ID`, `DEVICE_NAME` - Override DeviceID and DeviceName in the payload to duosecurity.com.
* `OS`, `OS_BUILD`, `OS_VERSION` - Override Os, OsBuild, and OsVersion in the payload to duosecurity.com.
* `DUO_CLIENT_VERSION` - Override DuoClientVersion in the payload to duosecurity.com.
