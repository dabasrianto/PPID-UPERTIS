package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	_ = godotenv.Load(".env")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "")
	if dbname == "" { log.Fatal("DB_NAME required") }
	sslmode := getEnv("DB_SSLMODE", "disable")

	db, _ := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode))
	defer db.Close()
	ctx := context.Background()

	fmt.Printf("=== Detail Rektor yang duplikat ===\n\n")
	rows, _ := db.QueryContext(ctx, `
		SELECT l.id, l.name, l.position, COALESCE(l.leadership_group,''), l.source, l.education, l.expertise, f.name as faculty
		FROM faculty_lecturers l
		JOIN faculties f ON f.id = l.faculty_id
		WHERE UPPER(l.name) IN ('YASLINA', 'YENDRIZAL JAFRI')
		ORDER BY l.name, l.source
	`)
	defer rows.Close()
	for rows.Next() {
		var id, name, pos, lg, src, edu, exp, fac string
		rows.Scan(&id, &name, &pos, &lg, &src, &edu, &exp, &fac)
		fmt.Printf("ID:     %s\n", id)
		fmt.Printf("Nama:   %s\n", name)
		fmt.Printf("Posisi: %s\n", pos)
		fmt.Printf("Group:  %s\n", lg)
		fmt.Printf("Source: %s\n", src)
		fmt.Printf("Edu:    %s\n", edu)
		fmt.Printf("Exp:    %s\n", exp)
		fmt.Printf("Fak:    %s\n", fac)
		fmt.Println("---")
	}

	// Delete the excel_import duplicates (keep pddikti ones which have leadership info)
	fmt.Println("\nMenghapus duplikat excel_import untuk rektor...")
	res, err := db.ExecContext(ctx, `
		DELETE FROM faculty_lecturers
		WHERE source = 'excel_import'
		AND pddikti_id IN ('1006037301', '1006116801')
	`)
	if err != nil {
		log.Fatal(err)
	}
	affected, _ := res.RowsAffected()
	fmt.Printf("✓ Dihapus: %d record duplikat dari excel_import\n", affected)

	// Find FK faculty ID dynamically
	var fkID string
	err = db.QueryRowContext(ctx, "SELECT id FROM faculties WHERE code = 'FK'").Scan(&fkID)
	if err != nil {
		log.Fatal("Failed to query FK faculty ID: ", err)
	}

	// Find the ID of the pddikti record for Yaslina dynamically
	var yaslinaID string
	err = db.QueryRowContext(ctx, "SELECT id FROM faculty_lecturers WHERE pddikti_id = '1006037301' OR (UPPER(name) LIKE '%YASLINA%' AND source = 'pddikti') LIMIT 1").Scan(&yaslinaID)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Failed to search Yaslina pddikti record: ", err)
	}

	// Update Yaslina's record if found
	if yaslinaID != "" {
		_, err = db.ExecContext(ctx, `
			UPDATE faculty_lecturers SET
				name = 'YASLINA, Dr,S.Kep,M.Kep,Ners,Sp.Kep.Kom',
				pddikti_id = '1006037301',
				faculty_id = $1,
				education = 'S3, Dr,S.Kep,M.Kep,Ners,Sp.Kep.Kom',
				expertise = 'Keperawatan (S1)'
			WHERE id = $2
		`, fkID, yaslinaID)
		if err != nil {
			log.Fatal("Failed to update Yaslina pddikti record: ", err)
		}
	} else {
		fmt.Println("⚠ Yaslina pddikti record not found to update")
	}

	// Find the ID of the pddikti record for Yendrizal Jafri dynamically
	var yendrizalID string
	err = db.QueryRowContext(ctx, "SELECT id FROM faculty_lecturers WHERE pddikti_id = '1006116801' OR (UPPER(name) LIKE '%YENDRIZAL JAFRI%' AND source = 'pddikti') LIMIT 1").Scan(&yendrizalID)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Failed to search Yendrizal Jafri pddikti record: ", err)
	}

	// Update Yendrizal Jafri's record if found
	if yendrizalID != "" {
		_, err = db.ExecContext(ctx, `
			UPDATE faculty_lecturers SET
				name = 'YENDRIZAL JAFRI, M.Biomed',
				pddikti_id = '1006116801',
				faculty_id = $1,
				education = 'S2, M.Biomed',
				expertise = 'Keperawatan (S1)'
			WHERE id = $2
		`, fkID, yendrizalID)
		if err != nil {
			log.Fatal("Failed to update Yendrizal Jafri pddikti record: ", err)
		}
	} else {
		fmt.Println("⚠ Yendrizal Jafri pddikti record not found to update")
	}
	fmt.Println("✓ Updated pddikti rector records with degrees, actual NIDN, and correct faculty")

	var total int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM faculty_lecturers").Scan(&total)
	fmt.Printf("✓ Total dosen sekarang: %d\n", total)
}
