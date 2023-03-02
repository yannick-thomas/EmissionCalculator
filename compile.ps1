$cd_to_workdir = "cd " + $(Get-Location) + ";"
$install_requirements = 'python -m pip install -r requirements.txt;'

$customtkinter_location = $(python -m pip show customtkinter)
$customtkinter_location = $customtkinter_location.split([Environment]::NewLine)[7].split(" ")[1] + "\customtkinter"
echo "Found customtkinter at:`n" + $customtkinter_location + "`n-> adding it to pyinstaller metadata..."


$build_files = ".\installer.exe", ".\build", ".\dist", ".\main.spec"
echo "Delete existing build..."
foreach($buildfile in $build_files)
{
    if (Test-Path $buildfile) {
        Remove-Item $buildfile -Recurse -Force
    } else {
        echo "Buildfile $buildfile doesn't exist..."
    }
}
$compile_exe = 'python -m PyInstaller --onedir --windowed --icon=favicon.ico --add-data "' + $customtkinter_location + '`;customtkinter/" --noconsole --noconfirm main.py;'
$copy_icon = 'cp favicon.ico ./dist/main/favicon.ico;'
$create_installer = 'makensis installer_config.nsi;'

$args = '-noexit -noprofile -command "' + $cd_to_workdir + 'echo "Checking requirements...`n;"' + $install_requirements + 'echo "`n`nCompiling...`n;"' + $compile_exe + 'echo "`n`nCopying Favicon...`n;"'+ $copy_icon + 'echo "`n`nCreating installer...`n;"' + $create_installer + '"'
Start-Process powershell -Verb RunAs -ArgumentList $args

# 'echo "`n`nDeleting existing build...`n;"'+ $delete_existing_build +