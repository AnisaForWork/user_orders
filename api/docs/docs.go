// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/auth": {
            "post": {
                "description": "sing in user if they have given valid credentials, returns access token(JWT)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sing in user",
                "parameters": [
                    {
                        "description": "login,password",
                        "name": "auth",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.Auth"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            }
        },
        "/auth/reg": {
            "post": {
                "description": "singUp with credentials user : login,full name,email,password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "login,full name,email,password",
                        "name": "reg",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.Registration"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            }
        },
        "/product/": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user provides products barcode, name, description and cost",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Create new product",
                "parameters": [
                    {
                        "description": "barcode,name,desc,cost",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/product.Created"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            }
        },
        "/product/all": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "configured to return set number of products(barcodes+name+cost) per request, user sends configures how many products per page be seen and offset, only owner of product can view it",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Returns all users products with pagnation",
                "parameters": [
                    {
                        "maximum": 50000,
                        "minimum": 1,
                        "type": "integer",
                        "description": "Next page to retrieve",
                        "name": "p",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "description": "Number of products info per page",
                        "name": "n",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            }
        },
        "/product/check/{checkName}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "sends user PDF check generated previously using product info, only owner of product can do it",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json",
                    "application/pdf"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Returns check of user product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Check file name",
                        "name": "checkName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            }
        },
        "/product/{barcode}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "returns full info about product user chose to view, only owner of product can view it",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Returns user product full info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product barcode",
                        "name": "barcode",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "archives product but don't delete it from storage, only owner of product can do it",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Delete user product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product barcode",
                        "name": "barcode",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            }
        },
        "/product/{barcode}/check": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "using info about given product generates PDF check using special PDF template",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json",
                    "application/pdf"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Generates check",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product barcode",
                        "name": "barcode",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.JSONResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.Auth": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 3,
                    "example": "Login123"
                },
                "password": {
                    "type": "string",
                    "default": "password",
                    "maxLength": 40,
                    "minLength": 6
                }
            }
        },
        "auth.Registration": {
            "type": "object",
            "required": [
                "email",
                "fullName",
                "login",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@test.com"
                },
                "fullName": {
                    "type": "string",
                    "maxLength": 75,
                    "minLength": 3,
                    "example": "Ivanov Ivan Ivanovich"
                },
                "login": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 3,
                    "example": "Login123"
                },
                "password": {
                    "type": "string",
                    "maxLength": 40,
                    "minLength": 6,
                    "example": "password"
                }
            }
        },
        "product.Created": {
            "type": "object",
            "required": [
                "barcode",
                "cost",
                "desc",
                "name"
            ],
            "properties": {
                "barcode": {
                    "type": "string",
                    "default": "1234567890"
                },
                "cost": {
                    "type": "integer",
                    "default": 100,
                    "minimum": 10
                },
                "desc": {
                    "type": "string",
                    "default": "Description",
                    "maxLength": 1000,
                    "minLength": 10
                },
                "name": {
                    "type": "string",
                    "default": "product name",
                    "maxLength": 60,
                    "minLength": 10
                }
            }
        },
        "response.JSONResult": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "User products service API",
	Description:      "User products service API using swagger 2.0.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
