package repository

import (
	"context"
	"time"

	"alumniproject/app/models/postgresql"
	"alumniproject/database/postgresql"
)

// GetAlumniByStatusPekerjaan retrieves alumni filtered by job status with more than 1 year of work
func GetAlumniByStatusPekerjaan(status string) ([]models.AlumniPekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := postgresql.DB.QueryContext(ctx, `
		SELECT 
			a.id, a.nama, a.jurusan, a.angkatan,
			p.bidang_industri, p.nama_perusahaan, p.posisi_jabatan,
			p.tanggal_mulai_kerja, p.gaji_range, p.status_pekerjaan,
			COUNT(*) OVER() AS total_bekerja_lebih_1_tahun
		FROM alumni a
		JOIN pekerjaan_alumni p ON a.id = p.alumni_id
		WHERE p.status_pekerjaan = $1
			AND AGE(CURRENT_DATE, p.tanggal_mulai_kerja) > INTERVAL '1 year'
		ORDER BY a.id
	`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AlumniPekerjaan
	var totalCount int
	for rows.Next() {
		var ap models.AlumniPekerjaan
		err := rows.Scan(
			&ap.ID, &ap.Nama, &ap.Jurusan, &ap.Angkatan,
			&ap.BidangIndustri, &ap.NamaPerusahaan, &ap.PosisiJabatan,
			&ap.TanggalMulaiKerja, &ap.GajiRange, &ap.StatusPekerjaan,
			&totalCount,
		)
		if err != nil {
			return nil, err
		}
		ap.TotalBekerjaLebih1Tahun = totalCount
		list = append(list, ap)
	}

	return list, nil
}

// GetAlumniWithLongTermJobs retrieves alumni with active jobs lasting more than 1 year
func GetAlumniWithLongTermJobs() ([]models.AlumniPekerjaan, error) {
	return GetAlumniByStatusPekerjaan("aktif")
}