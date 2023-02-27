$cd_to_workdir = "cd " + $(Get-Location) + ";"
$install_requirements = 'python -m pip install -r requirements.txt;'
$compile_exe = 'python -m PyInstaller --onedir --windowed --add-data "C:\Users\Krets\AppData\Local\Programs\Python\Python311\Lib\site-packages\customtkinter`;customtkinter/" --noconsole --noconfirm main.py'
$args = '-noexit -noprofile -command "' + $cd_to_workdir + 'echo "Checking requirements...`n;"' + $install_requirements + 'echo "`n`nCompiling...`n;"' + $compile_exe + '"'

Start-Process powershell -Verb RunAs -ArgumentList $args