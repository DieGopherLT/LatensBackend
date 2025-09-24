package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// Templates para cada archivo
const repositoryInterfaceTemplate = `
// {{.PascalCase}}Repository interface
type {{.PascalCase}}Repository interface {
	Create(ctx context.Context, {{.CamelCase}} *models.{{.PascalCase}}) error
	FindByID(ctx context.Context, id string) (*models.{{.PascalCase}}, error)
	FindAll(ctx context.Context) ([]*models.{{.PascalCase}}, error)
	Update(ctx context.Context, id string, update map[string]any) error
	Delete(ctx context.Context, id string) error
}
`

const repositoryImplTemplate = `package repository

import (
	"context"
	"errors"

	"github.com/DieGopherLT/mfc_backend/internal/database/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Mongo{{.PascalCase}}Repository struct {
	collection *mongo.Collection
}

func New{{.PascalCase}}Repository(client *mongo.Database) *Mongo{{.PascalCase}}Repository {
	collection := client.Collection("{{.SnakeCase}}")
	return &Mongo{{.PascalCase}}Repository{collection: collection}
}

func (r *Mongo{{.PascalCase}}Repository) Create(ctx context.Context, {{.CamelCase}} *models.{{.PascalCase}}) error {
	_, err := r.collection.InsertOne(ctx, {{.CamelCase}})
	return err
}

func (r *Mongo{{.PascalCase}}Repository) FindByID(ctx context.Context, id string) (*models.{{.PascalCase}}, error) {
	var {{.CamelCase}} models.{{.PascalCase}}

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	condition := map[string]any{"_id": objectID}
	err = r.collection.FindOne(ctx, condition).Decode(&{{.CamelCase}})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &{{.CamelCase}}, nil
}

func (r *Mongo{{.PascalCase}}Repository) FindAll(ctx context.Context) ([]*models.{{.PascalCase}}, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var {{.CamelCase}}s []*models.{{.PascalCase}}
	for cursor.Next(ctx) {
		var {{.CamelCase}} models.{{.PascalCase}}
		if err := cursor.Decode(&{{.CamelCase}}); err != nil {
			return nil, err
		}
		{{.CamelCase}}s = append({{.CamelCase}}s, &{{.CamelCase}})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return {{.CamelCase}}s, nil
}

func (r *Mongo{{.PascalCase}}Repository) Update(ctx context.Context, id string, update map[string]any) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	condition := map[string]any{"_id": objectID}
	_, err = r.collection.UpdateOne(ctx, condition, map[string]any{"$set": update})
	return err
}

func (r *Mongo{{.PascalCase}}Repository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	condition := map[string]any{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, condition)
	return err
}
`

const serviceTemplate = `package {{.KebabCase}}

import (
	"context"

	"github.com/DieGopherLT/mfc_backend/internal/database/models"
	"github.com/DieGopherLT/mfc_backend/internal/database/repository"
)

type {{.PascalCase}}Service struct {
	repo repository.{{.PascalCase}}Repository
}

func New{{.PascalCase}}Service(repo repository.{{.PascalCase}}Repository) *{{.PascalCase}}Service {
	return &{{.PascalCase}}Service{repo: repo}
}

func (s *{{.PascalCase}}Service) Create{{.PascalCase}}(ctx context.Context, {{.CamelCase}} *models.{{.PascalCase}}) error {
	return s.repo.Create(ctx, {{.CamelCase}})
}

func (s *{{.PascalCase}}Service) Get{{.PascalCase}}ByID(ctx context.Context, id string) (*models.{{.PascalCase}}, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *{{.PascalCase}}Service) GetAll{{.PascalCase}}s(ctx context.Context) ([]*models.{{.PascalCase}}, error) {
	return s.repo.FindAll(ctx)
}

func (s *{{.PascalCase}}Service) Update{{.PascalCase}}(ctx context.Context, id string, update map[string]any) error {
	return s.repo.Update(ctx, id, update)
}

func (s *{{.PascalCase}}Service) Delete{{.PascalCase}}(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
`

const controllerTemplate = `package controller

import (
	"github.com/DieGopherLT/mfc_backend/internal/database/models"
	"github.com/DieGopherLT/mfc_backend/internal/services/{{.KebabCase}}"
	"github.com/gofiber/fiber/v2"
)

type {{.PascalCase}}Handler struct {
	{{.CamelCase}}Service *{{.KebabCase}}.{{.PascalCase}}Service
}

func New{{.PascalCase}}Handler({{.CamelCase}}Service *{{.KebabCase}}.{{.PascalCase}}Service) *{{.PascalCase}}Handler {
	return &{{.PascalCase}}Handler{ {{.CamelCase}}Service: {{.CamelCase}}Service }
}

func (h *{{.PascalCase}}Handler) Create{{.PascalCase}}(c *fiber.Ctx) error {
	var {{.CamelCase}} models.{{.PascalCase}}

	if err := c.BodyParser(&{{.CamelCase}}); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := h.{{.CamelCase}}Service.Create{{.PascalCase}}(c.Context(), &{{.CamelCase}}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create {{.CamelCase}}",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "{{.PascalCase}} created successfully",
		"data": {{.CamelCase}},
	})
}

func (h *{{.PascalCase}}Handler) Get{{.PascalCase}}ByID(c *fiber.Ctx) error {
	id := c.Params("id")

	{{.CamelCase}}, err := h.{{.CamelCase}}Service.Get{{.PascalCase}}ByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch {{.CamelCase}}",
			"details": err.Error(),
		})
	}

	if {{.CamelCase}} == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "{{.PascalCase}} not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON({{.CamelCase}})
}

func (h *{{.PascalCase}}Handler) GetAll{{.PascalCase}}s(c *fiber.Ctx) error {
	{{.CamelCase}}s, err := h.{{.CamelCase}}Service.GetAll{{.PascalCase}}s(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch {{.CamelCase}}s",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON({{.CamelCase}}s)
}

func (h *{{.PascalCase}}Handler) Update{{.PascalCase}}(c *fiber.Ctx) error {
	id := c.Params("id")

	var update map[string]any
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := h.{{.CamelCase}}Service.Update{{.PascalCase}}(c.Context(), id, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update {{.CamelCase}}",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "{{.PascalCase}} updated successfully",
	})
}

func (h *{{.PascalCase}}Handler) Delete{{.PascalCase}}(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.{{.CamelCase}}Service.Delete{{.PascalCase}}(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete {{.CamelCase}}",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "{{.PascalCase}} deleted successfully",
	})
}
`

