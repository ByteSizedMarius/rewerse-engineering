# extracting certificates

The [powershell-script](./rewerse-engineering.ps1) contained in this directory automatically downloads the rewe apk from uptodown, extracts the certificate from it and splits the certificate into pem and key. 

The script has limited commandline-options, described in [#usage](#usage).

## download + run

In a powershell-window, run
```powershell
git clone https://github.com/ByteSizedMarius/rewerse-engineering
cd .\rewerse-engineering\docs
Set-ExecutionPolicy Bypass -Scope Process -Force
.\rewerse-engineering.ps1
```

or as a one-liner:
```powershell
git clone https://github.com/ByteSizedMarius/rewerse-engineering; Push-Location .\rewerse-engineering\docs; Set-ExecutionPolicy Bypass -Scope Process -Force; .\rewerse-engineering.ps1; Pop-Location
```

## manual extraction

You can also quite easily extract the pfx manually by first [downloading the apk](https://apkpure.com/de/rewe-supermarkt/de.rewe.app.mobile/download). Pay attention to which download button you click, as uptodown introduced dark patterns to get you to install their store instead. Next, rename the file from `.apk` or `.apkx` to `.zip`. If you had an `apkx`-file, the real apk is nested under the name `de.rewe.app.mobile.apk`. Extract it, and open it as a zip again. Next, navigate to `/res/raw`, where you will find the `mtls_prod.pfx`. [Torbens commands](https://github.com/torbenpfohl/rewe-discounts/blob/main/how%20to%20get%20private.pem%20and%20private.key.txt#L16) for extracting cert+key should work, but I have not tested them.

## usage 

```powershell
.\rewerse-engineering.ps1
```

- Optional: Provide a working directory (`-WorkingDirectory <Path>`); if not provided, the current working-directory is used
- Optional: Provide the apk-file to use (`-ApkFile <Path>`); if not provided:
	- the working directory is searched for an apk file
	- if not found, rewe apk ver. 4.0.2 is downloaded from uptodown.net

## misc

Apk versions tested with the script:

| version | tested  |
|---------|---------|
| 4.0.3   | ✅      |
| 4.0.2   | ✅      |
| 3.21.4  | ✅      |
| 3.20.0  | ✅      |
| 3.19.3  | ✅      |
| 3.18.6  | ✅      |
| 3.18.5  | ✅      |
| 3.18.4  | ✅      |
| 3.18.3  | ✅      |
| 3.18.2  | ✅      |
| 3.18.1  | ✅      |
| 3.18.0  | ✅      |
| 3.17.5  | ✅      |
| 3.16.6  | ✅      |
| 3.16.5  | ✅      |
| 3.16.2  | ✅      |

**Notes:**

- [Torben](https://github.com/torbenpfohl/rewe-discounts/blob/main/rewe_discounts/get_creds.py) also has a python helper for extracting the cert that's a bit less overengineered
- Starting with v3.19, the app now seems to be packaged as an xapk, which means the apk containing the certificate is nested. This requires unzipping twice
- The pfx bundled with the apk seems to be an old format (RC2-40-CBC) and may not be supported everywhere. You may have to convert it to a newer format if you are experiencing strange issues (ask me how I know)
- Currently trying to get better at powershell, feel free to criticise relentlessly