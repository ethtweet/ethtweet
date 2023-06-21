@echo off & setlocal enabledelayedexpansion
for /r .  %%i in (*.zip) do ( gpg --armor --detach-sign  %%i )
pause