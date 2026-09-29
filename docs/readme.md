# extracting certificates

Even though [manual extraction](#manual-extraction) is not necessarily difficult, I tried automating the whole process in PowerShell. The [script](./rewerse-engineering.ps1) contained in this directory extracts the certificate from a manually downloaded rewe apk and splits the certificate into pem and key. It has some commandline-options (including the ability to just extract pem and key from a manually obtained pfx) described in [#usage](#usage).

> [!CAUTION]
> The PS script opens the apk in memory as a zip archive to extract the pfx from it. Windows Defender does not seem to like this very much.


## download + run

1. Download the rewe xapk from [uptodown](https://rewe.en.uptodown.com/android/download/1220537953-x) in your browser. Read the buttons carefully: you want the xapk file, **not** the Uptodown installer or app store.
2. In a powershell-window, in the folder containing the apk, run this one-liner, that clones the repo, temporarily sets the execution policy and executes the script. The certificate files end up next to the apk.
```powershell
git clone https://github.com/ByteSizedMarius/rewerse-engineering; Set-ExecutionPolicy Bypass -Scope Process -Force; .\rewerse-engineering\docs\rewerse-engineering.ps1
```

## manual extraction

You can also quite easily extract the pfx manually.
1. [Download the apk](https://apkpure.com/de/rewe-supermarkt/de.rewe.app.mobile/download). Version does not really matter. However, pay attention to which download button you click, as many of these sites have been introducing dark patterns to get you to install their store instead. 
2. Rename the file from `.apk` or `.apkx` to `.zip` and open the zip.
	- If you had an `apkx`-file: Copy `de.rewe.app.mobile.apk` out of the zip and redo step 2 with this apk.
3. Navigate to `/res/raw`, where you will find the `mtls_prod.pfx`. Copy it out of the zip. 
4. Extract key and pem from the `.pfx`. 
	- Using the PowerShell-script: `./rewerse-engineering -PfxPath "/path/to/.pfx/"`
	- Or [Torbens openssl commands](https://github.com/torbenpfohl/rewe-discounts/blob/main/how%20to%20get%20private.pem%20and%20private.key.txt#L16)

## usage 

```powershell
.\rewerse-engineering.ps1
```

- Optional: Provide a working directory (`-WorkingDirectory <Path>`); if not provided, the current working-directory is used
- Optional: Provide the apk-file to use (`-ApkFile <Path>`); if not provided:
	- the working directory is searched for an .apk/.xapk file
	- if none is found, the script tells you where to download it
- Optional: Just extract key/pem from `mtls_prod.pfx` in current WorkingDirectory (`-Pfx`)
- Optional: Just extract key/pem from pfx at path (`-PfxPath <Path>`)

## misc

Apk versions tested with the script:

| version | tested |
|---------|--------|
| 5.18.1  | ✅      |
| 5.15.2  | ✅      |
| 5.11.0  | ✅      |
| 5.7.3   | ✅      |
| 4.1.0   | ✅      |
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