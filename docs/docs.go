// Package docs Product API.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/products": {
            "get": {
                "description": "Returns a list of products with optional pagination",
                "produces": ["application/json"],
                "summary": "List products",
                "operationId": "getAll",
                "parameters": [
                    {"type": "integer", "default": 100, "description": "Max items to return", "name": "limit", "in": "query"},
                    {"type": "integer", "default": 0, "description": "Offset", "name": "offset", "in": "query"}
                ],
                "responses": {
                    "200": {"description": "OK", "schema": {"type": "array", "items": {"$ref": "#/definitions/Product"}}},
                    "500": {"description": "Internal error", "schema": {"$ref": "#/definitions/APIError"}}
                }
            },
            "post": {
                "description": "Create a new product",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "summary": "Create product",
                "operationId": "create",
                "parameters": [
                    {"description": "Product body", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/ProductInput"}}
                ],
                "responses": {
                    "201": {"description": "Created", "schema": {"$ref": "#/definitions/Product"}},
                    "400": {"description": "Validation error", "schema": {"$ref": "#/definitions/APIError"}},
                    "500": {"description": "Internal error", "schema": {"$ref": "#/definitions/APIError"}}
                }
            }
        },
        "/products/{id}": {
            "get": {
                "description": "Returns a product by ID",
                "produces": ["application/json"],
                "summary": "Get product by ID",
                "operationId": "getByID",
                "parameters": [
                    {"type": "integer", "description": "Product ID", "name": "id", "in": "path", "required": true}
                ],
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/Product"}},
                    "400": {"description": "Invalid ID", "schema": {"$ref": "#/definitions/APIError"}},
                    "404": {"description": "Not found", "schema": {"$ref": "#/definitions/APIError"}},
                    "500": {"description": "Internal error", "schema": {"$ref": "#/definitions/APIError"}}
                }
            },
            "put": {
                "description": "Update a product",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "summary": "Update product",
                "operationId": "update",
                "parameters": [
                    {"type": "integer", "description": "Product ID", "name": "id", "in": "path", "required": true},
                    {"description": "Product body", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/ProductInput"}}
                ],
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/Product"}},
                    "400": {"description": "Validation error", "schema": {"$ref": "#/definitions/APIError"}},
                    "404": {"description": "Not found", "schema": {"$ref": "#/definitions/APIError"}},
                    "500": {"description": "Internal error", "schema": {"$ref": "#/definitions/APIError"}}
                }
            },
            "delete": {
                "description": "Delete a product",
                "summary": "Delete product",
                "operationId": "delete",
                "parameters": [
                    {"type": "integer", "description": "Product ID", "name": "id", "in": "path", "required": true}
                ],
                "responses": {
                    "204": {"description": "No Content"},
                    "400": {"description": "Invalid ID", "schema": {"$ref": "#/definitions/APIError"}},
                    "404": {"description": "Not found", "schema": {"$ref": "#/definitions/APIError"}},
                    "500": {"description": "Internal error", "schema": {"$ref": "#/definitions/APIError"}}
                }
            }
        }
    },
    "definitions": {
        "Product": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "name": {"type": "string"},
                "description": {"type": "string"},
                "price": {"type": "integer", "description": "Price in minor units (e.g. cents)"}
            }
        },
        "ProductInput": {
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {"type": "string", "maxLength": 500},
                "description": {"type": "string", "maxLength": 2000},
                "price": {"type": "integer", "minimum": 0}
            }
        },
        "APIError": {
            "type": "object",
            "properties": {
                "code": {"type": "string"},
                "message": {"type": "string"}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Product API",
	Description:      "REST API for product management",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
