## Duo Device Health

Bypass client-side checks for [Duo Device Health Application](https://duo.com/docs/device-health). Tested on macOS only.

**WARNING**: Multiple failed attempts will lock out you account. If you fail to auth with 2 consecutive attempts, wait for awhile before you try again. Alternatively reset your failed auth attempts by using another device (i.e. a phone) to successfully auth, and then try again.

**WARNING**: Testing and working on MacOS with Firefox only. 

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

* Use a web browser to login to your DUO protected account


