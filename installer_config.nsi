OutFile "installer.exe"

InstallDir "$PROGRAMFILES\emissionsrechner"

RequestExecutionLevel admin

Section 
    SetOutPath $INSTDIR
    File /r "dist\main\*"
    WriteUninstaller "uninstall.exe"
    CreateShortCut "$SMPROGRAMS\Emissionsrechner.lnk" "$INSTDIR\main.exe"
    CreateShortCut "$DESKTOP\Emissionsrechner.lnk" "$INSTDIR\main.exe"
SectionEnd

Section "uninstall"
    Delete "$SMPROGRAMS\Emissionsrechner.lnk"
    Delete "$DESKTOP\Emissionsrechner.lnk"
    Delete "$INSTDIR\uninstall.exe"  
    RMDir /r "$PROGRAMFILES\emissionsrechner"
SectionEnd