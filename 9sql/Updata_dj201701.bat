@echo off
echo ***************************************************
echo ********** 2017��1�³�̬������������������ű� **********
echo ***************************************************
set /p DJBSP="������DJBʵ��sys�û������룺"
set /p DJB2SP="������DJB2ʵ��sys�û������룺"
echo ***************************************************
echo 1����ͣADR�Ϳ��ӳ���....
echo ***************************************************
taskkill /f /t /im adr.exe
taskkill /f /t /im BJMain.exe
taskkill /f /im explorer.exe
taskkill /f /im explorer.exe
xcopy e:\oradata\*.* e:\oradata201701\ /e /y
xcopy e:\databak\*.* e:\databak201701\ /e /y
echo ***************************************************
echo 2�����ü�����DJB2���ݿ�ʵ��....
echo ***************************************************
net stop OracleServiceDJB2
rd /s /q I:\ORADATA
del E:\DATAFILE\*.* /f /q
xcopy e:\oradata\*.* I:\ORADATA\ /e /y
xcopy e:\databak\*.* e:\DATAFILE\ /e /y
net start OracleServiceDJB2
echo ***************************************************
echo 3��DJB2���ӱ�ռ�....
echo ***************************************************
sqlplus sys/%DJB2SP%@djb2 as sysdba @AddTABLEPACE.sql
echo ***************************************************
echo 4������DJB���ݿ�ʵ����ɾ��¼��������ش���....
echo ***************************************************

echo ***************************************************
echo 5�����±���DJB2���ݿ�ʵ��....
echo ***************************************************
net stop OracleServiceDJB2
xcopy I:\oradata\*.* e:\ORADATA\ /e /y
xcopy e:\DATAFILE\*.* e:\databak\ /e /y
rd /s /q E:\databak201701
rd /s /q E:\oradata201701
echo ***************************************************
echo 6����������Windows����....
echo ***************************************************
start explorer
echo ***************************************************
echo **********     ������ɣ���������˳�    **********
echo ***************************************************
pause>nul