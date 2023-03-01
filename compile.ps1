$cd_to_workdir = "cd " + $(Get-Location) + ";"
$install_requirements = 'python -m pip install -r requirements.txt;'
$compile_exe = 'python -m PyInstaller --onedir --windowed --icon=favicon.ico --add-data "C:\Users\Krets\AppData\Local\Programs\Python\Python311\Lib\site-packages\customtkinter`;customtkinter/" --noconsole --noconfirm main.py;'
$copy_icon = 'cp favicon.ico ./dist/main/favicon.ico;'
$create_installer = 'makensis installer_config.nsi;'
$args = '-noexit -noprofile -command "' + $cd_to_workdir + 'echo "Checking requirements...`n;"' + $install_requirements + 'echo "`n`nCompiling...`n;"' + $compile_exe + 'echo "`n`nCopying Favicon...`n;"'+ $copy_icon + 'echo "`n`nCreating installer...`n;"' + $create_installer + '"'

Start-Process powershell -Verb RunAs -ArgumentList $args