type ResourceData struct {
	PascalCase string // User, Product
	CamelCase  string // user, product
	KebabCase  string // user, product (for package names)
	SnakeCase  string // users, products (for collection names)
}

func toPascalCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func toSnakeCase(s string) string {
	return strings.ToLower(s) + "s"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate-resource.go <resource-name>")
		fmt.Println("Example: go run generate-resource.go Task")
		os.Exit(1)
	}

	resourceName := strings.TrimSpace(os.Args[1])
	if resourceName == "" {
		fmt.Println("Resource name cannot be empty")
		os.Exit(1)
	}

	data := ResourceData{
		PascalCase: toPascalCase(resourceName),
		CamelCase:  toCamelCase(resourceName),
		KebabCase:  strings.ToLower(resourceName),
		SnakeCase:  toSnakeCase(resourceName),
	}

	fmt.Printf("Generating resource: %s\n", data.PascalCase)

	// Create directories
	serviceDirPath := fmt.Sprintf("internal/services/%s", data.KebabCase)
	if err := os.MkdirAll(serviceDirPath, 0755); err != nil {
		fmt.Printf("Error creating service directory: %v\n", err)
		os.Exit(1)
	}

	// Generate files
	files := map[string]string{
		fmt.Sprintf("internal/database/repository/%s.go", data.KebabCase): repositoryImplTemplate,
		fmt.Sprintf("internal/services/%s/%s.go", data.KebabCase, data.KebabCase): serviceTemplate,
		fmt.Sprintf("internal/controller/%s_controller.go", data.KebabCase): controllerTemplate,
	}

	for filePath, templateContent := range files {
		if err := generateFile(filePath, templateContent, data); err != nil {
			fmt.Printf("Error generating %s: %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Generated: %s\n", filePath)
	}

	// Update interfaces.go
	if err := updateInterfaces(data); err != nil {
		fmt.Printf("Error updating interfaces: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Updated: internal/database/repository/interfaces.go\n")

	// Generate routes snippet
	generateRoutesSnippet(data)

	fmt.Printf("\n🎉 Resource '%s' generated successfully!\n", data.PascalCase)
	fmt.Println("\n📝 Don't forget to:")
	fmt.Printf("1. Add the model to internal/database/models/models.go\n")
	fmt.Printf("2. Add routes to cmd/api/routes.go (snippet generated below)\n")
	fmt.Printf("3. Import the service in routes.go\n")
}

func generateFile(filePath, templateContent string, data ResourceData) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Parse and execute template
	tmpl, err := template.New("resource").Parse(templateContent)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func updateInterfaces(data ResourceData) error {
	interfacePath := "internal/database/repository/interfaces.go"

	// Generate interface content
	tmpl, err := template.New("interface").Parse(repositoryInterfaceTemplate)
	if err != nil {
		return err
	}

	var interfaceContent strings.Builder
	if err := tmpl.Execute(&interfaceContent, data); err != nil {
		return err
	}

	// Append to file
	file, err := os.OpenFile(interfacePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(interfaceContent.String())
	return err
}

func generateRoutesSnippet(data ResourceData) {
	fmt.Printf("\n📋 Routes snippet for cmd/api/routes.go:\n")
	fmt.Printf("```go\n")
	fmt.Printf("// Add to imports:\n")
	fmt.Printf("\"github.com/DieGopherLT/mfc_backend/internal/services/%s\"\n\n", data.KebabCase)

	fmt.Printf("// Add to setupRoutes function:\n")
	fmt.Printf("%sRepo := repository.New%sRepository(db)\n", data.CamelCase, data.PascalCase)
	fmt.Printf("%sService := %s.New%sService(%sRepo)\n", data.CamelCase, data.KebabCase, data.PascalCase, data.CamelCase)
	fmt.Printf("%sHandler := controller.New%sHandler(%sService)\n\n", data.CamelCase, data.PascalCase, data.CamelCase)

	fmt.Printf("// Add routes:\n")
	fmt.Printf("%ss := v1.Group(\"/%ss\")\n", data.CamelCase, data.CamelCase)
	fmt.Printf("%ss.Post(\"/\", %sHandler.Create%s)\n", data.CamelCase, data.CamelCase, data.PascalCase)
	fmt.Printf("%ss.Get(\"/:id\", %sHandler.Get%sByID)\n", data.CamelCase, data.CamelCase, data.PascalCase)
	fmt.Printf("%ss.Get(\"/\", %sHandler.GetAll%ss)\n", data.CamelCase, data.CamelCase, data.PascalCase)
	fmt.Printf("%ss.Put(\"/:id\", %sHandler.Update%s)\n", data.CamelCase, data.CamelCase, data.PascalCase)
	fmt.Printf("%ss.Delete(\"/:id\", %sHandler.Delete%s)\n", data.CamelCase, data.CamelCase, data.PascalCase)
	fmt.Printf("```\n")
}