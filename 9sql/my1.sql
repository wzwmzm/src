SET ECHO OFF;
SET VERIFY OFF;


COLUMN '����' format a10;
select gjdqdm as "����",count(*) as "����" from BJ_YW_T_CRJRYDK WHERE CRRQ='&2' AND crkadm='055' and crbz='1' AND CRJSYDM<>'A' AND gjdqdm IN ('&1') GROUP BY gjdqdm;

column ���Թ� format a10;
select qwglzg as "���Թ�",count(*) as "������" from BJ_YW_T_CRJjtgjDK WHERE CRRQ='&2' AND crkadm='055' and crbz='1' AND substr(HC,1,2)<>'H/' AND qwglzg IN ('&1') GROUP BY qwglzg;


select count(*) as "�뾳�ÿ�" from BJ_YW_T_CRJRYDK WHERE crrq='&2' and crbz='1' and crkadm='055' and crjsydm<>'A' and jtgjbs in (select JWGJBS  from BJ_YW_T_CRJjtgjDK WHERE CRRQ='&2' AND crkadm='055' and crbz='1' AND substr(HC,1,2)<>'H/' AND qwglzg='&1' );

select count(*) as "�й���" from BJ_YW_T_CRJRYDK WHERE crrq='&2' AND GJDQDM IN ('CHN','TWN','HKG','MAC') and crbz='1' and crkadm='055' and crjsydm<>'A' and jtgjbs in (select JWGJBS  from BJ_YW_T_CRJjtgjDK WHERE CRRQ='&2' AND crkadm='055' and crbz='1' AND substr(HC,1,2)
<>'H/' AND qwglzg='&1' );

select count(*) as "������" from BJ_YW_T_CRJRYDK WHERE crrq='&2' AND GJDQDM  IN ('&1') and crbz='1' and crkadm='055' and crjsydm<>'A' and jtgjbs in (select JWGJBS  from BJ_YW_T_CRJjtgjDK WHERE CRRQ='&2' AND crkadm='055' and crbz='1' AND substr(HC,1,2)
<>'H/' AND qwglzg='&1' );

select COUNT(*) as "��������" from BJ_YW_T_CRJRYDK WHERE crrq='&2' AND GJDQDM IN ('&1') and crbz='1' and crkadm='055'  and jtgjbs not in (select JWGJBS  from BJ_YW_T_CRJjtgjDK WHERE CRRQ='&2' AND crkadm='055' and crbz='1' AND qwglzg='&1' and substr(hc,1,2)<>'H/' ) and substr(jtgjbs,1,2)<>'H/';





exit;