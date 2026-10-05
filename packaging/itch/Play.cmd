@echo off
setlocal
"%~dp0tdef.exe" %*
set "game_exit=%errorlevel%"
if not "%game_exit%"=="0" (
    echo.
    echo TDEF exited with error %game_exit%.
    pause
)
exit /b %game_exit%
