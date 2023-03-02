OutFile "installer.exe"

InstallDir "$PROGRAMFILES\Emissionsrechner"

RequestExecutionLevel admin
Icon "favicon.ico"

Section 
    SetOutPath $INSTDIR
    File /r "dist\main\*"
    WriteUninstaller "uninstall.exe"
    CreateShortCut "$SMPROGRAMS\Emissions rechner.lnk" "$INSTDIR\main.exe"
    CreateShortCut "$DESKTOP\Emissions rechner.lnk" "$INSTDIR\main.exe" 
SectionEnd

Section "uninstall"
    Delete "$SMPROGRAMS\Emissions rechner.lnk"
    Delete "$DESKTOP\Emissions rechner.lnk"
    Delete "$INSTDIR\uninstall.exe"  
    RMDir /r "$PROGRAMFILES\Emissionsrechner"
SectionEnd