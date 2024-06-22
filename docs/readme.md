# extracting certificates

This powershell-script contained in this directory extracts the certificates from the rewe apk. 

## usage 

- Optional: Provide a working directory (`-WorkingDirectory <Path>`); if not provided, the current working-directory is used
- Optional: Provide the apk-file to use (`-ApkFile <Path>`); if not provided:
	- the working directory is searched for an apk file
	- if not found, rewe apk ver. 3.18.5 is downloaded from uptodown.net

In a powershell-window, run
```powershell
.\rewerse-engineering.ps1
```

Currently trying to get better at powershell, feedback appreciated.

## download + run

In a powershell-window, run
```powershell
git clone https://github.com/ByteSizedMarius/rewerse-engineering
cd .\rewerse-engineering\docs

$OriginalExecutionPolicy = Get-ExecutionPolicy
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force

.\rewerse-engineering.ps1

Set-ExecutionPolicy -ExecutionPolicy $OriginalExecutionPolicy -Scope Process -Force
```

## manual extraction

You can also quite easily extract the pfx manually by [downloading the apk](https://rewe.en.uptodown.com/android/post-download/1014869773), opening it as a zip (or renaming the file to `.zip`) and then navigating to `/res/raw`. Then copy the pfx out of the zip. [Torbens commands](https://github.com/torbenpfohl/rewe-discounts/blob/main/how%20to%20get%20private.pem%20and%20private.key.txt#L16) for extracting cert+key should work, but I have not tested them.

## misc

[Torben](https://github.com/torbenpfohl/rewe-discounts) also has a python helper that's a bit less overengineered.