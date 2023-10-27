package handlers

import (
	"log"
	"rest-service/db"
	"rest-service/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func People(c *fiber.Ctx) error {
    log.Println("Received GET /people")

    var data map[string]string

    if err := c.BodyParser(&data); err != nil {
        log.Println(err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Could not parse the request",
        })
    }

    filter := models.Person{}

    if data["name"] != "" {
        filter.Name = data["name"]
    }
    if data["surname"] != "" {
        filter.Name = data["surname"]
    }
    if data["patronymic"] != "" {
        filter.Name = data["patronymic"]
    }

    var people []models.Person
    db.DB.Where(filter).Find(&people)

    return c.JSON(people)
}



func DeleteId(c *fiber.Ctx) error {
    log.Println("Received DELETE /people/delete")

    var data map[string]string

    if err := c.BodyParser(&data); err != nil {
        log.Println(err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Could not parse the request",
        })
    }

    id, err := strconv.Atoi(data["id"])
    if err != nil {
        log.Println("Couldn't parse the id to delete")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Could not parse the request",
        })
    }

    result := db.DB.Delete(&models.Person{}, id)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not delete",
        })
    }

    if result.RowsAffected == 0 {
        return c.Status(fiber.StatusNotModified).JSON(fiber.Map{
            "error": "No such record to delete",
        })
    }

    return c.JSON("Id " + strconv.Itoa(id) + " deleted")
}



func Add (c *fiber.Ctx) error {
    log.Println("Received POST /people/add")

    var data map[string]string

    if err := c.BodyParser(&data); err != nil {
        log.Println(err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Could not parse the request",
        })
    }

    // Check if a record with the same data already exists
    existingPerson := models.Person{}
    result := db.DB.Where(models.Person{Name: data["name"], Surname: data["surname"], Patronymic: data["patronymic"]}).First(&existingPerson)
    if result.Error == nil {
        // A matching record already exists, return an error or a message
        log.Println("Trying to add existing record, skipping")
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "error": "Record already exists",
        })
    }
    
    person := models.Person {
        Name: data["name"],
        Surname: data["surname"],
        Patronymic: data["patronymic"],
    }

    if person.Name == "" || person.Surname == "" {
        log.Println("Name or surname are missing.")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Name or surname are missing",
        })
    }

    err := enrich(&person)
    if err != nil {
        log.Println(err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not obtain information from external API",
        })
    }
    
    db.DB.Create(&person)

    log.Println("Added new record to DB.")
    return c.JSON(person)
}



func Modify (c *fiber.Ctx) error {
    log.Println("Received PUT /people/modify")

    var data map[string]string

    if err := c.BodyParser(&data); err != nil {
        log.Println(err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Could not parse the request",
        })
    }

    id, err := strconv.Atoi(data["id"])
    if err != nil {
        log.Println("Couldn't parse the id to delete")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Could not parse the request",
        })
    }

    person := models.Person{}

    result := db.DB.First(&person, id)
    if result.Error != nil {
        log.Println("Could not find the record to modify")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Could not find the record to modify",
        })
    }

    newData := map[string]interface{}{}

    if (data["name"] != "") {
        newData["name"] = data["name"]
    }
    if (data["surname"] != "") {
        newData["surname"] = data["surname"]
    }
    if (data["patronymic"] != "") {
        newData["patronymic"] = data["patronymic"]
    }
    if (data["age"] != "") {
        newData["age"] = data["age"]
    }
    if (data["gender"] != "") {
        newData["gender"] = data["gender"]
    }
    if (data["country"] != "") {
        newData["country"] = data["country"]
    }

    result = db.DB.Model(&person).Updates(newData)
    if result.Error != nil {
        log.Println("Error updating the record")
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Error updating the record",
        })
    }

    return c.JSON("Record Modified")
}
