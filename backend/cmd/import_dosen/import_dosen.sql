-- Import Dosen from Excel data
-- Generated automatically

BEGIN;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'AFRI ANDIKA', 'afri-andika', '', 'S2, S.M,M.M', 'Bisnis Digital (S1)', 'excel_import', NULL, true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'DINA HADIA', 'dina-hadia', 'Asisten Ahli', 'S2, S.E.,M.Si', 'Bisnis Digital (S1)', 'excel_import', '1025129202', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'NOFRIADI', 'nofriadi', 'Lektor', 'S2, Drs,M.M', 'Bisnis Digital (S1)', 'excel_import', '1018116301', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'RAFNELLY RAFKI', 'rafnelly-rafki', 'Lektor', 'S2, S.H,M.Kn.', 'Bisnis Digital (S1)', 'excel_import', '1003086602', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'RENI RAHMAWATI', 'reni-rahmawati', '', 'S2, S.E.,M.M.', 'Bisnis Digital (S1)', 'excel_import', '1013107201', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'RIO ANDIKA MALIK', 'rio-andika-malik', 'Lektor', 'S3, S.Kom,M.Kom,Dr', 'Bisnis Digital (S1)', 'excel_import', '1019058907', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'SRI MONA OCTAFIA', 'sri-mona-octafia', 'Lektor', 'S2, S.E.,M.M.', 'Bisnis Digital (S1)', 'excel_import', '1001109103', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'VICKY SETIA GUNAWAN', 'vicky-setia-gunawan', 'Lektor', 'S2, S.Kom,M.Kom', 'Bisnis Digital (S1)', 'excel_import', '1027089701', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'YOSI KURNIA', 'yosi-kurnia', 'Asisten Ahli', 'S2, S.E,M.Si', 'Bisnis Digital (S1)', 'excel_import', '1023116901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'B A MARTINUS', 'b-a-martinus', 'Lektor Kepala', 'S2, M.Si,Drs', 'Farmasi (S1)', 'excel_import', '0031056005', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'DEDI NOFIANDI', 'dedi-nofiandi', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1009117903', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'DIANA AGUSTIN', 'diana-agustin', '', 'S2, M.M,Apt,S.S.I', 'Farmasi (S1)', 'excel_import', '1017087702', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'DIZA SARTIKA', 'diza-sartika', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1024049006', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'ELMITRA', 'elmitra', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1025108501', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'EPI SUPRI WARDI', 'epi-supri-wardi', 'Lektor', 'S3, S.Si,M.Si,Dr', 'Farmasi (S1)', 'excel_import', '1031058901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'FARIDA RAHIM', 'farida-rahim', 'Lektor', 'S3, S.Si,Dr,M.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1015027802', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'HAZLI NURDIN', 'hazli-nurdin', 'Profesor', 'S3, Drs,Dr,M.Sc.', 'Farmasi (S1)', 'excel_import', '8868000016', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'IFMAILY', 'ifmaily', 'Lektor', 'S3, S.Si,M.Kes,Apt,Dr.', 'Farmasi (S1)', 'excel_import', '1001057202', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'INTAN KUMALASARI', 'intan-kumalasari', '', 'S2, M.Kes,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', NULL, true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'IRWANDI', 'irwandi', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1028088202', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'LOLA AZYENELA', 'lola-azyenela', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1022048703', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'MEGA YULIA', 'mega-yulia', 'Lektor', 'S3, S.H.,Dr,M.Farm,M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1012078703', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'MUTHIA MIRANDA ZAUNIT', 'muthia-miranda-zaunit', 'Lektor', 'S2, S.Pd,M.Si', 'Farmasi (S1)', 'excel_import', '1023119201', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'NESSA', 'nessa', 'Lektor', 'S2, S.Farm,Apt,M.Biomed', 'Farmasi (S1)', 'excel_import', '1015029102', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'NONI RAHAYU PUTRI', 'noni-rahayu-putri', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1023118703', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'PUSPA PAMESWARI', 'puspa-pameswari', 'Lektor', 'S2, M.Farm,M.Farm', 'Farmasi (S1)', 'excel_import', '1015028702', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'PUTRI DIAN AFRINDA', 'putri-dian-afrinda', 'Lektor', 'S2, S.Pd,M.Pd', 'Farmasi (S1)', 'excel_import', '1021048602', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'RIA AFRIANTI', 'ria-afrianti', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1005048101', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'SANDRA TRI JULI FENDRI', 'sandra-tri-juli-fendri', 'Lektor', 'S2, S.Si,M.Si', 'Farmasi (S1)', 'excel_import', '1005078902', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'TISA MANDALA SARI', 'tisa-mandala-sari', 'Lektor', 'S2, S.Pd,M.Si', 'Farmasi (S1)', 'excel_import', '1010088901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'VERAWATI', 'verawati', 'Lektor', 'S3, Dr,M.Farm,S.Farm,Apt,Dr', 'Farmasi (S1)', 'excel_import', '0007028103', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'WIDYASTUTI', 'widyastuti', 'Lektor Kepala', 'S3, S.Si,Dr,M.Farm,Apt', 'Farmasi (S1)', 'excel_import', '1012047502', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'ZULKARNI R', 'zulkarni-r', 'Lektor', 'S3, S.Si,Dr,M.M', 'Farmasi (S1)', 'excel_import', '1030126802', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ALYA MISDHAL RINI', 'alya-misdhal-rini', 'Asisten Ahli', 'S2, S.Gz,M.Biomed', 'Gizi (D3)', 'excel_import', '1001017604', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'HENDRA MUKHLIS', 'hendra-mukhlis', 'Lektor', 'S2, S.E.,M.Pd', 'Gizi (D3)', 'excel_import', '1029036701', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MARIA NOVA', 'maria-nova', 'Lektor', 'S2, S.KM,M.Kes', 'Gizi (D3)', 'excel_import', '1023118301', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'RISYA AHRIYASNA', 'risya-ahriyasna', 'Lektor', 'S2, M.Gz,S.Gz,Dietisien', 'Gizi (D3)', 'excel_import', '1016119201', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'SEPNI ASMIRA', 'sepni-asmira', 'Lektor', 'S2, S.TP,M.P', 'Gizi (D3)', 'excel_import', '1024097801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'TIKA DWITA ADFAR', 'tika-dwita-adfar', 'Lektor', 'S2, A.Md,S.ST,M.Biomed,Dietisien', 'Gizi (D3)', 'excel_import', '1018039001', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'YENSASNIDAR', 'yensasnidar', 'Lektor', 'S2, M.Pd,S.Gz,Dietisien', 'Gizi (D3)', 'excel_import', '1016076701', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'DEZI ILHAM', 'dezi-ilham', 'Asisten Ahli', 'S2, S.Pd,M.Biomed', 'Gizi (S1)', 'excel_import', '1014128901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'DWI SYAPUTRI YANTI', 'dwi-syaputri-yanti', '', 'S2, M.Gizi,S.Tr Gz', 'Gizi (S1)', 'excel_import', NULL, true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ERINA MASRI', 'erina-masri', 'Lektor', 'S2, M.Biomed,S.K.M.', 'Gizi (S1)', 'excel_import', '0007028204', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'HARLENI', 'harleni', 'Lektor', 'S2, S.Pd,M.Pd', 'Gizi (S1)', 'excel_import', '1004118801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'NURHAMIDAH', 'nurhamidah', 'Lektor', 'S2, S.KM,M.Biomed', 'Gizi (S1)', 'excel_import', '1020037701', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'RAHMITA YANTI', 'rahmita-yanti', 'Lektor', 'S2, S.KM,M.Kes', 'Gizi (S1)', 'excel_import', '1026098302', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'SRI INDRAYATI', 'sri-indrayati', 'Lektor', 'S2, S.Si,M.Si', 'Gizi (S1)', 'excel_import', '1012128901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'WIDIA DARA', 'widia-dara', 'Lektor', 'S2, S.P,M.P.', 'Gizi (S1)', 'excel_import', '1001026801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'WILDA LAILA', 'wilda-laila', 'Lektor', 'S2, M.Biomed,S.K.M.,Dietisien', 'Gizi (S1)', 'excel_import', '1017108302', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'YEYEP NATRIO', 'yeyep-natrio', 'Lektor', 'S2, S.S.,S.S.,M.Hum,M.Hum', 'Gizi (S1)', 'excel_import', '1024018903', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'AFRINUR ZAQIA', 'afrinur-zaqia', '', 'S2, S.Sos,M.I.Kom', 'Ilmu Komunikasi (S1)', 'excel_import', NULL, true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'ANNISA WERIFRAMAYENI', 'annisa-weriframayeni', 'Lektor', 'S2, S.I.Kom,M.I.Kom', 'Ilmu Komunikasi (S1)', 'excel_import', '1029039603', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'DANI PRAYOGA', 'dani-prayoga', 'Asisten Ahli', 'S2, M.I.Kom', 'Ilmu Komunikasi (S1)', 'excel_import', '1027109801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'DELPA', 'delpa', 'Lektor', 'S3, S.S,M.Soc.Sc.,PhD', 'Ilmu Komunikasi (S1)', 'excel_import', '1013067902', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'EDA ELYSIA', 'eda-elysia', 'Asisten Ahli', 'S2, A.Md,S.I.Kom,M.I.Kom', 'Ilmu Komunikasi (S1)', 'excel_import', '1008118801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'SHAFIRA ALISYA PUTRI MAULANA', 'shafira-alisya-putri-maulana', '', 'S2, S.I.Kom,S.I.Kom', 'Ilmu Komunikasi (S1)', 'excel_import', '1009079602', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'WAHYU FITRI', 'wahyu-fitri', 'Asisten Ahli', 'S2, S.Ds,M.I.Kom', 'Ilmu Komunikasi (S1)', 'excel_import', '1027039303', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('be39c54f-f925-4191-ba68-734604612e76', 'YOLANDA OKTARINA', 'yolanda-oktarina', 'Lektor', 'S2, S.Psi,M.M.', 'Ilmu Komunikasi (S1)', 'excel_import', '1007108803', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'AFRINITA EKA FITRI', 'afrinita-eka-fitri', '', 'S2, M.Keb,S.Tr. Keb', 'Kebidanan (D3)', 'excel_import', '0106049402', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'AMELYA PERMATA SARI', 'amelya-permata-sari', 'Asisten Ahli', 'S2, S.Tr.Keb,M.Keb', 'Kebidanan (D3)', 'excel_import', '1003109001', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ATHICA OVIANA', 'athica-oviana', 'Lektor', 'S2, S.ST,M.Keb', 'Kebidanan (D3)', 'excel_import', '1011118904', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'DIVA YULFERINA', 'diva-yulferina', 'Asisten Ahli', 'S2, S.ST,A.Md.Keb,M.Bio(Biomed)', 'Kebidanan (D3)', 'excel_import', '1029118901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'FENY WARTISA', 'feny-wartisa', 'Lektor', 'S3, M.KM,S.SiT,Doctor of Philosophy', 'Kebidanan (D3)', 'excel_import', '1018028801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MERI RAHMAYUNI', 'meri-rahmayuni', '', 'S2, S.ST,M.Biomed', 'Kebidanan (D3)', 'excel_import', NULL, true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'OKTI SATRIA', 'okti-satria', 'Lektor', 'S2, S.ST,S.ST,M.Keb', 'Kebidanan (D3)', 'excel_import', '1025108705', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'SISKA SYAFPUTRI', 'siska-syafputri', '', 'S2, S.KM,M.Biomed', 'Kebidanan (D3)', 'excel_import', '1027068502', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'SRI RAMADHANI FITRI', 'sri-ramadhani-fitri', 'Asisten Ahli', 'S2, S.Tr.Keb,MKM', 'Kebidanan (D3)', 'excel_import', '1027019801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'TRIVENI', 'triveni', 'Lektor', 'S2, S.ST,M.KM', 'Kebidanan (D3)', 'excel_import', '1017128902', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'WIENDHA SARI', 'wiendha-sari', '', 'S2, S.ST,M.Biomed', 'Kebidanan (D3)', 'excel_import', '1012018903', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'WIRA MEIRIZA', 'wira-meiriza', 'Lektor', 'S2, S.ST,M.Keb', 'Kebidanan (D3)', 'excel_import', '1003018901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'YENDA HASNITA', 'yenda-hasnita', 'Lektor', 'S2, S.Tr.Keb,M.Keb,M.Keb,S.Tr.Keb.', 'Kebidanan (D3)', 'excel_import', '1020039303', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ENDRA AMALIA', 'endra-amalia', 'Lektor', 'S2, S.Kep,M.Kep,Ners', 'Keperawatan (D3)', 'excel_import', '1023106901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'KALPANA KARTIKA', 'kalpana-kartika', 'Lektor', 'S2, S.Kep,Ners,M.Si', 'Keperawatan (D3)', 'excel_import', '1015108001', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'LISA FRADISA', 'lisa-fradisa', 'Lektor', 'S2, S.Si,M.Pd', 'Keperawatan (D3)', 'excel_import', '1006038402', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MAIDALIZA', 'maidaliza', 'Lektor', 'S2, S.Kep,M.Kep,Ners', 'Keperawatan (D3)', 'excel_import', '1019058005', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MARTA SURI', 'marta-suri', 'Lektor', 'S2, S.Kep,M.Kep', 'Keperawatan (D3)', 'excel_import', '1024038604', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MUHAMMAD ARIF', 'muhammad-arif', 'Lektor', 'S2, S.Kep,M.Kep,Ners', 'Keperawatan (D3)', 'excel_import', '1014098402', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'VERA SESRIANTY', 'vera-sesrianty', 'Lektor', 'S2, S.Kep,M.Kep,Ners', 'Keperawatan (D3)', 'excel_import', '1002117801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'DIA RESTI DEWI NANDA DEMUR', 'dia-resti-dewi-nanda-demur', 'Lektor', 'S2, S.Kep,M.Kep,Ners', 'Keperawatan (S1)', 'excel_import', '1008028602', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'IDA SURYATI', 'ida-suryati', 'Lektor', 'S2, S.Kep,M.Kep,Ners', 'Keperawatan (S1)', 'excel_import', '1030047503', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'LILISA MURNI', 'lilisa-murni', 'Asisten Ahli', 'S2, M.Pd,Dra', 'Keperawatan (S1)', 'excel_import', '1012106401', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'LISA MUSTIKA SARI', 'lisa-mustika-sari', 'Lektor', 'S2, S.Kep,M.Kep,Ners', 'Keperawatan (S1)', 'excel_import', '1014098501', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'TANTI ANGGREINIBOTI', 'tanti-anggreiniboti', 'Lektor', 'S2, S.Kep,MM', 'Keperawatan (S1)', 'excel_import', '1006067503', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'YASLINA', 'yaslina', 'Lektor', 'S3, Dr,S.Kep,M.Kep,Ners,Sp.Kep.Kom', 'Keperawatan (S1)', 'excel_import', '1006037301', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'YENDRIZAL JAFRI', 'yendrizal-jafri', 'Lektor', 'S2, M.Biomed', 'Keperawatan (S1)', 'excel_import', '1006116801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'EKA FITRIANDA', 'eka-fitrianda', 'Lektor', 'S3, M.Farm,S.Farm,Apt,Dr.', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '0023088004', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'JUNI FITRAH', 'juni-fitrah', '', 'S2, S.Si,M.Farm,Apt', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '8836033420', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'META EMILIA SURYA DHARMA', 'meta-emilia-surya-dharma', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '1012109202', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'MIMI ARIA', 'mimi-aria', 'Lektor', 'S2, M.Farm,S.Farm,Apt', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '1001078201', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'OKTA FERA', 'okta-fera', 'Lektor', 'S2, S.Si,M.Farm,Apt', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '1006107302', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'REVI YENTI', 'revi-yenti', 'Lektor', 'S2, S.Si,M.Si,Apt', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '0403027601', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'RINO WAHYUDI', 'rino-wahyudi', '', 'S2, S.Si,Apt,M.Farm.Klin.', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '8803133420', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'ROSLINDA RASYID', 'roslinda-rasyid', 'Lektor Kepala', 'S3, M.Si,Dra,Dra,Dr,Dr,Apt,Apt,MSi', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '8908630021', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'SANUBARI RELA TOBAT', 'sanubari-rela-tobat', 'Lektor', 'S3, Dr,M.Farm,S.Farm,Apt', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '1003028502', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('bcb407b8-2d79-48dd-bec4-24f99ba1bc20', 'SUHATRI', 'suhatri', 'Lektor Kepala', 'S3, Dra,Dra,Dr,Apt,Apt,M.S,MS,Dr', 'Pendidikan Profesi Apoteker (Profesi)', 'excel_import', '8949430021', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ALDO YULIANO MAS PUTRA', 'aldo-yuliano-mas-putra', 'Asisten Ahli', 'S2, S.Kep,M.M', 'Pendidikan Profesi Ners (Profesi)', 'excel_import', '1020078501', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ANDRYE FERNANDES', 'andrye-fernandes', 'Lektor', 'S2, S.Kep,M.Kep,Ners,Sp.Kep.An', 'Pendidikan Profesi Ners (Profesi)', 'excel_import', '1015079002', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'FALERI SISKAYUNERE', 'faleri-siskayunere', 'Lektor', 'S2, S.Kep,M.Kep', 'Pendidikan Profesi Ners (Profesi)', 'excel_import', '1025028003', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MERA DELIMA', 'mera-delima', 'Lektor', 'S2, S.Kep,S.Kep,M.Kep,Ners', 'Pendidikan Profesi Ners (Profesi)', 'excel_import', '1001107202', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'RINAWATI KASRIN', 'rinawati-kasrin', 'Lektor', 'S2, S.Kep,M.Kep', 'Pendidikan Profesi Ners (Profesi)', 'excel_import', '1005057005', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'YESSI ANDRIANI', 'yessi-andriani', 'Lektor', 'S2, S.Kep,M.Kep,Ners,Sp.Kep.Mat', 'Pendidikan Profesi Ners (Profesi)', 'excel_import', '1016078603', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ALI ASMUL', 'ali-asmul', 'Asisten Ahli', 'S2, S.Pd.I,M.Pd', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1007098705', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ENDANG SURIANI', 'endang-suriani', 'Lektor', 'S2, M.Kes', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1005107604', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MARISA', 'marisa', 'Lektor', 'S2, S.Pd,M.Pd', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1003038601', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'NOVA MUSTIKA', 'nova-mustika', 'Lektor', 'S2, S.Pd,M.Pd', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1006118801', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'PUTRA RAHMADEA UTAMI', 'putra-rahmadea-utami', 'Lektor', 'S2, S.Si,M.Biomed', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1017019001', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'RENOWATI', 'renowati', 'Lektor', 'S2, M.Biomed', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1001077301', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'RINDA LESTARI', 'rinda-lestari', 'Lektor', 'S2, S.Pd,M.Pd', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1012037604', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'RITA PERMATASARI', 'rita-permatasari', 'Lektor', 'S2, S.ST,M.Biotek.', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1013039302', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'SURAINI', 'suraini', 'Lektor', 'S2, M.Si,Dra', 'Teknologi Laboratorium Medis (D3)', 'excel_import', '1020116503', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'ANGGUN SOPHIA', 'anggun-sophia', 'Lektor', 'S2, S.Pd,M.Pd', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1005079301', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'BETTI ROSITA', 'betti-rosita', 'Lektor', 'S2, M.Si', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1004128001', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'CHAIRANI', 'chairani', 'Lektor', 'S2, M.Biomed,S.S.T', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1016128401', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'DEF PRIMAL', 'def-primal', 'Lektor Kepala', 'S2, S.Kep,M.Biomed', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1026128401', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'DEWI YUDIANA SHINTA', 'dewi-yudiana-shinta', 'Lektor', 'S3, S.Si,M.Si,Apt,Dr', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1016017602', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'IKHWAN RESMALA SUDJI', 'ikhwan-resmala-sudji', 'Lektor', 'S3, S.Si,M.Si,Dr.rer.nat.', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1023097901', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MELLY SISKA SURYANI', 'melly-siska-suryani', 'Lektor', 'S2, M.Hum,S.S', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1003088203', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MERI WULANDARI', 'meri-wulandari', 'Asisten Ahli', 'S2, S.ST,M.Biotek.', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1014099201', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'MUHAMMAD DIKI JULIANDI', 'muhammad-diki-juliandi', 'Lektor', 'S2, S.ST,M.Biotek.', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1010079501', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'SUDIYANTO', 'sudiyanto', '', 'S2, S.E.,M.P.h', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '8968001024', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

INSERT INTO faculty_lecturers (faculty_id, name, slug, position, education, expertise, source, pddikti_id, active, sort_order)
VALUES ('00b491d9-2adb-4436-a368-68f6410def65', 'VETRA SUSANTO', 'vetra-susanto', 'Lektor', 'S2, S.S.T,M.K.M.', 'Teknologi Laboratorium Medis (D4)', 'excel_import', '1008098101', true, 0)
ON CONFLICT (pddikti_id) DO UPDATE SET
  name = EXCLUDED.name,
  position = EXCLUDED.position,
  education = EXCLUDED.education,
  expertise = EXCLUDED.expertise,
  faculty_id = EXCLUDED.faculty_id;

COMMIT;
