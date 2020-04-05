@echo off
echo ***************************************************
echo ********** 2017年1月常态化单机版升级批处理脚本 **********
echo ***************************************************
set /p DJBSP="请输入DJB实例sys用户的密码："
set /p DJB2SP="请输入DJB2实例sys用户的密码："
echo ***************************************************
echo 1、关停ADR和壳子程序....
echo ***************************************************
taskkill /f /t /im adr.exe
taskkill /f /t /im BJMain.exe
taskkill /f /im explorer.exe
taskkill /f /im explorer.exe
xcopy e:\oradata\*.* e:\oradata201701\ /e /y
xcopy e:\databak\*.* e:\databak201701\ /e /y
echo ***************************************************
echo 2、重置及重启DJB2数据库实例....
echo ***************************************************
net stop OracleServiceDJB2
rd /s /q I:\ORADATA
del E:\DATAFILE\*.* /f /q
xcopy e:\oradata\*.* I:\ORADATA\ /e /y
xcopy e:\databak\*.* e:\DATAFILE\ /e /y
net start OracleServiceDJB2
echo ***************************************************
echo 3、DJB2增加表空间....
echo ***************************************************
sqlplus sys/%DJB2SP%@djb2 as sysdba @AddTABLEPACE.sql
echo ***************************************************
echo 4、调整DJB数据库实例，删记录并增加相关代码....
echo ***************************************************

echo ***************************************************
echo 5、重新备份DJB2数据库实例....
echo ***************************************************
net stop OracleServiceDJB2
xcopy I:\oradata\*.* e:\ORADATA\ /e /y
xcopy e:\DATAFILE\*.* e:\databak\ /e /y
rd /s /q E:\databak201701
rd /s /q E:\oradata201701
echo ***************************************************
echo 6、重新启动Windows桌面....
echo ***************************************************
start explorer
echo ***************************************************
echo **********     升级完成！按任意键退出    **********
echo ***************************************************
pause>nul