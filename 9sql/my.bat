@echo off
set /p GJ="请输入国籍:"
set /p RQ="请输入日期:"

sqlplus -S qgtg/qgtg@6.0.96.8 @my1.sql %GJ% %RQ%

pause
