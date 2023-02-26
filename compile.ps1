$cd_to_workdir = "cd " + $(Get-Location) + ";"
$install_requirements = 'pip install -r requirements.txt;'
$compile_exe = 'pyinstaller --onedir --windowed --add-data "C:\Python311\Lib\site-packages\customtkinter`;customtkinter/" --noconsole --noconfirm main.py'
$args = '-noexit -noprofile -command "' + $cd_to_workdir + 'echo "Checking requirements...`n;"' + $install_requirements + 'echo "`n`nCompiling...`n;"' + $compile_exe + '"'

Start-Process powershell -Verb RunAs -ArgumentList $args