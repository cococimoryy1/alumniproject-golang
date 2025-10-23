package service

import (
	"strconv"
	"strings"
	"time"

	"alumniproject/app/models/postgresql"
	"alumniproject/app/repository/postgresql"
	"github.com/gofiber/fiber/v2"

)

func GetAllAlumni(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	list, err := repository.GetAllAlumni(role, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}


func GetAlumniByIDService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	data, err := repository.GetAlumniByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data alumni tidak ditemukan"})
	}

	return c.JSON(data)
}


func CreateAlumniService(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)


	var req models.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Semua field wajib diisi"})
	}

	alumni := models.Alumni{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
		CreatedBy:  userID,
	}

	if err := repository.CreateAlumni(&alumni); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data"})
	}

	return c.JSON(fiber.Map{
		"message": "Data alumni berhasil ditambahkan",
		"data":    alumni,
	})
}

func UpdateAlumniService(c *fiber.Ctx) error {
    id, _ := strconv.Atoi(c.Params("id"))
    userID := c.Locals("user_id").(int)
    role := c.Locals("role").(string)

    var req models.UpdateAlumniRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Input tidak valid",
        })
    }

    if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Nama, jurusan, dan email wajib diisi",
        })
    }

    // Ambil data dari repository
    data, err := repository.GetAlumniByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Data alumni tidak ditemukan",
        })
    }

    // Validasi akses
    if role != "admin" && data.CreatedBy != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Tidak boleh ubah data ini",
        })
    }

    // Update field
    data.Nama = req.Nama
    data.Jurusan = req.Jurusan
    data.Angkatan = req.Angkatan
    data.TahunLulus = req.TahunLulus
    data.Email = req.Email
    data.NoTelepon = req.NoTelepon
    data.Alamat = req.Alamat
    data.UpdatedAt = time.Now()

    // Simpan ke repository
    err = repository.UpdateAlumni(&data, userID, role)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Gagal memperbarui data alumni",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Data alumni berhasil diperbarui",
        "data":    data,
    })
}


func DeleteAlumniService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	data, err := repository.GetAlumniByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data alumni tidak ditemukan"})
	}

	if role != "admin" && data.CreatedBy != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh hapus data ini"})
	}

	err = repository.DeleteAlumni(id, userID, role)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data alumni berhasil dihapus"})
}
func RestoreAlumniService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	err = repository.RestoreAlumni(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data alumni berhasil dikembalikan"})
}



// func DeleteAlumniService(c *fiber.Ctx) error {
//     id, err := strconv.Atoi(c.Params("id"))
//     if err != nil {
//         return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
//     }

//     var userID int
//     switch v := c.Locals("user_id").(type) {
//     case int:
//         userID = v
//     case float64:
//         userID = int(v)
//     case string:
//         idInt, _ := strconv.Atoi(v)
//         userID = idInt
//     default:
//         return c.Status(401).JSON(fiber.Map{"error": "User ID tidak valid di token"})
//     }

//     role := c.Locals("role").(string)

//     // Pastikan data ada
//     data, err := repository.GetAlumniByID(id)
//     if err != nil {
//         return c.Status(404).JSON(fiber.Map{"error": "Data alumni tidak ditemukan"})
//     }

//     // Validasi akses
//     if role != "admin" && data.CreatedBy != userID {
//         return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh hapus data ini"})
//     }

//     // Jalankan delete (repository)
//     err = repository.DeleteAlumni(id, userID, role)
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{"error": err.Error()})
//     }

//     return c.JSON(fiber.Map{"message": "Data alumni berhasil dihapus (soft delete)"})
// }

// func RestoreAlumniService(c *fiber.Ctx) error {
//     id, err := strconv.Atoi(c.Params("id"))
//     if err != nil {
//         return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
//     }

//     // Panggil repository restore
//     err = repository.RestoreAlumni(id)
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{"error": err.Error()})
//     }

//     return c.JSON(fiber.Map{"message": "Data alumni berhasil direstore"})
// }




// GetAlumniPaginated -> ambil data alumni dengan pagination, sorting, dan search
func GetAlumniService(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "id")
    order := c.Query("order", "asc")
    search := c.Query("search", "")

    offset := (page - 1) * limit

    whitelist := map[string]bool{
        "id": true, "nim": true, "nama": true,
        "jurusan": true, "angkatan": true, "tahun_lulus": true,
    }
    if !whitelist[sortBy] {
        sortBy = "id"
    }
    if strings.ToLower(order) != "desc" {
        order = "asc"
    }

    alumni, err := repository.GetAlumniRepo(search, sortBy, order, limit, offset)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data alumni"})
    }

    total, err := repository.CountAlumniRepo(search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal hitung data alumni"})
    }

	response := &models.AlumniResponse{
		Data: alumni,
		Meta: &models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

    return c.JSON(response)
}