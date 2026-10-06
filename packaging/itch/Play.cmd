@echo off
setlocal
"%~dp0termtd.exe" %*
set "game_exit=%errorlevel%"
if not "%game_exit%"=="0" (
    echo.
    echo TERMTD exited with error %game_exit%.
    pause
)
exit /b %game_exit%
